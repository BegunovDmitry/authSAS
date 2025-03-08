package redis

import (
	"authSAS/internal/utils"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type TemporaryStorage struct {
	client *redis.Client
	codeTTL time.Duration
}

func NewStorage(client *redis.Client, codeTTL time.Duration) (*TemporaryStorage) {
	return &TemporaryStorage{client: client, codeTTL: codeTTL}
}

func (s *TemporaryStorage) KeepTwoFACode(ctx context.Context, email string, code int) (err error) {
	key := fmt.Sprintf("2fa_code_key: %s", email)

	err = s.client.Set(ctx, key, code, s.codeTTL).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *TemporaryStorage) GetTwoFACode(ctx context.Context, email string) (code int, err error) {
	key := fmt.Sprintf("2fa_code_key: %s", email)

	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, utils.Err2FACodeNotFound
		}
		return 0, err
	}

	if val == "" {
		return 0, utils.Err2FACodeNotFound
	}

	code, err = strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return code, nil
}

func (s *TemporaryStorage) KeepEmailVerifyCode(ctx context.Context, email string, code int) (err error) {
	key := fmt.Sprintf("email_verify_key: %s", email)

	err = s.client.Set(ctx, key, code, s.codeTTL).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *TemporaryStorage) GetEmailVerifyCode(ctx context.Context, email string) (code int, err error) {
	key := fmt.Sprintf("email_verify_key: %s", email)

	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if val == "" {
		return 0, utils.ErrEmailVerifyCodeNotFound
	}

	code, err = strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return code, nil
}

func (s *TemporaryStorage) KeepPassRecoverCode(ctx context.Context, email string, code int) (err error) {
	key := fmt.Sprintf("pass_recover_key: %s", email)

	err = s.client.Set(ctx, key, code, s.codeTTL).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *TemporaryStorage) GetPassRecoverCode(ctx context.Context, email string) (code int, err error) {
	key := fmt.Sprintf("pass_recover_key: %s", email)

	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if val == "" {
		return 0, utils.ErrPassRecoverCodeNotFound
	}

	code, err = strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return code, nil
}