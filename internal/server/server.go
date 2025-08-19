package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/morozoffnor/home-storage/internal/config"
	"github.com/morozoffnor/home-storage/internal/handler"
)

func rootRouter(h *handler.APIHandler, m *Middleware) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	// auth endpoints do not need auth middleware
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.User.Register)
		r.Post("/login", h.User.Login)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(m.Auth)
		r.Route("/user", func(r chi.Router) {
			r.Get("/", h.User.GetAll)

			// user
			r.Route("/{user_id}", func(r chi.Router) {
				r.Use(m.userCtx)
				r.Get("/", h.User.Get)

				// user's homes
				r.Route("/home", func(r chi.Router) {
					r.Get("/", h.User.GetHomes)
					r.Post("/", h.User.AddHome)
				})
			})
		})
		// home management
		r.Route("/home", func(r chi.Router) {
			r.Get("/", h.Home.GetAll)
			r.Post("/", h.Home.Create)

			// distinct home
			r.Route("/{home_id}", func(r chi.Router) {
				r.Use(m.homeCtx)
				r.Get("/", h.Home.Get)
				r.Put("/", h.Home.Update)
				r.Delete("/", h.Home.Delete)
			})
		})
	})

	return r
}

func New(cfg *config.Config, h *handler.APIHandler, m *Middleware) *http.Server {
	router := rootRouter(h, m)
	return &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: router,
	}
}
