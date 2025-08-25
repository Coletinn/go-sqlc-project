package main

import (
	"database/sql"
	"log"
	"sqlc-testing/api"
	"sqlc-testing/services"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://gustavo:1910@localhost:5432/postgres?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	// Connect to the database
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatal("Database unreachable:", err)
	}

	userService := services.NewUserService(conn)

	server := api.NewServer(userService)

	log.Printf("Server running on %s\n", serverAddress)
	if err := server.Start(serverAddress); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
