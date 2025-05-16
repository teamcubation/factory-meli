package memory

import (
	"hex-todo/core/model"
)

type InMemoryRepo struct {
	tasks []*model.Task
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		tasks: make([]*model.Task, 0),
	}
}

func (r *InMemoryRepo) Save(task *model.Task) error {
	r.tasks = append(r.tasks, task)
	return nil
}

func (r *InMemoryRepo) FindAll() ([]*model.Task, error) {
	return r.tasks, nil
}
