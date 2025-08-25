package api

import (
	"net/http"
	"sqlc-testing/db"
	"sqlc-testing/utils"

	"github.com/gin-gonic/gin"
)

type createStoreRequest struct {
	Name	string	`json:"name" binding:"required"`
	Address	string	`json:"address" binding:"required"`
	Phone	string	`json:"phone"`
}

func (server *Server) createStore(ctx *gin.Context) {
	var req createStoreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.CreateStoreParams{
		Name: req.Name,
		Address: req.Address,
		Phone: utils.NullString(req.Phone),
	}

	store, err := server.services.Store.CreateStore(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, store)
}

func (server *Server) getAllStores(ctx *gin.Context) {
	stores, err := server.services.Store.GetStores(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "can't fetch stores"})
        return
	}

	ctx.JSON(200, stores)
}
