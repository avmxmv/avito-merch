package model

import "time"

const (
	InitialCoins      = 1000
	MaxUsernameLength = 50
	MinPasswordLength = 6

	TransactionTypeTransfer = "transfer"
	TransactionTypePurchase = "purchase"

	TokenExpiration = 24 * time.Hour
)
