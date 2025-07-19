package auth

import (
	"errors"
	"github.com/melnik-dev/go_todo_jwt/pkg/logger"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/configs"
	"github.com/melnik-dev/go_todo_jwt/pkg/jwt"
	"github.com/melnik-dev/go_todo_jwt/pkg/response"
)

type HandlerDeps struct {
	AuthService IService
	*configs.Config
}

type Handler struct {
	AuthService IService
	*configs.Config
}

func NewHandler(r *gin.Engine, deps *HandlerDeps) {
	handler := &Handler{
		AuthService: deps.AuthService,
		Config:      deps.Config,
	}
	auth := r.Group("/auth")
	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)
}

func (h *Handler) Register(c *gin.Context) {
	logHandle := handlerLogger(c)
	logHandle.Debug("Received request to Register")

	var input RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		logHandle.WithError(err).Warn("Failed to bind JSON in Register")
		response.BadRequest(c, "Invalid input data")
		return
	}
	logHandle = logHandle.WithField("user_name", input.Name)

	userId, err := h.AuthService.Register(input.Name, input.Password)
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			logHandle.Warn(ErrUserExists.Error())
			response.BadRequest(c, ErrUserExists.Error())
			return
		}
		logHandle.WithError(err).Error("Failed to register user")
		response.InternalServerError(c, "Failed to register user")
		return
	}
	logHandle = logHandle.WithField("user_id", userId)

	token, err := jwt.NewJWT(h.Config.JWT.Secret).Create(jwt.Data{
		UserId:   userId,
		TokenTTL: h.Config.JWT.TokenTTL,
	})
	if err != nil {
		logHandle.WithError(err).Error("Failed to create JWT for user")
		response.InternalServerError(c, "Failed to create authentication token")
		return
	}

	logHandle.Debug("Register successfully")
	response.Success(c, http.StatusOK, RegisterResponse{Token: token})
}

func (h *Handler) Login(c *gin.Context) {
	logHandle := handlerLogger(c)
	logHandle.Debug("Received request to Login")

	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		logHandle.WithError(err).Warn("Failed to bind JSON in Login")
		response.BadRequest(c, "Invalid input data")
		return
	}
	logHandle = logHandle.WithField("user_name", input.Name)

	userId, err := h.AuthService.Login(input.Name, input.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidLogin) {
			logHandle.Warn(ErrInvalidLogin.Error())
			response.Unauthorized(c, ErrInvalidLogin.Error())
			return
		}
		logHandle.WithError(err).Error("Failed to login user")
		response.InternalServerError(c, "Failed to login")
		return
	}

	token, err := jwt.NewJWT(h.Config.JWT.Secret).Create(jwt.Data{
		UserId:   userId,
		TokenTTL: h.Config.JWT.TokenTTL,
	})
	if err != nil {
		logHandle.WithError(err).Error("Failed to create JWT for user")
		response.InternalServerError(c, "Failed to create authentication token")
		return
	}

	logHandle.Debug("Login successfully")
	response.Success(c, http.StatusOK, LoginResponse{Token: token})
}

func handlerLogger(c *gin.Context) *logrus.Entry {
	return logger.FromContext(c).WithField("layer", "Handler auth layer")
}
