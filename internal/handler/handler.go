package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/internal/service"
)

var (
	ErrUserNotAuth = errors.New("internal server error (user not authenticated)")
)

type Auth interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type Task interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Get(c *gin.Context)
	GetAll(c *gin.Context)
}

type Handler struct {
	Auth *AuthHandler
	Task *TaskHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Auth: NewAuthHandler(service.Auth),
		Task: NewTaskHandler(service.Task),
	}
}
