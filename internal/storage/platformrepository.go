package storage

type PlatformRepository interface {
	Find() ([]string, error)
}
