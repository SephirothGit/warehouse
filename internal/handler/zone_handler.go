package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SephirothGit/warehouse/internal/service"
)

type ZoneHandler struct {
	zoneService service.ZoneService
}

type createZoneRequest struct {
	WarehouseID int    `json:"warehouse_id"`
	Name        string `json:"name"`
}

func NewZoneHandler(zoneService service.ZoneService) *ZoneHandler {
	return &ZoneHandler{
		zoneService: zoneService,
	}
}

func (z *ZoneHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var req createZoneRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondError(w, err, "invalid request body")
		return
	}

	id, err := z.zoneService.CreateZone(r.Context(), req.WarehouseID, req.Name)
	if err != nil {
		respondError(w, err, "unable to create zone")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (z *ZoneHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	warehouseIDStr := r.URL.Query().Get("warehouse_id")
	warehouseID, err := strconv.Atoi(warehouseIDStr)
	if err != nil {
		respondError(w, err, "invalid warehouse_id")
		return
	}

	zones, err := z.zoneService.ListZonesByWarehouse(r.Context(), warehouseID)
	if err != nil {
		respondError(w, err, "unable to fetch zones")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zones)
}
