package model

type Resume struct {
	ID        int    `db:"id"`
	ProfileID int    `db:"profile_id"`
	Name      string `db:"name"`
	UID       string `db:"uid"`
}
