package global

import (
	"Yuri/twitter/core/ports/services"
	"Yuri/twitter/infrastructure/api/handlers"

	"github.com/gin-gonic/gin"
)

type GlobalHandler struct {
	UserHandler  *handlers.UserHandler
	TweetHandler *handlers.TweetHandler
}

func NewRouters(userService services.IUserService, tweetService services.ITweetService) *gin.Engine {
	userHandler := &handlers.UserHandler{Service: userService}
	tweetHandler := &handlers.TweetHandler{Service: tweetService}

	router := gin.Default()

	//USER ROUTES
	router.POST("/user/create", userHandler.CreateUser)

	//TWEET ROUTES
	router.POST("/tweet/create", tweetHandler.CreateTweet)

	return router
}
