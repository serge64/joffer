package storage

import "guthub.com/serge64/joffer/internal/model"

type UserRepository interface {
	Create(*model.User) (int, error)
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}
