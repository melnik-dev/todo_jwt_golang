package task

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/melnik-dev/go_todo_jwt/internal/user"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) Create(task *Task) (*Task, error) {
	var id int
	query := `INSERT INTO tasks (user_id, title, description)
				VALUES ($1, $2, $3) 
				RETURNING id`

	row := repo.db.QueryRow(query, task.UserID, task.Title, task.Description)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	task.ID = id
	return task, nil
}

func (repo *Repository) Update(task *Task) error {
	query := `UPDATE tasks 
				SET title = $1, description = $2, completed = $3 
				WHERE id = $4 AND user_id = $5`

	result, err := repo.db.Exec(query, task.Title, task.Description, task.Completed, task.ID, task.UserID)
	if err != nil {
		return err
	}

	row, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (repo *Repository) DeleteById(task *Task) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`

	result, err := repo.db.Exec(query, task.ID, task.UserID)
	if err != nil {
		return err
	}

	row, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (repo *Repository) GetById(task *Task) (*Task, error) {
	query := `SELECT * FROM tasks WHERE id = $1 AND user_id = $2`

	err := repo.db.Get(task, query, task.ID, task.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return task, nil
}

func (repo *Repository) GetAll(user *user.User) ([]Task, error) {
	tasks := make([]Task, 0)
	query := `SELECT * FROM tasks WHERE user_id = $1`

	err := repo.db.Select(&tasks, query, user.ID)
	if err != nil {
		return nil, err
	}
	return tasks, err
}
