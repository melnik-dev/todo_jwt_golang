package task

import (
	"github.com/melnik-dev/go_todo_jwt/internal/user"
	"github.com/sirupsen/logrus"
)

type IService interface {
	Create(userID int, title, desc string) (int, error)
	Update(userID, taskID int, title, desc string, completed bool) error
	Delete(userID, taskID int) error
	GetById(userID, taskID int) (*Task, error)
	GetAll(userID int) ([]Task, error)
}

type Service struct {
	taskRepo IRepository
	logger   *logrus.Logger
}

func NewService(repo IRepository, logger *logrus.Logger) *Service {
	return &Service{
		taskRepo: repo,
		logger:   logger,
	}
}

func (s *Service) Create(userID int, title, desc string) (int, error) {
	logServ := serviceLogger(s.logger).WithFields(logrus.Fields{
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
	logServ := serviceLogger(s.logger).WithFields(logrus.Fields{
		"user_id": userID,
		"task_id": taskID,
	})
	logServ.Debug("Attempting to Update")

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
	logServ := serviceLogger(s.logger).WithFields(logrus.Fields{
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
	logServ := serviceLogger(s.logger).WithFields(logrus.Fields{
		"user_id": userID,
		"task_id": taskID,
	})
	logServ.Debug("Attempting to Delete")

	task, err := s.taskRepo.GetById(&Task{
		ID:     taskID,
		UserID: userID,
	})
	if err != nil {
		logServ.WithError(err).Error("Failed to GetById")
		return nil, err
	}

	logServ.Debug("GetById successfully")
	return task, nil
}

func (s *Service) GetAll(userID int) ([]Task, error) {
	logServ := serviceLogger(s.logger).WithFields(logrus.Fields{
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

func serviceLogger(l *logrus.Logger) *logrus.Entry {
	return l.WithField("layer", "Service task layer")
}
