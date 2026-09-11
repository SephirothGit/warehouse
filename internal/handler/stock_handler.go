package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SephirothGit/warehouse/internal/service"
)

type StockHandler struct {
	stockService service.StockService
}

type moveStockRequest struct {
	FromShelfID int `json:"from_shelf_id"`
	ToShelfID   int `json:"to_shelf_id"`
	ProductID   int `json:"product_id"`
	Quantity    int `json:"quantity"`
}
type addStockRequest struct {
	ShelfID   int `json:"shelf_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func NewStockHandler(stockService service.StockService) *StockHandler {
	return &StockHandler{
		stockService: stockService,
	}
}

func (h *StockHandler) MoveHandler(w http.ResponseWriter, r *http.Request) {
	var req moveStockRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userIDStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusInternalServerError)
		return
	}

	err = h.stockService.MoveStock(r.Context(), req.FromShelfID, req.ToShelfID, req.ProductID, req.Quantity, userID)
	if err != nil {
		http.Error(w, "unable to move stock", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *StockHandler) AddHandler(w http.ResponseWriter, r *http.Request) {
	var req addStockRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userIDStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "user id not found in context", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusInternalServerError)
		return
	}

	err = h.stockService.AddStock(r.Context(), req.ShelfID, req.ProductID, req.Quantity, userID)
	if err != nil {
		http.Error(w, "unable to add stock", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *StockHandler) GetByShelfHandler(w http.ResponseWriter, r *http.Request) {
	shelfIDStr := r.URL.Query().Get("shelf_id")
	shelfID, err := strconv.Atoi(shelfIDStr)
	if err != nil {
		http.Error(w, "invalid shelf_id", http.StatusBadRequest)
		return
	}

	items, err := h.stockService.GetByShelf(r.Context(), shelfID)
	if err != nil {
		http.Error(w, "unable to fetch stock items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
