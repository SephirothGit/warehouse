package service

import "github.com/SephirothGit/warehouse/internal/repository"

type zoneService struct {
	repo repository.ZoneRepo
}

type ZoneService interface {
	CreateZone(warehouseID int, name string) (int, error)
	ListZonesByWarehouse(warehouseID int) ([]repository.Zone, error)
}

func NewZoneService(repo repository.ZoneRepo) ZoneService {
	return &zoneService{
		repo: repo,
	}
}

func (z *zoneService) CreateZone(warehouseID int, name string) (int, error) {
	return z.repo.Create(warehouseID, name)
}

func (z *zoneService) ListZonesByWarehouse(warehouseID int) ([]repository.Zone, error) {
	return z.repo.GetByWarehouse(warehouseID)
}
