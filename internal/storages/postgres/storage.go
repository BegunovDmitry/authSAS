package postgres

import (
	"context"
	"errors"

	"authSAS/internal/models"
	"authSAS/internal/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"
)

type PermanentStorage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) (*PermanentStorage) {
	return &PermanentStorage{pool: pool}
}

// For Session Service 

func (s *PermanentStorage) GetUserByEmail(ctx context.Context, email string) (user models.User, err error) {
	query := `SELECT id, email, password_hash, is_verified, use_2fa, is_admin 
	FROM users 
	WHERE email = $1`

	err = s.pool.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.PassHash,
		&user.IsVerified,
		&user.Use2FA,
		&user.IsAdmin,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, utils.ErrUserNotFound
		}

		return models.User{}, err
	}

	return user, nil
}

func (s *PermanentStorage) KeepLogoutJWT(ctx context.Context, uid int64, token string) (err error) {
	query := `INSERT INTO bad_jwts (user_id, token) 
	VALUES ($1, $2);`

	_, err = s.pool.Exec(ctx, query, uid, token)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return utils.ErrJWTAlreadyAdded
		}
		return err
	}

	return nil
}

// For Account Service 

func (s *PermanentStorage) CreateUser(ctx context.Context, email string, passHash []byte) (userId int64, err error) {
	query := `INSERT INTO users (email, password_hash) 
	VALUES ($1, $2) 
	RETURNING id;`

	err = s.pool.QueryRow(ctx, query, email, passHash).Scan(&userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, utils.ErrUserAlreadyExists
		}
		return 0, err
	}

	return userId, nil
}

func (s *PermanentStorage) VerifyEmail(ctx context.Context, email string) (err error) {
	query := `UPDATE users 
	SET is_verified = true 
	WHERE email = $1`

	result, err := s.pool.Exec(ctx, query, email)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}

func (s *PermanentStorage) ChangePassword(ctx context.Context, email string, newPassHash []byte) (err error) {
	query := `UPDATE users 
	SET password_hash = $1 
	WHERE email = $2`

	result, err := s.pool.Exec(ctx, query, newPassHash, email)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}