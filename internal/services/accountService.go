package services

import (
	"authSAS/internal/utils"
	emailsender "authSAS/internal/utils/emailSender"
	utils_random "authSAS/internal/utils/randomCode"
	"context"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	logger *slog.Logger
	tokenTTL time.Duration
	emailSender *emailsender.EmailSender
	userGetter UserGetter
	userCreator UserCreator
	emailVerificator 	EmailVerificator
	passChanger 	PassChanger
	emailVerifyCodeKeeper 	EmailVerifyCodeKeeper
	emailVerifyCodeGetter 	EmailVerifyCodeGetter
	passRecoverCodeKeeper 	PassRecoverCodeKeeper
	passRecoverCodeGetter 	PassRecoverCodeGetter
}

func NewAccountService(logger *slog.Logger, tokenTTL time.Duration, emailSender *emailsender.EmailSender, permanentStorage PermanentStorage, temporaryStorage TemporaryStorage) *AccountService {
	return &AccountService{
		logger: logger,
		tokenTTL: tokenTTL,
		emailSender: emailSender,
		userGetter: permanentStorage,
		userCreator: permanentStorage,
		emailVerificator: permanentStorage,
		passChanger: permanentStorage,
		emailVerifyCodeKeeper: temporaryStorage,
		emailVerifyCodeGetter: temporaryStorage,
		passRecoverCodeKeeper: temporaryStorage,
		passRecoverCodeGetter: temporaryStorage,
	}
}

func (a *AccountService) Register(ctx context.Context, email string, password string) (userId int64, err error) {
	
	a.logger.Debug("Trying to register user", "email", email)

	if email == "" {
		a.logger.Debug("Register user error", "email", email, "err", utils.ErrEmptyEmail)
		return 0, utils.ErrInvalidCredentials
	}

	if password == "" {
		a.logger.Debug("Register user error", "email", email, "err", utils.ErrEmptyPassword)
		return 0, utils.ErrInvalidCredentials
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Debug("Register user error", "email", email, "err", err.Error())
		return 0, utils.ErrInternalServer
	}

	userId, err = a.userCreator.CreateUser(ctx, email, passHash)
	if err != nil {
		a.logger.Debug("Register user error", "email", email, "err", err.Error())
		if err == utils.ErrUserAlreadyExists {
			return 0, err
		}
		return 0, utils.ErrInternalServer
	}

	a.logger.Debug("User registered", "email", email, "uid", userId)

	return userId, nil
}

func (a *AccountService) EmailVerifySendCode(ctx context.Context, email string) (msg string, err error) {

	a.logger.Debug("Trying to send email verify code", "email", email)

	if email == "" {
		a.logger.Debug("Sending email verify code user error", "email", email, "err", utils.ErrEmptyEmail)
		return "Error", utils.ErrInvalidCredentials
	}

	_, err = a.userGetter.GetUserByEmail(ctx, email)
	if err != nil {
		a.logger.Debug("Sending email verify code user error", "email", email, "err", err.Error())
		if err == utils.ErrUserNotFound {
			return "Error", utils.ErrInvalidCredentials
		}
		return "Error", utils.ErrInternalServer
	}

	randCode := utils_random.RandRange(1000, 9999)
	
	a.emailSender.SendEmail(email, randCode)

	if err := a.emailVerifyCodeKeeper.KeepEmailVerifyCode(ctx, email, randCode); err != nil {
		a.logger.Debug("Sending email verify code user error", "email", email, "err", err.Error())
		return "Error", utils.ErrInternalServer
	}

	a.logger.Debug("Email verify code sended", "email", email, "code", randCode)

	return "Code sended", nil
}

