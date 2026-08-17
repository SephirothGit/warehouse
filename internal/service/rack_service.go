package service

import "github.com/SephirothGit/warehouse/internal/repository"

type rackService struct {
	repo repository.RackRepo
}

type RackService interface {
	CreateRack(zoneID int, code string) (int, error)
	GetByZone(zoneID int) ([]repository.Rack, error)
}

func NewRackService(repo repository.RackRepo) RackService {
	return &rackService{
		repo: repo,
	}
}

func (r *rackService) CreateRack(zoneID int, code string) (int, error) {
	return r.repo.Create(zoneID, code)
}

func (r *rackService) GetByZone(zoneID int) ([]repository.Rack, error) {
	return r.repo.GetByZone(zoneID)
}
