package repository

import (
	"context"
	"database/sql"
)

type BaseRepository struct {
	db *sql.DB
}

func (r *BaseRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *BaseRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
