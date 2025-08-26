package services

import (
	"context"
	"database/sql"
	"sqlc-testing/db"
)

// Struct representing Product services
type ProductService struct {
	queries *db.Queries
}

// Constructor for Product
func NewProductService(conn *sql.DB) *ProductService {
	return &ProductService{
		queries: db.New(conn),
	}
}

func (p *ProductService) CreateProduct(ctx context.Context, params db.CreateProductParams) (db.Product, error) {
	return p.queries.CreateProduct(ctx, params)
}

func (p *ProductService) GetProducts(ctx context.Context) ([]db.Product, error) {
	return p.queries.ListProducts(ctx)
}

func (p *ProductService) GetProductByID(ctx context.Context, id int32) (db.Product, error) {
	return p.queries.GetProductByID(ctx, id)
}

func (p *ProductService) GetProductBySKU(ctx context.Context, sku string) (db.Product, error) {
	return p.queries.GetProductBySKU(ctx, sku)
}

func (p *ProductService) DeleteProduct(ctx context.Context, id int32) error {
	return p.queries.DeleteProduct(ctx, id)
}
