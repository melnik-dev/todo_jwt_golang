package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/configs"
	"github.com/melnik-dev/go_todo_jwt/pkg/jwt"
	"github.com/melnik-dev/go_todo_jwt/pkg/response"
)

type HandlerDeps struct {
	AuthService *Service
	*configs.Config
}

type Handler struct {
	AuthService *Service
	*configs.Config
}

func NewHandler(r *gin.Engine, deps HandlerDeps) {
	handler := &Handler{
		AuthService: deps.AuthService,
		Config:      deps.Config,
	}
	auth := r.Group("/auth")
	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)
}

func (h *Handler) Register(c *gin.Context) {
	var input RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Failed to bind JSON in Register: %v", err)
		response.BadRequest(c, "Invalid input data")
		return
	}

	userId, err := h.AuthService.Register(input.Name, input.Password)
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			response.BadRequest(c, ErrUserExists.Error())
			return
		}
		log.Printf("Failed to register user %s: %v", input.Name, err)
		response.InternalServerError(c, "Failed to register user")
		return
	}

	token, err := jwt.NewJWT(h.Config.JWT.Secret).Create(jwt.Data{
		UserId:   userId,
		TokenTTL: h.Config.JWT.TokenTTL,
	})
	if err != nil {
		log.Printf("Failed to create JWT for user %d: %v", userId, err)
		response.InternalServerError(c, "Failed to create authentication token")
		return
	}

	response.Success(c, http.StatusOK, RegisterResponse{Token: token})
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Failed to bind JSON in Login: %v", err)
		response.BadRequest(c, "Invalid input data")
		return
	}

	userId, err := h.AuthService.Login(input.Name, input.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidLogin) {
			response.Unauthorized(c, ErrInvalidLogin.Error())
			return
		}
		log.Printf("Failed to login user %s: %v", input.Name, err)
		response.InternalServerError(c, "Failed to login")
		return
	}

	token, err := jwt.NewJWT(h.Config.JWT.Secret).Create(jwt.Data{
		UserId:   userId,
		TokenTTL: h.Config.JWT.TokenTTL,
	})
	if err != nil {
		log.Printf("Failed to create JWT for user %d: %v", userId, err)
		response.InternalServerError(c, "Failed to create authentication token")
		return
	}

	response.Success(c, http.StatusOK, LoginResponse{Token: token})
}
