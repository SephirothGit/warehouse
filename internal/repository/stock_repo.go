package repository

import (
	"database/sql"
)

type StockRepo interface {
	GetByShelf(shelfID int) ([]StockItem, error)
	MoveStock(fromShelfID int, toSHelfID int, productID int, quantity int, userID int) error
	AddStock(shelfID int, productID int, quantity int, userID int) error
}

type StockItem struct {
	ID        int
	ShelfID   int
	ProductID int
	Quantity  int
}

type stockRepo struct {
	db *sql.DB
}

func NewStockRepo(db *sql.DB) StockRepo {
	return &stockRepo{
		db: db,
	}
}

func (s *stockRepo) GetByShelf(shelfID int) ([]StockItem, error) {
	rows, err := s.db.Query("SELECT id, shelf_id, product_id, quantity FROM stock_items WHERE shelf_id = $1", shelfID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []StockItem
	for rows.Next() {
		var item StockItem
		err := rows.Scan(&item.ID, &item.ShelfID, &item.ProductID, &item.Quantity)
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

func (s *stockRepo) MoveStock(fromShelfID int, toSHelfID int, productID int, quantity int, userID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE stock_items SET quantity = quantity - $1 WHERE shelf_id = $2 AND product_id = $3",
		quantity, fromShelfID, productID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO shelf_items (shelf_id, product_id, quantity)
VALUES ($!, $2, $3)
ON CONFLICT (shelf_id, product_id) 
DO UPDATE SET quantity = stock_items.quantity + EXCLUDED.quantity`, toSHelfID, productID, quantity)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO stock_movements (product_id, from_shelf_id, to_shelf_id, quantity, moved_by)
VALUES ($1, $2, $3, $4, $5)`, productID, fromShelfID, toSHelfID, quantity, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *stockRepo) AddStock(shelfID int, productID int, quantity int, userID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO stock_items (shelf_id, product_id, quantity)
	VALUES ($1, $2, $3)
	ON CONFLICT (shelf_id, product_id)
	DO UPDATE SET quantity = stock_items.quantity + EXCLUDED.quantity`, shelfID, productID, quantity)

	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO stock_movements (product_id, from_shelf_id, to_shelf_id, quantity, moved_by)
	VALUES ($1, NULL, $2, $3, $4)`, productID, shelfID, quantity, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
