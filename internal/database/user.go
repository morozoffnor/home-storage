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
	var homeIDs []int

	query := `SELECT * FROM users WHERE email=$1`
	err := u.conn.QueryRow(u.ctx, query, email).
		Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PassHash,
			&user.CreatedAt,
			&user.UpdatedAt,
			&homeIDs,
		)

	if err != nil {
		fmt.Printf("user.get error: %v\n", err)
		return nil, err
	}
	homes, err := u.GetHomes(user.ID)
	if err != nil {
		fmt.Printf("user.get (homes) error: %v\n", err)
		return nil, err
	}
	user.Homes = homes
	return &user, nil
}

func (u *User) GetByID(id int) (*types.User, error) {
	var user types.User
	var homeIDs []int

	query := `SELECT * FROM users WHERE id=$1`
	err := u.conn.QueryRow(u.ctx, query, id).
		Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PassHash,
			&user.CreatedAt,
			&user.UpdatedAt,
			&homeIDs)
	if err != nil {
		return nil, err
	}
	homes, err := u.GetHomes(user.ID)
	if err != nil {
		fmt.Printf("user.getbyid (homes) error: %v\n", err)
		return nil, err
	}
	user.Homes = homes
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
		var homeIDs []int
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PassHash,
			&user.CreatedAt,
			&user.UpdatedAt,
			&homeIDs)
		if err != nil {
			return nil, err
		}
		homes, err := u.GetHomes(user.ID)
		if err != nil {
			fmt.Printf("user.getall (homes) error: %v\n", err)
			return nil, err
		}
		user.Homes = homes
		users = append(users, &user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *User) Update(id int, name string) error {
	tx, err := u.conn.Begin(u.ctx)
	if err != nil {
		return err
	}

	query := `UPDATE users SET username = $1, updated_at = NOW() WHERE id = $2`
	_, err = tx.Exec(u.ctx, query, name, id)
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

func (u *User) AddHome(userID int, homeID int) error {
	tx, err := u.conn.Begin(u.ctx)
	if err != nil {
		return err
	}
	query := `UPDATE users SET homes = homes || $1::INTEGER WHERE id = $2`
	_, err = tx.Exec(u.ctx, query, homeID, userID)
	if err != nil {
		fmt.Println("error while executing user home transaction")
		return err
	}

	err = tx.Commit(u.ctx)
	if err != nil {
		fmt.Println("error while commiting user home transaction")
		return err
	}
	return nil
}

func (u *User) GetHomes(userID int) ([]*types.Home, error) {
	var homes []*types.Home

	query := `SELECT h.* FROM homes h WHERE h.id = ANY (SELECT UNNEST(homes) FROM users WHERE id = $1);`
	rows, err := u.conn.Query(u.ctx, query, userID)
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
