package services

import (
	"context"
	"database/sql"
	"sqlc-testing/db"
)

// Struct representing Store Inventory services
type StoreInventoryService struct {
	queries *db.Queries
}

// Constructor for Store Inventory
func NewStoreInventoryService(conn *sql.DB) *StoreInventoryService {
	return &StoreInventoryService{
		queries: db.New(conn),
	}
}

func (si *StoreInventoryService) CreateStoreInventoryItem(ctx context.Context, params db.CreateInventoryItemParams) (db.StoreInventory, error) {
	return si.queries.CreateInventoryItem(ctx, params)
}

func (si *StoreInventoryService) GetStoreInventoryByStore(ctx context.Context, storeID int32) ([]db.GetInventoryByStoreRow, error) {
	return si.queries.GetInventoryByStore(ctx, storeID)
}
