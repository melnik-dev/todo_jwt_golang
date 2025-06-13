package task

type Task struct {
	ID          int    `db:"id" json:"id"`
	UserID      int    `db:"user_id" json:"user_id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Completed   bool   `db:"completed" json:"completed"`
}
