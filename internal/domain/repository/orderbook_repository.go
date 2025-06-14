package repository

import (
	"broker-backend/internal/domain/model"
	"context"
)

type OrderBookRepository interface {
	Create(ctx context.Context, o *model.OrderBookEntry) error
	ListByUser(ctx context.Context, userID int64) ([]model.OrderBookEntry, error)
}
