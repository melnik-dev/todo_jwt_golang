package user

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) Create(user *User) (*User, error) {
	var id int
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`

	row := repo.db.QueryRow(query, user.Name, user.Password)
	if err := row.Scan(&id); err != nil {
		return user, err
	}
	return user, nil
}

func (repo *Repository) Get(username string) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE username = $1`

	err := repo.db.Get(&user, query, username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
