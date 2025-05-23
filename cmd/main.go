package main

import (
	"tweet-tq-rafael/core/ports"
	"tweet-tq-rafael/internal/adapters/databases"
	"tweet-tq-rafael/internal/adapters/http"
	"tweet-tq-rafael/internal/services"
)

func main() {

	//Creating repositories
	db, err := databases.InitPostgresDB()
	if err != nil {
		panic("Error connecting DB")
	}
	repos := databases.RepositoriesFactory(db)

	// Creating services
	userService := services.UserService{Repo: repos.UserRepository}
	tweetService := services.TweetService{Repo: repos.TweetRepository}
	services := ports.Services{
		UserService:  userService,
		TweetService: tweetService,
	}

	//Creating router
	router := http.NewRouter(services)
	router.Run(":8080")
}
