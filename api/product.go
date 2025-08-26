package api

import (
	"net/http"
	"sqlc-testing/db"
	"sqlc-testing/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createProductRequest struct {
	Name		string			`json:"name" binding:"required"`
	Description	string			`json:"description"`
	Price		float64			`json:"price" binding:"required"`
	Sku			string			`json:"sku" binding:"required"`
	Category	string			`json:"category"`
}

func (server *Server) createProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.CreateProductParams{
		Name: req.Name,
		Description: utils.NullString(req.Description),
		Price: req.Price,
		Sku: req.Sku,
		Category: utils.NullString(req.Category),
	}

	product, err := server.services.Product.CreateProduct(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (server *Server) getAllProducts(ctx *gin.Context) {
	products, err := server.services.Product.GetProducts(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "can't fetch products"})
        return
	}

	ctx.JSON(200, products)
}

func (server *Server) getProductByID(ctx *gin.Context) {
	idParam := ctx.Param("id")

	requestId, err := strconv.Atoi(idParam)
    if err != nil {
        ctx.JSON(400, gin.H{"error": "invalid product ID"})
        return
    }

	products, err := server.services.Product.GetProductByID(ctx, int32(requestId))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "can't fetch products"})
        return
	}

	ctx.JSON(200, products)
}

func (server *Server) getProductBySKU(ctx *gin.Context) {
	idParam := ctx.Param("sku")
	
	product, err := server.services.Product.GetProductBySKU(ctx, idParam)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "can't fetch product SKU"})
        return
	}

	ctx.JSON(200, product)
}

func (server *Server) deleteProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")

	requestId, err := strconv.Atoi(idParam)
    if err != nil {
        ctx.JSON(400, gin.H{"error": "invalid product ID"})
        return
    }

	err = server.services.Product.DeleteProduct(ctx, int32(requestId))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "could not delete product"})
        return
	}

	ctx.JSON(200, gin.H{"message": "Product succesfully deleted"})
}
