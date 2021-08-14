package storage

import "guthub.com/serge64/joffer/internal/model"

type GroupRepository interface {
	Create(*model.Group) (int, error)
	Find(int) ([]model.Group, error)
	FindList(int) ([]string, error)
	Update(*model.Group) error
	Delete(int) error
	Response(int, int) error
}
