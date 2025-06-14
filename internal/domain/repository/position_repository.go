package repository

import (
	"broker-backend/internal/domain/model"
	"context"
)

type PositionRepository interface {
	Create(ctx context.Context, p *model.Position) error
	ListByUser(ctx context.Context, userID int64) ([]model.Position, error)
}
