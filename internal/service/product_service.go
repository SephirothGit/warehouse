package service

import "github.com/SephirothGit/warehouse/internal/repository"

type productService struct {
	repo repository.ProductRepo
}

type ProductService interface {
	Create(sku, name, unit string) (int, error)
	GetAll() ([]repository.Product, error)
}

func NewProductService(repo repository.ProductRepo) ProductService {
	return &productService{
		repo: repo,
	}
}

func (p *productService) Create(sku, name, unit string) (int, error) {
	return p.repo.Create(sku, name, unit)
}

func (p *productService) GetAll() ([]repository.Product, error) {
	return p.repo.GetAll()
}
