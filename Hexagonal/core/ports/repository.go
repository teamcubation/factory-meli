package ports

import "hex-todo/core/model"

type TaskRepository interface {
	Save(task *model.Task) error
	FindAll() ([]*model.Task, error)
}
