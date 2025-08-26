package services

import (
	"context"
	"database/sql"
	"sqlc-testing/db"
)

type StoreService struct {
    queries *db.Queries
}

func NewStoreService(conn *sql.DB) *StoreService {
    return &StoreService{
        queries: db.New(conn),
    }
}

func (s *StoreService) CreateStore(ctx context.Context, params db.CreateStoreParams) (db.Store, error) {
    if params.Phone.String != "" {
        params.Phone.Valid = true
    } else {
        params.Phone.Valid = false
    }

    return s.queries.CreateStore(ctx, params)
}

func (s *StoreService) GetStores(ctx context.Context) ([]db.Store, error) {
    return s.queries.ListStores(ctx)
}
