package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func NewPool(ctx context.Context, url string, min, max int32) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	cfg.MinConns = min
	cfg.MaxConns = max
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute
	return pgxpool.NewWithConfig(ctx, cfg)
}
