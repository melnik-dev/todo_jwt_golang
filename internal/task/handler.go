package task

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/configs"
	"github.com/melnik-dev/go_todo_jwt/pkg/middleware"
	"net/http"
)

type HandlerDeps struct {
	TaskService *Service
	*configs.Config
}

type Handler struct {
	TaskService *Service
	*configs.Config
}

func NewHandler(r *gin.Engine, deps HandlerDeps) {
	handler := &Handler{
		TaskService: deps.TaskService,
		Config:      deps.Config,
	}
	task := r.Group("/task")
	task.Use(middleware.IsAuthed(handler.Config.JWT.Secret))
	task.POST("/create", handler.Create)
	task.PUT("/:id", handler.Update)
	task.DELETE("/:id", handler.Delete)
	task.GET("/:id", handler.Get)
	task.GET("/", handler.GetAll)
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var input CreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": http.StatusText(http.StatusBadRequest)})
		return
	}

	taskId, err := h.TaskService.Create(userID, input.Title, input.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "failed  to create task"})
		return
	}

	data := CreateResponse{
		ID: taskId,
	}
	c.JSON(http.StatusOK, data)
}

func (h *Handler) Update(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var uri URIParam
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": http.StatusText(http.StatusBadRequest)})
		return
	}

	var input UpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": http.StatusText(http.StatusBadRequest)})
		return
	}

	err := h.TaskService.Update(userID, uri.ID, input.Title, input.Description, input.Completed)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"err": ErrTaskNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": ErrFailedGetTask.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "task updated"})
}

func (h *Handler) Delete(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var uri URIParam
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": http.StatusText(http.StatusBadRequest)})
		return
	}

	err := h.TaskService.Delete(userID, uri.ID)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"err": ErrTaskNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": ErrFailedGetTask.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "task deleted"})
}

func (h *Handler) Get(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var uri URIParam
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": http.StatusText(http.StatusBadRequest)})
		return
	}

	task, err := h.TaskService.GetById(userID, uri.ID)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"err": ErrTaskNotFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": ErrFailedGetTask.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (h *Handler) GetAll(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	tasks, err := h.TaskService.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": ErrFailedGetTask.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func getUserID(c *gin.Context) (int, bool) {
	strId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"err": ErrUserNotFound})
		return 0, false
	}
	userID, ok := strId.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type"})
		return 0, false
	}
	return userID, true
}
