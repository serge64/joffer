package storagepg

type PlatformRepository struct {
	store *Store
}

func (r *PlatformRepository) Find() ([]string, error) {
	platforms := []string{}

	if err := r.store.db.Select(
		&platforms,
		"SELECT full_name FROM platforms;",
	); err != nil {
		return nil, err
	}

	return platforms, nil
}
