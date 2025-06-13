package task

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/configs"
	"github.com/melnik-dev/go_todo_jwt/pkg/middleware"
	"github.com/melnik-dev/go_todo_jwt/pkg/response"
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
		log.Printf("Failed to bind JSON in Create task: %v", err)
		response.BadRequest(c, "Invalid input data")
		return
	}

	taskId, err := h.TaskService.Create(userID, input.Title, input.Description)
	if err != nil {
		log.Printf("Failed to create task for user %d: %v", userID, err)
		response.InternalServerError(c, "Failed to create task")
		return
	}

	response.Success(c, http.StatusOK, CreateResponse{ID: taskId})
}

func (h *Handler) Update(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var uri URIParam
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Printf("Failed to bind URI in Update task: %v", err)
		response.BadRequest(c, "Invalid task ID")
		return
	}

	var input UpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Failed to bind JSON in Update task: %v", err)
		response.BadRequest(c, "Invalid input data")
		return
	}

	err := h.TaskService.Update(userID, uri.ID, input.Title, input.Description, input.Completed)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			response.NotFound(c, ErrTaskNotFound.Error())
			return
		}
		log.Printf("Failed to update task %d for user %d: %v", uri.ID, userID, err)
		response.InternalServerError(c, "Failed to update task")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func (h *Handler) Delete(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var uri URIParam
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Printf("Failed to bind URI in Delete task: %v", err)
		response.BadRequest(c, "Invalid task ID")
		return
	}

	err := h.TaskService.Delete(userID, uri.ID)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			response.NotFound(c, ErrTaskNotFound.Error())
			return
		}
		log.Printf("Failed to delete task %d for user %d: %v", uri.ID, userID, err)
		response.InternalServerError(c, "Failed to delete task")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (h *Handler) Get(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var uri URIParam
	if err := c.ShouldBindUri(&uri); err != nil {
		log.Printf("Failed to bind URI in Get task: %v", err)
		response.BadRequest(c, "Invalid task ID")
		return
	}

	task, err := h.TaskService.GetById(userID, uri.ID)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			response.NotFound(c, ErrTaskNotFound.Error())
			return
		}
		log.Printf("Failed to get task %d for user %d: %v", uri.ID, userID, err)
		response.InternalServerError(c, "Failed to get task")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"task": task})
}

func (h *Handler) GetAll(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	tasks, err := h.TaskService.GetAll(userID)
	if err != nil {
		log.Printf("Failed to get all tasks for user %d: %v", userID, err)
		response.InternalServerError(c, "Failed to get tasks")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"tasks": tasks})
}

func getUserID(c *gin.Context) (int, bool) {
	strId, ok := c.Get("user_id")
	if !ok {
		log.Printf("User ID not found in context")
		response.InternalServerError(c, ErrUserNotFound.Error())
		return 0, false
	}
	userID, ok := strId.(int)
	if !ok {
		log.Printf("Invalid user_id type in context")
		response.InternalServerError(c, "Invalid user ID type")
		return 0, false
	}
	return userID, true
}
