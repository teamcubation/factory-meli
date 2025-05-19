package global

import (
	"Yuri/twitter/core/ports/services"
	"Yuri/twitter/internal/tweet"
	"Yuri/twitter/internal/user"

	"github.com/gin-gonic/gin"
)

type GlobalHandler struct {
	UserHandler  *user.Handler
	TweetHandler *tweet.Handler
}

func NewRouters(userService services.IUserService, tweetService services.ITweetService) *gin.Engine {
	userHandler := &user.Handler{Service: userService}
	tweetHandler := &tweet.Handler{Service: tweetService}

	router := gin.Default()

	//USER ROUTES
	router.POST("/user/create", userHandler.CreateUser)

	//TWEET ROUTES
	router.POST("/tweet/create", tweetHandler.CreateTweet)

	return router
}
