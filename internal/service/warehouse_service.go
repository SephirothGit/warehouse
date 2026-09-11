package service

import (
	"context"

	"github.com/SephirothGit/warehouse/internal/repository"
)

type warehouseService struct {
	repo repository.WarehouseRepo
}

type WarehouseService interface {
	CreateWarehouse(ctx context.Context, name, address string) (int, error)
	ListWarehouses(ctx context.Context) ([]repository.Warehouse, error)
}

func NewWarehouseService(repo repository.WarehouseRepo) WarehouseService {
	return &warehouseService{
		repo: repo,
	}
}

func(w *warehouseService) CreateWarehouse(ctx context.Context, name string, address string) (int, error) {
	return w.repo.Create(ctx, name, address)
}

func(w *warehouseService) ListWarehouses(ctx context.Context) ([]repository.Warehouse, error) {
	return w.repo.GetAll(ctx)
}