package task

import (
	"github.com/melnik-dev/go_todo_jwt/internal/user"
	"github.com/melnik-dev/go_todo_jwt/pkg/logger"
	"github.com/sirupsen/logrus"
)

type Service struct {
	taskRepo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{taskRepo: repo}
}

func (s *Service) Create(userID int, title, desc string) (int, error) {
	logServ := serviceLogger().WithFields(logrus.Fields{
		"user_id": userID,
	})
	logServ.Debug("Attempting to Create")

	task := &Task{
		UserID:      userID,
		Title:       title,
		Description: desc,
	}
	_, err := s.taskRepo.Create(task)
	if err != nil {
		logServ.WithError(err).Error("Failed to Create")
		return 0, err
	}

	logServ.Debug("Create successfully")
	return task.ID, nil
}

func (s *Service) Update(userID, taskID int, title, desc string, completed bool) error {
	logServ := serviceLogger().WithFields(logrus.Fields{
		"user_id": userID,
		"task_id": taskID,
	})
	logServ.Debug("Attempting to Create")

	task := &Task{
		ID:          taskID,
		UserID:      userID,
		Title:       title,
		Description: desc,
		Completed:   completed,
	}
	err := s.taskRepo.Update(task)
	if err != nil {
		logServ.WithError(err).Error("Failed to Update")
		return err
	}

	logServ.Debug("Update successfully")
	return nil
}

func (s *Service) Delete(userID, taskID int) error {
	logServ := serviceLogger().WithFields(logrus.Fields{
		"user_id": userID,
		"task_id": taskID,
	})
	logServ.Debug("Attempting to Delete")

	task := &Task{
		ID:     taskID,
		UserID: userID,
	}

	err := s.taskRepo.DeleteById(task)
	if err != nil {
		logServ.WithError(err).Error("Failed to Delete")
		return err
	}

	logServ.Debug("Delete successfully")
	return nil
}

func (s *Service) GetById(userID, taskID int) (*Task, error) {
	logServ := serviceLogger().WithFields(logrus.Fields{
		"user_id": userID,
		"task_id": taskID,
	})
	logServ.Debug("Attempting to Delete")

	task := &Task{
		ID:     taskID,
		UserID: userID,
	}

	_, err := s.taskRepo.GetById(task)
	if err != nil {
		logServ.WithError(err).Error("Failed to GetById")
		return nil, err
	}

	logServ.Debug("GetById successfully")
	return task, nil
}

func (s *Service) GetAll(userID int) ([]Task, error) {
	logServ := serviceLogger().WithFields(logrus.Fields{
		"user_id": userID,
	})
	logServ.Debug("Attempting to GetAll")

	us := &user.User{ID: userID}

	tasks, err := s.taskRepo.GetAll(us)
	if err != nil {
		logServ.WithError(err).Error("Failed to GetAll")
		return nil, err
	}

	logServ.Debug("GetAll successfully")
	return tasks, nil
}

func serviceLogger() *logrus.Entry {
	return logger.GetLogger().WithField("layer", "Service task layer")
}
