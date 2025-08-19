package database

import (
	"fmt"

	"github.com/morozoffnor/home-storage/internal/types"
)

func (u *User) Create(username string, email string, passHash string) error {
	tx, err := u.conn.Begin(u.ctx)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`
	_, err = tx.Exec(u.ctx, query, username, email, passHash)
	if err != nil {
		fmt.Println("error while executing user transaction")
		return err
	}

	err = tx.Commit(u.ctx)
	if err != nil {
		fmt.Println("error while committing user transaction")
		return err
	}
	return nil
}

func (u *User) Exists(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`
	var exists bool

	err := u.conn.QueryRow(u.ctx, query, email).Scan(&exists)
	if err != nil {
		// returning true if err so we don't create a user by mistake
		return true, err
	}
	return exists, nil
}

func (u *User) Get(email string) (*types.User, error) {
	var user types.User

	query := `SELECT * FROM users WHERE email=$1`
	err := u.conn.QueryRow(u.ctx, query, email).
		Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PassHash,
			&user.CreatedAt,
			&user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) GetAll() ([]*types.User, error) {
	var users []*types.User

	query := `SELECT * FROM users`
	rows, err := u.conn.Query(u.ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user types.User
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PassHash,
			&user.CreatedAt,
			&user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}
