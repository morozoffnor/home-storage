package database

import (
	"fmt"

	"github.com/morozoffnor/home-storage/internal/types"
)

func (h *Home) Create(homeName string, description string) (int, error) {
	tx, err := h.conn.Begin(h.ctx)
	if err != nil {
		return 0, err
	}

	var homeID int
	query := `INSERT INTO homes (name, description) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRow(h.ctx, query, homeName, description).Scan(&homeID)
	if err != nil {
		tx.Rollback(h.ctx)
		fmt.Println("error while executing home transaction")
		return 0, err
	}

	err = tx.Commit(h.ctx)
	if err != nil {
		fmt.Println("error while committing home transaction")
		return 0, err
	}
	return homeID, nil
}

func (h *Home) Get(id int) (*types.Home, error) {
	var home types.Home

	query := `SELECT * FROM homes WHERE id=$1`
	err := h.conn.QueryRow(h.ctx, query, id).
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

func (h *Home) GetAll() ([]*types.Home, error) {
	var homes []*types.Home

	query := `SELECT * FROM homes`
	rows, err := h.conn.Query(h.ctx, query)
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

func (h *Home) Update(home *types.Home) error {
	tx, err := h.conn.Begin(h.ctx)
	if err != nil {
		return err
	}

	query := `UPDATE homes SET name = $1, description = $2, updated_at = NOW() WHERE id = $3`
	_, err = tx.Exec(h.ctx, query, home.Name, home.Description, home.ID)
	if err != nil {
		fmt.Println("error while executing home transaction")
		return err
	}
	err = tx.Commit(h.ctx)
	if err != nil {
		fmt.Println("error while committing home transaction")
		return err
	}
	return nil
}

func (h *Home) Delete(id int) error {
	tx, err := h.conn.Begin(h.ctx)
	if err != nil {
		return err
	}

	query := `DELETE FROM homes WHERE id = $1`
	_, err = tx.Exec(h.ctx, query, id)
	if err != nil {
		fmt.Println("error while executing delete home transaction")
		return err
	}

	err = tx.Commit(h.ctx)
	if err != nil {
		fmt.Println("error while committing delete home transaction")
		return err
	}
	return nil
}

func (h *Home) ContainersCount(id int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM containers WHERE home_id = $1`
	err := h.conn.QueryRow(h.ctx, query, id).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (h *Home) ItemsCount(id int) (int, error) {
	// count items in all containers in home
	var count int
	query := `SELECT COUNT(*) FROM items WHERE container_id IN (SELECT id FROM containers WHERE home_id = $1)`
	err := h.conn.QueryRow(h.ctx, query, id).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
