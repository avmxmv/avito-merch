package repository

import (
	"avito-merch/internal/model"
	"context"
	"database/sql"
	"errors"
)

var ErrMerchNotFound = errors.New("merch not found")

type MerchRepository interface {
	GetByName(ctx context.Context, name string) (*model.Merch, error)
}

type MerchPostgres struct {
	db *sql.DB
}

func NewMerchPostgres(db *sql.DB) MerchRepository {
	return &MerchPostgres{db: db}
}

func (r *MerchPostgres) GetByName(ctx context.Context, name string) (*model.Merch, error) {
	query := `SELECT id, name, price FROM merch WHERE name = $1`
	var merch model.Merch
	err := r.db.QueryRowContext(ctx, query, name).Scan(&merch.ID, &merch.Name, &merch.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMerchNotFound
		}
		return nil, err
	}
	return &merch, nil
}
