package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/morozoffnor/home-storage/internal/config"
	api "github.com/morozoffnor/home-storage/internal/handler/api"
	frontend "github.com/morozoffnor/home-storage/internal/handler/frontend"
)

func rootRouter(api *api.APIHandler, frontend *frontend.FrontendHandler, m *Middleware) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	// auth endpoints do not need auth middleware
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", api.User.Register)
		r.Post("/login", api.User.Login)
	})

	r.Get("/login", frontend.LoginPage)

	r.Group(func(r chi.Router) {
		r.Use(m.FrontendAuth)
		r.Get("/", frontend.HomePage)
		r.Get("/homes", frontend.Home.GetAll)

		r.Route("/home", func(r chi.Router) {

			r.Route("/{home_id}", func(r chi.Router) {
				r.Use(m.homeCtx)
				r.Get("/", frontend.Container.GetAllInHome)
			})
		})
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(m.Auth)
		r.Route("/user", func(r chi.Router) {
			r.Get("/", api.User.GetAll)

			// user
			r.Route("/{user_id}", func(r chi.Router) {
				r.Use(m.userCtx)
				r.Get("/", api.User.Get)

				// user's homes
				r.Route("/home", func(r chi.Router) {
					r.Get("/", api.User.GetHomes)
					r.Post("/", api.User.AddHome)
				})
			})
		})
		// home management
		r.Route("/home", func(r chi.Router) {
			r.Get("/", api.Home.GetAll)
			r.Post("/", api.Home.Create)

			// distinct home
			r.Route("/{home_id}", func(r chi.Router) {
				r.Use(m.homeCtx)
				r.Get("/", api.Home.Get)
				r.Put("/", api.Home.Update)
				r.Delete("/", api.Home.Delete)

				r.Route("/container", func(r chi.Router) {
					r.Get("/", api.Container.GetAllInHome)
					r.Post("/", api.Container.Create)

					r.Route("/{container_id}", func(r chi.Router) {
						r.Use(m.containerCtx)
						r.Get("/", api.Container.Get)

						r.Route("/item", func(r chi.Router) {
							r.Get("/", api.Item.GetAllInContainer)
							r.Post("/", api.Item.Create)

							r.Route("/{item_id}", func(r chi.Router) {
								r.Use(m.itemCtx)
								r.Get("/", api.Item.Get)
								r.Delete("/", api.Item.Delete)
							})
						})

					})

				})
			})
		})
	})

	return r
}

func New(cfg *config.Config, h *api.APIHandler, f *frontend.FrontendHandler, m *Middleware) *http.Server {
	router := rootRouter(h, f, m)
	return &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: router,
	}
}
