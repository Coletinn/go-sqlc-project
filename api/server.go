package api

import (
	"github.com/gin-gonic/gin"
	"sqlc-testing/services"
)

type Server struct {
	router *gin.Engine
	userService *services.UserService
}

func NewServer(userService *services.UserService) *Server {
	server := &Server{
		router: gin.Default(),
		userService: userService,
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("/", s.getAllUsers)
			users.GET("/:id", s.getUserByID)
			users.POST("/", s.createUser)
			users.PATCH("/:id", s.updateUser)
			users.DELETE("/:id", s.deleteUser)
		}
	}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
