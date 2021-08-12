package storagepg

import "guthub.com/serge64/joffer/internal/model"

type LetterRepository struct {
	store *Store
}

func (r *LetterRepository) Create(letter *model.Letter) (int, error) {
	if err := r.store.db.QueryRow(
		"INSERT INTO letters (user_id, name, body) VALUES ($1, $2, $3) RETURNING id",
		letter.UserID,
		letter.Name,
		letter.Body,
	).Scan(&letter.ID); err != nil {
		return 0, err
	}
	return letter.ID, nil
}

func (r *LetterRepository) Update(letter *model.Letter) error {
	if _, err := r.store.db.NamedExec(
		"UPDATE letters SET name = :name, body = :body WHERE id = :id;",
		letter,
	); err != nil {
		return err
	}
	return nil
}

func (r *LetterRepository) Find(userID int) ([]model.Letter, error) {
	letters := []model.Letter{}
	if err := r.store.db.Select(
		&letters,
		"SELECT * FROM letters WHERE user_id = $1",
		userID,
	); err != nil {
		return nil, err
	}
	return letters, nil
}

func (r *LetterRepository) Delete(userID int, letterID int) error {
	if _, err := r.store.db.Exec(
		"DELETE FROM letters WHERE id = $1 AND user_id = $2;",
		letterID,
		userID,
	); err != nil {
		return err
	}
	return nil
}
