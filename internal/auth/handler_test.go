package auth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/melnik-dev/go_todo_jwt/configs"
	"github.com/melnik-dev/go_todo_jwt/internal/auth"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockAuthService struct {
	RegisterMock func(username, password string) (int, error)
	LoginMock    func(username, password string) (int, error)
}

func (m *MockAuthService) Register(username, password string) (int, error) {
	return m.RegisterMock(username, password)
}

func (m *MockAuthService) Login(username, password string) (int, error) {
	return m.LoginMock(username, password)
}

func mockGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		logger := logrus.New()
		logger.SetOutput(io.Discard)
		c.Set("logger", logrus.NewEntry(logger))
		c.Next()
	})
	return r
}

type Options struct {
	h    *auth.Handler
	name string
}

func requestRegisterHelper(t *testing.T, opts Options) *httptest.ResponseRecorder {
	t.Helper()
	r := mockGin()
	r.POST("/auth/register", opts.h.Register)

	name := opts.name
	if name == "" {
		name = "testname"
	}

	body, _ := json.Marshal(map[string]string{
		"username": name,
		"password": "test_pass",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func requestLoginHelper(t *testing.T, opts Options) *httptest.ResponseRecorder {
	t.Helper()
	r := mockGin()
	r.POST("/auth/login", opts.h.Login)

	name := opts.name
	if name == "" {
		name = "testname"
	}

	body, _ := json.Marshal(map[string]string{
		"username": name,
		"password": "test_pass",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestHandler_Register_Success(t *testing.T) {
	handler := &auth.Handler{
		AuthService: &MockAuthService{
			RegisterMock: func(username, password string) (int, error) {
				return 42, nil
			},
		},
		Config: &configs.Config{
			JWT: configs.ConfJWT{Secret: "secret", TokenTTL: time.Hour},
		},
	}

	w := requestRegisterHelper(t, Options{h: handler})

	if w.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
	}

	type ResponseWrapper struct {
		Status int                   `json:"status"`
		Data   auth.RegisterResponse `json:"data"`
	}

	var res ResponseWrapper
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Fatal(err)
	}

	if res.Data.Token == "" {
		t.Errorf("expected non-empty token, got: '%s'", res.Data.Token)
	}
}

func TestHandler_Register_Fail(t *testing.T) {
	handler := &auth.Handler{
		AuthService: &MockAuthService{
			RegisterMock: func(username, password string) (int, error) {
				return 0, fmt.Errorf("test error")
			},
		},
		Config: &configs.Config{
			JWT: configs.ConfJWT{Secret: "secret", TokenTTL: time.Hour},
		},
	}

	w := requestRegisterHelper(t, Options{h: handler})

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandler_Register_FailExists(t *testing.T) {
	handler := &auth.Handler{
		AuthService: &MockAuthService{
			RegisterMock: func(username, password string) (int, error) {
				return 0, auth.ErrUserExists
			},
		},
		Config: &configs.Config{
			JWT: configs.ConfJWT{Secret: "secret", TokenTTL: time.Hour},
		},
	}

	w := requestRegisterHelper(t, Options{h: handler})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_Register_FailInvalid(t *testing.T) {
	handler := &auth.Handler{
		AuthService: &MockAuthService{
			RegisterMock: func(username, password string) (int, error) {
				return 0, fmt.Errorf("invalid input data")
			},
		},
		Config: &configs.Config{
			JWT: configs.ConfJWT{Secret: "secret", TokenTTL: time.Hour},
		},
	}

	w := requestRegisterHelper(t, Options{h: handler, name: "test_name"})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandler_Login_Success(t *testing.T) {
	handler := &auth.Handler{
		AuthService: &MockAuthService{
			LoginMock: func(username, password string) (int, error) {
				return 42, nil
			},
		},
		Config: &configs.Config{
			JWT: configs.ConfJWT{Secret: "secret", TokenTTL: time.Hour},
		},
	}

	w := requestLoginHelper(t, Options{h: handler})

	if w.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
	}

	type ResponseWrapper struct {
		Status int                `json:"status"`
		Data   auth.LoginResponse `json:"data"`
	}

	var res ResponseWrapper
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Fatal(err)
	}

	if res.Data.Token == "" {
		t.Errorf("expected non-empty token, got: '%s'", res.Data.Token)
	}
}

func TestHandler_Login_Fail(t *testing.T) {
	handler := &auth.Handler{
		AuthService: &MockAuthService{
			LoginMock: func(username, password string) (int, error) {
				return 0, fmt.Errorf("test error")
			},
		},
		Config: &configs.Config{
			JWT: configs.ConfJWT{Secret: "secret", TokenTTL: time.Hour},
		},
	}

	w := requestLoginHelper(t, Options{h: handler})

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandler_Login_FailUnauthorized(t *testing.T) {
	handler := &auth.Handler{
		AuthService: &MockAuthService{
			LoginMock: func(username, password string) (int, error) {
				return 0, auth.ErrInvalidLogin
			},
		},
		Config: &configs.Config{
			JWT: configs.ConfJWT{Secret: "secret", TokenTTL: time.Hour},
		},
	}

	w := requestLoginHelper(t, Options{h: handler})

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestHandler_Login_FailInvalid(t *testing.T) {
	handler := &auth.Handler{
		AuthService: &MockAuthService{
			LoginMock: func(username, password string) (int, error) {
				return 0, fmt.Errorf("invalid input data")
			},
		},
		Config: &configs.Config{
			JWT: configs.ConfJWT{Secret: "secret", TokenTTL: time.Hour},
		},
	}

	w := requestLoginHelper(t, Options{h: handler, name: "test_name"})

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}
