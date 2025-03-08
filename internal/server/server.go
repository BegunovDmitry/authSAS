package server

import (
	"context"

	sasv1 "github.com/BegunovDmitry/authSASproto/result/go"

	"google.golang.org/grpc"
)

type SessionService interface {
	Login(ctx context.Context, email string, password string) (token string, msg string, err error)
	Logout(ctx context.Context, token string) (msg string, err error)
	LoginWith2FACode(ctx context.Context, email string, code int) (token string, err error)
}

type AccountService interface {
	Register(ctx context.Context, email string, password string) (userId int64, err error)
	EmailVerifySendCode(ctx context.Context, email string) (msg string, err error)
	EmailVerify(ctx context.Context, email string, code int) (msg string, err error)
	PasswordRecoverSendCode(ctx context.Context, email string) (msg string, err error)
	PasswordRecover(ctx context.Context, email string, newPassword string, code int) (msg string, err error)
}

type Server struct {
	sasv1.UnimplementedAuthServer
	sessionService SessionService
	accountService AccountService
}

func RegisterServer(grpc *grpc.Server, sessionService SessionService, accountService AccountService) {
	sasv1.RegisterAuthServer(grpc, &Server{sessionService: sessionService, accountService: accountService})
}

// Logic

func (s *Server) Login(ctx context.Context, req *sasv1.LoginRequest) (*sasv1.LoginResponce, error) {

	email := req.GetEmail()
	password := req.GetPassword()

	token, msg, err := s.sessionService.Login(ctx, email, password)

	return &sasv1.LoginResponce{
		Token: token,
		Msg: msg,
	}, err

}

func (s *Server) Logout(ctx context.Context, req *sasv1.LogoutRequest) (*sasv1.LogoutResponce, error) {

	token := req.GetToken()

	msg, err := s.sessionService.Logout(ctx, token)

	return &sasv1.LogoutResponce{
		Msg: msg,
	}, err

}

func (s *Server) LoginWith2FACode(ctx context.Context, req *sasv1.LoginWith2FACodeRequest) (*sasv1.LoginWith2FACodeResponce, error) {

	email := req.GetEmail()
	code := req.GetCode()

	token, err := s.sessionService.LoginWith2FACode(ctx, email, int(code))

	return &sasv1.LoginWith2FACodeResponce{
		Token: token,
	}, err

}

func (s *Server) Register(ctx context.Context, req *sasv1.RegisterRequest) (*sasv1.RegisterResponce, error) {
	
	email := req.GetEmail()
	password := req.GetPassword()

	userId, err := s.accountService.Register(ctx, email, password)

	return &sasv1.RegisterResponce{
		UserId: userId,
	}, err
}

func (s *Server) EmailVerifySendCode(ctx context.Context, req *sasv1.EmailVerifySendCodeRequest) (*sasv1.EmailVerifySendCodeResponce, error) {

	email := req.GetEmail()

	msg, err := s.accountService.EmailVerifySendCode(ctx, email)

	return &sasv1.EmailVerifySendCodeResponce{
		Msg: msg,
	}, err
}

func (s *Server) EmailVerify(ctx context.Context, req *sasv1.EmailVerifyRequest) (*sasv1.EmailVerifyResponce, error) {

	email := req.GetEmail()
	code := req.GetCode()

	msg, err := s.accountService.EmailVerify(ctx, email, int(code))

	return &sasv1.EmailVerifyResponce{
		Msg: msg,
	}, err
}

func (s *Server) PasswordRecoverSendCode(ctx context.Context, req *sasv1.PasswordRecoverSendCodeRequest) (*sasv1.PasswordRecoverSendCodeResponce, error) {

	email := req.GetEmail()

	msg, err := s.accountService.PasswordRecoverSendCode(ctx, email)

	return &sasv1.PasswordRecoverSendCodeResponce{
		Msg: msg,
	}, err
}

func (s *Server) PasswordRecover(ctx context.Context, req *sasv1.PasswordRecoverRequest) (*sasv1.PasswordRecoverResponce, error) {

	email := req.GetEmail()
	newPassword := req.GetNewPassword()
	code := req.GetCode()

	msg, err := s.accountService.PasswordRecover(ctx, email, newPassword, int(code))

	return &sasv1.PasswordRecoverResponce{
		Msg: msg,
	}, err
}

