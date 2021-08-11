package sessionstorage

import (
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/boj/redistore.v1"
	"guthub.com/serge64/joffer/internal/config"
)

type Session struct {
	Redis *redistore.RediStore
}

func New(config *config.Config) (*Session, error) {
	redis, err := redistore.NewRediStore(
		10,
		"tcp",
		config.SessionAddr,
		"",
		[]byte(config.SessionKey),
	)
	if err != nil {
		return nil, err
	}

	redis.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		Secure:   true,
		MaxAge:   config.SessionMaxAge,
		HttpOnly: true,
	}

	session := &Session{
		Redis: redis,
	}

	logrus.Println("open redis connection")

	return session, nil
}

func (s *Session) Close() {
	if err := s.Redis.Close(); err != nil {
		logrus.Println(errors.Wrap(err, "err closing redis connection"))
	} else {
		logrus.Println("close redis connection")
	}
}
