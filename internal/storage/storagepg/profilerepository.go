package storagepg

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"guthub.com/serge64/joffer/internal/config"
	"guthub.com/serge64/joffer/internal/model"
)

type ProfileRepository struct {
	store *Store
	cron  *gocron.Scheduler
	tasks map[int]*gocron.Job
}

func (r *ProfileRepository) Create(profile *model.Profile, code string, config *config.Config) error {
	if err := profile.Authorization(code, config); err != nil {
		return err
	}

	if _, err := profile.Me(); err != nil {
		return err
	}

	if err := r.store.db.QueryRow(
		"INSERT INTO profiles (user_id, platform_id, name, email, access_token, refresh_token, expiry) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;",
		profile.UserID,
		profile.PlatformID,
		profile.Name,
		profile.Email,
		profile.AccessToken,
		profile.RefreshToken,
		profile.Expiry,
	).Scan(&profile.ID); err != nil {
		return err
	}

	resumes, err := profile.Resumes()
	if err != nil {
		return nil
	}

	for _, v := range resumes {
		if err := r.store.Resume().Create(&v); err != nil {
			return err
		}
	}

	t := profile.Expiry
	j, err := r.cron.Every(t).Seconds().StartAt(time.Now().Add(time.Duration(t*int(time.Second)))).Do(r.update, r.store.db, profile)
	if err != nil {
		return err
	}

	r.tasks[profile.ID] = j
	logrus.WithFields(logrus.Fields{
		"profile_id": profile.ID,
	}).Warnf("profile created")

	return nil
}

func (r *ProfileRepository) update(db *sqlx.DB, p model.Profile) {
	logger := logrus.WithFields(logrus.Fields{
		"profile_id": p.ID,
	})

	logger.Warnf("profile start updated")

	if err := p.UpdateToken(); err != nil {
		logrus.Error(err)
		return
	}

	if _, err := db.NamedExec(
		"UPDATE profiles SET access_token = :access_token, refresh_token = :refresh_token, expiry = :expiry WHERE id = :id;",
		p,
	); err != nil {
		logrus.Error(err)
		return
	}

	logrus.Warnf("profile end updated")
}

func (r *ProfileRepository) Find(userID int, platformID int) (*model.Profile, error) {
	profile := model.Profile{}

	if err := r.store.db.Get(
		&profile,
		"SELECT * FROM profiles WHERE user_id = $1 AND platform_id = $2;",
		userID,
		platformID,
	); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *ProfileRepository) Delete(userID int, platformID int) error {
	ids := []int{}

	if err := r.store.db.Select(
		&ids,
		"SELECT g.id FROM groups g INNER JOIN profiles p ON p.id = g.profile_id WHERE p.user_id = $1",
		userID,
	); err != nil {
		return err
	}

	if _, err := r.store.db.Exec(
		"DELETE FROM profiles WHERE user_id = $1 AND platform_id = $2;",
		userID,
		platformID,
	); err != nil {
		return err
	}

	for _, v := range ids {
		r.store.Group().Delete(v)
	}

	return nil
}
