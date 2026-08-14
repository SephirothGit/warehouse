package repository

import "database/sql"

type WarehouseRepo interface {
	Create(name, address string) (int, error)
	GetAll() ([]Warehouse, error)
}
type warehouseRepo struct {
	db *sql.DB
}

type Warehouse struct {
	ID      int
	Name    string
	Address string
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

func (w *warehouseRepo) GetAll() ([]Warehouse, error) {
	rows, err := w.db.Query("SELECT id, name, address FROM warehouses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Warehouse
	for rows.Next() {
		var item Warehouse
		err := rows.Scan(&item.ID, &item.Name, &item.Address)
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
