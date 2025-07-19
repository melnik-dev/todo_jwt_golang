package auth

import (
	"fmt"
	"github.com/melnik-dev/go_todo_jwt/internal/user"
	"github.com/melnik-dev/go_todo_jwt/pkg/crypto"
	"github.com/sirupsen/logrus"
)

type IService interface {
	Register(username, password string) (int, error)
	Login(username, password string) (int, error)
}

type Service struct {
	userRepo user.IRepository
	logger   *logrus.Logger
}

func NewService(repo user.IRepository, logger *logrus.Logger) *Service {
	return &Service{
		userRepo: repo,
		logger:   logger,
	}
}

func (s *Service) Register(username, password string) (int, error) {
	logServ := serviceLogger(s.logger).WithField("user_name", username)
	logServ.Debug("Attempting to Register new user")

	existedUser, _ := s.userRepo.Get(username)
	if existedUser != nil {
		logServ.Warn(ErrUserExists.Error())
		return 0, ErrUserExists
	}

	hashPassword, err := crypto.HashPassword(password)
	if err != nil {
		logServ.WithError(err).Error("failed to hash password")
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	u := &user.User{
		Name:     username,
		Password: hashPassword,
	}
	_, err = s.userRepo.Create(u)
	if err != nil {
		logServ.WithError(err).Error("failed to Create")
		return 0, err
	}

	logServ.Debug("User Register successfully")
	return u.ID, nil
}

func (s *Service) Login(username, password string) (int, error) {
	logServ := serviceLogger(s.logger).WithField("user_name", username)
	logServ.Debug("Attempting to Login new user")

	existedUser, err := s.userRepo.Get(username)
	if err != nil {
		logServ.WithError(err).Warn("failed to fetch user")
		return 0, ErrInvalidLogin
	}

	if !crypto.ComparePasswords(password, existedUser.Password) {
		logServ.Warn(ErrInvalidLogin.Error())
		return 0, ErrInvalidLogin
	}

	logServ.Debug("User Login successfully")
	return existedUser.ID, nil
}

func serviceLogger(l *logrus.Logger) *logrus.Entry {
	return l.WithField("layer", "Service auth layer")
}
