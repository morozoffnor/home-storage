package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/morozoffnor/home-storage/internal/config"
	"github.com/morozoffnor/home-storage/internal/handler"
)

func rootRouter(h *handler.APIHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Route("/home", func(r chi.Router) {
			r.Get("/", h.GetAllHomes)
			r.Post("/", h.CreateHome)

			r.Route("/{home_id}", func(r chi.Router) {
				r.Use(homeCtx)
				r.Get("/", h.GetHome)
				r.Put("/", h.UpdateHome)
				r.Delete("/", h.DeleteHome)
			})
		})
	})

	return r
}

func New(cfg *config.Config, h *handler.APIHandler) *http.Server {
	router := rootRouter(h)
	return &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: router,
	}
}
