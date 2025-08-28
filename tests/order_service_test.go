package tests

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"sqlc-testing/db"
	"sqlc-testing/services"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// setupTestDB connects to the test database
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	connStr := os.Getenv("DB_DSN")
	dbConn, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	require.NoError(t, dbConn.Ping())
	log.Println("Connected to test database")
	return dbConn
}

// CreateRandomOrder seeds user, store, product, inventory, and returns order params
func CreateRandomOrder(t *testing.T, q *db.Queries) (db.CreateOrderParams, []db.CreateOrderItemParams, db.Product) {
	t.Helper()
	ctx := context.Background()

	// Seed user
	user, err := q.CreateUser(ctx, db.CreateUserParams{
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
	})
	require.NoError(t, err)

	// Seed store
	store, err := q.CreateStore(ctx, db.CreateStoreParams{
		Name:    gofakeit.Company(),
		Address: gofakeit.Address().Address,
	})
	require.NoError(t, err)

	// Seed product
	product, err := q.CreateProduct(ctx, db.CreateProductParams{
		Name:  gofakeit.ProductName(),
		Sku:   gofakeit.UUID(),
		Price: gofakeit.Price(10, 100),
	})
	require.NoError(t, err)

	// Seed inventory
	initialQty := int32(50)
	_, err = q.CreateInventoryItem(ctx, db.CreateInventoryItemParams{
		StoreID:   store.StoreID,
		ProductID: product.ProductID,
		Quantity:  initialQty,
	})
	require.NoError(t, err)

	// Create order params
	orderParams := db.CreateOrderParams{
		UserID:          user.ID,
		StoreID:         store.StoreID,
		DeliveryAddress: gofakeit.Address().Address,
	}

	// Create one order item
	items := []db.CreateOrderItemParams{
		{
			ProductID: product.ProductID,
			Quantity:  5,
			UnitPrice: product.Price,
		},
	}

	return orderParams, items, product
}

func TestOrderTransaction(t *testing.T) {
	dbConn := setupTestDB(t)
	defer dbConn.Close()

	ctx := context.Background()

	// Wrap the test in a transaction and rollback at the end
	tx, err := dbConn.Begin()
	require.NoError(t, err)
	defer tx.Rollback() // ensures no permanent changes

	service := services.NewOrderService(dbConn) // pass DB; execTx will create its own transaction

	// Create random order
	orderParams, items, product := CreateRandomOrder(t, service.Queries)

	// Run transaction
	result, err := service.OrderTransaction(ctx, orderParams, items)
	require.NoError(t, err)
	require.NotZero(t, result.OrderID)
	require.Equal(t, float64(items[0].Quantity)*product.Price, result.TotalAmount)
	require.Equal(t, "confirmed", result.Status)

	// Check inventory was decremented
	inv, err := service.Queries.GetInventoryByStoreAndProduct(ctx, db.GetInventoryByStoreAndProductParams{
		StoreID:   orderParams.StoreID,
		ProductID: product.ProductID,
	})
	require.NoError(t, err)
	expectedQty := int32(50) - items[0].Quantity
	require.Equal(t, expectedQty, inv.Quantity)
}
