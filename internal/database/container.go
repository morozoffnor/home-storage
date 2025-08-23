package database

import (
	"fmt"

	"github.com/morozoffnor/home-storage/internal/types"
)

func (c *Container) Create(name string, description string, category string, location string, homeID int) (int, error) {
	tx, err := c.conn.Begin(c.ctx)
	if err != nil {
		return 0, err
	}
	var containerID int
	query := `INSERT INTO containers (name, description, category, location, home_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = tx.QueryRow(c.ctx, query, name, description, category, location, homeID).Scan(&containerID)
	if err != nil {
		tx.Rollback(c.ctx)
		fmt.Println("error while executing container transaction")
		return 0, err
	}
	err = tx.Commit(c.ctx)
	if err != nil {
		fmt.Println("error while committing container transaction")
		return 0, err
	}
	return containerID, nil
}

func (c *Container) Get(id int) (*types.Container, error) {
	var container types.Container

	query := `SELECT * FROM containers WHERE id=$1`
	err := c.conn.QueryRow(c.ctx, query, id).
		Scan(&container.ID,
			&container.Name,
			&container.Description,
			&container.Category,
			&container.Location,
			&container.CreatedAt,
			&container.UpdatedAt,
			&container.HomeID)
	if err != nil {
		return nil, err
	}
	return &container, nil
}

func (c *Container) GetAll() ([]*types.Container, error) {
	var containers []*types.Container

	query := `SELECT * FROM containers`
	rows, err := c.conn.Query(c.ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var container types.Container
		err = rows.
			Scan(&container.ID,
				&container.Name,
				&container.Description,
				&container.Category,
				&container.Location,
				&container.CreatedAt,
				&container.UpdatedAt,
				&container.HomeID)
		if err != nil {
			return nil, err
		}
		containers = append(containers, &container)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func (c *Container) GetAllInHome(homeID int) ([]*types.Container, error) {
	var containers []*types.Container

	query := `SELECT * FROM containers WHERE home_id=$1`
	rows, err := c.conn.Query(c.ctx, query, homeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var container types.Container
		err = rows.
			Scan(&container.ID,
				&container.Name,
				&container.Description,
				&container.Category,
				&container.Location,
				&container.CreatedAt,
				&container.UpdatedAt,
				&container.HomeID)
		if err != nil {
			return nil, err
		}
		containers = append(containers, &container)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func (c *Container) Update(container *types.Container) error {
	tx, err := c.conn.Begin(c.ctx)
	if err != nil {
		return err
	}

	query := `UPDATE containers SET name = $1, description = $2, category = $3, location = $4, updated_at = NOW() WHERE id = $5`
	_, err = tx.Exec(c.ctx, query, container.Name, container.Description, container.Category, container.Location, container.ID)
	if err != nil {
		fmt.Println("error while executing update container transaction")
		return err
	}
	err = tx.Commit(c.ctx)
	if err != nil {
		fmt.Println("error while committing update container transaction")
		return err
	}
	return nil
}

// TODO: delete all items as well
func (c *Container) Delete(id int) error {
	tx, err := c.conn.Begin(c.ctx)
	if err != nil {
		return err
	}

	query := `DELETE FROM containers WHERE id = $1`
	_, err = tx.Exec(c.ctx, query, id)
	if err != nil {
		fmt.Println("error while executing delete container transaction")
		return err
	}

	err = tx.Commit(c.ctx)
	if err != nil {
		fmt.Println("error while committing delete container transaction")
		return err
	}
	return nil
}

func (c *Container) ItemsCount(containerID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM items WHERE container_id = $1`
	err := c.conn.QueryRow(c.ctx, query, containerID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
