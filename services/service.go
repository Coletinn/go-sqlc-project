package services

import (
	"database/sql"
)

type Services struct {
	User  *UserService
	Store *StoreService
}

func NewServices(conn *sql.DB) *Services {
	return &Services{
		User:  NewUserService(conn),
		Store: NewStoreService(conn),
	}
}
