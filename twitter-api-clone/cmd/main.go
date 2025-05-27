package main

import (
	"database/sql"
	"twitter-api-clone/internal/adapter/http"
	mysqlrepo "twitter-api-clone/internal/repository/mysql"
	"twitter-api-clone/internal/service"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/db?parseTime=true")

	if err != nil {
		panic(err.Error()) // handle error properly instead of panic in real app
	}

	defer db.Close()

	userRepo := mysqlrepo.NewUserRepository(db)
	tweetRepo := mysqlrepo.NewTweetRepository(db)
	followRepo := mysqlrepo.NewFollowRepository(db)

	tweetSvc := service.NewTweetService(tweetRepo)
	userSvc := service.NewUserService(userRepo)
	followSvc := service.NewFollowService(followRepo)

	router := http.NewRouter(tweetSvc, userSvc, followSvc)

	router.Run(":8080")
}
