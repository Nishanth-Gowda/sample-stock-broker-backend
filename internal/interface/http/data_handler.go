package handler

import (
	"broker-backend/internal/domain/model"
	"encoding/json"
	"net/http"
)

type DataHandler struct{}

func NewDataHandler() *DataHandler { return &DataHandler{} }

func (h *DataHandler) Holdings(w http.ResponseWriter, r *http.Request) {
	holdings := []model.Holding{
		{Symbol: "AAPL", Quantity: 10, Average: 150.0, Value: 1600},
		{Symbol: "GOOGL", Quantity: 5, Average: 2500.0, Value: 13000},
	}
	respondJSON(w, holdings)
}

func (h *DataHandler) OrderBook(w http.ResponseWriter, r *http.Request) {
	orders := []model.OrderBookEntry{
		{Symbol: "AAPL", Quantity: 10, Side: "BUY", Price: 150, UnrealizedPNL: 100, RealizedPNL: 0},
		{Symbol: "MSFT", Quantity: 5, Side: "SELL", Price: 300, UnrealizedPNL: -50, RealizedPNL: 20},
	}
	respondJSON(w, orders)
}

func (h *DataHandler) Positions(w http.ResponseWriter, r *http.Request) {
	positions := []model.Position{
		{Symbol: "AAPL", Quantity: 10, Entry: 150, Current: 160, PNL: 100, PNLPct: 6.67},
		{Symbol: "GOOGL", Quantity: 5, Entry: 2500, Current: 2600, PNL: 500, PNLPct: 4},
	}
	respondJSON(w, positions)
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
