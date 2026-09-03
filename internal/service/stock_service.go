package service

import "github.com/SephirothGit/warehouse/internal/repository"

type stockService struct {
	repo repository.StockRepo
}

type StockService interface {
	GetByShelf(shelfID int) ([]repository.StockItem, error)
	MoveStock(fromShelfID, toShelfID, productID, quantity, userID int) error
	AddStock(shelfID, productID, quantity, userID int) error
}

func NewStockService(repo repository.StockRepo) StockService {
	return &stockService{
		repo: repo,
	}
}

func (s *stockService) GetByShelf(shelfID int) ([]repository.StockItem, error) {
	return s.repo.GetByShelf(shelfID)
}

func (s *stockService) MoveStock(fromShelfID, toShelfID, productID, quantity, userID int) error {
	return s.repo.MoveStock(fromShelfID, toShelfID, productID, quantity, userID)
}

func (s *stockService) AddStock(shelfID, productID, quantity, userID int) error {
	return s.repo.AddStock(shelfID, productID, quantity, userID)
}
