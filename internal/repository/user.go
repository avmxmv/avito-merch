package repository

import (
	"avito-merch/internal/model"
	"context"
	"database/sql"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	UpdateCoins(ctx context.Context, tx *sql.Tx, id int, coins int) error
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(db *sql.DB) UserRepository {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Create(ctx context.Context, tx *sql.Tx, user *model.User) error {
	query := `INSERT INTO users (username, password_hash, coins) VALUES ($1, $2, $3)`
	_, err := tx.ExecContext(ctx, query, user.Username, user.PasswordHash, user.Coins)
	return err
}

func (r *UserPostgres) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `SELECT id, username, password_hash, coins FROM users WHERE id = $1`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserPostgres) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT id, username, password_hash, coins FROM users WHERE username = $1`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Coins)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserPostgres) UpdateCoins(ctx context.Context, tx *sql.Tx, id int, coins int) error {
	query := `UPDATE users SET coins = $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, coins, id)
	return err
}

func (r *UserPostgres) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
