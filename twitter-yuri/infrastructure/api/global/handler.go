package global

import (
	"Yuri/twitter/core/ports/services"
	"Yuri/twitter/infrastructure/api/handlers"

	"github.com/gin-gonic/gin"
)

type GlobalHandler struct {
	UserHandler     *handlers.UserHandler
	TweetHandler    *handlers.TweetHandler
	TimelineHandler *handlers.TimelineHandler
}

func NewRouters(userService services.IUserService, tweetService services.ITweetService, timelineService services.ITimelineService) *gin.Engine {
	userHandler := &handlers.UserHandler{Service: userService}
	tweetHandler := &handlers.TweetHandler{Service: tweetService}
	timelineHandler := &handlers.TimelineHandler{Service: timelineService}

	router := gin.Default()

	//USER ROUTES
	router.POST("/user/create", userHandler.CreateUser)
	router.PATCH("/user/follow", userHandler.FollowUser)
	router.PATCH("/user/unfollow", userHandler.UnfollowUser)

	//TWEET ROUTES
	router.POST("/tweet/create", tweetHandler.CreateTweet)
	router.GET("/tweet", tweetHandler.FindTweetsByUserId)

	//TIMELINE ROUTES
	router.GET("/timeline", timelineHandler.GetTimeline)

	return router
}
