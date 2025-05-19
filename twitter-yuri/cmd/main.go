package main

import (
	"Yuri/twitter/db"
	"Yuri/twitter/internal/global"
	"Yuri/twitter/internal/tweet"
	"Yuri/twitter/internal/user"
	"fmt"
)

func main() {
	_, err := db.Connect()
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	fmt.Println("Connected to MongoDB successfully")

	//REPOSITORIES
	userRepo := user.NewUserRepo()
	tweetRepo := tweet.NewTweetRepo()

	//SERVICES
	userSvc := user.NewUserService(userRepo)
	tweetSvc := tweet.NewTweetService(tweetRepo)

	//ROUTERS
	router := global.NewRouters(userSvc, tweetSvc)

	router.Run(":8080")

}
