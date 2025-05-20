package http

import "github.com/gin-gonic/gin"

func NewRouter(handler *Handler) *gin.Engine {
	router := gin.Default()

	router.POST("/users", handler.CreateUserHandler)
	router.GET("/users/:id", handler.GetUserById)
	router.POST("/users/follow", handler.Following)

	return router
}
