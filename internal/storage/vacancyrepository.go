package storage

import "guthub.com/serge64/joffer/internal/model"

type VacancyRepository interface {
	Find(*model.Filter) ([]model.Vacancy, error)
	CountSelected(int) (int, error)
	Response(int) error
	Toggle(int) error
}
