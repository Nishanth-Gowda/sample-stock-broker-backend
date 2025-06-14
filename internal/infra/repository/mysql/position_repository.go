package mysql

import (
	"broker-backend/internal/domain/model"
	"broker-backend/internal/domain/repository"
	"context"

	"github.com/jmoiron/sqlx"
)

var _ repository.PositionRepository = (*PositionRepository)(nil)

type PositionRepository struct {
	db *sqlx.DB
}

func NewPositionRepository(db *sqlx.DB) *PositionRepository {
	return &PositionRepository{db: db}
}

func (r *PositionRepository) Create(ctx context.Context, p *model.Position) error {
	query := `INSERT INTO positions(user_id, symbol, quantity, entry_price, current_price, pnl, pnl_pct) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, p.UserID, p.Symbol, p.Quantity, p.Entry, p.Current, p.PNL, p.PNLPct)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err == nil {
		p.ID = id
	}
	return err
}

func (r *PositionRepository) ListByUser(ctx context.Context, userID int64) ([]model.Position, error) {
	var positions []model.Position
	query := `SELECT id, user_id, symbol, quantity, entry_price, current_price, pnl, pnl_pct, created_at FROM positions WHERE user_id=?`
	if err := r.db.SelectContext(ctx, &positions, query, userID); err != nil {
		return nil, err
	}
	return positions, nil
}
