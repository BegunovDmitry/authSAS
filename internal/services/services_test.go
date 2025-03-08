package services_test

import (
	"log/slog"
	"os"

	test_config "authSAS/config"
	emailsender "authSAS/internal/utils/emailSender"
)

var cfg = test_config.TestsConfig
var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
var emailSender = emailsender.NewEmailSender(logger, cfg["email"], cfg["password"])