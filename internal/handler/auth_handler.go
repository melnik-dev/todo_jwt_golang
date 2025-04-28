package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/internal/model"
	"github.com/melnik-dev/go_todo_jwt/internal/service"
	"net/http"
)

type AuthHandler struct {
	AuthService service.Auth
}

func NewAuthHandler(service service.Auth) *AuthHandler {
	return &AuthHandler{AuthService: service}
}

func (ah *AuthHandler) Register(c *gin.Context) {
	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid input format"})
		return
	}

	token, err := ah.AuthService.Register(input.Name, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid input format"})
		return
	}

	token, err := ah.AuthService.Login(input.Name, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
