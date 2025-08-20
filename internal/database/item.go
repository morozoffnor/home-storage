package database

import (
	"fmt"

	"github.com/morozoffnor/home-storage/internal/types"
)

func (i *Item) Create(name string, description string, category string, containerID int) (int, error) {
	tx, err := i.conn.Begin(i.ctx)
	if err != nil {
		return 0, err
	}
	var itemID int
	query := `INSERT INTO items (name, description, category, container_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err = tx.QueryRow(i.ctx, query, name, description, category, containerID).Scan(&itemID)
	if err != nil {
		tx.Rollback(i.ctx)
		fmt.Println("error while executing item transaction")
		return 0, err
	}
	err = tx.Commit(i.ctx)
	if err != nil {
		fmt.Println("error while committing item transaction")
		return 0, err
	}
	return containerID, nil
}

func (i *Item) Get(id int) (*types.Item, error) {
	var item types.Item

	query := `SELECT * FROM items WHERE id=$1`
	err := i.conn.QueryRow(i.ctx, query, id).
		Scan(&item.ID,
			&item.Name,
			&item.Description,
			&item.Category,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.ContainerID)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (i *Item) GetAll() ([]*types.Item, error) {
	var items []*types.Item

	query := `SELECT * FROM items`
	rows, err := i.conn.Query(i.ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item types.Item
		err = rows.
			Scan(&item.ID,
				&item.Name,
				&item.Description,
				&item.Category,
				&item.CreatedAt,
				&item.UpdatedAt,
				&item.ContainerID)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (i *Item) GetAllInContainer(homeID int) ([]*types.Item, error) {
	var items []*types.Item

	query := `SELECT * FROM items WHERE container_id=$1`
	rows, err := i.conn.Query(i.ctx, query, homeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item types.Item
		err = rows.
			Scan(&item.ID,
				&item.Name,
				&item.Description,
				&item.Category,
				&item.CreatedAt,
				&item.UpdatedAt,
				&item.ContainerID)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (i *Item) Update(item *types.Item) error {
	tx, err := i.conn.Begin(i.ctx)
	if err != nil {
		return err
	}

	query := `UPDATE items SET name = $1, description = $2, category = $3, updated_at = NOW() WHERE id = $4`
	_, err = tx.Exec(i.ctx, query, item.Name, item.Description, item.Category, item.ID)
	if err != nil {
		fmt.Println("error while executing update item transaction")
		return err
	}
	err = tx.Commit(i.ctx)
	if err != nil {
		fmt.Println("error while committing update item transaction")
		return err
	}
	return nil
}

func (i *Item) Delete(id int) error {
	tx, err := i.conn.Begin(i.ctx)
	if err != nil {
		return err
	}

	query := `DELETE FROM items WHERE id = $1`
	_, err = tx.Exec(i.ctx, query, id)
	if err != nil {
		fmt.Println("error while executing delete item transaction")
		return err
	}

	err = tx.Commit(i.ctx)
	if err != nil {
		fmt.Println("error while committing delete item transaction")
		return err
	}
	return nil
}
