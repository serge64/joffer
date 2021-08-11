package storagepg

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sqlx.DB
}

func New(databaseURL string) (*Store, error) {
	db, err := createCient(databaseURL)
	if err != nil {
		return nil, err
	}

	s := &Store{
		db: db,
	}

	return s, nil
}

func createCient(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	logrus.Println("open db connection")

	return db, nil
}

func (s *Store) Close() {
	if err := s.db.Close(); err != nil {
		logrus.Println(errors.Wrap(err, "err closing db connection"))
	} else {
		logrus.Println("close db connection")
	}
}
