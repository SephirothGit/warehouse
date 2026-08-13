package repository

import "database/sql"

type ShelfRepo interface {
	Create(rackID int, level int) (int, error)
	GetByRack(rackID int) ([]Shelf, error)
}

type Shelf struct {
	ID     int
	RackID int
	Level  int
}

type shelfRepo struct {
	db *sql.DB
}

func NewShelfRepo(db *sql.DB) ShelfRepo {
	return &shelfRepo{
		db: db,
	}
}

func (s *shelfRepo) Create(rackID int, level int) (int, error) {
	row := s.db.QueryRow("INSERT INTO shelves (rack_id, level) VALUES ($1, $2) RETURNING id", rackID, level)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *shelfRepo) GetByRack(rackID int) ([]Shelf, error) {
	rows, err := s.db.Query("SELECT id, rack_id, level FROM shelves WHERE rack_id = $1", rackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Shelf
	for rows.Next() {
		var item Shelf
		err := rows.Scan(&item.ID, &item.RackID, &item.Level)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}
