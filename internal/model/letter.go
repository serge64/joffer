package model

type Letter struct {
	ID     int    `json:"id" db:"id"`
	UserID int    `json:"-" db:"user_id"`
	Name   string `json:"name" db:"name"`
	Body   string `json:"body" db:"body"`
}
