package handlers

import (
	"broker-backend/internal/domain/model"
	"broker-backend/internal/services"
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

type PositionHandler struct{ svc services.PositionService }

func NewPositionHandler(svc services.PositionService) *PositionHandler {
	return &PositionHandler{svc: svc}
}

type createPositionRequest struct {
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	Entry    float64 `json:"entry_price"`
	Current  float64 `json:"current_price"`
}

func (h *PositionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createPositionRequest
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

	pos := model.Position{Symbol: req.Symbol, Quantity: req.Quantity, Entry: req.Entry, Current: req.Current}
	created, err := h.svc.Create(r.Context(), userID, pos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *PositionHandler) List(w http.ResponseWriter, r *http.Request) {
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
