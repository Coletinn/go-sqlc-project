package api

import (
	"net/http"
	"sqlc-testing/db"
	"sqlc-testing/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Name	string	`json:"name" binding:"required"`
	Email	string	`json:"email" binding:"required"`
	Phone	string	`json:"phone"`
}

type updateUserRequest struct {
    Name  *string `json:"name"`
    Email *string `json:"email"`
    Phone *string `json:"phone"`
}

func (server *Server) getUserByID(ctx *gin.Context) {
	// Get "id" parameter from the URL
    idParam := ctx.Param("id")
    requestId, err := strconv.Atoi(idParam)
    if err != nil {
        ctx.JSON(400, gin.H{"error": "invalid user ID"})
        return
    }

    user, err := server.services.User.GetUserByID(ctx, int32(requestId))
    if err != nil {
        ctx.JSON(404, gin.H{"error": "user not found"})
        return
    }

    ctx.JSON(200, user)
}

func (server *Server) getAllUsers(ctx *gin.Context) {
	users, err := server.services.User.ListUsers(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "can't fetch users"})
        return
	}

	ctx.JSON(200, users)
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := db.CreateUserParams{
		Name: req.Name,
		Email: req.Email,
		Phone: utils.NullString(req.Phone),
	}

	user, err := server.services.User.CreateUser(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (server *Server) updateUser(ctx *gin.Context) {
    idParam := ctx.Param("id")
    requestId, err := strconv.Atoi(idParam)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    var req updateUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    user, err := server.services.User.GetUserByID(ctx, int32(requestId))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    params := db.UpdateUserParams{
        ID:    int32(requestId),
        Name:  user.Name,
        Email: user.Email,
        Phone: user.Phone,
    }

    if req.Name != nil {
        params.Name = *req.Name
    }
    if req.Email != nil {
        params.Email = *req.Email
    }
    if req.Phone != nil {
        params.Phone = utils.NullString(*req.Phone)
    }

    updatedUser, err := server.services.User.UpdateUser(ctx, params)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, updatedUser)
}

func (server *Server) deleteUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
    requestId, err := strconv.Atoi(idParam)
    if err != nil {
        ctx.JSON(400, gin.H{"error": "invalid user ID"})
        return
    }

	err = server.services.User.DeleteUser(ctx, int32(requestId))
    if err != nil {
        ctx.JSON(404, gin.H{"error": "user not found"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
