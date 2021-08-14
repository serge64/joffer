package storagepg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"guthub.com/serge64/joffer/internal/model"
)

type GroupRepository struct {
	store *Store
}

func (r *GroupRepository) Create(g *model.Group) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx := r.store.db.MustBeginTx(ctx, nil)

	if err := tx.QueryRow(
		"INSERT INTO groups (profile_id, name, resume, letter) VALUES ($1, $2, $3, $4) RETURNING id;",
		g.ProfileID,
		g.Name,
		g.Resume,
		g.Letter,
	).Scan(&g.ID); err != nil {
		tx.Rollback()
		logrus.Error(err)
		return 0, err
	}

	for _, v := range g.Positions {
		task := &model.Task{
			GroupID: g.ID,
			Name:    v,
		}
		if err := r.store.Task().Create(task, tx); err != nil {
			tx.Rollback()
			logrus.Error(err)
			return 0, err
		}
	}

	tx.Commit()
	return g.ID, nil
}

func (r *GroupRepository) Find(profileID int) ([]model.Group, error) {
	groups := []model.Group{}

	if err := r.store.db.Select(
		&groups,
		"SELECT * FROM groups WHERE profile_id = $1;",
		profileID,
	); err != nil {
		return nil, err
	}

	for i := 0; i < len(groups); i++ {
		v := &groups[i]

		tasks, err := r.store.Task().Find(v.ID)
		if err != nil {
			return nil, err
		}

		count, err := r.store.Vacancy().CountSelected(v.ID)
		if err != nil {
			return nil, err
		}

		v.Positions = r.tasksToArray(tasks)
		v.Count = count
	}

	return groups, nil
}

func (r *GroupRepository) FindList(profileID int) ([]string, error) {
	groups := []string{}

	if err := r.store.db.Select(
		&groups,
		"SELECT name FROM groups WHERE profile_id = $1;",
		profileID,
	); err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *GroupRepository) Update(g *model.Group) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx := r.store.db.MustBeginTx(ctx, nil)

	if _, err := tx.ExecContext(
		ctx,
		"UPDATE groups SET name = $1, resume = $2, letter = $3 WHERE id = $4;",
		g.Name,
		g.Resume,
		g.Letter,
		g.ID,
	); err != nil {
		tx.Rollback()
		return err
	}

	tasks, err := r.store.Task().Find(g.ID)
	if err != nil {
		tx.Commit()
		return err
	}

	slice := r.tasksToArray(tasks)
	diff := r.difference(slice, g.Positions)

	for _, v := range diff {
		if err := r.store.Task().Delete(g.ID, v); err != nil {
			t := &model.Task{
				GroupID: g.ID,
				Name:    v,
			}
			if err := r.store.Task().Create(t, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	return nil
}

func (r *GroupRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx := r.store.db.MustBeginTx(ctx, nil)

	if _, err := tx.ExecContext(
		ctx,
		"DELETE FROM groups WHERE id = $1;",
		id,
	); err != nil {
		tx.Rollback()
		return err
	}

	tasks, err := r.store.Task().Find(id)
	if err != nil {
		tx.Commit()
		return err
	}

	for _, v := range tasks {
		r.store.Task().Delete(id, v.Name)
	}

	tx.Commit()
	return nil
}

func (r *GroupRepository) Response(userID int, id int) error {
	var (
		vacancyID   []string
		resumeID    string
		body        string
		accessToken string
		b           bytes.Buffer
	)

	if err := r.store.db.Select(
		&vacancyID,
		"SELECT number FROM vacancies v INNER JOIN tasks t ON t.id = v.task_id INNER JOIN groups g ON g.id = t.group_id WHERE v.selected = true AND t.group_id = $1;",
		id,
	); err != nil {
		return err
	}

	if err := r.store.db.Get(
		&resumeID,
		"SELECT uid FROM resumes r INNER JOIN groups g ON g.resume = r.name WHERE g.id = $1;",
		id,
	); err != nil {
		return err
	}

	if err := r.store.db.Get(
		&body,
		"SELECT body FROM letters l INNER JOIN groups g ON l.name = g.letter WHERE g.id = $1;",
		id,
	); err != nil {
		return err
	}

	if err := r.store.db.Get(
		&accessToken,
		"SELECT access_token FROM profiles WHERE user_id = $1;",
		userID,
	); err != nil {
		return err
	}

	for _, v := range vacancyID {
		values := map[string]io.Reader{
			"vacancy_id": strings.NewReader(v),
			"resume_id":  strings.NewReader(resumeID),
			"message":    strings.NewReader(body),
		}

		w := multipart.NewWriter(&b)
		for key, r := range values {
			var fw io.Writer
			if x, ok := r.(io.Closer); ok {
				defer x.Close()
			}
			fw, err := w.CreateFormField(key)
			if err != nil {
				return nil
			}
			if _, err := io.Copy(fw, r); err != nil {
				return err
			}
		}
		w.Close()

		req, err := http.NewRequest("POST", "https://api.hh.ru/negotiations", &b)
		if err != nil {
			return err
		}

		req.Header.Add("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", w.FormDataContentType())

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		if res.StatusCode != http.StatusCreated {
			return fmt.Errorf("bad request: %s", res.Status)
		}

		if _, err = r.store.db.Exec(
			"UPDATE vacancies SET responsed = true WHERE number = $1",
			v,
		); err != nil {
			return err
		}
	}

	return nil
}

func (r *GroupRepository) difference(slice1 []string, slice2 []string) []string {
	diff := []string{}

	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}
		}
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

func (r *GroupRepository) tasksToArray(t []model.Task) []string {
	tasks := []string{}
	for _, v := range t {
		tasks = append(tasks, v.Name)
	}
	return tasks
}
