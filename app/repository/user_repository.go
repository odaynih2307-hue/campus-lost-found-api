package repository

import (
	"campus-lost-found-api/app/model"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

var ErrNotFound = errors.New("data tidak ditemukan")
var ErrDuplicate = errors.New("data sudah digunakan")

type UserRepository interface {
	FindByID(context.Context, int64) (model.User, error)
	FindByUsername(context.Context, string) (model.User, error)
	Create(context.Context, model.User) (model.User, error)
	UpdateRole(context.Context, int64, string) error
}
type userPostgresRepository struct{ pool *pgxpool.Pool }

func NewUserRepository(p *pgxpool.Pool) UserRepository { return &userPostgresRepository{p} }
func (r *userPostgresRepository) FindByID(ctx context.Context, id int64) (model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx, `SELECT id,username,email,password,role,is_active,created_at FROM users WHERE id=$1`, id).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return u, ErrNotFound
	}
	return u, err
}
func (r *userPostgresRepository) FindByUsername(ctx context.Context, username string) (model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx, `SELECT id,username,email,password,role,is_active,created_at FROM users WHERE lower(username)=lower($1)`, username).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return u, ErrNotFound
	}
	return u, err
}
func (r *userPostgresRepository) Create(ctx context.Context, u model.User) (model.User, error) {
	err := r.pool.QueryRow(ctx, `INSERT INTO users(username,email,password,role,is_active) VALUES($1,$2,$3,'user',true) RETURNING id,username,email,password,role,is_active,created_at`, u.Username, u.Email, u.Password).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.CreatedAt)
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		return u, ErrDuplicate
	}
	return u, err
}
func (r *userPostgresRepository) UpdateRole(ctx context.Context, id int64, role string) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE users SET role=$1 WHERE id=$2`, role, id)
	if err == nil && cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}
