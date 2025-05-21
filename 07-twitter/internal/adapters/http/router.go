package http

import (
	handlers "07-twitter/internal/adapters/http/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/users", handler.CreateUserHandler)
	router.GET("/users/:id", handler.GetUserById)
	router.POST("/users/follow", handler.Follow)
	router.DELETE("/users/unfollow", handler.Unfollow)

	return router
}
