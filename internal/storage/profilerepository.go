package storage

import (
	"guthub.com/serge64/joffer/internal/config"
	"guthub.com/serge64/joffer/internal/model"
)

type ProfileRepository interface {
	Create(*model.Profile, string, *config.Config) error
	Find(int, int) (*model.Profile, error)
	Delete(int, int) error
}
