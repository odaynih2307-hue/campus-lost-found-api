package repository

import (
	"campus-lost-found-api/app/model"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type StoredToken struct {
	UserID    int64
	ExpiresAt time.Time
}
type TokenRepository interface {
	Create(context.Context, int64, string, time.Time) error
	FindActive(context.Context, string) (StoredToken, error)
	Revoke(context.Context, string) error
}
type tokenRepo struct{ pool *pgxpool.Pool }

func NewTokenRepository(p *pgxpool.Pool) TokenRepository { return &tokenRepo{p} }
func (r *tokenRepo) Create(ctx context.Context, uid int64, h string, e time.Time) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO refresh_tokens(user_id,token_hash,expires_at) VALUES($1,$2,$3)`, uid, h, e)
	return err
}
func (r *tokenRepo) FindActive(ctx context.Context, h string) (StoredToken, error) {
	var s StoredToken
	err := r.pool.QueryRow(ctx, `SELECT user_id,expires_at FROM refresh_tokens WHERE token_hash=$1 AND revoked_at IS NULL AND expires_at>NOW()`, h).Scan(&s.UserID, &s.ExpiresAt)
	if err != nil {
		return s, ErrNotFound
	}
	return s, nil
}
func (r *tokenRepo) Revoke(ctx context.Context, h string) error {
	_, err := r.pool.Exec(ctx, `UPDATE refresh_tokens SET revoked_at=NOW() WHERE token_hash=$1 AND revoked_at IS NULL`, h)
	return err
}

var _ model.User
