package handler

import (
	"net/http"

	"github.com/morozoffnor/home-storage/internal/config"
	"github.com/morozoffnor/home-storage/internal/database"
)

type APIHandler struct {
	cfg *config.Config
	db  *database.Database
}

func New(cfg *config.Config, db *database.Database) *APIHandler {
	return &APIHandler{
		cfg: cfg,
		db:  db,
	}
}

func (h *APIHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
