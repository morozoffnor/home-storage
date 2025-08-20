package types

import "time"

type Container struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"creaated_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	HomeID      int       `json:"home_id"`
}
