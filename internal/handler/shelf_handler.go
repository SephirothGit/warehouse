package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SephirothGit/warehouse/internal/service"
)

type ShelfHandler struct {
	shelfService service.ShelfService
}

type createShelfRequest struct {
	RackID int `json:"rack_id"`
	Level  int `json:"level"`
}

func NewShelfHandler(shelfService service.ShelfService) *ShelfHandler {
	return &ShelfHandler{
		shelfService: shelfService,
	}
}

func (h *ShelfHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var req createShelfRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondError(w, err, "invalid request body")
		return
	}

	id, err := h.shelfService.CreateShelf(r.Context(), req.RackID, req.Level)
	if err != nil {
		respondError(w, err, "unable to create a shelf")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *ShelfHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	rackIDStr := r.URL.Query().Get("rack_id")
	rackID, err := strconv.Atoi(rackIDStr)
	if err != nil {
		respondError(w, err, "unable to fetch shelves")
		return
	}

	shelves, err := h.shelfService.GetByRack(r.Context(), rackID)
	if err != nil {
		respondError(w, err, "unable to fetch shelves")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shelves)
}
