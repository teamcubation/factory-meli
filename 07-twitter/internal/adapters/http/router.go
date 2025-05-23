package http

import (
	handlers "07-twitter/internal/adapters/http/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter(userHandler *handlers.UserHandler, tweetHandler *handlers.TweetHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/users", userHandler.CreateUserHandler)
	router.GET("/users/:id", userHandler.GetUserById)
	router.POST("/users/follow", userHandler.Follow)
	router.DELETE("/users/unfollow", userHandler.Unfollow)

	router.POST("/tweet", tweetHandler.CreateTweet)
	router.GET("/tweet/:id", tweetHandler.GetTweetsByUser)
	router.GET("/tweet/timeline/:id", tweetHandler.GetTweetTimelineByUser)

	return router
}
