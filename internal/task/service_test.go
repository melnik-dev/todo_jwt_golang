package task_test

import (
	"fmt"
	"github.com/melnik-dev/go_todo_jwt/internal/task"
	"github.com/melnik-dev/go_todo_jwt/internal/user"
	"github.com/sirupsen/logrus"
	"io"
	"testing"
)

type MockTaskRepository struct {
	CreateMock     func(task *task.Task) (*task.Task, error)
	UpdateMock     func(task *task.Task) error
	DeleteByIdMock func(task *task.Task) error
	GetByIdMock    func(task *task.Task) (*task.Task, error)
	GetAllMock     func(user *user.User) ([]task.Task, error)
}

func (m *MockTaskRepository) Create(task *task.Task) (*task.Task, error) {
	return m.CreateMock(task)
}

func (m *MockTaskRepository) Update(task *task.Task) error {
	return m.UpdateMock(task)
}

func (m *MockTaskRepository) DeleteById(task *task.Task) error {
	return m.DeleteByIdMock(task)
}

func (m *MockTaskRepository) GetById(task *task.Task) (*task.Task, error) {
	return m.GetByIdMock(task)
}

func (m *MockTaskRepository) GetAll(user *user.User) ([]task.Task, error) {
	return m.GetAllMock(user)
}

func mockLogger() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(io.Discard)
	return l
}

func TestService_Create_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{
		CreateMock: func(task *task.Task) (*task.Task, error) {
			task.ID = 1
			return task, nil
		},
	}

	service := task.NewService(mockRepo, mockLogger())

	expId, err := service.Create(42, "test_title", "test_desc")
	if err != nil {
		t.Fatal(err)
	}

	if expId != 1 {
		t.Errorf("task id should be 1, was %d", expId)
	}
}

func TestService_Create_Fail(t *testing.T) {
	mockRepo := &MockTaskRepository{
		CreateMock: func(task *task.Task) (*task.Task, error) {
			task.ID = 1
			return nil, fmt.Errorf("test error")
		},
	}

	service := task.NewService(mockRepo, mockLogger())

	_, err := service.Create(42, "test_title", "test_desc")
	if err == nil {
		t.Fatal(err)
	}
}

func TestService_Update_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{
		UpdateMock: func(task *task.Task) error {
			return nil
		},
	}

	service := task.NewService(mockRepo, mockLogger())

	err := service.Update(42, 1, "test_title", "test_desc", true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_Update_Fail(t *testing.T) {
	mockRepo := &MockTaskRepository{
		UpdateMock: func(task *task.Task) error {
			return fmt.Errorf("test error")
		},
	}

	service := task.NewService(mockRepo, mockLogger())

	err := service.Update(42, 1, "test_title", "test_desc", true)
	if err == nil {
		t.Fatal(err)
	}
}

func TestService_Delete_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{
		DeleteByIdMock: func(task *task.Task) error {
			return nil
		},
	}

	service := task.NewService(mockRepo, mockLogger())

	err := service.Delete(42, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_Delete_Fail(t *testing.T) {
	mockRepo := &MockTaskRepository{
		DeleteByIdMock: func(task *task.Task) error {
			return fmt.Errorf("test error")
		},
	}

	service := task.NewService(mockRepo, mockLogger())

	err := service.Delete(42, 1)
	if err == nil {
		t.Fatal(err)
	}
}

func TestService_GetById_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{
		GetByIdMock: func(task *task.Task) (*task.Task, error) {
			task.ID = 1
			task.UserID = 42
			task.Title = "test_title"
			task.Description = "test_desc"
			task.Completed = true
			return task, nil
		},
	}

	service := task.NewService(mockRepo, mockLogger())

	exp, err := service.GetById(42, 1)
	if err != nil {
		t.Fatal(err)
	}

	if exp.ID != 1 ||
		exp.UserID != 42 ||
		exp.Title != "test_title" ||
		exp.Description != "test_desc" ||
		exp.Completed != true {
		t.Errorf("Unexpected result: %+v", exp)
	}
}

func TestService_GetById_Fail(t *testing.T) {
	mockRepo := &MockTaskRepository{
		GetByIdMock: func(task *task.Task) (*task.Task, error) {
			return nil, fmt.Errorf("test error")
		},
	}

	service := task.NewService(mockRepo, mockLogger())

	_, err := service.GetById(42, 1)
	if err == nil {
		t.Fatal(err)
	}
}

func TestService_GetAll_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{
		GetAllMock: func(user *user.User) ([]task.Task, error) {
			return []task.Task{
				{
					ID:          1,
					UserID:      42,
					Title:       "test_title_1",
					Description: "test_desc_1",
					Completed:   true,
				},
				{
					ID:          2,
					UserID:      42,
					Title:       "test_title_2",
					Description: "test_desc_2",
					Completed:   false,
				},
			}, nil
		},
	}

	service := task.NewService(mockRepo, mockLogger())

	exp, err := service.GetAll(42)
	if err != nil {
		t.Fatal(err)
	}

	if len(exp) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(exp))
	}

	if exp[0].ID != 1 ||
		exp[0].UserID != 42 ||
		exp[0].Title != "test_title_1" ||
		exp[0].Description != "test_desc_1" ||
		exp[0].Completed != true {
		t.Errorf("Unexpected first task: %+v", exp[0])
	}

	if exp[1].ID != 2 ||
		exp[1].UserID != 42 ||
		exp[1].Title != "test_title_2" ||
		exp[1].Description != "test_desc_2" ||
		exp[1].Completed != false {
		t.Errorf("Unexpected second task: %+v", exp[1])
	}
}

func TestService_GetAll_Fail(t *testing.T) {
	mockRepo := &MockTaskRepository{
		GetAllMock: func(user *user.User) ([]task.Task, error) {
			return nil, fmt.Errorf("test error")
		},
	}

	service := task.NewService(mockRepo, mockLogger())

	_, err := service.GetAll(42)
	if err == nil {
		t.Fatal(err)
	}
}