func (a *AccountService) EmailVerify(ctx context.Context, email string, code int) (msg string, err error) {

	a.logger.Debug("Trying to verify user's email", "email", email, "code", code)

	if email == "" {
		a.logger.Debug("Verifying user's email error", "email", email, "err", utils.ErrEmptyEmail)
		return "Error", utils.ErrInvalidCredentials
	}

	if code == 0 {
		a.logger.Debug("Verifying user's email error", "email", email, "err", utils.ErrWrongVerificationCode)
		return "Error", utils.ErrInvalidCredentials
	}

	sendedCode, err := a.emailVerifyCodeGetter.GetEmailVerifyCode(ctx, email)
	if err != nil {
		a.logger.Debug("Verifying user's email error", "email", email, "err", err.Error())
		if err == utils.ErrEmailVerifyCodeNotFound {
			return "Error", utils.ErrInvalidCredentials
		}
		return "Error", utils.ErrInternalServer
	}

	if sendedCode != code {
		a.logger.Debug("Verifying user's email error", "email", email, "err", utils.ErrWrongVerificationCode)
		return "Error", utils.ErrInvalidCredentials
	}

	if err := a.emailVerificator.VerifyEmail(ctx, email); err != nil {
		a.logger.Debug("Verifying user's email error", "email", email, "err", err.Error())
		if err == utils.ErrUserNotFound {
			return "Error", utils.ErrInvalidCredentials
		}
		return "Error", utils.ErrInternalServer
	}

	a.logger.Debug("User's email verified succesfully", "email", email)

	return "Success", nil
}

func (a *AccountService) PasswordRecoverSendCode(ctx context.Context, email string) (msg string, err error) {
	
	a.logger.Debug("Trying to send pass recover code", "email", email)

	if email == "" {
		a.logger.Debug("Sending pass recover code error", "email", email, "err", utils.ErrEmptyEmail)
		return "Error", utils.ErrInvalidCredentials
	}

	_, err = a.userGetter.GetUserByEmail(ctx, email)
	if err != nil {
		a.logger.Debug("Sending pass recover code error", "email", email, "err", err.Error())
		if err == utils.ErrUserNotFound {
			return "Error", utils.ErrInvalidCredentials
		}
		return "Error", utils.ErrInternalServer
	}

	randCode := utils_random.RandRange(1000, 9999)

	a.emailSender.SendEmail(email, randCode)

	if err := a.passRecoverCodeKeeper.KeepPassRecoverCode(ctx, email, randCode); err != nil {
		a.logger.Debug("Sending pass recover code error", "email", email, "err", err.Error())
		return "Error", utils.ErrInternalServer
	}

	a.logger.Debug("Pass recover code sended", "email", email, "code", randCode)

	return "Code sended", nil
}

func (a *AccountService) PasswordRecover(ctx context.Context, email string, newPassword string, code int) (msg string, err error) {

	a.logger.Debug("Trying to change user's password", "email", email, "code", code)

	if email == "" {
		a.logger.Debug("Changing user's password error", "email", email, "err", utils.ErrEmptyEmail)
		return "Error", utils.ErrInvalidCredentials
	}

	if newPassword == "" {
		a.logger.Debug("Changing user's password error", "email", email, "err", utils.ErrEmptyPassword)
		return "Error", utils.ErrInvalidCredentials
	}

	if code == 0 {
		a.logger.Debug("Changing user's password error", "email", email, "err", utils.ErrWrongPasswordRecoverCode)
		return "Error", utils.ErrInvalidCredentials
	}

	sendedCode, err := a.passRecoverCodeGetter.GetPassRecoverCode(ctx, email)
	if err != nil {
		a.logger.Debug("Changing user's password error", "email", email, "err", err.Error())
		if err == utils.ErrPassRecoverCodeNotFound {
			return "Error", utils.ErrInvalidCredentials
		}
		return "Error", utils.ErrInternalServer
	}

	if sendedCode != code {
		a.logger.Debug("Changing user's password error", "email", email, "err", utils.ErrWrongPasswordRecoverCode)
		return "Error", utils.ErrInvalidCredentials
	}

	newPassHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Debug("Changing user's password error", "email", email, "err", err.Error())
		return "Error", utils.ErrInternalServer
	}

	if err := a.passChanger.ChangePassword(ctx, email, newPassHash); err != nil {
		a.logger.Debug("Changing user's password error", "email", email, "err", err.Error())
		if err == utils.ErrUserNotFound {
			return "Error", err
		}
		return "Error", utils.ErrInternalServer
	}

	a.logger.Debug("User's password changed succsefully", "email", email)

	return "Success", nil
}