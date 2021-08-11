package storagepg

import (
	"guthub.com/serge64/joffer/internal/model"
	"guthub.com/serge64/joffer/internal/storage"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) (int, error) {
	if err := u.Validate(); err != nil {
		return 0, storage.ErrUserNotValid
	}

	u.BeforeCreate()

	if err := r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id;",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID); err != nil {
		return 0, storage.ErrEmailNotUnique
	}

	return u.ID, nil
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	u := model.User{}
	if err := r.store.db.Get(
		&u,
		"SELECT * FROM users WHERE id = $1;",
		id,
	); err != nil {
		return nil, storage.ErrRecordNotFound
	}

	return &u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := model.User{}
	if err := r.store.db.Get(
		&u,
		"SELECT * FROM users WHERE email = $1;",
		email,
	); err != nil {
		return nil, storage.ErrRecordNotFound
	}

	return &u, nil
}
