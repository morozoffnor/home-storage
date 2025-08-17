package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/morozoffnor/home-storage/internal/config"
	"github.com/morozoffnor/home-storage/internal/handler"
)

func newRouter(h *handler.APIHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/api/note", h.CreateItem)
	return r
}

func New(cfg *config.Config, h *handler.APIHandler) *http.Server {
	router := newRouter(h)
	return &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: router,
	}
}
