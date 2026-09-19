package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SephirothGit/warehouse/internal/repository"
	"github.com/SephirothGit/warehouse/internal/service"
)

type ProductHandler struct {
	productService service.ProductService
}

type createProductRequest struct {
	SKU  string `json:"sku"`
	Name string `json:"name"`
	Unit string `json:"unit"`
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.productService.Create(r.Context(), req.SKU, req.Name, req.Unit)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			http.Error(w, "sku already exists", http.StatusConflict)
			return
		}
		http.Error(w, "unable to create product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *ProductHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetAll(r.Context())
	if err != nil {
		http.Error(w, "unable to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
