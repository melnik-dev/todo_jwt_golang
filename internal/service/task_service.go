package service

import (
	"fmt"
	"github.com/melnik-dev/go_todo_jwt/internal/model"
	"github.com/melnik-dev/go_todo_jwt/internal/repository"
)

type TaskService struct {
	repo repository.Task
}

func NewTaskService(taskRepo repository.Task) *TaskService {
	return &TaskService{repo: taskRepo}
}

func (ts *TaskService) CreateTask(userID int, title, desc string) (int, error) {
	if title == "" {
		return 0, fmt.Errorf("service: %w", ErrRequiredField)
	}
	taskId, err := ts.repo.Create(userID, title, desc)
	if err != nil {
		return 0, err
	}

	return taskId, nil
}

func (ts *TaskService) UpdateTask(userID, taskIdD int, title, desc string, completed *bool) error {
	if userID == 0 || taskIdD == 0 || title == "" || desc == "" {
		return fmt.Errorf("service: %w", ErrRequiredField)
	}

	defaultCompleted := false
	if completed != nil {
		defaultCompleted = *completed
	}

	err := ts.repo.Update(userID, taskIdD, title, desc, defaultCompleted)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) DeleteTask(userID, taskIdD int) error {
	if userID == 0 || taskIdD == 0 {
		return fmt.Errorf("service: %w", ErrRequiredField)
	}

	err := ts.repo.DeleteById(userID, taskIdD)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) GetTaskById(userID, taskIdD int) (model.Task, error) {
	var task model.Task
	if userID == 0 || taskIdD == 0 {
		return task, fmt.Errorf("service: %w", ErrRequiredField)
	}

	task, err := ts.repo.GetById(userID, taskIdD)
	if err != nil {
		return task, err
	}

	return task, nil
}

func (ts *TaskService) GetTasks(userID int) ([]model.Task, error) {
	tasks := make([]model.Task, 0)
	if userID == 0 {
		return tasks, fmt.Errorf("service: %w", ErrRequiredField)
	}

	tasks, err := ts.repo.GetAll(userID)
	if err != nil {
		return tasks, err
	}

	return tasks, nil
}
