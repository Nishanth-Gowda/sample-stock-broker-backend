package services

import (
	"broker-backend/internal/domain/model"
	"broker-backend/internal/domain/repository"
	"context"
	"log"
)

type HoldingService interface {
	Create(ctx context.Context, userID int64, h model.Holding) (*model.Holding, error)
	List(ctx context.Context, userID int64) ([]model.Holding, error)
}

type holdingService struct {
	repo repository.HoldingRepository
}

func NewHoldingService(repo repository.HoldingRepository) HoldingService {
	return &holdingService{repo: repo}
}

func (s *holdingService) Create(ctx context.Context, userID int64, h model.Holding) (*model.Holding, error) {
	h.UserID = userID
	h.Value = h.Average * h.Quantity // simple calculation
	if err := s.repo.Create(ctx, &h); err != nil {
		log.Printf("holding_service: failed to create holding for user %d: %v", userID, err)
		return nil, err
	}
	log.Printf("holding_service: created holding %d for user %d (%s x%.2f)", h.ID, userID, h.Symbol, h.Quantity)
	return &h, nil
}

func (s *holdingService) List(ctx context.Context, userID int64) ([]model.Holding, error) {
	list, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		log.Printf("holding_service: failed to list holdings for user %d: %v", userID, err)
		return nil, err
	}
	log.Printf("holding_service: fetched %d holdings for user %d", len(list), userID)
	return list, nil
}
