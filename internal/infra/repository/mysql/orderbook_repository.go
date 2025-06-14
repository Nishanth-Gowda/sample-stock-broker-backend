package mysql

import (
	"broker-backend/internal/domain/model"
	"broker-backend/internal/domain/repository"
	"context"

	"github.com/jmoiron/sqlx"
)

var _ repository.OrderBookRepository = (*OrderBookRepository)(nil)

type OrderBookRepository struct{ db *sqlx.DB }

func NewOrderBookRepository(db *sqlx.DB) *OrderBookRepository { return &OrderBookRepository{db: db} }

func (r *OrderBookRepository) Create(ctx context.Context, o *model.OrderBookEntry) error {
	q := `INSERT INTO orderbook(user_id, symbol, quantity, side, price, unrealized_pnl, realized_pnl) VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, q, o.UserID, o.Symbol, o.Quantity, o.Side, o.Price, o.UnrealizedPNL, o.RealizedPNL)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		o.ID = id
	}
	return err
}

func (r *OrderBookRepository) ListByUser(ctx context.Context, userID int64) ([]model.OrderBookEntry, error) {
	var list []model.OrderBookEntry
	q := `SELECT id, user_id, symbol, quantity, side, price, unrealized_pnl, realized_pnl, created_at FROM orderbook WHERE user_id=?`
	if err := r.db.SelectContext(ctx, &list, q, userID); err != nil {
		return nil, err
	}
	return list, nil
}
