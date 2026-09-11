package service

import (
	"context"

	"github.com/SephirothGit/warehouse/internal/repository"
)

type productService struct {
	repo repository.ProductRepo
}

type ProductService interface {
	Create(ctx context.Context, sku, name, unit string) (int, error)
	GetAll(ctx context.Context) ([]repository.Product, error)
}

func NewProductService(repo repository.ProductRepo) ProductService {
	return &productService{
		repo: repo,
	}
}

func (p *productService) Create(ctx context.Context, sku, name, unit string) (int, error) {
	return p.repo.Create(ctx, sku, name, unit)
}

func (p *productService) GetAll(ctx context.Context) ([]repository.Product, error) {
	return p.repo.GetAll(ctx)
}
