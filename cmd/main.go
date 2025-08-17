package main

import (
	"context"

	"github.com/morozoffnor/home-storage/internal/config"
	"github.com/morozoffnor/home-storage/internal/database"
	"github.com/morozoffnor/home-storage/internal/handler"
	"github.com/morozoffnor/home-storage/internal/server"
)

func main() {
	cfg := config.New()
	db, err := database.New(cfg, context.Background())
	if err != nil {
		panic(err)
	}
	h := handler.New(cfg, db)
	s := server.New(cfg, h)

	s.ListenAndServe()
}
