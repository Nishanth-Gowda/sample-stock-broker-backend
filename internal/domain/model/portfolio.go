package model

// Holding represents a user's asset holding.
type Holding struct {
	ID        int64   `json:"id" db:"id"`
	UserID    int64   `json:"user_id" db:"user_id"`
	Symbol    string  `json:"symbol" db:"symbol"`
	Quantity  float64 `json:"quantity" db:"quantity"`
	Average   float64 `json:"average_price" db:"average_price"`
	Value     float64 `json:"value" db:"value"`
	CreatedAt string  `json:"created_at" db:"created_at"`
}

// OrderBookEntry represents an order with PNL.
type OrderBookEntry struct {
	ID            int64   `json:"id" db:"id"`
	UserID        int64   `json:"user_id" db:"user_id"`
	Symbol        string  `json:"symbol" db:"symbol"`
	Quantity      float64 `json:"quantity" db:"quantity"`
	Side          string  `json:"side" db:"side"` // BUY or SELL
	Price         float64 `json:"price" db:"price"`
	UnrealizedPNL float64 `json:"unrealized_pnl" db:"unrealized_pnl"`
	RealizedPNL   float64 `json:"realized_pnl" db:"realized_pnl"`
	CreatedAt     string  `json:"created_at" db:"created_at"`
}

// Position represents an open position with PNL card.
type Position struct {
	ID        int64   `json:"id" db:"id"`
	UserID    int64   `json:"user_id" db:"user_id"`
	Symbol    string  `json:"symbol" db:"symbol"`
	Quantity  float64 `json:"quantity" db:"quantity"`
	Entry     float64 `json:"entry_price" db:"entry_price"`
	Current   float64 `json:"current_price" db:"current_price"`
	PNL       float64 `json:"pnl" db:"pnl"`
	PNLPct    float64 `json:"pnl_pct" db:"pnl_pct"`
	CreatedAt string  `json:"created_at" db:"created_at"`
}
