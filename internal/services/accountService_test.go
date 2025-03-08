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

func TestRegister(t *testing.T) {

	permStor := mockups.NewPermStorMokup()
	tempStor := mockups.NewTempStorMokup()
	accService := services.NewAccountService(logger, time.Minute*5, emailSender, permStor, tempStor)
	ctx := context.Background()

	cases := []struct {
		desc string
		inEmail string
		inPassword string
		outUserId int64
		mustFail bool
		fail error
	}{
		{
			desc: "case 1 - right reg.",
			inEmail: "test@mail.ru",
			inPassword: "admin",
			outUserId: 1,
			mustFail: false,
		},
		{
			desc: "case 2 - reg. again with same data",
			inEmail: "test@mail.ru",
			inPassword: "admin",
			outUserId: 0,
			mustFail: true,
			fail: utils.ErrUserAlreadyExists,
		},
		{
			desc: "case 3 - empty email",
			inEmail: "",
			inPassword: "admin",
			outUserId: 0,
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 4 - empty password",
			inEmail: "123@mail.ru",
			inPassword: "",
			outUserId: 0,
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
	}

	for _, tC := range cases {
		userId, err := accService.Register(ctx, tC.inEmail, tC.inPassword)

		if !tC.mustFail {
			require.NoError(t, err)
			require.Equal(t, tC.outUserId, userId)
		} else {
			require.ErrorIs(t, err, tC.fail)
			require.Equal(t, tC.outUserId, userId)
		}
	}
}

func TestEmailVerifySendCode(t *testing.T) {

	permStor := mockups.NewPermStorMokup()
	tempStor := mockups.NewTempStorMokup()
	accService := services.NewAccountService(logger, time.Minute*5, emailSender, permStor, tempStor)
	ctx := context.Background()

	// registering email for (case 1) test
	accService.Register(ctx, "test@mail.ru", "admin")

	cases := []struct {
		desc string
		inEmail string
		outMsg string
		mustFail bool
		fail error
	}{
		{
			desc: "case 1 - right sending",
			inEmail: "test@mail.ru",
			outMsg: "Code sended",
			mustFail: false,
		},
		{
			desc: "case 2 - sending on not registered email",
			inEmail: "123@mail.ru",
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 3 - empty email",
			inEmail: "",
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
	}

	for _, tC := range cases {
		msg, err := accService.EmailVerifySendCode(ctx, tC.inEmail)

		if !tC.mustFail {
			require.NoError(t, err)
			require.Equal(t, tC.outMsg, msg)

			code, err := tempStor.GetEmailVerifyCode(ctx, tC.inEmail)
			require.NoError(t, err)
			require.NotEmpty(t, code)

		} else {
			require.ErrorIs(t, err, tC.fail)
			require.Equal(t, tC.outMsg, msg)
		}
	}
}

func TestEmailVerify(t *testing.T) {

	permStor := mockups.NewPermStorMokup()
	tempStor := mockups.NewTempStorMokup()
	accService := services.NewAccountService(logger, time.Minute*5, emailSender, permStor, tempStor)
	ctx := context.Background()

	// prepare for (case 1) test
	accService.Register(ctx, "test@mail.ru", "admin")
	tempStor.KeepEmailVerifyCode(ctx, "test@mail.ru", 1234)

	cases := []struct {
		desc string
		inEmail string
		inCode int
		outMsg string
		mustFail bool
		fail error
	}{
		{
			desc: "case 1 - right verification",
			inEmail: "test@mail.ru",
			inCode: 1234,
			outMsg: "Success",
			mustFail: false,
		},
		{
			desc: "case 2 - wrong email",
			inEmail: "123@mail.ru",
			inCode: 1234,
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 2 - wrong code",
			inEmail: "test@mail.ru",
			inCode: 5555,
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 3 - empty email",
			inEmail: "",
			inCode: 1234,
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 4 - empty code",
			inEmail: "test@mail.ru",
			inCode: 0,
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
	}

	for _, tC := range cases {
		msg, err := accService.EmailVerify(ctx, tC.inEmail, tC.inCode)

		if !tC.mustFail {
			require.NoError(t, err)
			require.Equal(t, tC.outMsg, msg)
		} else {
			require.ErrorIs(t, err, tC.fail)
			require.Equal(t, tC.outMsg, msg)
		}
	}
}

func TestPasswordRecoverSendCode(t *testing.T) {

	permStor := mockups.NewPermStorMokup()
	tempStor := mockups.NewTempStorMokup()
	accService := services.NewAccountService(logger, time.Minute*5, emailSender, permStor, tempStor)
	ctx := context.Background()

	// registering email for (case 1) test
	accService.Register(ctx, "test@mail.ru", "admin")

	cases := []struct {
		desc string
		inEmail string
		outMsg string
		mustFail bool
		fail error
	}{
		{
			desc: "case 1 - right sending",
			inEmail: "test@mail.ru",
			outMsg: "Code sended",
			mustFail: false,
		},
		{
			desc: "case 2 - sending on not registered email",
			inEmail: "123@mail.ru",
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 3 - empty email",
			inEmail: "",
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
	}

	for _, tC := range cases {
		msg, err := accService.PasswordRecoverSendCode(ctx, tC.inEmail)

		if !tC.mustFail {
			require.NoError(t, err)
			require.Equal(t, tC.outMsg, msg)

			code, err := tempStor.GetPassRecoverCode(ctx, tC.inEmail)
			require.NoError(t, err)
			require.NotEmpty(t, code)

		} else {
			require.ErrorIs(t, err, tC.fail)
			require.Equal(t, tC.outMsg, msg)
		}
	}
}

func TestPasswordRecover(t *testing.T) {

	permStor := mockups.NewPermStorMokup()
	tempStor := mockups.NewTempStorMokup()
	accService := services.NewAccountService(logger, time.Minute*5, emailSender, permStor, tempStor)
	ctx := context.Background()

	// prepare for (case 1) test
	accService.Register(ctx, "test@mail.ru", "admin")
	tempStor.KeepPassRecoverCode(ctx, "test@mail.ru", 1234)

	cases := []struct {
		desc string
		inEmail string
		inNewPassword string
		inCode int
		outMsg string
		mustFail bool
		fail error
	}{
		{
			desc: "case 1 - right verification",
			inEmail: "test@mail.ru",
			inNewPassword: "admin123",
			inCode: 1234,
			outMsg: "Success",
			mustFail: false,
		},
		{
			desc: "case 2 - wrong email",
			inEmail: "123@mail.ru",
			inNewPassword: "admin123",
			inCode: 1234,
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 2 - wrong code",
			inEmail: "test@mail.ru",
			inNewPassword: "admin123",
			inCode: 5555,
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 3 - empty email",
			inEmail: "",
			inNewPassword: "admin123",
			inCode: 1234,
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 3 - empty email",
			inEmail: "test@mail.ru",
			inNewPassword: "",
			inCode: 1234,
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
		{
			desc: "case 4 - empty code",
			inEmail: "test@mail.ru",
			inNewPassword: "admin123",
			inCode: 0,
			outMsg: "Error",
			mustFail: true,
			fail: utils.ErrInvalidCredentials,
		},
	}

	for _, tC := range cases {
		msg, err := accService.PasswordRecover(ctx, tC.inEmail, tC.inNewPassword, tC.inCode)

		if !tC.mustFail {
			require.NoError(t, err)
			require.Equal(t, tC.outMsg, msg)
		} else {
			require.ErrorIs(t, err, tC.fail)
			require.Equal(t, tC.outMsg, msg)
		}
	}
}