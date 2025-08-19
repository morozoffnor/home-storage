package handler

import (
	"github.com/morozoffnor/home-storage/internal/auth"
	"github.com/morozoffnor/home-storage/internal/config"
	"github.com/morozoffnor/home-storage/internal/database"
)

type APIHandler struct {
	cfg  *config.Config
	db   *database.Database
	auth *auth.Auth
	User User
	Home Home
}

func New(cfg *config.Config, db *database.Database, a *auth.Auth) *APIHandler {
	return &APIHandler{
		cfg:  cfg,
		db:   db,
		auth: a,
		User: User{
			auth: a,
			db:   db,
		},
		Home: Home{
			db: db,
		},
	}
}
