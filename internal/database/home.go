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

func (db *Database) GetHome(id int) (*types.Home, error) {
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

func (db *Database) GetAllHomes() ([]*types.Home, error) {
	var homes []*types.Home

	query := `SELECT * FROM homes`
	rows, err := db.conn.Query(db.ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var home types.Home
		err = rows.Scan(
			&home.ID,
			&home.Name,
			&home.Description,
			&home.CreatedAt,
			&home.UpdatedAt)
		if err != nil {
			return nil, err
		}
		homes = append(homes, &home)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return homes, nil
}

func (db *Database) UpdateHome(home *types.Home) error {
	tx, err := db.conn.Begin(db.ctx)
	if err != nil {
		return err
	}

	query := `UPDATE homes SET name = $1, description = $2, updated_at = NOW() WHERE id = $3`
	_, err = tx.Exec(db.ctx, query, home.Name, home.Description, home.ID)
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

func (db *Database) DeleteHome(id int) error {
	tx, err := db.conn.Begin(db.ctx)
	if err != nil {
		return err
	}

	query := `DELETE FROM homes WHERE id = $1`
	_, err = tx.Exec(db.ctx, query, id)
	if err != nil {
		fmt.Println("error while executing delete home transaction")
		return err
	}

	err = tx.Commit(db.ctx)
	if err != nil {
		fmt.Println("error while committing delete home transaction")
		return err
	}
	return nil
}
