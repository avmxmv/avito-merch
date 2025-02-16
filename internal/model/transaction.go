package model

import "time"

type Transaction struct {
	ID        int       `json:"id"`
	FromUser  *int      `json:"from_user,omitempty"` // nil для системных операций
	ToUser    int       `json:"to_user" validate:"required"`
	Amount    int       `json:"amount" validate:"min=1"`
	Type      string    `json:"type" validate:"oneof=transfer purchase"` // transfer или purchase
	CreatedAt time.Time `json:"created_at"`

	// Для отображения в истории
	FromUsername string `json:"from_username,omitempty"`
	ToUsername   string `json:"to_username,omitempty"`
}
