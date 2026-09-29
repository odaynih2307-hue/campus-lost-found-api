package repository

import (
	"campus-lost-found-api/app/model"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"
	"time"
)

type ItemRepository interface {
	List(context.Context, model.ListQuery) ([]model.Item, bool, error)
	FindByID(context.Context, int64) (model.Item, error)
	Create(context.Context, model.Item) (model.Item, error)
	Update(context.Context, model.Item) (model.Item, error)
	Delete(context.Context, int64) error
}
type itemRepo struct{ pool *pgxpool.Pool }

func NewItemRepository(p *pgxpool.Pool) ItemRepository { return &itemRepo{p} }
func (r *itemRepo) List(ctx context.Context, q model.ListQuery) ([]model.Item, bool, error) {
	limit := q.Limit
	args := []any{}
	where := []string{"1=1"}
	n := 1
	if q.CursorCreated != nil {
		where = append(where, `(created_at,id)<($`+fmtInt(n)+`,$`+fmtInt(n+1)+`)`)
		args = append(args, *q.CursorCreated, q.CursorID)
		n += 2
	}
	if q.Search != "" {
		where = append(where, `(title ILIKE $`+fmtInt(n)+` OR description ILIKE $`+fmtInt(n)+` OR location ILIKE $`+fmtInt(n)+`)`)
		args = append(args, "%"+q.Search+"%")
		n++
	}
	if q.Status != "" {
		where = append(where, `status=$`+fmtInt(n))
		args = append(args, q.Status)
		n++
	}
	if q.Category != "" {
		where = append(where, `category=$`+fmtInt(n))
		args = append(args, q.Category)
		n++
	}
	if q.OwnerID > 0 {
		where = append(where, `reported_by=$`+fmtInt(n))
		args = append(args, q.OwnerID)
		n++
	}
	sql := `SELECT id,title,description,category,location,status,reported_by,created_at,updated_at FROM items WHERE ` + strings.Join(where, " AND ") + ` ORDER BY created_at DESC,id DESC LIMIT $` + fmtInt(n)
	args = append(args, limit+1)
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()
	out := []model.Item{}
	for rows.Next() {
		var i model.Item
		if err := rows.Scan(&i.ID, &i.Title, &i.Description, &i.Category, &i.Location, &i.Status, &i.ReportedBy, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, false, err
		}
		out = append(out, i)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}
	more := len(out) > limit
	if more {
		out = out[:limit]
	}
	return out, more, nil
}
func fmtInt(i int) string {
	if i < 1 {
		return "1"
	}
	return strconv.Itoa(i)
}
func (r *itemRepo) FindByID(ctx context.Context, id int64) (model.Item, error) {
	var i model.Item
	err := r.pool.QueryRow(ctx, `SELECT id,title,description,category,location,status,reported_by,created_at,updated_at FROM items WHERE id=$1`, id).Scan(&i.ID, &i.Title, &i.Description, &i.Category, &i.Location, &i.Status, &i.ReportedBy, &i.CreatedAt, &i.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return i, ErrNotFound
	}
	return i, err
}
func (r *itemRepo) Create(ctx context.Context, i model.Item) (model.Item, error) {
	err := r.pool.QueryRow(ctx, `INSERT INTO items(title,description,category,location,status,reported_by) VALUES($1,$2,$3,$4,$5,$6) RETURNING id,title,description,category,location,status,reported_by,created_at,updated_at`, i.Title, i.Description, i.Category, i.Location, i.Status, i.ReportedBy).Scan(&i.ID, &i.Title, &i.Description, &i.Category, &i.Location, &i.Status, &i.ReportedBy, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}
func (r *itemRepo) Update(ctx context.Context, i model.Item) (model.Item, error) {
	err := r.pool.QueryRow(ctx, `UPDATE items SET title=$1,description=$2,category=$3,location=$4,status=$5,updated_at=NOW() WHERE id=$6 RETURNING id,title,description,category,location,status,reported_by,created_at,updated_at`, i.Title, i.Description, i.Category, i.Location, i.Status, i.ID).Scan(&i.ID, &i.Title, &i.Description, &i.Category, &i.Location, &i.Status, &i.ReportedBy, &i.CreatedAt, &i.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return i, ErrNotFound
	}
	return i, err
}
func (r *itemRepo) Delete(ctx context.Context, id int64) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM items WHERE id=$1`, id)
	if err == nil && cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

var _ = time.Now
