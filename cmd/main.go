package main

import (
	"context"

	"github.com/morozoffnor/home-storage/internal/auth"
	"github.com/morozoffnor/home-storage/internal/config"
	"github.com/morozoffnor/home-storage/internal/database"
	api "github.com/morozoffnor/home-storage/internal/handler/api"
	"github.com/morozoffnor/home-storage/internal/handler/frontend"
	"github.com/morozoffnor/home-storage/internal/server"
)

func main() {
	cfg := config.New()
	db, err := database.New(cfg, context.Background())
	if err != nil {
		panic(err)
	}
	a := auth.New(cfg)
	api := api.New(cfg, db, a)
	frontend := frontend.New(db)
	m := server.NewMiddleware(a, db)
	s := server.New(cfg, api, frontend, m)

	s.ListenAndServe()
}
