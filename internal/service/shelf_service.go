package service

import (
	"context"

	"github.com/SephirothGit/warehouse/internal/repository"
)

type shelfService struct {
	repo repository.ShelfRepo
}

type ShelfService interface {
	CreateShelf(ctx context.Context, rackID int, level int) (int, error)
	GetByRack(ctx context.Context, rackID int) ([]repository.Shelf, error)
}

func NewShelfService(repo repository.ShelfRepo) ShelfService {
	return &shelfService{
		repo: repo,
	}
}

func (s *shelfService) CreateShelf(ctx context.Context, rackID int, level int) (int, error) {
	return s.repo.Create(ctx, rackID, level)
}

func (s *shelfService) GetByRack(ctx context.Context, rackID int) ([]repository.Shelf, error) {
	return s.repo.GetByRack(ctx, rackID)
}
