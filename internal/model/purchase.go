package model

import "time"

type Purchase struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id" validate:"required"`
	MerchID   int       `json:"merch_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"min=1"`
	CreatedAt time.Time `json:"created_at"`

	// Для отображения в инвентаре
	MerchName string `json:"merch_name,omitempty"`
	Price     int    `json:"price,omitempty"`
}
