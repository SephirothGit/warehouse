package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type ProductRepo interface {
	Create(ctx context.Context, sku, name, unit string) (int, error)
	GetAll(ctx context.Context) ([]Product, error)
}

type Product struct {
	ID   int
	SKU  string
	Name string
	Unit string
}

type productRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) ProductRepo {
	return &productRepo{
		db: db,
	}
}

func (p *productRepo) Create(ctx context.Context, sku, name, unit string) (int, error) {
	row := p.db.QueryRowContext(ctx, "INSERT INTO products (sku, name, unit) VALUES ($1, $2, $3) RETURNING id", sku, name, unit)

	var id int
	err := row.Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, ErrAlreadyExists
		}
		return 0, err
	}
	return id, nil
}

func (p *productRepo) GetAll(ctx context.Context) ([]Product, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, sku, name, unit FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Product
	for rows.Next() {
		var item Product
		err := rows.Scan(&item.ID, &item.SKU, &item.Name, &item.Unit)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
