package main

import (
	"database/sql"
	"log"
	"net"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"grpc/internal/config"
	productv1 "grpc/internal/grpc/gen/productv1"
	grpcserver "grpc/internal/grpc/server"
	postgresrepo "grpc/internal/storage/postgres"
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

	repository := postgresrepo.NewProductRepository(db)

	listener, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	server := grpc.NewServer()
	productv1.RegisterProductServiceServer(server, grpcserver.NewProductServer(repository))
	reflection.Register(server)

	log.Printf("gRPC API listening on %s", cfg.GRPCAddress)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
