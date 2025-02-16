package repository

import (
	"avito-merch/internal/model"
	"context"
	"database/sql"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *sql.Tx, transaction *model.Transaction) error
	GetUserHistory(ctx context.Context, userID int) ([]model.Transaction, error)
}

type TransactionPostgres struct {
	db *sql.DB
}

func NewTransactionPostgres(db *sql.DB) TransactionRepository {
	return &TransactionPostgres{db: db}
}

func (r *TransactionPostgres) Create(ctx context.Context, tx *sql.Tx, transaction *model.Transaction) error {
	query := `
		INSERT INTO transactions(from_user, to_user, amount, type) 
		VALUES($1, $2, $3, $4) 
		RETURNING id, created_at`
	return tx.QueryRowContext(ctx, query,
		transaction.FromUser,
		transaction.ToUser,
		transaction.Amount,
		transaction.Type,
	).Scan(&transaction.ID, &transaction.CreatedAt)
}

func (r *TransactionPostgres) GetUserHistory(ctx context.Context, userID int) ([]model.Transaction, error) {
	query := `
		SELECT t.id, t.from_user, t.to_user, t.amount, t.type, t.created_at,
		       u_from.username, u_to.username
		FROM transactions t
		LEFT JOIN users u_from ON t.from_user = u_from.id
		JOIN users u_to ON t.to_user = u_to.id
		WHERE t.from_user = $1 OR t.to_user = $1
		ORDER BY t.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		var fromUser sql.NullInt64
		err := rows.Scan(
			&t.ID,
			&fromUser,
			&t.ToUser,
			&t.Amount,
			&t.Type,
			&t.CreatedAt,
			&t.FromUsername,
			&t.ToUsername,
		)
		if err != nil {
			return nil, err
		}
		if fromUser.Valid {
			t.FromUser = new(int)
			*t.FromUser = int(fromUser.Int64)
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}
