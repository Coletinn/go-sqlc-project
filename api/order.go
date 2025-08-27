package api

import (
	"net/http"
	"sqlc-testing/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createOrderRequest struct {
	UserID			int32			`json:"user_id" binding:"required"`
	StoreID			int32			`json:"store_id" binding:"required"`
	TotalAmount		float64			`json:"total_amount" binding:"required"`
	DeliveryAddress	string			`json:"delivery_address" binding:"required"`
}

type createOrderEnhancedRequest struct {
	UserID          int32                   `json:"user_id" binding:"required"`
	StoreID         int32                   `json:"store_id" binding:"required"`
	DeliveryAddress string                  `json:"delivery_address" binding:"required"`
	Items           []createOrderItemRequest `json:"items" binding:"required,dive"`
}

type createOrderItemRequest struct {
	ProductID int32   `json:"product_id" binding:"required"`
	Quantity  int32   `json:"quantity" binding:"required"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
}

func (server *Server) createOrder(ctx *gin.Context) {
	var req createOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.CreateOrderParams{
		UserID: req.UserID,
		StoreID: req.StoreID,
		TotalAmount: req.TotalAmount,
		DeliveryAddress: req.DeliveryAddress,
	}

	order, err := server.services.Order.CreateOrder(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (server *Server) createOrderEnhanced(ctx *gin.Context) {
	var req createOrderEnhancedRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Map request items to db.CreateOrderItemParams
	items := make([]db.CreateOrderItemParams, len(req.Items))
	for i, item := range req.Items {
		items[i] = db.CreateOrderItemParams{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			TotalPrice: float64(item.Quantity) * item.UnitPrice,
		}
	}

	// Call service
	orderResult, err := server.services.Order.OrderTransaction(ctx, db.CreateOrderParams{
		UserID:          req.UserID,
		StoreID:         req.StoreID,
		TotalAmount:     0,
		DeliveryAddress: req.DeliveryAddress,
	}, items)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, orderResult)
}

func (server *Server) getOrdersByUser(ctx *gin.Context) {
	userIdParam := ctx.Param("userId")

	requestId, err := strconv.Atoi(userIdParam)
    if err != nil {
        ctx.JSON(400, errorResponse(err))
        return
    }

	orders, err := server.services.Order.GetOrdersByUser(ctx, int32(requestId))
	if err != nil {
		ctx.JSON(400, errorResponse(err))
        return
	}

	ctx.JSON(200, orders)
}

func (server *Server) getOrdersByID(ctx *gin.Context) {
	IdParam := ctx.Param("id")

	requestId, err := strconv.Atoi(IdParam)
    if err != nil {
        ctx.JSON(400, errorResponse(err))
        return
    }

	orders, err := server.services.Order.GetOrderByID(ctx, int32(requestId))
	if err != nil {
		ctx.JSON(400, errorResponse(err))
        return
	}

	ctx.JSON(200, orders)
}
