package mysql

import (
	"broker-backend/internal/domain/model"
	"broker-backend/internal/domain/repository"
	"context"

	"github.com/jmoiron/sqlx"
)

var _ repository.HoldingRepository = (*HoldingRepository)(nil)

type HoldingRepository struct {
	db *sqlx.DB
}

func NewHoldingRepository(db *sqlx.DB) *HoldingRepository {
	return &HoldingRepository{db: db}
}

func (r *HoldingRepository) Create(ctx context.Context, h *model.Holding) error {
	query := `INSERT INTO holdings(user_id, symbol, quantity, average_price, value) VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, h.UserID, h.Symbol, h.Quantity, h.Average, h.Value)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err == nil {
		h.ID = id
	}
	return err
}

func (r *HoldingRepository) ListByUser(ctx context.Context, userID int64) ([]model.Holding, error) {
	var holdings []model.Holding
	query := `SELECT id, user_id, symbol, quantity, average_price, value, created_at FROM holdings WHERE user_id=?`
	if err := r.db.SelectContext(ctx, &holdings, query, userID); err != nil {
		return nil, err
	}
	return holdings, nil
}
