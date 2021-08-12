package storage

import (
	"github.com/jmoiron/sqlx"
	"guthub.com/serge64/joffer/internal/model"
)

type TaskRepository interface {
	Create(*model.Task, *sqlx.Tx) error
	Find(int) ([]model.Task, error)
	Delete(int, string) error
}
