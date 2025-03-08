package services

import (
	"authSAS/internal/models"
	"context"
)

// SessionService storage interfaces

type UserGetter interface {
	GetUserByEmail(ctx context.Context, email string) (user models.User, err error)
}

type LogoutJWTKeeper interface {
	KeepLogoutJWT(ctx context.Context, uid int64, token string) (err error)
}

type TwoFACodeKeeper interface {
	KeepTwoFACode(ctx context.Context, email string, code int) (err error)
}

type TwoFACodeGetter interface {
	GetTwoFACode(ctx context.Context, email string) (code int, err error)
}

// AccountService storage interfaces

type UserCreator interface {
	CreateUser(ctx context.Context, email string, passHash []byte) (userId int64, err error)
}

type EmailVerificator interface {
	VerifyEmail(ctx context.Context, email string) (err error)
}

type PassChanger interface {
	ChangePassword(ctx context.Context, email string, newPassHash []byte) (err error)
}

type EmailVerifyCodeKeeper interface {
	KeepEmailVerifyCode(ctx context.Context, email string, code int) (err error)
}

type EmailVerifyCodeGetter interface {
	GetEmailVerifyCode(ctx context.Context, email string) (code int, err error)
}

type PassRecoverCodeKeeper interface {
	KeepPassRecoverCode(ctx context.Context, email string, code int) (err error)
}

type PassRecoverCodeGetter interface {
	GetPassRecoverCode(ctx context.Context, email string) (code int, err error)
}




type PermanentStorage interface {
	UserGetter
	LogoutJWTKeeper

	UserCreator
	EmailVerificator
	PassChanger
}

type TemporaryStorage interface {
	TwoFACodeKeeper
	TwoFACodeGetter

	EmailVerifyCodeKeeper
	EmailVerifyCodeGetter
	PassRecoverCodeKeeper
	PassRecoverCodeGetter
}