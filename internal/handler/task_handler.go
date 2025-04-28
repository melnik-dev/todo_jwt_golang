package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/internal/model"
	repoErrors "github.com/melnik-dev/go_todo_jwt/internal/repository"
	"github.com/melnik-dev/go_todo_jwt/internal/service"
	serviceErrors "github.com/melnik-dev/go_todo_jwt/internal/service"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	TaskService service.Task
}

func NewTaskHandler(service service.Task) *TaskHandler {
	return &TaskHandler{TaskService: service}
}

func (th *TaskHandler) Create(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": ErrUserNotAuth})
		return
	}

	var input model.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid input format"})
		return
	}

	taskId, err := th.TaskService.CreateTask(userID, input.Title, input.Description)
	if err != nil {
		if errors.Is(err, serviceErrors.ErrRequiredField) {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		// if errors.Is(err, serviceErrors.ErrValidationFailed)

		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to create task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": taskId})
}

func (th *TaskHandler) Update(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": ErrUserNotAuth})
		return
	}

	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid task ID format"})
		return
	}

	var input model.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid input format"})
		return
	}
	// TODO struct UpdateTask
	completed := input.Completed
	err = th.TaskService.UpdateTask(userID, taskId, input.Title, input.Description, &completed)
	if err != nil {
		if errors.Is(err, repoErrors.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"err": "Task not found"})
			return
		}
		if errors.Is(err, serviceErrors.ErrRequiredField) {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to get task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "task updated"})
}

func (th *TaskHandler) Delete(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": ErrUserNotAuth})
		return
	}

	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid task ID format"})
		return
	}

	err = th.TaskService.DeleteTask(userID, taskId)
	if err != nil {
		if errors.Is(err, repoErrors.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"err": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to get task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "task deleted"})
}

func (th *TaskHandler) Get(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": ErrUserNotAuth})
		return
	}

	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid task ID format"})
		return
	}

	task, err := th.TaskService.GetTaskById(userID, taskId)
	if err != nil {
		if errors.Is(err, repoErrors.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"err": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to get task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (th *TaskHandler) GetAll(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": ErrUserNotAuth})
		return
	}

	tasks, err := th.TaskService.GetTasks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to get task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func getUserID(c *gin.Context) (int, error) {
	userID, ok := c.Get("user_id")
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	idInt, ok := userID.(int)
	if !ok {
		return 0, errors.New("user ID in context is not of type int")
	}
	return idInt, nil
}
