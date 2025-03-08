package models

type User struct {
	Id int64
	Email string
	PassHash []byte
	IsVerified bool
	Use2FA bool
	IsAdmin bool
}