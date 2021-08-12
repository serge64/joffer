package storage

import "guthub.com/serge64/joffer/internal/model"

type GroupRepository interface {
	Create(*model.Group) (int, error)
	Find(int) ([]model.Group, error)
	Update(*model.Group) error
	Delete(int) error
}
