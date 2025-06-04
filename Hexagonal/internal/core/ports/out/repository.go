package out

import "hex-todo/internal/core/domain"

type TaskRepository interface {
	Save(task *domain.Task) error
	FindAll() ([]*domain.Task, error)
}
