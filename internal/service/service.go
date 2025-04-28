package service

import (
	"errors"
	"github.com/melnik-dev/go_todo_jwt/internal/model"
	"github.com/melnik-dev/go_todo_jwt/internal/repository"
	"time"
)

var (
	ErrRequiredField = errors.New("required field is missing")
)

type Auth interface {
	Register(username, password string) (string, error)
	Login(username, password string) (string, error)
	GenerateToken(userId int) (string, error)
	HashPassword(password string) (string, error)
	ComparePasswords(hash, password string) bool
}

type Task interface {
	CreateTask(userID int, title, desc string) (int, error)
	UpdateTask(userID, taskIdD int, title, desc string, completed *bool) error
	DeleteTask(userID, taskID int) error
	GetTaskById(userID, taskID int) (model.Task, error)
	GetTasks(userID int) ([]model.Task, error)
}

type Service struct {
	Auth *AuthService
	Task *TaskService
}

func NewService(repo *repository.Repository, jwtSecret string, time time.Duration) *Service {
	return &Service{
		Auth: NewAuthService(repo.Auth, jwtSecret, time),
		Task: NewTaskService(repo.Task),
	}
}
