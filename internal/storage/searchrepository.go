package storage

type SearchRepository interface {
	Position(int, string) ([]string, error)
	Company(int, string) ([]string, error)
	Area(int, string) ([]string, error)
}
