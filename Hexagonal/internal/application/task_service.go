package application

import (
	"errors"
	"hex-todo/internal/core/domain"
	portIn "hex-todo/internal/core/ports/in"
	portOut "hex-todo/internal/core/ports/out"
)

type TaskServiceImpl struct {
	repo   portOut.TaskRepository
	nextID int
}

func NewTaskService(repo portOut.TaskRepository) portIn.TaskService {
	return &TaskServiceImpl{
		repo:   repo,
		nextID: 1,
	}
}

func (s *TaskServiceImpl) CreateTask(title string) (*domain.Task, error) {
	if title == "" {
		return nil, errors.New("el título no puede estar vacío")
	}

	task := &domain.Task{
		ID:    s.nextID,
		Title: title,
	}
	s.nextID++

	err := s.repo.Save(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskServiceImpl) ListTasks() ([]*domain.Task, error) {
	return s.repo.FindAll()
}
