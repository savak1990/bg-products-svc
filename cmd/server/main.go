package main

import (
	"log"
	"net/http"
	"time"

	"github.com/savak1990/bg-products-svc/internal/config"
	"github.com/savak1990/bg-products-svc/internal/httpserver"
	"github.com/savak1990/bg-products-svc/internal/products"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	repo := products.NewInMemoryStore()
	s := httpserver.NewServer(repo)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      s.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("starting %s server on %s\n", cfg.Env, cfg.Addr)
	log.Fatal(srv.ListenAndServe())
}
