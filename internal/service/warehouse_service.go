package service

import "github.com/SephirothGit/warehouse/internal/repository"

type warehouseService struct {
	repo repository.WarehouseRepo
}

type WarehouseService interface {
	CreateWarehouse(name, address string) (int, error)
	ListWarehouses() ([]repository.Warehouse, error)
}

func NewWarehouseService(repo repository.WarehouseRepo) WarehouseService {
	return &warehouseService{
		repo: repo,
	}
}

func(w *warehouseService) CreateWarehouse(name string, address string) (int, error) {
	return w.repo.Create(name, address)
}

func(w *warehouseService) ListWarehouses() ([]repository.Warehouse, error) {
	return w.repo.GetAll()
}