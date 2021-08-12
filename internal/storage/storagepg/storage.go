package storagepg

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"guthub.com/serge64/joffer/internal/storage"

	_ "github.com/lib/pq"
)

type Store struct {
	db                *sqlx.DB
	userRepository    *UserRepository
	letterRepository  *LetterRepository
	groupRepository   *GroupRepository
	vacancyRepository *VacancyRepository
	taskRepository    *TaskRepository
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

func (s *Store) User() storage.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
		}
	}
	return s.userRepository
}

func (s *Store) Letter() storage.LetterRepository {
	if s.letterRepository == nil {
		s.letterRepository = &LetterRepository{
			store: s,
		}
	}
	return s.letterRepository
}

func (s *Store) Group() storage.GroupRepository {
	if s.groupRepository == nil {
		s.groupRepository = &GroupRepository{
			store: s,
		}
	}
	return s.groupRepository
}

func (s *Store) Vacancy() storage.VacancyRepository {
	if s.vacancyRepository == nil {
		s.vacancyRepository = &VacancyRepository{
			store: s,
		}
	}
	return s.vacancyRepository
}

func (s *Store) Task() storage.TaskRepository {
	if s.taskRepository == nil {
		s.taskRepository = &TaskRepository{
			store: s,
			cron:  gocron.NewScheduler(time.UTC),
			tasks: make(map[int]*gocron.Job),
		}

		s.taskRepository.cron.SetMaxConcurrentJobs(1, 1)
		s.taskRepository.cron.StartAsync()
	}
	return s.taskRepository
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
