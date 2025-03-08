package utils

import "errors"

var (
	ErrInternalServer = errors.New("internal server error")

	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrEmptyEmail = errors.New("email is required")
	ErrEmptyPassword = errors.New("password is required")
	ErrEmptyJWT = errors.New("token is required")

	ErrJWTAlreadyAdded = errors.New("jwt already added")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrUserNotFound = errors.New("user not found")
	Err2FACodeNotFound = errors.New("2 factor auth code not found in temp. storage")
	ErrEmailVerifyCodeNotFound = errors.New("email verify code not found in temp. storage")
	ErrPassRecoverCodeNotFound = errors.New("password recover code not found in temp. storage")

	ErrWrong2FACode = errors.New("wrong 2 factor auth code")
	ErrWrongVerificationCode = errors.New("wrong email verification code")
	ErrWrongPasswordRecoverCode = errors.New("wrong password recover code")
)