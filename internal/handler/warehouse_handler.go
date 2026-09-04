package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SephirothGit/warehouse/internal/service"
)

type WarehouseHandler struct {
	warehouseService service.WarehouseService
}

type createWarehouseRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func NewWarehouseHandler(warehouseService service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{
		warehouseService: warehouseService,
	}
}

func (h *WarehouseHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var req createWarehouseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.warehouseService.CreateWarehouse(req.Name, req.Address)
	if err != nil {
		http.Error(w, "unable to create warehouse", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *WarehouseHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	warehouses, err := h.warehouseService.ListWarehouses()
	if err != nil {
		http.Error(w, "unable to fetch warehouses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(warehouses)
}
