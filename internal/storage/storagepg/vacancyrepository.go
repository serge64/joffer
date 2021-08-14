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

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"guthub.com/serge64/joffer/internal/model"
)

type VacancyRepository struct {
	store *Store
}

func (r *VacancyRepository) search(db *sqlx.DB, t *model.Task) {
	logger := logrus.WithFields(logrus.Fields{
		"task_id":   t.ID,
		"task_name": t.Name,
		"group_id":  t.GroupID,
	})

	logger.Warnf("task started")

	vacancies := model.SearchVacancies(t.Name, t.ID)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx := db.MustBeginTx(ctx, nil)

	for _, v := range vacancies {
		if _, err := tx.NamedExecContext(
			ctx,
			"INSERT INTO vacancies (task_id, number, link, name, salary_from, salary_to, company, area, description, at_published, responsed, selected) VALUES (:task_id, :number, :link, :name, :salary_from, :salary_to, :company, :area, :description, :at_published, :responsed, :selected) ON CONFLICT(number) DO UPDATE SET at_published = :at_published, salary_from = :salary_from, salary_to = :salary_to WHERE vacancies.task_id = :task_id;",
			&v,
		); err != nil {
			tx.Rollback()
			logrus.Error(err)
			return
		}
	}

	logger.Warnf("task ended: found vacancies %d", len(vacancies))

	if err := tx.Commit(); err != nil {
		logrus.Error(err)
		return
	}
}

func (r *VacancyRepository) Find(filter *model.Filter) ([]model.Vacancy, error) {
	vacancies := []model.Vacancy{}
	selectQuery := "SELECT pl.name AS site, v.id, g.name AS group_name, number, link, v.name, salary_from, salary_to, company, area, description, to_char(at_published, 'YYYY-MM-DD HH24:MI:SS') AS at_published, selected"
	joinQuery := "FROM vacancies v INNER JOIN tasks t ON t.id = v.task_id INNER JOIN groups g ON g.id = t.group_id INNER JOIN profiles r ON r.id = g.profile_id INNER JOIN platforms pl ON pl.id = r.platform_id"
	query := fmt.Sprintf("%s %s %s ORDER BY at_published DESC LIMIT 20;", selectQuery, joinQuery, filter.ToString())

	if err := r.store.db.Select(
		&vacancies,
		query,
	); err != nil {
		return nil, err
	}

	for i := 0; i < len(vacancies); i++ {
		vacancies[i].ConvertSalary()
	}

	return vacancies, nil
}

func (r *VacancyRepository) CountSelected(groupID int) (int, error) {
	var count int

	if err := r.store.db.QueryRow(
		"SELECT COUNT(selected) FROM vacancies INNER JOIN tasks t ON t.id = task_id WHERE t.group_id = $1 AND selected = true AND responsed = false;",
		groupID,
	).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *VacancyRepository) Response(userID int, id int) error {
	var (
		vacancyID   string
		resumeID    string
		body        string
		accessToken string
		b           bytes.Buffer
	)

	if err := r.store.db.Get(
		&vacancyID,
		"SELECT number FROM vacancies WHERE id = $1;",
		id,
	); err != nil {
		return err
	}

	if err := r.store.db.Get(
		&resumeID,
		"SELECT uid FROM resumes r INNER JOIN profiles p ON p.id = r.profile_id INNER JOIN groups g ON p.id = g.profile_id INNER JOIN tasks t ON g.id = t.group_id INNER JOIN vacancies v ON t.id = v.task_id WHERE v.id = $1;",
		id,
	); err != nil {
		return err
	}

	if err := r.store.db.Get(
		&body,
		"SELECT body FROM letters l INNER JOIN groups g ON l.name = g.letter INNER JOIN tasks t ON g.id = t.group_id INNER JOIN vacancies v ON v.task_id = t.id WHERE v.id = $1;",
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

	values := map[string]io.Reader{
		"vacancy_id": strings.NewReader(vacancyID),
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
		"UPDATE vacancies SET responsed = true WHERE id = $1",
		id,
	); err != nil {
		return err
	}

	return nil
}

func (r *VacancyRepository) Toggle(id int) error {
	_, err := r.store.db.Exec("UPDATE vacancies SET selected = NOT selected WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
