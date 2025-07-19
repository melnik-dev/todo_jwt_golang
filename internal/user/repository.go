package user

import (
	"github.com/melnik-dev/go_todo_jwt/pkg/db"
	"github.com/sirupsen/logrus"
)

type IRepository interface {
	Create(user *User) (*User, error)
	Get(username string) (*User, error)
}

type Repository struct {
	db     *db.Db
	logger *logrus.Logger
}

func NewRepository(db *db.Db, logger *logrus.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) Create(user *User) (*User, error) {
	userLogger := repositoryLogger(r.logger)
	userLogger.Debug("Attempting to Create user")

	var id int
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`

	row := r.db.QueryRow(query, user.Name, user.Password)
	if err := row.Scan(&id); err != nil {
		userLogger.WithError(err).Error("Failed to create user in database")
		return user, err
	}
	user.ID = id

	userLogger.WithField("user_id", user.ID).Debug("User Create successfully")
	return user, nil
}

func (r *Repository) Get(username string) (*User, error) {
	userLogger := repositoryLogger(r.logger)
	userLogger.Debug("Attempting to Get user")

	var user User
	query := `SELECT * FROM users WHERE username = $1`

	err := r.db.Get(&user, query, username)
	if err != nil {
		userLogger.WithError(err).Error("Failed to get user in database")
		return nil, err
	}

	userLogger.WithField("user_id", user.ID).Debug("User Get successfully")
	return &user, nil
}

func repositoryLogger(l *logrus.Logger) *logrus.Entry {
	return l.WithField("layer", "Repository user layer")
}
