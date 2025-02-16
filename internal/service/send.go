package service

import (
	"avito-merch/internal/model"
	"avito-merch/internal/repository"
	"context"
)

type sendService struct {
	userRepo  repository.UserRepository
	transRepo repository.TransactionRepository
}

func NewSendService(
	userRepo repository.UserRepository,
	transRepo repository.TransactionRepository,
) SendService {
	return &sendService{
		userRepo:  userRepo,
		transRepo: transRepo,
	}
}

func (s *sendService) SendCoins(ctx context.Context, senderID int, recipientUsername string, amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	recipient, err := s.userRepo.GetByUsername(ctx, recipientUsername)
	if err != nil {
		return ErrInvalidRecipient
	}

	if senderID == recipient.ID {
		return ErrSameUser
	}

	tx, err := s.userRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sender, err := s.userRepo.GetByID(ctx, senderID)
	if err != nil {
		return err
	}

	if sender.Coins < amount {
		return ErrInsufficientFunds
	}

	if err := s.userRepo.UpdateCoins(ctx, tx, sender.ID, sender.Coins-amount); err != nil {
		return err
	}

	if err := s.userRepo.UpdateCoins(ctx, tx, recipient.ID, recipient.Coins+amount); err != nil {
		return err
	}

	transaction := &model.Transaction{
		FromUser: &sender.ID,
		ToUser:   recipient.ID,
		Amount:   amount,
		Type:     model.TransactionTypeTransfer,
	}

	return s.transRepo.Create(ctx, tx, transaction)
}
