package repository

import "database/sql"

type ZoneRepo interface {
	Create(warehouseID int, name string) (int, error)
	GetByWarehouse(warehouseID int) ([]Zone, error)
}

type Zone struct {
	ID          int
	WarehouseID int
	Name        string
}

type zoneRepo struct {
	db *sql.DB
}

func NewZoneRepo(db *sql.DB) ZoneRepo {
	return &zoneRepo{
		db: db,
	}
}

func (z *zoneRepo) Create(warehouseID int, name string) (int, error) {
	row := z.db.QueryRow("INSERT INTO zones (warehouse_id, name) VALUES ($1, $2) RETURNING id", warehouseID, name)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (z *zoneRepo) GetByWarehouse(warehouseID int) ([]Zone, error) {
	rows, err := z.db.Query("SELECT id, warehouse_id, name FROM zones WHERE warehouse_id = $1", warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Zone
	for rows.Next() {
		var item Zone
		err := rows.Scan(&item.ID, &item.WarehouseID, &item.Name)
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
