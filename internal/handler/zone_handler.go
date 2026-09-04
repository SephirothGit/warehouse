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
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}

	id, err := z.zoneService.CreateZone(req.WarehouseID, req.Name)
	if err != nil {
		http.Error(w, "unable to create zone", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (z *ZoneHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	warehouseIDStr := r.URL.Query().Get("warehouse_id")
	warehouseID, err := strconv.Atoi(warehouseIDStr)
	if err != nil {
		http.Error(w, "invalid warehouse_id", http.StatusBadRequest)
		return
	}

	zones, err := z.zoneService.ListZonesByWarehouse(warehouseID)
	if err != nil {
		http.Error(w, "unable to fetch zones", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zones)
}
