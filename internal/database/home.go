package database

import (
	"fmt"

	"github.com/morozoffnor/home-storage/internal/types"
)

func (db *Database) CreateHome(homeName string, description string) error {
	tx, err := db.conn.Begin(db.ctx)
	if err != nil {
		return err
	}

	query := `INSERT INTO homes (name, description) VALUES ($1, $2)`
	_, err = tx.Exec(db.ctx, query, homeName, description)
	if err != nil {
		fmt.Println("error while executing home transaction")
		return err
	}

	err = tx.Commit(db.ctx)
	if err != nil {
		fmt.Println("error while committing home transaction")
		return err
	}
	return nil
}

func (db *Database) GetHomeByID(id int) (*types.Home, error) {
	var home types.Home

	query := `SELECT * FROM homes WHERE id=$1`
	err := db.conn.QueryRow(db.ctx, query, id).
		Scan(
			&home.ID,
			&home.Name,
			&home.Description,
			&home.CreatedAt,
			&home.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &home, nil
}
