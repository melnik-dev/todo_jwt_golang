package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/melnik-dev/go_todo_jwt/internal/model"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Auth interface {
	CreateUser(username, password string) (int, error)
	GetUser(username string) (model.User, error)
}

type Task interface {
	Create(userID int, title, description string) (int, error)
	Update(userID, taskID int, title, description string, completed bool) error
	DeleteById(userID, taskID int) error
	GetById(userID, taskID int) (model.Task, error)
	GetAll(userID int) ([]model.Task, error)
}

type Repository struct {
	Auth *AuthRepo
	Task *TaskRepo
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Auth: NewAuthRepo(db),
		Task: NewTaskRepo(db),
	}
}
