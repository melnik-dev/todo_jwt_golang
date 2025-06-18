package task

import (
	"database/sql"
	"errors"
	"github.com/melnik-dev/go_todo_jwt/internal/user"
	"github.com/melnik-dev/go_todo_jwt/pkg/db"
	"github.com/melnik-dev/go_todo_jwt/pkg/logger"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	db *db.Db
}

func NewRepository(db *db.Db) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(task *Task) (*Task, error) {
	logRepo := repositoryLogger().WithFields(logrus.Fields{
		"user_id": task.UserID,
		"title":   task.Title,
	})
	logRepo.Debug("Attempting to Create")

	var id int
	query := `INSERT INTO tasks (user_id, title, description)
				VALUES ($1, $2, $3) 
				RETURNING id`

	row := r.db.QueryRow(query, task.UserID, task.Title, task.Description)
	if err := row.Scan(&id); err != nil {
		logRepo.WithError(err).Error("Failed to insert database")
		return nil, err
	}
	task.ID = id

	logRepo.Debug("Insert database successfully")
	return task, nil
}

func (r *Repository) Update(task *Task) error {
	logRepo := repositoryLogger().WithFields(logrus.Fields{
		"user_id": task.UserID,
		"task_id": task.ID,
	})
	logRepo.Debug("Attempting to Update")

	query := `UPDATE tasks 
				SET title = $1, description = $2, completed = $3 
				WHERE id = $4 AND user_id = $5`

	result, err := r.db.Exec(query, task.Title, task.Description, task.Completed, task.ID, task.UserID)
	if err != nil {
		logRepo.WithError(err).Error("Failed to Update database")
		return err
	}

	row, err := result.RowsAffected()
	if err != nil {
		logRepo.WithError(err).Error("Failed rows affected by Update database")
		return err
	}

	if row == 0 {
		logRepo.WithError(err).Warn(ErrTaskNotFound.Error())
		return ErrTaskNotFound
	}

	logRepo.Debug("Update database successfully")
	return nil
}

func (r *Repository) DeleteById(task *Task) error {
	logRepo := repositoryLogger().WithFields(logrus.Fields{
		"user_id": task.UserID,
		"task_id": task.ID,
	})
	logRepo.Debug("Attempting to Delete")

	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, task.ID, task.UserID)
	if err != nil {
		logRepo.WithError(err).Error("Failed to Delete database")
		return err
	}

	row, err := result.RowsAffected()
	if err != nil {
		logRepo.WithError(err).Error("Failed rows affected by Delete database")
		return err
	}

	if row == 0 {
		logRepo.WithError(err).Warn(ErrTaskNotFound.Error())
		return ErrTaskNotFound
	}

	logRepo.Debug("Delete database successfully")
	return nil
}

func (r *Repository) GetById(task *Task) (*Task, error) {
	logRepo := repositoryLogger().WithFields(logrus.Fields{
		"user_id": task.UserID,
		"task_id": task.ID,
	})
	logRepo.Debug("Attempting to GetById")

	query := `SELECT * FROM tasks WHERE id = $1 AND user_id = $2`

	err := r.db.Get(task, query, task.ID, task.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logRepo.WithError(err).Warn(ErrTaskNotFound.Error())
			return nil, ErrTaskNotFound
		}
		logRepo.WithError(err).Error("Failed to GetById database")
		return nil, err
	}

	logRepo.Debug("GetById database successfully")
	return task, nil
}

func (r *Repository) GetAll(user *user.User) ([]Task, error) {
	logRepo := repositoryLogger().WithFields(logrus.Fields{
		"user_id": user.ID,
	})
	logRepo.Debug("Attempting to GetAll")

	tasks := make([]Task, 0)
	query := `SELECT * FROM tasks WHERE user_id = $1`

	err := r.db.Select(&tasks, query, user.ID)
	if err != nil {
		logRepo.WithError(err).Error("Failed to GetAll database")
		return nil, err
	}

	logRepo.Debug("GetAll database successfully")
	return tasks, err
}

func repositoryLogger() *logrus.Entry {
	return logger.GetLogger().WithField("layer", "Repository task layer")
}
