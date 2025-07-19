package auth_test

import (
	"errors"
	"fmt"
	"github.com/melnik-dev/go_todo_jwt/internal/auth"
	"github.com/melnik-dev/go_todo_jwt/internal/user"
	"github.com/melnik-dev/go_todo_jwt/pkg/crypto"
	"github.com/sirupsen/logrus"
	"io"
	"testing"
)

type MockUserRepository struct {
	GetMock    func(username string) (*user.User, error)
	CreateMock func(user *user.User) (*user.User, error)
}

func (m *MockUserRepository) Get(username string) (*user.User, error) {
	return m.GetMock(username)
}

func (m *MockUserRepository) Create(user *user.User) (*user.User, error) {
	return m.CreateMock(user)
}

func mockLogger() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(io.Discard)
	return l
}

func TestService_Register_Success(t *testing.T) {
	mockRepo := &MockUserRepository{
		CreateMock: func(u *user.User) (*user.User, error) {
			u.ID = 42
			return u, nil
		},
		GetMock: func(username string) (*user.User, error) {
			return nil, nil
		},
	}

	service := auth.NewService(mockRepo, mockLogger())

	expId, err := service.Register("test_user", "test_pass")
	if err != nil {
		t.Fatal(err)
	}

	if expId != 42 {
		t.Fatalf("expected ID 42, got %d", expId)
	}
}

func TestService_Register_Fail(t *testing.T) {
	mockRepo := &MockUserRepository{
		CreateMock: func(u *user.User) (*user.User, error) {
			return nil, fmt.Errorf("failed to create")
		},
		GetMock: func(username string) (*user.User, error) {
			return nil, nil
		},
	}

	service := auth.NewService(mockRepo, mockLogger())

	_, err := service.Register("test_user", "test_pass")
	if err == nil {
		t.Fatal(err)
	}
}

func TestService_Register_FailExist(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetMock: func(username string) (*user.User, error) {
			return &user.User{ID: 42}, nil
		},
	}

	service := auth.NewService(mockRepo, mockLogger())

	_, err := service.Register("test_user", "test_pass")
	if !errors.Is(err, auth.ErrUserExists) {
		t.Fatalf("expected ErrUserExists, got %v", err)
	}
}

func TestService_Login_Success(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetMock: func(username string) (*user.User, error) {
			pass, _ := crypto.HashPassword("test_pass")
			return &user.User{
				ID:       42,
				Password: pass,
			}, nil
		},
	}

	service := auth.NewService(mockRepo, mockLogger())

	expId, err := service.Login("test_user", "test_pass")
	if err != nil {
		t.Fatal(err)
	}

	if expId != 42 {
		t.Fatalf("expected ID 42, got %d", expId)
	}
}

func TestService_Login_Fail(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetMock: func(username string) (*user.User, error) {
			return nil, auth.ErrInvalidLogin
		},
	}

	service := auth.NewService(mockRepo, mockLogger())

	_, err := service.Login("test_user", "test_pass")
	if err == nil {
		t.Fatal(err)
	}
}
