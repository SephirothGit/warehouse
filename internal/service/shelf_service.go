package service

import "github.com/SephirothGit/warehouse/internal/repository"

type shelfService struct {
	repo repository.ShelfRepo
}

type ShelfService interface {
	CreateShelf(rackID int, level int) (int, error)
	GetByRack(rackID int) ([]repository.Shelf, error)
}

func NewShelfService(repo repository.ShelfRepo) ShelfService {
	return &shelfService{
		repo: repo,
	}
}

func (s *shelfService) CreateShelf(rackID int, level int) (int, error) {
	return s.repo.Create(rackID, level)
}

func (s *shelfService) GetByRack(rackID int) ([]repository.Shelf, error) {
	return s.repo.GetByRack(rackID)
}
