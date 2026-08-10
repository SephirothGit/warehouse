package repository

import "database/sql"

type WarehouseRepo interface {
	Create(name, address string) (int, error)
}
type warehouseRepo struct {
	db *sql.DB
}

func NewWarehouseRepo(db *sql.DB) WarehouseRepo {
	return &warehouseRepo{
		db: db,
	}
}

func (w *warehouseRepo) Create(name, address string) (int, error) {
	row := w.db.QueryRow("INSERT INTO warehouses (name, address) VALUES ($1, $2) RETURNING id", name, address)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
