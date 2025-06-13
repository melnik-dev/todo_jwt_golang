package auth

import (
	"fmt"
	"github.com/melnik-dev/go_todo_jwt/internal/user"
	"github.com/melnik-dev/go_todo_jwt/pkg/crypto"
)

type Service struct {
	UserRepo *user.Repository
}

func NewService(repo *user.Repository) *Service {
	return &Service{
		UserRepo: repo,
	}
}

func (service *Service) Register(username, password string) (int, error) {
	existedUser, _ := service.UserRepo.Get(username)
	if existedUser != nil {
		return 0, ErrUserExists
	}

	hashPassword, err := crypto.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("service: failed to hash password: %w", err)
	}

	u := &user.User{
		Name:     username,
		Password: hashPassword,
	}
	_, err = service.UserRepo.Create(u)
	if err != nil {
		return 0, err
	}

	return u.ID, nil
}

func (service *Service) Login(username, password string) (int, error) {
	existedUser, _ := service.UserRepo.Get(username)
	if existedUser == nil {
		return 0, ErrInvalidLogin
	}

	if !crypto.ComparePasswords(password, existedUser.Password) {
		return 0, ErrInvalidLogin
	}

	return existedUser.ID, nil
}
