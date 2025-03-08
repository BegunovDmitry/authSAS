package app

import (
	"fmt"
	"log/slog"
	"net"

	"authSAS/internal/config"
	authServer "authSAS/internal/server"
	"authSAS/internal/services"
	emailsender "authSAS/internal/utils/emailSender"

	"google.golang.org/grpc"
)

type App struct {
	logger *slog.Logger
	grpsServer *grpc.Server
	config *config.Config
}

func NewApp(logger *slog.Logger, config *config.Config, permanentStorage services.PermanentStorage, temporaryStorage services.TemporaryStorage) *App {

	sender := emailsender.NewEmailSender(logger, config.EmailSender.Email, config.EmailSender.Password)

	sessionService := services.NewSessionService(logger, config.JWTTokenTTL, config.JWTSecret, sender, permanentStorage, temporaryStorage)
	accountService := services.NewAccountService(logger, config.JWTTokenTTL, sender, permanentStorage, temporaryStorage)
	logger.Info("All services initialized")

	grpsServer := grpc.NewServer()

	authServer.RegisterServer(grpsServer, sessionService, accountService)
	logger.Info("gRPC server registered")

	return &App{
		logger: logger,
		grpsServer: grpsServer,
		config: config,
	}
}

func (a *App) MustRun() {
	if err := a.runApp(); err != nil {
		panic(err)
	}
}

func (a *App) StopApp() {
	a.grpsServer.GracefulStop()
}

func (a *App) runApp() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.config.Grpc.Domain, a.config.Grpc.Port))
	if err != nil {
		return fmt.Errorf("listen failed: - err: %w", err)
	}

	if err := a.grpsServer.Serve(l); err != nil {
		return fmt.Errorf("grpcServer serve failed: - err: %w", err)
	}

	return nil
}
