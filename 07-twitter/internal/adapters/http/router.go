package http

import "github.com/gin-gonic/gin"

func NewRouter(handler *Handler) *gin.Engine {
	router := gin.Default()

	router.POST("/users", handler.CreateUserHandler)

	return router
}
