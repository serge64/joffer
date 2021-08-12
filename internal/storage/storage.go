package storage

type Store interface {
	User() UserRepository
	Group() GroupRepository
	Letter() LetterRepository
	Vacancy() VacancyRepository
	Task() TaskRepository
}
