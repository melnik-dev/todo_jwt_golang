package task

import (
	"github.com/melnik-dev/go_todo_jwt/internal/user"
)

type Service struct {
	taskRepo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{taskRepo: repo}
}

func (service *Service) Create(userID int, title, desc string) (int, error) {
	task := &Task{
		UserID:      userID,
		Title:       title,
		Description: desc,
	}
	_, err := service.taskRepo.Create(task)
	if err != nil {
		return 0, err
	}

	return task.ID, nil
}

func (service *Service) Update(userID, taskID int, title, desc string, completed bool) error {
	task := &Task{
		ID:          taskID,
		UserID:      userID,
		Title:       title,
		Description: desc,
		Completed:   completed,
	}
	err := service.taskRepo.Update(task)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) Delete(userID, taskIdD int) error {
	task := &Task{
		ID:     taskIdD,
		UserID: userID,
	}

	err := service.taskRepo.DeleteById(task)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) GetById(userID, taskIdD int) (*Task, error) {
	task := &Task{
		ID:     taskIdD,
		UserID: userID,
	}

	_, err := service.taskRepo.GetById(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (service *Service) GetAll(userID int) ([]Task, error) {
	us := &user.User{ID: userID}

	tasks, err := service.taskRepo.GetAll(us)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
