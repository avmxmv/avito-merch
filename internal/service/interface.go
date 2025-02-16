package service

import (
	"avito-merch/internal/model"
	"context"
)

type AuthService interface {
	Authenticate(ctx context.Context, username, password string) (string, error)
	ParseToken(tokenString string) (int, error)
}

type BuyService interface {
	BuyItem(ctx context.Context, userID int, merchName string) error
}

type SendService interface {
	SendCoins(ctx context.Context, senderID int, recipientUsername string, amount int) error
}

type InfoService interface {
	GetUserInfo(ctx context.Context, userID int) (*model.InfoResponse, error)
}
