package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/melnik-dev/go_todo_jwt/internal/model"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (au *AuthRepo) CreateUser(username, password string) (int, error) {
	var id int
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`

	row := au.db.QueryRow(query, username, password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (au *AuthRepo) GetUser(username string) (model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE username = $1`

	err := au.db.Get(&user, query, username)
	return user, err
}
