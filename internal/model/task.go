package model

type Task struct {
	ID          int    `db:"id" json:"id"`
	UserID      string `db:"user_id" json:"user_id"`
	Title       string `db:"title" json:"title" binding:"required"`
	Description string `db:"description" json:"description"`
	Completed   bool   `db:"completed" json:"completed"`
}
