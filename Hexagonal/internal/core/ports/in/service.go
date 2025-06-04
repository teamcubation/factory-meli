package in

import "hex-todo/internal/core/domain"

type TaskService interface {
	CreateTask(title string) (*domain.Task, error)
	ListTasks() ([]*domain.Task, error)
}
