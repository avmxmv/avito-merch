package model

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username" validate:"required,min=3,max=50"`
	PasswordHash string    `json:"-"`
	Coins        int       `json:"coins" validate:"min=0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
