package service

import "errors"

var (
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrInvalidRecipient   = errors.New("invalid recipient")
	ErrSameUser           = errors.New("cannot send to yourself")
	ErrItemNotFound       = errors.New("item not found")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)
