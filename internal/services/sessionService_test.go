package services_test

import (
	"context"
	"testing"
	"time"

	"authSAS/internal/services"
	"authSAS/internal/storages/mockups"
	"authSAS/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {

	permStor := mockups.NewPermStorMokup()
	tempStor := mockups.NewTempStorMokup()
	accService := services.NewAccountService(logger, time.Minute*5, emailSender, permStor, tempStor)
	sesService := services.NewSessionService(logger, time.Minute*5, "test", emailSender, permStor, tempStor)
	ctx := context.Background()

	// registering email for (case 1) test
	accService.Register(ctx, "test@mail.ru", "admin")

	// registering email with 2fa for (case 2) test
	accService.Register(ctx, "test2@mail.ru", "admin")
	user := permStor.UsersStorage["test2@mail.ru"]
	user.Use2FA= true
	permStor.UsersStorage["test2@mail.ru"] = user

	cases := []struct {
		desc string
		inEmail string
		inPassword string
		outToken string
		outMsg string
		mustFail bool
		fail error
	}{
		{
			desc: "case 1 - right login",
			inEmail: "test@mail.ru",
			inPassword: "admin",
			outMsg: "Authorized",
			mustFail: false,
		},
		{
			desc: "case 2 - login with 2fa",
			inEmail: "test2@mail.ru",
			inPassword: "admin",
			outMsg: "2FA code sended",
			mustFail: false,
		},
		{
			desc: "case 3 - unregistered email",
			inEmail: "123@mail.ru",
			inPassword: "admin",
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 4 - invalid password",
			inEmail: "test@mail.ru",
			inPassword: "admin123",
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 5 - empty email",
			inEmail: "",
			inPassword: "admin",
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 6 - empty password",
			inEmail: "test@mail.ru",
			inPassword: "",
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
	}

	for _, tC := range cases {
		token, msg, err := sesService.Login(ctx, tC.inEmail, tC.inPassword)

		if !tC.mustFail {
			require.NoError(t, err)
			require.Equal(t, tC.outMsg, msg)
			
			if msg != "2FA code sended" {
				require.NotEmpty(t, token)
			} else {
				require.Empty(t, token)

				code, err := tempStor.GetTwoFACode(ctx, tC.inEmail)
				require.NoError(t, err)
				require.NotEmpty(t, code)
			}
			
		} else {
			require.ErrorIs(t, err, tC.fail)
			require.Equal(t, tC.outMsg, msg)
			require.Empty(t, token)
		}
	}
}

func TestLogout(t *testing.T) {

	permStor := mockups.NewPermStorMokup()
	tempStor := mockups.NewTempStorMokup()
	accService := services.NewAccountService(logger, time.Minute*5, emailSender, permStor, tempStor)
	sesService := services.NewSessionService(logger, time.Minute*5, "test", emailSender, permStor, tempStor)
	ctx := context.Background()

	// preparing for (case 1) test
	accService.Register(ctx, "test@mail.ru", "admin")
	validToken,_,_ := sesService.Login(ctx, "test@mail.ru", "admin")

	cases := []struct {
		desc string
		inToken string
		outMsg string
		mustFail bool
		fail error
	}{
		{
			desc: "case 1 - logout with valid token",
			inToken: validToken,
			outMsg: "Success",
			mustFail: false,
		},
		{
			desc: "case 2 - logout with INVALID token",
			inToken: "invalid",
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 3 - emty token string",
			inToken: "",
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 4 - logout with already added token",
			inToken: validToken,
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrJWTAlreadyAdded,
		},
	}

	for _, tC := range cases {
		msg, err := sesService.Logout(ctx, tC.inToken)

		if !tC.mustFail {
			require.NoError(t, err)
			require.Equal(t, tC.outMsg, msg)		
		} else {
			require.ErrorIs(t, err, tC.fail)
			require.Equal(t, tC.outMsg, msg)
		}
	}
}

func TestLoginWith2FACode(t *testing.T) {

	permStor := mockups.NewPermStorMokup()
	tempStor := mockups.NewTempStorMokup()
	accService := services.NewAccountService(logger, time.Minute*5, emailSender, permStor, tempStor)
	sesService := services.NewSessionService(logger, time.Minute*5, "test", emailSender, permStor, tempStor)
	ctx := context.Background()

	// preparing for (case 1) test
	accService.Register(ctx, "test@mail.ru", "admin")
	tempStor.KeepTwoFACode(ctx, "test@mail.ru", 1234)
	

	cases := []struct {
		desc string
		inEmail string
		inCode int
		outToken string
		mustFail bool
		fail error
	}{
		{
			desc: "case 1 - right 2FA login",
			inEmail: "test@mail.ru",
			inCode: 1234,
			mustFail: false,
		},
		{
			desc: "case 2 - 2fa login with wrong email",
			inEmail: "123@mail.ru",
			inCode: 1234,
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 3 - 2fa login with wrong code",
			inEmail: "test@mail.ru",
			inCode: 5555,
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 4 - empty email",
			inEmail: "",
			inCode: 5555,
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 5 - empty code",
			inEmail: "test@mail.ru",
			inCode: 0,
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
	}

	for _, tC := range cases {
		token, err := sesService.LoginWith2FACode(ctx, tC.inEmail, tC.inCode)

		if !tC.mustFail {
			require.NoError(t, err)
			require.NotEmpty(t, token)
		} else {
			require.ErrorIs(t, err, tC.fail)
			require.Empty(t, token)
		}
	}
}