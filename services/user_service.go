package service

import (
    "context"
    "database/sql"
    "sqlc-testing/db"
)

type UserService struct {
    queries *db.Queries
}

func NewUserService(conn *sql.DB) *UserService {
    return &UserService{
        queries: db.New(conn),
    }
}

func (s *UserService) CreateUser(ctx context.Context, name, email, phone string) (db.User, error) {
    params := db.CreateUserParams{
        Name:  name,
        Email: email,
    }
    
    // Adiciona phone se fornecido
    if phone != "" {
        params.Phone = sql.NullString{String: phone, Valid: true}
    }
    
    return s.queries.CreateUser(ctx, params)
}

func (s *UserService) GetUserByID(ctx context.Context, id int32) (db.User, error) {
    return s.queries.GetUserByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]db.User, error) {
    return s.queries.ListUsers(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, id int32, name, email, phone string) (db.User, error) {
    params := db.UpdateUserParams{
        ID:    id,
        Name:  name,
        Email: email,
    }
    
    if phone != "" {
        params.Phone = sql.NullString{String: phone, Valid: true}
    }
    
    return s.queries.UpdateUser(ctx, params)
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
    return s.queries.DeleteUser(ctx, id)
}