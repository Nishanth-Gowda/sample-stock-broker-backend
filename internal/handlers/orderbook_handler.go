package handlers

import (
	"broker-backend/internal/domain/model"
	"broker-backend/internal/services"
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

type OrderBookHandler struct{ svc services.OrderBookService }

func NewOrderBookHandler(svc services.OrderBookService) *OrderBookHandler {
	return &OrderBookHandler{svc: svc}
}

type createOrderRequest struct {
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Side     string  `json:"side"` // BUY or SELL
	Price    float64 `json:"price"`
}

func (h *OrderBookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	uidAny := claims["user_id"]
	var userID int64
	switch v := uidAny.(type) {
	case float64:
		userID = int64(v)
	case int64:
		userID = v
	}

	entry := model.OrderBookEntry{Symbol: req.Symbol, Quantity: req.Quantity, Side: req.Side, Price: req.Price, UnrealizedPNL: 0, RealizedPNL: 0}
	created, err := h.svc.Create(r.Context(), userID, entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *OrderBookHandler) List(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	uidAny := claims["user_id"]
	var userID int64
	switch v := uidAny.(type) {
	case float64:
		userID = int64(v)
	case int64:
		userID = v
	}

	list, err := h.svc.List(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
