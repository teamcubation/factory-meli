package main

import (
	"hex-todo/internal/adapter/http"
	"hex-todo/internal/adapter/memory"
	"hex-todo/internal/service"
)

func main() {
	repo := memory.NewInMemoryRepo()
	svc := service.NewTaskService(repo)

	router := http.NewRouter(svc)

	router.Run(":8080")
}
