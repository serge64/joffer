package model

type Group struct {
	ID        int      `json:"id" db:"id"`
	UserID    int      `json:"-" db:"user_id"`
	Name      string   `json:"name" db:"name"`
	Resume    string   `json:"resume" db:"resume"`
	Letter    string   `json:"letter" db:"letter"`
	Positions []string `json:"positions" db:"-"`
	Count     int      `json:"count" db:"-"`
}
