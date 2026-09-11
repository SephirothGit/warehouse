package service

import (
	"context"

	"github.com/SephirothGit/warehouse/internal/repository"
)

type zoneService struct {
	repo repository.ZoneRepo
}

type ZoneService interface {
	CreateZone(ctx context.Context, warehouseID int, name string) (int, error)
	ListZonesByWarehouse(ctx context.Context, warehouseID int) ([]repository.Zone, error)
}

func NewZoneService(repo repository.ZoneRepo) ZoneService {
	return &zoneService{
		repo: repo,
	}
}

func (z *zoneService) CreateZone(ctx context.Context, warehouseID int, name string) (int, error) {
	return z.repo.Create(ctx, warehouseID, name)
}

func (z *zoneService) ListZonesByWarehouse(ctx context.Context, warehouseID int) ([]repository.Zone, error) {
	return z.repo.GetByWarehouse(ctx, warehouseID)
}
