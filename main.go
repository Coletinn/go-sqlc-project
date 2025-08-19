package main

import (
    "context"
    "database/sql"
    "fmt"
    "math/rand"
    "os"
    "time"

    _ "github.com/lib/pq"
    "sqlc-testing/services"
)

func generateRandomUser() (string, string, string) {
    names := []string{"Alice", "Bob", "Carol", "David", "Eva"}
    emails := []string{"alice@test.com", "bob@test.com", "carol@test.com", "david@test.com", "eva@test.com"}
    phones := []string{"111-222-333", "222-333-444", "333-444-555", "444-555-666", "555-666-777"}

    return names[rand.Intn(len(names))],
        emails[rand.Intn(len(emails))],
        phones[rand.Intn(len(phones))]
}

func main() {
    rand.Seed(time.Now().UnixNano())

    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "postgres://gustavo:1910@localhost:5432/postgres?sslmode=disable"
    }

    conn, err := sql.Open("postgres", dsn)
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    if err := conn.Ping(); err != nil {
        panic(fmt.Sprintf("database not reachable: %v", err))
    }

    userService := service.NewUserService(conn)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    name, email, phone := generateRandomUser()

    user, err := userService.CreateUser(ctx, name, email, phone)
    if err != nil {
        panic(err)
    }

    fmt.Println("Created user:", user)

    // Buscar o mesmo usuário
    found, err := userService.GetUserByID(ctx, user.ID)
    if err != nil {
        panic(err)
    }

    fmt.Println("Fetched user:", found)

    // Listar todos usuários
    users, err := userService.ListUsers(ctx)
    if err != nil {
        panic(err)
    }

    fmt.Println("All users:", users)
}
