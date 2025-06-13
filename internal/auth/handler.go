package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/configs"
	"github.com/melnik-dev/go_todo_jwt/pkg/jwt"
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"err": http.StatusText(http.StatusBadRequest)})
		return
	}

	userId, err := h.AuthService.Register(input.Name, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}
	token, err := jwt.NewJWT(h.Config.JWT.Secret).Create(jwt.Data{
		UserId:   userId,
		TokenTTL: h.Config.JWT.TokenTTL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	data := RegisterResponse{
		Token: token,
	}
	c.JSON(http.StatusOK, data)
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": http.StatusText(http.StatusBadRequest)})
		return
	}

	userId, err := h.AuthService.Login(input.Name, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
		return
	}

	token, err := jwt.NewJWT(h.Config.JWT.Secret).Create(jwt.Data{
		UserId:   userId,
		TokenTTL: h.Config.JWT.TokenTTL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	data := LoginResponse{
		Token: token,
	}
	c.JSON(http.StatusOK, data)
}
