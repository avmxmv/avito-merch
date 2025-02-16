package service

import (
	"avito-merch/internal/model"
	"avito-merch/internal/repository"
	"context"
)

type infoService struct {
	userRepo     repository.UserRepository
	purchaseRepo repository.PurchaseRepository
	transRepo    repository.TransactionRepository
}

func NewInfoService(
	userRepo repository.UserRepository,
	purchaseRepo repository.PurchaseRepository,
	transRepo repository.TransactionRepository,
) InfoService {
	return &infoService{
		userRepo:     userRepo,
		purchaseRepo: purchaseRepo,
		transRepo:    transRepo,
	}
}

func (s *infoService) GetUserInfo(ctx context.Context, userID int) (*model.InfoResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	purchases, err := s.purchaseRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transRepo.GetUserHistory(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := &model.InfoResponse{
		Coins:     user.Coins,
		Inventory: make([]model.InventoryItem, 0),
		History: model.TransactionHistory{
			Received: make([]model.ReceivedTransaction, 0),
			Sent:     make([]model.SentTransaction, 0),
		},
	}

	itemCounts := make(map[string]*model.InventoryItem)
	for _, p := range purchases {
		if item, exists := itemCounts[p.MerchName]; exists {
			item.Quantity += p.Quantity
		} else {
			itemCounts[p.MerchName] = &model.InventoryItem{
				Type:     p.MerchName,
				Quantity: p.Quantity,
				Price:    p.Price,
			}
		}
	}

	for _, item := range itemCounts {
		response.Inventory = append(response.Inventory, *item)
	}

	for _, t := range transactions {
		switch {
		case t.FromUser != nil && *t.FromUser == userID:
			response.History.Sent = append(response.History.Sent, model.SentTransaction{
				ToUser: t.ToUsername,
				Amount: t.Amount,
				SentAt: t.CreatedAt,
			})
		case t.ToUser == userID:
			response.History.Received = append(response.History.Received, model.ReceivedTransaction{
				FromUser:   t.FromUsername,
				Amount:     t.Amount,
				ReceivedAt: t.CreatedAt,
			})
		}
	}

	return response, nil
}
