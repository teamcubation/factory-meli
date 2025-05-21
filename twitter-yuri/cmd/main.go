package main

import (
	timelineService "Yuri/twitter/application/timeline/service"
	tweetService "Yuri/twitter/application/tweet/service"
	userService "Yuri/twitter/application/user/service"
	"Yuri/twitter/infrastructure/api/global"
	repos "Yuri/twitter/infrastructure/db/repositories"
)

func main() {

	//REPOSITORIES
	userRepo := repos.NewUserRepo()
	tweetRepo := repos.NewTweetRepo()

	//SERVICES
	userSvc := userService.NewUserService(userRepo)
	tweetSvc := tweetService.NewTweetService(tweetRepo)
	timelineSvc := timelineService.NewTimelineService(tweetRepo, userRepo)

	//ROUTERS
	router := global.NewRouters(userSvc, tweetSvc, timelineSvc)

	router.Run(":8080")

}
