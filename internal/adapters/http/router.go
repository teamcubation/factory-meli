package http

import (
	"tweet-tq-rafael/core/ports"

	"github.com/gin-gonic/gin"
)

func NewRouter(services ports.Services) *gin.Engine {
	userHandler := &UserHandler{Service: services.UserService}
	tweetHandler := &TweetHandler{Service: services.TweetService}

	router := gin.Default()

	//Profile Routes
	router.POST("/user", userHandler.Create)
	router.GET("/timeline/:id", userHandler.Timeline)
	router.POST("/follow", userHandler.Follow)
	router.POST("/unfollow", userHandler.Unfollow)

	//Tweet Routes
	router.POST("/tweet", tweetHandler.Create)

	return router
}
