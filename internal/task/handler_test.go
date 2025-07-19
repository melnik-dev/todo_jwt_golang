package task_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/internal/task"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockTaskService struct {
	CreateMock  func(userID int, title, desc string) (int, error)
	UpdateMock  func(userID, taskID int, title, desc string, completed bool) error
	DeleteMock  func(userID, taskID int) error
	GetByIdMock func(userID, taskID int) (*task.Task, error)
	GetAllMock  func(userID int) ([]task.Task, error)
}

func (m *MockTaskService) Create(userID int, title, desc string) (int, error) {
	return m.CreateMock(userID, title, desc)
}

func (m *MockTaskService) Update(userID, taskID int, title, desc string, completed bool) error {
	return m.UpdateMock(userID, taskID, title, desc, completed)
}

func (m *MockTaskService) Delete(userID, taskID int) error {
	return m.DeleteMock(userID, taskID)
}

func (m *MockTaskService) GetById(userID, taskID int) (*task.Task, error) {
	return m.GetByIdMock(userID, taskID)
}

func (m *MockTaskService) GetAll(userID int) ([]task.Task, error) {
	return m.GetAllMock(userID)
}

func mockGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		logger := logrus.New()
		logger.SetOutput(io.Discard)
		c.Set("logger", logrus.NewEntry(logger))
		c.Set("user_id", 42)
		c.Next()
	})
	return r
}

type Options struct {
	h      *task.Handler
	title  string
	taskID any
}

func requestCreateHelper(t *testing.T, opts Options) *httptest.ResponseRecorder {
	t.Helper()
	r := mockGin()
	r.POST("/task/create", opts.h.Create)

	if opts.title == "empty" {
		opts.title = ""
	} else {
		opts.title = "test_title"
	}

	body, _ := json.Marshal(map[string]string{
		"title":       opts.title,
		"description": "test_desc",
	})

	req := httptest.NewRequest(http.MethodPost, "/task/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func requestUpdateHelper(t *testing.T, opts Options) *httptest.ResponseRecorder {
	t.Helper()
	r := mockGin()
	r.PUT("/task/:id", opts.h.Update)

	if opts.title == "empty" {
		opts.title = ""
	} else {
		opts.title = "test_title"
	}
	if opts.taskID == nil {
		opts.taskID = 1
	}

	body, _ := json.Marshal(map[string]any{
		"title":       opts.title,
		"description": "test_desc",
		"completed":   true,
	})

	url := fmt.Sprintf("/task/%s", fmt.Sprint(opts.taskID))
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func requestDeleteHelper(t *testing.T, opts Options) *httptest.ResponseRecorder {
	t.Helper()
	r := mockGin()
	r.DELETE("/task/:id", opts.h.Delete)

	if opts.taskID == nil {
		opts.taskID = 1
	}

	url := fmt.Sprintf("/task/%s", fmt.Sprint(opts.taskID))
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func requestGetHelper(t *testing.T, opts Options) *httptest.ResponseRecorder {
	t.Helper()
	r := mockGin()
	r.GET("/task/:id", opts.h.Get)

	if opts.taskID == nil {
		opts.taskID = 1
	}

	url := fmt.Sprintf("/task/%s", fmt.Sprint(opts.taskID))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func requestGetAllHelper(t *testing.T, opts Options) *httptest.ResponseRecorder {
	t.Helper()
	r := mockGin()
	r.GET("/task", opts.h.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/task", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestHandler_Create_Success(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			CreateMock: func(userID int, title, desc string) (int, error) {
				return 1, nil
			},
		},
	}

	w := requestCreateHelper(t, Options{h: handler})

	if w.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
	}

	type ResponseWrapper struct {
		Status int                 `json:"status"`
		Data   task.CreateResponse `json:"data"`
	}

	var res ResponseWrapper
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Fatal(err)
	}

	if res.Data.ID != 1 {
		t.Errorf("expected 1, got %d", res.Data.ID)
	}
}

