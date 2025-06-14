package model

// Holding represents a user's asset holding.
type Holding struct {
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Average  float64 `json:"average_price"`
	Value    float64 `json:"value"`
}

// OrderBookEntry represents an order with PNL.
type OrderBookEntry struct {
	Symbol        string  `json:"symbol"`
	Quantity      float64 `json:"quantity"`
	Side          string  `json:"side"` // BUY or SELL
	Price         float64 `json:"price"`
	UnrealizedPNL float64 `json:"unrealized_pnl"`
	RealizedPNL   float64 `json:"realized_pnl"`
}

// Position represents an open position with PNL card.
type Position struct {
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Entry    float64 `json:"entry_price"`
	Current  float64 `json:"current_price"`
	PNL      float64 `json:"pnl"`
	PNLPct   float64 `json:"pnl_pct"`
}
