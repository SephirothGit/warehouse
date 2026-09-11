package service

import (
	"context"

	"github.com/SephirothGit/warehouse/internal/repository"
)

type rackService struct {
	repo repository.RackRepo
}

type RackService interface {
	CreateRack(ctx context.Context, zoneID int, code string) (int, error)
	GetByZone(ctx context.Context, zoneID int) ([]repository.Rack, error)
}

func NewRackService(repo repository.RackRepo) RackService {
	return &rackService{
		repo: repo,
	}
}

func (r *rackService) CreateRack(ctx context.Context, zoneID int, code string) (int, error) {
	return r.repo.Create(ctx, zoneID, code)
}

func (r *rackService) GetByZone(ctx context.Context, zoneID int) ([]repository.Rack, error) {
	return r.repo.GetByZone(ctx, zoneID)
}
