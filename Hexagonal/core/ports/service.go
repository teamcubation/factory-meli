package ports

import "hex-todo/core/model"

type TaskService interface {
	CreateTask(title string) (*model.Task, error)
	ListTasks() ([]*model.Task, error)
}
