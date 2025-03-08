package mockups

import (
	"authSAS/internal/utils"
	"context"
	"fmt"
	"sync"
)

type TempStorMockup struct {
	codeStorage map[string] int
	sync.RWMutex
}

func NewTempStorMokup() (*TempStorMockup) {
	return &TempStorMockup{codeStorage: make(map[string] int)}
}

func (s *TempStorMockup) KeepTwoFACode(ctx context.Context, email string, code int) (err error) {
	key := fmt.Sprintf("2fa_code_key: %s", email)

	s.RWMutex.Lock()
	s.codeStorage[key] = code
	s.RWMutex.Unlock()

	return nil
}

func (s *TempStorMockup) GetTwoFACode(ctx context.Context, email string) (code int, err error) {
	key := fmt.Sprintf("2fa_code_key: %s", email)

	s.RWMutex.RLock()
	result , ok := s.codeStorage[key]
	s.RWMutex.RUnlock()

	if !ok {
		return 0, utils.Err2FACodeNotFound
	}

	return result, nil
}

func (s *TempStorMockup) KeepEmailVerifyCode(ctx context.Context, email string, code int) (err error) {
	key := fmt.Sprintf("email_verify_key: %s", email)

	s.RWMutex.Lock()
	s.codeStorage[key] = code
	s.RWMutex.Unlock()

	return nil
}

func (s *TempStorMockup) GetEmailVerifyCode(ctx context.Context, email string) (code int, err error) {
	key := fmt.Sprintf("email_verify_key: %s", email)

	s.RWMutex.RLock()
	result , ok := s.codeStorage[key]
	s.RWMutex.RUnlock()

	if !ok {
		return 0, utils.ErrEmailVerifyCodeNotFound
	}

	return result, nil
}

func (s *TempStorMockup) KeepPassRecoverCode(ctx context.Context, email string, code int) (err error) {
	key := fmt.Sprintf("pass_recover_key: %s", email)

	s.RWMutex.Lock()
	s.codeStorage[key] = code
	s.RWMutex.Unlock()

	return nil
}

func (s *TempStorMockup) GetPassRecoverCode(ctx context.Context, email string) (code int, err error) {
	key := fmt.Sprintf("pass_recover_key: %s", email)

	s.RWMutex.RLock()
	result , ok := s.codeStorage[key]
	s.RWMutex.RUnlock()

	if !ok {
		return 0, utils.ErrPassRecoverCodeNotFound
	}

	return result, nil
}