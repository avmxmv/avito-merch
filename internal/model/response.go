package model

import "time"

type InfoResponse struct {
	Coins     int                `json:"coins"`
	Inventory []InventoryItem    `json:"inventory"`
	History   TransactionHistory `json:"history"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type TransactionHistory struct {
	Received []ReceivedTransaction `json:"received"`
	Sent     []SentTransaction     `json:"sent"`
}

type ReceivedTransaction struct {
	FromUser   string    `json:"from_user"`
	Amount     int       `json:"amount"`
	ReceivedAt time.Time `json:"received_at"`
}

type SentTransaction struct {
	ToUser string    `json:"to_user"`
	Amount int       `json:"amount"`
	SentAt time.Time `json:"sent_at"`
}
