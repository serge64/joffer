package storage

type Store interface {
	User() UserRepository
	Group() GroupRepository
}
