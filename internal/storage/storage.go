package storage

type Store interface {
	User() UserRepository
	Profile() ProfileRepository
	Resume() ResumeRepository
	Group() GroupRepository
	Letter() LetterRepository
	Vacancy() VacancyRepository
	Task() TaskRepository
	Platform() PlatformRepository
	Search() SearchRepository
}
