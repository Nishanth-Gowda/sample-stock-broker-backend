package services

import (
	"broker-backend/internal/domain/model"
	"broker-backend/internal/domain/repository"
	"context"
	"log"
)

type OrderBookService interface {
	Create(ctx context.Context, userID int64, o model.OrderBookEntry) (*model.OrderBookEntry, error)
	List(ctx context.Context, userID int64) ([]model.OrderBookEntry, error)
}

type orderBookService struct {
	repo repository.OrderBookRepository
}

func NewOrderBookService(repo repository.OrderBookRepository) OrderBookService {
	return &orderBookService{repo: repo}
}

func (s *orderBookService) Create(ctx context.Context, userID int64, o model.OrderBookEntry) (*model.OrderBookEntry, error) {
	o.UserID = userID
	// simple PNL calc placeholder (0)
	if err := s.repo.Create(ctx, &o); err != nil {
		log.Printf("orderbook_service: create failed user %d: %v", userID, err)
		return nil, err
	}
	log.Printf("orderbook_service: created order %d for user %d (%s %s x%.2f)", o.ID, userID, o.Side, o.Symbol, o.Quantity)
	return &o, nil
}

func (s *orderBookService) List(ctx context.Context, userID int64) ([]model.OrderBookEntry, error) {
	list, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		log.Printf("orderbook_service: list error user %d: %v", userID, err)
		return nil, err
	}
	return list, nil
}
