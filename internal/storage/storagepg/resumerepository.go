package storagepg

import "guthub.com/serge64/joffer/internal/model"

type ResumeRepository struct {
	store *Store
}

func (r *ResumeRepository) Create(resume *model.Resume) error {
	if _, err := r.store.db.Exec(
		"INSERT INTO resumes (profile_id, name, uid) VALUES ($1, $2, $3);",
		resume.ProfileID,
		resume.Name,
		resume.UID,
	); err != nil {
		return err
	}
	return nil
}

func (r *ResumeRepository) Find(profileID int) ([]model.Resume, error) {
	resumes := []model.Resume{}

	if err := r.store.db.Select(
		&resumes,
		"SELECT * FROM resumes WHERE profile_id = $1;",
		profileID,
	); err != nil {
		return nil, err
	}

	return resumes, nil
}

func (r *ResumeRepository) Delete(id int) error {
	if _, err := r.store.db.Exec(
		"DELETE FROM resumes WHERE id = $1;",
		id,
	); err != nil {
		return err
	}

	return nil
}
