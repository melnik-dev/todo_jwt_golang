package task

import (
	"errors"
	"github.com/melnik-dev/go_todo_jwt/pkg/logger"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/configs"
	"github.com/melnik-dev/go_todo_jwt/pkg/middleware"
	"github.com/melnik-dev/go_todo_jwt/pkg/response"
)

type HandlerDeps struct {
	TaskService IService
	*configs.Config
}

type Handler struct {
	TaskService IService
	*configs.Config
}

func NewHandler(r *gin.Engine, deps *HandlerDeps) {
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
	logHandle := handlerLogger(c)
	logHandle.Debug("Received request to Create")

	userID, ok := middleware.GetUserID(c)
	if !ok {
		return
	}

	var input CreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		logHandle.WithError(err).Warn("Failed to bind JSON in Create")
		response.BadRequest(c, "Invalid input")
		return
	}

	taskId, err := h.TaskService.Create(userID, input.Title, input.Description)
	if err != nil {
		logHandle.WithError(err).Error("Failed to Create")
		response.InternalServerError(c, "Failed to create task")
		return
	}

	logHandle.Debug("Create successfully")
	response.Success(c, http.StatusOK, CreateResponse{ID: taskId})
}

func (h *Handler) Update(c *gin.Context) {
	logHandle := handlerLogger(c)
	logHandle.Debug("Received request to Update")

	userID, ok := middleware.GetUserID(c)
	if !ok {
		return
	}

	var uri URIParam
	if err := c.ShouldBindUri(&uri); err != nil {
		logHandle.WithError(err).Warn("Failed to bind ID")
		response.BadRequest(c, "Invalid task ID")
		return
	}
	logHandle = logHandle.WithField("task_id", uri.ID)

	var input UpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		logHandle.WithError(err).Warn("Failed to bind JSON in Update")
		response.BadRequest(c, "Invalid input data")
		return
	}

	err := h.TaskService.Update(userID, uri.ID, input.Title, input.Description, input.Completed)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			logHandle.Warn(ErrTaskNotFound.Error())
			response.NotFound(c, ErrTaskNotFound.Error())
			return
		}
		logHandle.WithError(err).Error("Failed to Update")
		response.InternalServerError(c, "Failed to update task")
		return
	}

	logHandle.Debug("Update successfully")
	response.Success(c, http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func (h *Handler) Delete(c *gin.Context) {
	logHandle := handlerLogger(c)
	logHandle.Debug("Received request to Delete")

	userID, ok := middleware.GetUserID(c)
	if !ok {
		return
	}

	var uri URIParam
	if err := c.ShouldBindUri(&uri); err != nil {
		logHandle.WithError(err).Warn("Failed to bind ID")
		response.BadRequest(c, "Invalid task ID")
		return
	}
	logHandle = logHandle.WithField("task_id", uri.ID)

	err := h.TaskService.Delete(userID, uri.ID)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			logHandle.WithError(err).Warn(ErrTaskNotFound.Error())
			response.NotFound(c, ErrTaskNotFound.Error())
			return
		}
		logHandle.WithError(err).Error("Failed to delete")
		response.InternalServerError(c, "Failed to delete task")
		return
	}

	logHandle.Debug("Delete successfully")
	response.Success(c, http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (h *Handler) Get(c *gin.Context) {
	logHandle := handlerLogger(c)
	logHandle.Debug("Received request to Get task")

	userID, ok := middleware.GetUserID(c)
	if !ok {
		return
	}

	var uri URIParam
	if err := c.ShouldBindUri(&uri); err != nil {
		logHandle.WithError(err).Warn("Failed to bind task ID")
		response.BadRequest(c, "Invalid task ID")
		return
	}
	logHandle = logHandle.WithField("task_id", uri.ID)

	task, err := h.TaskService.GetById(userID, uri.ID)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			logHandle.WithError(err).Warn(ErrTaskNotFound.Error())
			response.NotFound(c, ErrTaskNotFound.Error())
			return
		}
		logHandle.WithError(err).Error("Failed to GetById")
		response.InternalServerError(c, "Failed to get task")
		return
	}

	logHandle.Debug("GetById successfully")
	response.Success(c, http.StatusOK, gin.H{"task": task})
}

func (h *Handler) GetAll(c *gin.Context) {
	logHandle := handlerLogger(c)
	logHandle.Debug("Received request to GetAll")

	userID, ok := middleware.GetUserID(c)
	if !ok {
		return
	}

	tasks, err := h.TaskService.GetAll(userID)
	if err != nil {
		logHandle.WithError(err).Error("Failed to get all tasks")
		response.InternalServerError(c, "Failed to get tasks")
		return
	}

	logHandle.Debug("GetAll successfully")
	response.Success(c, http.StatusOK, gin.H{"tasks": tasks})
}

func handlerLogger(c *gin.Context) *logrus.Entry {
	return logger.FromContext(c).WithField("layer", "Handler task layer")
}
