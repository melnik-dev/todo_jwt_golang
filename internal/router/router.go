package router

import (
	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/internal/handler"
	"github.com/melnik-dev/go_todo_jwt/internal/middleware"
	"net/http"
)

func NewRouter(h *handler.Handler, jwtSecret string) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", ping)
	{
		auth := r.Group("/auth")
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
	}
	{
		task := r.Group("/task")
		task.Use(middleware.ValidateToken(jwtSecret))
		task.POST("/create", h.Task.Create)
		task.PUT("/:id", h.Task.Update)
		task.DELETE("/:id", h.Task.Delete)
		task.GET("/:id", h.Task.Get)
		task.GET("/", h.Task.GetAll)
	}
	return r
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
