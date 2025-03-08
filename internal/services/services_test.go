package services_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"authSAS/internal/config"
	"authSAS/internal/services"
	"authSAS/internal/storages/mockups"
	emailsender "authSAS/internal/utils/emailSender"
)

type Tester struct {
	T *testing.T
	cfg *config.Config
	logger *slog.Logger
	permStor *mockups.PermStorMockup
	tempStor *mockups.TempStorMockup
	accService *services.AccountService
	sesService *services.SessionService
	emailSender *emailsender.EmailSender
}

func NewTester(t *testing.T) (context.Context, *Tester) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../../config/config.yaml")
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	emailSender := emailsender.NewEmailSender(logger, cfg.EmailSender.Email, cfg.EmailSender.Password)

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.Grpc.RequestTimeout)

	permStor := mockups.NewPermStorMokup()
	tempStor := mockups.NewTempStorMokup()
	accService := services.NewAccountService(logger, cfg.JWTTokenTTL, emailSender, permStor, tempStor)
	sesService := services.NewSessionService(logger, cfg.JWTTokenTTL, cfg.JWTSecret, emailSender, permStor, tempStor)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	return ctx, &Tester{
		T: t,
		cfg: cfg,
		logger: logger,
		permStor: permStor,
		tempStor: tempStor,
		accService: accService,
		sesService: sesService,
		emailSender: emailSender,
	}
}