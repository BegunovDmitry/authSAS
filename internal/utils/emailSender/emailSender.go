package emailsender

import (
	"authSAS/internal/utils"
	"log/slog"
	"net/smtp"
	"strconv"
)

type EmailSender struct {
	logger *slog.Logger
	email string
	password string
}

func NewEmailSender(logger *slog.Logger, email string, password string)  *EmailSender{
	return &EmailSender{
		logger: logger,
		email: email,
		password: password,
	}
}



func (s *EmailSender) SendEmail(userEmail string, code int) error {

	if s.email == "" || s.password == "" {
		s.logger.Debug("Email sender error", "err", "Inavlid config")
		return utils.ErrInternalServer
	}

	if userEmail == "" || code == 0 {
		s.logger.Debug("Email sender error", "email", userEmail, "err", "Inavlid credentials")
		return utils.ErrInvalidCredentials
	}

	to := userEmail
	from := s.email
	mesage := []byte("Hello from MyApp!\r\n"+
					"Here is your code: "+strconv.Itoa(code))

	addr := "smtp.yandex.ru:587"
	host := "smtp.yandex.ru"

	user := s.email
	auth := smtp.PlainAuth("", user, s.password, host)

	err := smtp.SendMail(addr, auth, from, []string{to}, mesage)

	if err != nil {
		s.logger.Debug("Send email code error", "email", userEmail, "err", err.Error())
		return utils.ErrInternalServer
	}

	return nil
}