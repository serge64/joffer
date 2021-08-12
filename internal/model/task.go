package model

type Task struct {
	ID      int    `db:"id"`
	GroupID int    `db:"group_id"`
	Name    string `db:"name"`
}
