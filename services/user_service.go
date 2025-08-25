package services

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

func (s *UserService) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
    if params.Phone.String != "" {
        params.Phone.Valid = true
    } else {
        params.Phone.Valid = false
    }

    return s.queries.CreateUser(ctx, params)
}

func (s *UserService) GetUserByID(ctx context.Context, id int32) (db.User, error) {
    return s.queries.GetUserByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]db.User, error) {
    return s.queries.ListUsers(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
    if params.Phone.String != "" {
        params.Phone.Valid = true
    } else {
        params.Phone.Valid = false
    }

    return s.queries.UpdateUser(ctx, params)
}


func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
    return s.queries.DeleteUser(ctx, id)
}