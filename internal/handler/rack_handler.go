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
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.rackService.CreateRack(req.ZoneID, req.Code)
	if err != nil {
		http.Error(w, "unable to create rack", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id)
}

func (h *RackHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	zoneIDStr := r.URL.Query().Get("zone_id")
	zoneID, err := strconv.Atoi(zoneIDStr)
	if err != nil {
		http.Error(w, "invalid zone_id", http.StatusBadRequest)
		return
	}

	racks, err := h.rackService.GetByZone(zoneID)
	if err != nil {
		http.Error(w, "unable to fetch racks", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(racks)
}
