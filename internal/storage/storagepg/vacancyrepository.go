package storagepg

import (
	"context"
	"fmt"
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

	vacancies := model.SearchVacancies(t.Name, t.GroupID)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx := db.MustBeginTx(ctx, nil)

	for _, v := range vacancies {
		if err := r.upsert(ctx, tx, v); err != nil {
			tx.Rollback()
			logrus.Error(err)
			return
		}
	}

	logger.Warnf("task ended")

	if err := tx.Commit(); err != nil {
		logrus.Error(err)
		return
	}
}

func (r *VacancyRepository) upsert(ctx context.Context, tx *sqlx.Tx, v model.Vacancy) error {
	if _, err := tx.NamedExecContext(
		ctx,
		"INSERT INTO vacancies (task_id, number, link, name, salary_from, salary_to, company, area, description, at_published, responsed, selected) VALUES (:task_id, :number, :link, :name, :salary_from, :salary_to, :company, :area, :description, :at_published, :responsed, :selected) ON CONFLICT (number) DO NOTHING;",
		&v,
	); err != nil {
		return err
	}

	return nil
}

func (r *VacancyRepository) Find(filter *model.Filter) ([]model.Vacancy, error) {
	vacancies := []model.Vacancy{}
	selectQuery := "SELECT p.name AS site, id, g.name AS group_name, number, link, name, salary_from, salary_to, company, area, desription, at_published, selected"
	joinQuery := "FROM vacancies INNER JOIN platforms p ON p.id = platform_id INNER JOIN tasks t ON t.id = task_id INNER JOIN groups g ON g.id = t.group_id"
	query := fmt.Sprintf("%s %s %s ORDER BY at_published_at ASC LIMIT 20;", selectQuery, joinQuery, filter.ToString())
	fmt.Println(query)
	if err := r.store.db.Select(
		&vacancies,
		query,
	); err != nil {
		return nil, err
	}

	for _, v := range vacancies {
		v.ConvertSalary()
	}

	return vacancies, nil
}

func (r *VacancyRepository) CountSelected(groupID int) (int, error) {
	var count int

	if err := r.store.db.QueryRow(
		"SELECT COUNT(selected) FROM vacancies INNER JOIN tasks t ON t.id = task_id WHERE t.group_id = $1 AND selected = true;",
		groupID,
	).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *VacancyRepository) Response(vacancyID int) error {
	_, err := r.store.db.Exec("UPDATE vacancies SET responsed = true WHERE id = $1", vacancyID)
	if err != nil {
		return err
	}
	return nil
}

func (r *VacancyRepository) Toggle(vacancyID int) error {
	_, err := r.store.db.Exec("UPDATE vacancies SET selected = NOT selected WHERE id = $1", vacancyID)
	if err != nil {
		return err
	}
	return nil
}
