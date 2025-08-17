package main

import (
	"github.com/morozoffnor/home-storage/internal/config"
	"github.com/morozoffnor/home-storage/internal/server"
)

func main() {
	cfg := config.New(":3000")
	s := server.New(cfg)

	s.ListenAndServe()
}
