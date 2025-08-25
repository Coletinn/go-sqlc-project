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
	}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