func TestHandler_Create_Fail(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			CreateMock: func(userID int, title, desc string) (int, error) {
				return 0, fmt.Errorf("test error")
			},
		},
	}

	w := requestCreateHelper(t, Options{h: handler})

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandler_Register_FailInvalid(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			CreateMock: func(userID int, title, desc string) (int, error) {
				return 0, fmt.Errorf("invalid input data")
			},
		},
	}

	w := requestCreateHelper(t, Options{h: handler, title: "empty"})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_Update_Success(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			UpdateMock: func(userID, taskID int, title, desc string, completed bool) error {
				return nil
			},
		},
	}

	w := requestUpdateHelper(t, Options{h: handler})

	if w.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandler_Update_Fail(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			UpdateMock: func(userID, taskID int, title, desc string, completed bool) error {
				return fmt.Errorf("test error")
			},
		},
	}

	w := requestUpdateHelper(t, Options{h: handler})

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandler_Update_FailNotFound(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			UpdateMock: func(userID, taskID int, title, desc string, completed bool) error {
				return task.ErrTaskNotFound
			},
		},
	}

	w := requestUpdateHelper(t, Options{h: handler})

	if w.Code != http.StatusNotFound {
		t.Errorf("expected %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandler_Update_FailInvalidData(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			UpdateMock: func(userID, taskID int, title, desc string, completed bool) error {
				return fmt.Errorf("invalid input data")
			},
		},
	}

	w := requestUpdateHelper(t, Options{h: handler, title: "empty"})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_Update_FailInvalidID(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			UpdateMock: func(userID, taskID int, title, desc string, completed bool) error {
				return fmt.Errorf("invalid id")
			},
		},
	}

	w := requestUpdateHelper(t, Options{h: handler, taskID: "invalid_id"})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_Delete_Success(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			DeleteMock: func(userID, taskID int) error {
				return nil
			},
		},
	}

	w := requestDeleteHelper(t, Options{h: handler})

	if w.Code != http.StatusOK {
		t.Errorf("expected %d, got %v", http.StatusOK, w.Code)
	}
}

func TestHandler_Delete_Fail(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			DeleteMock: func(userID, taskID int) error {
				return fmt.Errorf("test error")
			},
		},
	}

	w := requestDeleteHelper(t, Options{h: handler})

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandler_Delete_FailNotFound(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			DeleteMock: func(userID, taskID int) error {
				return task.ErrTaskNotFound
			},
		},
	}

	w := requestDeleteHelper(t, Options{h: handler})

	if w.Code != http.StatusNotFound {
		t.Errorf("expected %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandler_Delete_FailInvalidID(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			DeleteMock: func(userID, taskID int) error {
				return fmt.Errorf("invalid id")
			},
		},
	}

	w := requestDeleteHelper(t, Options{h: handler, taskID: "invalid_id"})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_Get_Success(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			GetByIdMock: func(userID, taskID int) (*task.Task, error) {
				return &task.Task{}, nil
			},
		},
	}

	w := requestGetHelper(t, Options{h: handler})

	if w.Code != http.StatusOK {
		t.Errorf("expected %d, got %v", http.StatusOK, w.Code)
	}
}

func TestHandler_Get_Fail(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			GetByIdMock: func(userID, taskID int) (*task.Task, error) {
				return nil, fmt.Errorf("test error")
			},
		},
	}

	w := requestGetHelper(t, Options{h: handler})

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandler_Get_FailNotFound(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			GetByIdMock: func(userID, taskID int) (*task.Task, error) {
				return nil, task.ErrTaskNotFound
			},
		},
	}

	w := requestGetHelper(t, Options{h: handler})

	if w.Code != http.StatusNotFound {
		t.Errorf("expected %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandler_Get_FailInvalidID(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			GetByIdMock: func(userID, taskID int) (*task.Task, error) {
				return nil, fmt.Errorf("invalid id")
			},
		},
	}

	w := requestDeleteHelper(t, Options{h: handler, taskID: "invalid_id"})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_GetAll_Success(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			GetAllMock: func(userID int) ([]task.Task, error) {
				tasks := make([]task.Task, 0)
				return tasks, nil
			},
		},
	}

	w := requestGetAllHelper(t, Options{h: handler})

	if w.Code != http.StatusOK {
		t.Errorf("expected %d, got %v", http.StatusOK, w.Code)
	}
}

func TestHandler_GetAll_Fail(t *testing.T) {
	handler := &task.Handler{
		TaskService: &MockTaskService{
			GetAllMock: func(userID int) ([]task.Task, error) {
				return nil, fmt.Errorf("test error")
			},
		},
	}

	w := requestGetAllHelper(t, Options{h: handler})

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
