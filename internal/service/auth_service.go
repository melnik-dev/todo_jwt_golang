package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/melnik-dev/go_todo_jwt/internal/model"
	"github.com/melnik-dev/go_todo_jwt/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	repo      repository.Auth
	jwtSecret string
	tokenTTL  time.Duration
}

func NewAuthService(authRepo repository.Auth, jwtSecret string, time time.Duration) *AuthService {
	return &AuthService{
		repo:      authRepo,
		jwtSecret: jwtSecret,
		tokenTTL:  time,
	}
}

func (as *AuthService) Register(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", fmt.Errorf("service: %w", ErrRequiredField)
	}

	hashPassword, err := as.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("service: failed to hash password: %w", err)
	}

	userID, err := as.repo.CreateUser(username, hashPassword)
	if err != nil {
		return "", err
	}

	token, err := as.GenerateToken(userID)
	if err != nil {
		return "", fmt.Errorf("service: %w", err)
	}

	return token, nil
}

func (as *AuthService) Login(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", fmt.Errorf("service: %w", ErrRequiredField)
	}

	var user model.User
	user, err := as.repo.GetUser(username)
	if err != nil {
		return "", err
	}

	if !as.ComparePasswords(password, user.Password) {
		return "", fmt.Errorf("service: invalid login or password")
	}

	token, err := as.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("service: %w", err)
	}

	return token, nil
}

func (as *AuthService) GenerateToken(userId int) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(as.tokenTTL).Unix(),
	})

	return claims.SignedString([]byte(as.jwtSecret))
}

func (as *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (as *AuthService) ComparePasswords(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
