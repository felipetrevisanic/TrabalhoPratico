package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	appservice "graphql/src/application/service"
	"graphql/src/config"
	repositories "graphql/src/infrastructure/repositories"
	graphapi "graphql/src/interfaces/graphql"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	repository := repositories.NewProductRepository(db)
	service := appservice.NewProductService(repository)
	server, err := graphapi.NewServer(service)
	if err != nil {
		log.Fatalf("create graphql server: %v", err)
	}

	httpServer := &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           server,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("GraphQL API listening on %s", cfg.HTTPAddress)

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen and serve: %v", err)
	}
}
