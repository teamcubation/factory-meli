package memory

import (
	"hex-todo/internal/core/domain"
)

type InMemoryRepo struct {
	tasks []*domain.Task
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		tasks: make([]*domain.Task, 0),
	}
}

func (r *InMemoryRepo) Save(task *domain.Task) error {
	r.tasks = append(r.tasks, task)
	return nil
}

func (r *InMemoryRepo) FindAll() ([]*domain.Task, error) {
	return r.tasks, nil
}
