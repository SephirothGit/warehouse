package repository

import "database/sql"

type RackRepo interface {
	Create(zoneID int, code string) (int, error)
	GetByZone(zoneID int) ([]Rack, error)
}

type Rack struct {
	ID     int
	ZoneID int
	Code   string
}

type rackRepo struct {
	db *sql.DB
}

func NewRackRepo(db *sql.DB) RackRepo {
	return &rackRepo{
		db: db,
	}
}

func (r *rackRepo) Create(zoneID int, code string) (int, error) {
	row := r.db.QueryRow("INSERT INTO racks (zone_id, code) VALUES ($1, $2) RETURNING id", zoneID, code)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *rackRepo) GetByZone(zoneID int) ([]Rack, error) {
	rows, err := r.db.Query("SELECT id, zone_id, code FROM racks WHERE zone_id = $1", zoneID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Rack
	for rows.Next() {
		var item Rack
		err := rows.Scan(&item.ID, &item.ZoneID, &item.Code)
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
