package storage

import "guthub.com/serge64/joffer/internal/model"

type ResumeRepository interface {
	Create(*model.Resume) error
	Find(int) ([]model.Resume, error)
	Delete(int) error
}
