package storagepg

import (
	"errors"

	"github.com/go-co-op/gocron"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"guthub.com/serge64/joffer/internal/model"
)

type TaskRepository struct {
	store *Store
	cron  *gocron.Scheduler
	tasks map[int]*gocron.Job
}

func (r *TaskRepository) Create(t *model.Task, tx *sqlx.Tx) error {
	if err := tx.QueryRow(
		"INSERT INTO tasks (group_id, name) VALUES ($1, $2) RETURNING id;",
		t.GroupID,
		t.Name,
	).Scan(&t.ID); err != nil {
		return err
	}

	j, err := r.cron.Every(2).Minutes().Do(r.store.vacancyRepository.search, r.store.db, t)
	if err != nil {
		return err
	}

	r.tasks[t.ID] = j
	logrus.WithFields(logrus.Fields{
		"task_id":   t.ID,
		"task_name": t.Name,
		"group_id":  t.GroupID,
	}).Warnf("task created")

	return nil
}

func (r *TaskRepository) Find(groupID int) ([]model.Task, error) {
	tasks := []model.Task{}

	if err := r.store.db.Select(
		&tasks,
		"SELECT * FROM tasks WHERE group_id = $1;",
		groupID,
	); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) Delete(groupID int, name string) error {
	id := 0
	if err := r.store.db.QueryRow(
		"DELETE FROM tasks WHERE group_id = $1 AND name = $2 RETURNING id;",
		groupID,
		name,
	).Scan(&id); err != nil {
		return err
	}

	j, ok := r.tasks[id]
	if !ok {
		err := errors.New("task not found")
		logrus.Error(err)
		return err
	}

	r.cron.RemoveByReference(j)
	delete(r.tasks, id)

	logrus.WithFields(logrus.Fields{
		"task_id":   id,
		"task_name": name,
		"group_id":  groupID,
	}).Warnf("task removed")

	return nil
}
