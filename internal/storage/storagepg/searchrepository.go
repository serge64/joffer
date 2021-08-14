package storagepg

import (
	"fmt"
)

type SearchRepository struct {
	store *Store
}

func (s *SearchRepository) Position(userID int, pattern string) ([]string, error) {
	query := fmt.Sprintf(
		"SELECT DISTINCT LOWER(v.name) FROM vacancies v INNER JOIN tasks t ON t.id = v.task_id INNER JOIN groups g ON g.id = t.group_id INNER JOIN profiles p ON p.id = g.profile_id WHERE LOWER(v.name) ~ '.*%s.*' AND p.user_id = %d ORDER BY LOWER(v.name) ASC LIMIT 20;",
		pattern,
		userID,
	)

	return s.search(query)
}

func (s *SearchRepository) Company(userID int, pattern string) ([]string, error) {
	query := fmt.Sprintf(
		"SELECT DISTINCT LOWER(company) FROM vacancies v INNER JOIN tasks t ON t.id = v.task_id INNER JOIN groups g ON g.id = t.group_id INNER JOIN profiles p ON p.id = g.profile_id WHERE LOWER(company) ~ '.*%s.*' AND p.user_id = %d ORDER BY LOWER(company) ASC LIMIT 20;",
		pattern,
		userID,
	)

	return s.search(query)
}

func (s *SearchRepository) Area(userID int, pattern string) ([]string, error) {
	query := fmt.Sprintf(
		"SELECT DISTINCT LOWER(area) FROM vacancies v INNER JOIN tasks t ON t.id = v.task_id INNER JOIN groups g ON g.id = t.group_id INNER JOIN profiles p ON p.id = g.profile_id WHERE LOWER(area) ~ '.*%s.*' AND p.user_id = %d ORDER BY LOWER(area) ASC LIMIT 20;",
		pattern,
		userID,
	)

	return s.search(query)
}

func (s *SearchRepository) search(query string) ([]string, error) {
	found := []string{}

	if err := s.store.db.Select(&found, query); err != nil {
		return nil, err
	}

	return found, nil
}
