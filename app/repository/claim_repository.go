package repository

import (
	"campus-lost-found-api/app/model"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClaimRepository interface {
	Create(context.Context, model.Claim) (model.Claim, error)
	ListByItem(context.Context, int64) ([]model.Claim, error)
	FindByID(context.Context, int64) (model.Claim, error)
	UpdateStatus(context.Context, int64, string) (model.Claim, error)
}
type claimRepo struct{ pool *pgxpool.Pool }

func NewClaimRepository(p *pgxpool.Pool) ClaimRepository { return &claimRepo{p} }
func (r *claimRepo) Create(ctx context.Context, c model.Claim) (model.Claim, error) {
	err := r.pool.QueryRow(ctx, `INSERT INTO claims(item_id,claimant_id,note) VALUES($1,$2,$3) RETURNING id,item_id,claimant_id,note,status,created_at,updated_at`, c.ItemID, c.ClaimantID, c.Note).Scan(&c.ID, &c.ItemID, &c.ClaimantID, &c.Note, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}
func (r *claimRepo) ListByItem(ctx context.Context, id int64) ([]model.Claim, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,item_id,claimant_id,note,status,created_at,updated_at FROM claims WHERE item_id=$1 ORDER BY created_at DESC,id DESC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Claim{}
	for rows.Next() {
		var c model.Claim
		if err := rows.Scan(&c.ID, &c.ItemID, &c.ClaimantID, &c.Note, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
func (r *claimRepo) FindByID(ctx context.Context, id int64) (model.Claim, error) {
	var c model.Claim
	err := r.pool.QueryRow(ctx, `SELECT id,item_id,claimant_id,note,status,created_at,updated_at FROM claims WHERE id=$1`, id).Scan(&c.ID, &c.ItemID, &c.ClaimantID, &c.Note, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return c, ErrNotFound
	}
	return c, err
}
func (r *claimRepo) UpdateStatus(ctx context.Context, id int64, status string) (model.Claim, error) {
	var c model.Claim
	err := r.pool.QueryRow(ctx, `UPDATE claims SET status=$1,updated_at=NOW() WHERE id=$2 RETURNING id,item_id,claimant_id,note,status,created_at,updated_at`, status, id).Scan(&c.ID, &c.ItemID, &c.ClaimantID, &c.Note, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return c, ErrNotFound
	}
	return c, err
}
