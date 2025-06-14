package services

import (
	"broker-backend/internal/domain/model"
	"broker-backend/internal/domain/repository"
	"context"
	"log"
)

type PositionService interface {
	Create(ctx context.Context, userID int64, p model.Position) (*model.Position, error)
	List(ctx context.Context, userID int64) ([]model.Position, error)
}

type positionService struct{ repo repository.PositionRepository }

func NewPositionService(repo repository.PositionRepository) PositionService {
	return &positionService{repo: repo}
}

func (s *positionService) Create(ctx context.Context, userID int64, p model.Position) (*model.Position, error) {
	p.UserID = userID
	// Calculate PNL simple example
	p.PNL = (p.Current - p.Entry) * p.Quantity
	p.PNLPct = (p.Current - p.Entry) / p.Entry * 100
	if err := s.repo.Create(ctx, &p); err != nil {
		log.Printf("position_service: failed create position for user %d: %v", userID, err)
		return nil, err
	}
	log.Printf("position_service: created position %d for user %d (%s)", p.ID, userID, p.Symbol)
	return &p, nil
}

func (s *positionService) List(ctx context.Context, userID int64) ([]model.Position, error) {
	list, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		log.Printf("position_service: list error user %d: %v", userID, err)
		return nil, err
	}
	return list, nil
}
