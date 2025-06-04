package main

import (
	"hex-todo/internal/adapter/http"
	"hex-todo/internal/adapter/memory"
	"hex-todo/internal/application"
)

func main() {
	repo := memory.NewInMemoryRepo()
	svc := application.NewTaskService(repo)

	router := http.NewRouter(svc)

	router.Run(":8080")
}
