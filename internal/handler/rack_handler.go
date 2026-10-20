package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SephirothGit/warehouse/internal/service"
)

type RackHandler struct {
	rackService service.RackService
}

type createRackRequest struct {
	ZoneID int    `json:"zone_id"`
	Code   string `json:"code"`
}

func NewRackHandler(rackService service.RackService) *RackHandler {
	return &RackHandler{
		rackService: rackService,
	}
}

func (h *RackHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var req createRackRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondError(w, err, "invalid request body")
		return
	}

	id, err := h.rackService.CreateRack(r.Context(), req.ZoneID, req.Code)
	if err != nil {
		respondError(w, err, "unable to create rack")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *RackHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	zoneIDStr := r.URL.Query().Get("zone_id")
	zoneID, err := strconv.Atoi(zoneIDStr)
	if err != nil {
		respondError(w, err, "invalid zone_id")
		return
	}

	racks, err := h.rackService.GetByZone(r.Context(), zoneID)
	if err != nil {
		respondError(w, err, "unable to fetch racks")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(racks)
}
