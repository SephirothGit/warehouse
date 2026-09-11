package service

import (
	"context"

	"github.com/SephirothGit/warehouse/internal/repository"
)

type stockService struct {
	repo repository.StockRepo
}

type StockService interface {
	GetByShelf(ctx context.Context, shelfID int) ([]repository.StockItem, error)
	MoveStock(ctx context.Context, fromShelfID, toShelfID, productID, quantity, userID int) error
	AddStock(ctx context.Context, shelfID, productID, quantity, userID int) error
}

func NewStockService(repo repository.StockRepo) StockService {
	return &stockService{
		repo: repo,
	}
}

func (s *stockService) GetByShelf(ctx context.Context, shelfID int) ([]repository.StockItem, error) {
	return s.repo.GetByShelf(ctx, shelfID)
}

func (s *stockService) MoveStock(ctx context.Context, fromShelfID, toShelfID, productID, quantity, userID int) error {
	return s.repo.MoveStock(ctx, fromShelfID, toShelfID, productID, quantity, userID)
}

func (s *stockService) AddStock(ctx context.Context, shelfID, productID, quantity, userID int) error {
	return s.repo.AddStock(ctx, shelfID, productID, quantity, userID)
}
