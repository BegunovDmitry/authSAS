package services

import (
	"authSAS/internal/utils"
	emailsender "authSAS/internal/utils/emailSender"
	"authSAS/internal/utils/jwt"
	"context"
	"errors"
	"log/slog"
	"time"

	utils_random "authSAS/internal/utils/randomCode"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type SessionService struct {
	logger *slog.Logger
	tokenTTL time.Duration
	jwtSecret string
	emailSender *emailsender.EmailSender
	userGetter UserGetter
	logoutJWTKeeper LogoutJWTKeeper
	twoFACodeKeeper TwoFACodeKeeper
	twoFACodeGetter TwoFACodeGetter
}

func NewSessionService(logger *slog.Logger, tokenTTL time.Duration, secret string, emailSender *emailsender.EmailSender, permanentStorage PermanentStorage, temporaryStorage TemporaryStorage) *SessionService {
	return &SessionService{
		logger: logger,
		tokenTTL: tokenTTL,
		jwtSecret: secret, 
		emailSender: emailSender,
		userGetter: permanentStorage,
		logoutJWTKeeper: permanentStorage,
		twoFACodeKeeper: temporaryStorage,
		twoFACodeGetter: temporaryStorage,
	}
}


func (s *SessionService) Login(ctx context.Context, email string, password string) (token string, msg string, err error) {

	s.logger.Debug("Trying to login user", "email", email)

	if email == "" {
		s.logger.Debug("User login error", "email", email, "err", utils.ErrEmptyEmail)
		return "", "Error", utils.ErrInvalidCredentials
	}

	if password == "" {
		s.logger.Debug("User login error", "email", email, "err", utils.ErrEmptyPassword)
		return "", "Error", utils.ErrInvalidCredentials
	}

	user, err := s.userGetter.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Debug("User login error", "email", email, "err", err.Error())
		if err == utils.ErrUserNotFound {
			return "", "Error", utils.ErrInvalidCredentials
		}
		return "", "Error", utils.ErrInternalServer
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.logger.Debug("User login error", "email", email, "err", "invalid password (not null)")
		return "", "Error", utils.ErrInvalidCredentials
	}

	if user.Use2FA {
		s.logger.Debug("Trying to send 2FA code", "email", email)

		randCode := utils_random.RandRange(1000, 9999)

		s.emailSender.SendEmail(email, randCode)

		if err = s.twoFACodeKeeper.KeepTwoFACode(ctx, email, randCode); err != nil {
			s.logger.Debug("Sending 2FA code error", "email", email, "err", err.Error())
			return "", "Error", utils.ErrInternalServer
		}

		s.logger.Debug("2FA code sended", "email", email, "code", randCode)

		return "", "2FA code sended", nil
	}

	token, err = utils_jwt.NewToken(user, s.tokenTTL, s.jwtSecret)
	if err != nil {
		s.logger.Debug("User login error", "email", email, "err", err.Error())
		return "", "Error", utils.ErrInternalServer
	}

	s.logger.Debug("User logined succesfully", "email", email)

	return token, "Authorized", nil
}

func (s *SessionService) Logout(ctx context.Context, tokenString string) (msg string, err error) {

	s.logger.Debug("Trying to logout user", "token", tokenString)

	if tokenString == "" {
		s.logger.Debug("Logout user error", "token", tokenString, "err", utils.ErrEmptyJWT)
		return "Error", utils.ErrInvalidCredentials
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		s.logger.Debug("Trying to logout user", "token", tokenString, "err", "invalid token")
		return "Error", utils.ErrInvalidCredentials
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Debug("Trying to logout user", "token", tokenString, "err", "invalid claims in token")
		return "Error", utils.ErrInvalidCredentials
	}
	uid := int64(claims["uid"].(float64))

	if err := s.logoutJWTKeeper.KeepLogoutJWT(ctx, uid, tokenString); err != nil {
		s.logger.Debug("Trying to logout user", "token", tokenString, "err", err.Error())
		if errors.Is(err, utils.ErrJWTAlreadyAdded) {
			return "Error", utils.ErrJWTAlreadyAdded
		}
		return "Error", utils.ErrInternalServer
	}

	s.logger.Debug("User logouted succesfully", "token", tokenString, "uid", uid)

	return "Success", nil
}

func (s *SessionService) LoginWith2FACode(ctx context.Context, email string, code int) (token string, err error) {

	s.logger.Debug("Trying to 2FA login user", "email", email, "code", code)

	if email == "" {
		s.logger.Debug("User 2FA login error", "email", email, "err", utils.ErrEmptyEmail)
		return "", utils.ErrInvalidCredentials
	}

	if code == 0 {
		s.logger.Debug("User 2FA login error", "email", email, "err", "null 2FA code")
		return "", utils.ErrInvalidCredentials
	}

	sendedCode, err := s.twoFACodeGetter.GetTwoFACode(ctx, email)
	if err != nil {
		s.logger.Debug("User 2FA login error", "email", email, "err", err.Error())
		if err == utils.Err2FACodeNotFound {
			return "", utils.ErrInvalidCredentials
		}
		return "", utils.ErrInternalServer
	}

	if sendedCode != code {
		s.logger.Debug("User 2FA login error", "email", email, "err", "invalid 2FA code")
		return "", utils.ErrInvalidCredentials
	}

	user, err := s.userGetter.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Debug("User 2FA login error", "email", email, "err", err.Error())
		if err == utils.ErrUserNotFound {
			return "", utils.ErrInvalidCredentials
		}
		return "", utils.ErrInternalServer
	}

	token, err = utils_jwt.NewToken(user, s.tokenTTL, s.jwtSecret)
	if err != nil {
		s.logger.Debug("User 2FA login error", "email", email, "err", err.Error())
		return "", utils.ErrInternalServer
	}

	s.logger.Debug("User logined with 2FA succesfully", "email", email)

	return token, nil

}