package main

import (
	"database/sql"
	"log"
	"os"
	"sqlc-testing/api"
	"sqlc-testing/services"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	dbDsn := os.Getenv("DB_DSN")
	if dbDsn == "" {
		log.Fatal("No DB_DSN env variable found")
	}
	// Connect to the database
	conn, err := sql.Open(dbDriver, dbDsn)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatal("Database unreachable:", err)
	}

	services := services.NewServices(conn)

	server := api.NewServer(services)

	log.Printf("Server running on %s\n", serverAddress)
	if err := server.Start(serverAddress); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
