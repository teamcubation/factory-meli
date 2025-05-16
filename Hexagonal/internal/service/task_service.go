package service

import (
	"errors"
	"hex-todo/core/model"
	"hex-todo/core/ports"
)

type TaskServiceImpl struct {
	repo   ports.TaskRepository
	nextID int
}

func NewTaskService(repo ports.TaskRepository) ports.TaskService {
	return &TaskServiceImpl{
		repo:   repo,
		nextID: 1,
	}
}

func (s *TaskServiceImpl) CreateTask(title string) (*model.Task, error) {
	if title == "" {
		return nil, errors.New("el título no puede estar vacío")
	}

	task := &model.Task{
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

func (s *TaskServiceImpl) ListTasks() ([]*model.Task, error) {
	return s.repo.FindAll()
}
