package services

import (
	"context"
	"database/sql"
	"fmt"
	"sqlc-testing/db"
)

// Struct representing Order services
type OrderService struct {
	*db.Queries
	db *sql.DB
}

// Constructor for Order
func NewOrderService(conn *sql.DB) *OrderService {
	return &OrderService{
		Queries: db.New(conn),
		db: conn,
	}
}

func (order *OrderService) execTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := order.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := db.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type OrderItemResult struct {
	OrderItemID int32   `json:"order_item_id"`
	ProductID   int32   `json:"product_id"`
	ProductName string  `json:"product_name,omitempty"`
	SKU         string  `json:"sku,omitempty"`
	Quantity    int32   `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

type OrderTxResult struct {
	OrderID         int32             `json:"order_id"`
	UserID          int32             `json:"user_id"`
	StoreID         int32             `json:"store_id"`
	TotalAmount     float64           `json:"total_amount"`
	Status          string            `json:"status"`
	DeliveryAddress string            `json:"delivery_address"`
	OrderDate       string            `json:"order_date"`
	Items           []OrderItemResult `json:"items"`
}

func (o *OrderService) OrderTransaction(ctx context.Context, orderParams db.CreateOrderParams, items []db.CreateOrderItemParams) (OrderTxResult, error) {
	var result OrderTxResult

	err := o.execTx(ctx, func(q *db.Queries) error {
		// Create order
		order, err := q.CreateOrder(ctx, orderParams)
		if err != nil {
			return err
		}

		totalAmount := 0.0
		orderItems := make([]OrderItemResult, 0, len(items))

		for i := range items {
			item := &items[i]
			item.OrderID = order.OrderID

			inv, err := q.GetInventoryByStoreAndProduct(ctx, db.GetInventoryByStoreAndProductParams{
				StoreID:   order.StoreID,
				ProductID: item.ProductID,
			})
			if err != nil {
				return fmt.Errorf("error fetching inventory for product %d: %w", item.ProductID, err)
			}
			if inv.Quantity < item.Quantity {
				return fmt.Errorf("not enough inventory for product %d", item.ProductID)
			}

			item.TotalPrice = float64(item.Quantity) * item.UnitPrice

			createdItem, err := q.CreateOrderItem(ctx, *item)
			if err != nil {
				return fmt.Errorf("error creating order item: %w", err)
			}

			_, err = q.UpdateInventoryQuantity(ctx, db.UpdateInventoryQuantityParams{
				StoreID:   order.StoreID,
				ProductID: item.ProductID,
				Quantity:  inv.Quantity - item.Quantity,
			})
			if err != nil {
				return fmt.Errorf("error updating inventory: %w", err)
			}

			totalAmount += item.TotalPrice

			// Append item to result
			orderItems = append(orderItems, OrderItemResult{
				OrderItemID: createdItem.OrderItemID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
				TotalPrice:  item.TotalPrice,
			})
		}

		// Update order total amount and status
		_, err = q.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
			OrderID: order.OrderID,
			Status:  "confirmed",
		})
		if err != nil {
			return fmt.Errorf("error updating order status: %w", err)
		}

		order.TotalAmount = totalAmount

		// Populate result
		result = OrderTxResult{
			OrderID:         order.OrderID,
			UserID:          order.UserID,
			StoreID:         order.StoreID,
			TotalAmount:     totalAmount,
			Status:          "confirmed",
			DeliveryAddress: order.DeliveryAddress,
			OrderDate:       order.OrderDate.Time.Format("2006-01-02 15:04:05"),
			Items:           orderItems,
		}

		return nil
	})

	return result, err
} 

func (o *OrderService) GetOrderByID(ctx context.Context, orderID int32) (db.Order, error) {
	return o.Queries.GetOrderByID(ctx, orderID)
}

func (o *OrderService) GetOrdersByUser(ctx context.Context, userID int32) ([]db.ListOrdersByUserRow, error) {
	return o.Queries.ListOrdersByUser(ctx, userID)
}

func (o *OrderService) GetOrdersByStore(ctx context.Context, storeID int32) ([]db.ListOrdersByStoreRow, error) {
	return o.Queries.ListOrdersByStore(ctx, storeID)
}

func (o *OrderService) DeleteOrder(ctx context.Context, storeID int32) error {
	return o.Queries.DeleteOrder(ctx, storeID)
}
