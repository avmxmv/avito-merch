package repository

import (
	"avito-merch/internal/model"
	"context"
	"database/sql"
)

type PurchaseRepository interface {
	Create(ctx context.Context, tx *sql.Tx, purchase *model.Purchase) error
	GetByUserID(ctx context.Context, userID int) ([]model.Purchase, error)
}

type PurchasePostgres struct {
	db *sql.DB
}

func NewPurchasePostgres(db *sql.DB) PurchaseRepository {
	return &PurchasePostgres{db: db}
}

func (r *PurchasePostgres) Create(ctx context.Context, tx *sql.Tx, purchase *model.Purchase) error {
	query := `
		INSERT INTO purchases(user_id, merch_id, quantity) 
		VALUES($1, $2, $3) 
		RETURNING id, created_at`
	return tx.QueryRowContext(ctx, query,
		purchase.UserID,
		purchase.MerchID,
		purchase.Quantity,
	).Scan(&purchase.ID, &purchase.CreatedAt)
}

func (r *PurchasePostgres) GetByUserID(ctx context.Context, userID int) ([]model.Purchase, error) {
	query := `
		SELECT p.id, p.user_id, p.merch_id, p.quantity, p.created_at,
		       m.name, m.price
		FROM purchases p
		JOIN merch m ON p.merch_id = m.id
		WHERE p.user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []model.Purchase
	for rows.Next() {
		var p model.Purchase
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.MerchID,
			&p.Quantity,
			&p.CreatedAt,
			&p.MerchName,
			&p.Price,
		)
		if err != nil {
			return nil, err
		}
		purchases = append(purchases, p)
	}
	return purchases, nil
}
