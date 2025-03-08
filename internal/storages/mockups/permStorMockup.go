package mockups

import (
	"authSAS/internal/models"
	"authSAS/internal/utils"
	"context"
	"sync"
)

type PermStorMockup struct {
	UsersStorage map[string] models.User
	JwtStore map[int64] string
	usersCnt int
	sync.RWMutex
 
}

func NewPermStorMokup() (*PermStorMockup) {
	return &PermStorMockup{
		UsersStorage: make(map[string] models.User), 
		JwtStore: make(map[int64] string),
		usersCnt: 0,
	}
}

func (s *PermStorMockup) GetUserByEmail(ctx context.Context, email string) (user models.User, err error) {
	s.RWMutex.RLock()
	result, ok := s.UsersStorage[email]
	s.RWMutex.RUnlock()

	if !ok {
		return models.User{}, utils.ErrUserNotFound
	}
	
	return result, nil
}

func (s *PermStorMockup) KeepLogoutJWT(ctx context.Context, uid int64, token string) (err error) {
	s.RWMutex.RLock()
	_, ok := s.JwtStore[uid]
	s.RWMutex.RUnlock()

	if ok {
		return utils.ErrJWTAlreadyAdded
	}

	s.RWMutex.Lock()
	s.JwtStore[uid] = token
	s.RWMutex.Unlock()

	return nil
}

func (s *PermStorMockup) CreateUser(ctx context.Context, email string, passHash []byte) (userId int64, err error) {
	s.RWMutex.RLock()
	_, ok := s.UsersStorage[email]
	s.RWMutex.RUnlock()

	if ok {
		return 0, utils.ErrUserAlreadyExists
	}

	user := models.User{
		Id: int64(s.usersCnt), 
		Email: email,
		PassHash: passHash,
		IsVerified: false,
		Use2FA: false,
		IsAdmin: false,
	}

	s.RWMutex.Lock()
	s.UsersStorage[email] = user
	s.usersCnt++
	s.RWMutex.Unlock()

	return int64(s.usersCnt), nil
}

func (s *PermStorMockup) VerifyEmail(ctx context.Context, email string) (err error) {
	s.RWMutex.RLock()
	result, ok := s.UsersStorage[email]
	s.RWMutex.RUnlock()

	if !ok {
		return utils.ErrUserNotFound
	}

	result.IsVerified = true

	s.RWMutex.Lock()
	s.UsersStorage[email] = result
	s.RWMutex.Unlock()

	return nil
}

func (s *PermStorMockup) ChangePassword(ctx context.Context, email string, newPassHash []byte) (err error) {
	s.RWMutex.RLock()
	result, ok := s.UsersStorage[email]
	s.RWMutex.RUnlock()

	if !ok {
		return utils.ErrUserNotFound
	}

	result.PassHash = newPassHash

	s.RWMutex.Lock()
	s.UsersStorage[email] = result
	s.RWMutex.Unlock()

	return nil
}