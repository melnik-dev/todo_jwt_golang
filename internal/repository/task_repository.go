package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/melnik-dev/go_todo_jwt/internal/model"
)

type TaskRepo struct {
	db *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (tr *TaskRepo) Create(userID int, title, desc string) (int, error) {
	var id int
	query := `INSERT INTO tasks (user_id, title, description)
				VALUES ($1, $2, $3) 
				RETURNING id`

	row := tr.db.QueryRow(query, userID, title, desc)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (tr *TaskRepo) Update(userID, taskId int, title, desc string, completed bool) error {
	query := `UPDATE tasks 
				SET title = $1, description = $2, completed = $3 
				WHERE id = $4 AND user_id = $5`

	result, err := tr.db.Exec(query, title, desc, completed, taskId, userID)
	if err != nil {
		return err
	}

	row, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return fmt.Errorf("repository: %w", ErrTaskNotFound)
	}

	return nil
}

func (tr *TaskRepo) DeleteById(userID, taskId int) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`

	result, err := tr.db.Exec(query, taskId, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("repository: %w", ErrTaskNotFound)
	}

	return nil
}

func (tr *TaskRepo) GetById(userID, taskID int) (model.Task, error) {
	var task model.Task
	query := `SELECT * FROM tasks WHERE id = $1 AND user_id = $2`

	err := tr.db.Get(&task, query, taskID, userID)
	return task, err
}

func (tr *TaskRepo) GetAll(userID int) ([]model.Task, error) {
	tasks := make([]model.Task, 0)
	query := `SELECT * FROM tasks WHERE user_id = $1`

	err := tr.db.Select(&tasks, query, userID)
	return tasks, err
}
