package types

import "time"

type User struct {
	ID        int
	Username  string
	Email     string
	PassHash  string `json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
