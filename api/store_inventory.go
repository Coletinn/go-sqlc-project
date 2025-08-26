package api

import (
	"net/http"
	"sqlc-testing/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createInventoryRequest struct {
	StoreID   int32  `json:"store_id" binding:"required"`
	ProductID int32  `json:"product_id" binding:"required"`
	Quantity  int32  `json:"quantity" binding:"required"`
}

func (server *Server) createInventoryItem(ctx *gin.Context) {
	var req createInventoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.CreateInventoryItemParams{
		StoreID:   req.StoreID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	inv, err := server.services.StoreInventory.CreateStoreInventoryItem(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inv)
}

func (server *Server) getInventoryByStore(ctx *gin.Context) {
	idParam := ctx.Param("id")

	requestId, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	inv, err := server.services.StoreInventory.GetStoreInventoryByStore(ctx, int32(requestId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, inv)
}
