package model

type User struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"username" json:"username" binding:"required,min=3"`
	Password string `db:"password" json:"password" binding:"required,min=3"`
}
