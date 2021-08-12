package storage

import "guthub.com/serge64/joffer/internal/model"

type LetterRepository interface {
	Create(*model.Letter) (int, error)
	Update(*model.Letter) error
	Find(int) ([]model.Letter, error)
	Delete(int, int) error
}
