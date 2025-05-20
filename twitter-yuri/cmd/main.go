package main

import (
	tweetService "Yuri/twitter/application/tweet/service"
	userService "Yuri/twitter/application/user/service"
	"Yuri/twitter/infrastructure/api/global"
	repos "Yuri/twitter/infrastructure/db/repositories"
)

func main() {
	// _, err := db.Connect()
	// if err != nil {
	// 	fmt.Println("Error connecting to MongoDB:", err)
	// 	return
	// }
	// fmt.Println("Connected to MongoDB successfully")

	//REPOSITORIES
	userRepo := repos.NewUserRepo()
	tweetRepo := repos.NewTweetRepo()

	//SERVICES
	userSvc := userService.NewUserService(userRepo)
	tweetSvc := tweetService.NewTweetService(tweetRepo)

	//ROUTERS
	router := global.NewRouters(userSvc, tweetSvc)

	router.Run(":8080")

}
