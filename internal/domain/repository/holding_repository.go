package repository

import (
	"broker-backend/internal/domain/model"
	"context"
)

type HoldingRepository interface {
	Create(ctx context.Context, h *model.Holding) error
	ListByUser(ctx context.Context, userID int64) ([]model.Holding, error)
}
