package services

import (
	"database/sql"
)

type Services struct {
	User  *UserService
	Store *StoreService
	Product *ProductService
	StoreInventory *StoreInventoryService
	Order *OrderService
}

func NewServices(conn *sql.DB) *Services {
	return &Services{
		User:  NewUserService(conn),
		Store: NewStoreService(conn),
		Product: NewProductService(conn),
		StoreInventory: NewStoreInventoryService(conn),
		Order: NewOrderService(conn),
	}
}
