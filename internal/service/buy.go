package service

import (
	"avito-merch/internal/model"
	"avito-merch/internal/repository"
	"context"
)

type buyService struct {
	userRepo     repository.UserRepository
	merchRepo    repository.MerchRepository
	purchaseRepo repository.PurchaseRepository
	transRepo    repository.TransactionRepository
}

func NewBuyService(
	userRepo repository.UserRepository,
	merchRepo repository.MerchRepository,
	purchaseRepo repository.PurchaseRepository,
	transRepo repository.TransactionRepository,
) BuyService {
	return &buyService{
		userRepo:     userRepo,
		merchRepo:    merchRepo,
		purchaseRepo: purchaseRepo,
		transRepo:    transRepo,
	}
}

func (s *buyService) BuyItem(ctx context.Context, userID int, merchName string) error {
	merch, err := s.merchRepo.GetByName(ctx, merchName)
	if err != nil {
		return ErrItemNotFound
	}

	tx, err := s.userRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Coins < merch.Price {
		return ErrInsufficientFunds
	}

	// Обновляем баланс
	if err := s.userRepo.UpdateCoins(ctx, tx, user.ID, user.Coins-merch.Price); err != nil {
		return err
	}

	// Создаем запись о покупке
	purchase := &model.Purchase{
		UserID:  user.ID,
		MerchID: merch.ID,
	}
	if err := s.purchaseRepo.Create(ctx, tx, purchase); err != nil {
		return err
	}

	// Записываем транзакцию
	transaction := &model.Transaction{
		ToUser: user.ID,
		Amount: merch.Price,
		Type:   model.TransactionTypePurchase,
	}
	return s.transRepo.Create(ctx, tx, transaction)
}
