package api

import (
	"github.com/gin-gonic/gin"
	"sqlc-testing/services"
)

type Server struct {
	router *gin.Engine
	services *services.Services
}

func NewServer(svcs *services.Services) *Server {
	server := &Server{
		router: gin.Default(),
		services: svcs,
	}
	server.setupRoutes()
	return server
}

func (server *Server) setupRoutes() {
	api := server.router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("/", server.getAllUsers)
			users.GET("/:id", server.getUserByID)
			users.POST("/", server.createUser)
			users.PATCH("/:id", server.updateUser)
			users.DELETE("/:id", server.deleteUser)
		}
		stores := api.Group("/stores")
		{
			stores.GET("/", server.getAllStores)
			stores.POST("/", server.createStore)
		}
		products := api.Group("/products")
		{
			products.GET("/", server.getAllProducts)
			products.GET("/id/:id", server.getProductByID)
			products.GET("/sku/:sku", server.getProductBySKU)
			products.POST("/", server.createProduct)
			products.DELETE("/:id", server.deleteProduct)
		}
		orders := api.Group("/orders")
		{
			orders.GET("/user/:userId", server.getOrdersByUser)
			orders.GET("/id/:id", server.getOrdersByID)
			orders.POST("/", server.createOrder)
			// Enhanced create order with items and inventory handling
			orders.POST("/enhanced", server.createOrderEnhanced)
		}
		inventory := api.Group("/inventory")
		{
			inventory.GET("/:id", server.getInventoryByStore)
			inventory.POST("/", server.createInventoryItem)
		}
	}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
