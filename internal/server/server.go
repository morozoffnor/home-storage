package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/morozoffnor/home-storage/internal/config"
)

func newRouter() *chi.Mux {
	r := chi.NewRouter()
	return r
}

func New(cfg *config.Config) *http.Server {
	router := newRouter()
	return &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: router,
	}
}
