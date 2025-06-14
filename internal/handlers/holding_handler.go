package handlers

import (
	"broker-backend/internal/domain/model"
	"broker-backend/internal/services"
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

type HoldingHandler struct {
	svc services.HoldingService
}

func NewHoldingHandler(svc services.HoldingService) *HoldingHandler {
	return &HoldingHandler{svc: svc}
}

type createHoldingRequest struct {
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Average  float64 `json:"average_price"`
}

func (h *HoldingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createHoldingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// get user_id from jwt claims
	_, claims, _ := jwtauth.FromContext(r.Context())
	uidAny := claims["user_id"]
	var userID int64
	switch v := uidAny.(type) {
	case float64:
		userID = int64(v)
	case int64:
		userID = v
	}

	holding := model.Holding{
		Symbol:   req.Symbol,
		Quantity: req.Quantity,
		Average:  req.Average,
	}
	created, err := h.svc.Create(r.Context(), userID, holding)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *HoldingHandler) List(w http.ResponseWriter, r *http.Request) {
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
