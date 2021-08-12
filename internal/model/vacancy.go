package model

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/valyala/fastjson"
)

type Vacancy struct {
	ID          int     `json:"id" db:"id"`
	GroupName   string  `json:"group" db:"group_name"`
	TaskID      int     `json:"-" db:"task_id"`
	Number      string  `json:"vacancy_id" db:"number"`
	Link        string  `json:"vacancy_link" db:"link"`
	Name        string  `json:"vacancy_name" db:"name"`
	Salary      string  `json:"salary" db:"-"`
	SalaryFrom  float64 `json:"-" db:"salary_from"`
	SalaryTo    float64 `json:"-" db:"salary_to"`
	Company     string  `json:"company" db:"company"`
	Area        string  `json:"area" db:"area"`
	Description string  `json:"description" db:"description"`
	AtPublished string  `json:"at_published" db:"at_published"`
	Responsed   bool    `json:"-" db:"responsed"`
	Selected    bool    `json:"selected" db:"selected"`
}

func (v *Vacancy) ConvertSalary() {
	var salary []string

	if v.SalaryFrom > 0 {
		salary[0] = fmt.Sprintf("от %v", v.SalaryFrom)
	}
	if v.SalaryTo > 0 {
		salary[1] = fmt.Sprintf("до %v", v.SalaryTo)
	}

	v.Salary = strings.Join(salary, " ")
}

func SearchVacancies(pattern string, taskID int) []Vacancy {
	urls := getURLs(pattern)
	ch := make(chan []Vacancy)
	vacancies := []Vacancy{}
	counter := 0

	for _, url := range urls {
		go func(url string) {
			r, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
			}
			ch <- getVacancies(r, taskID)
		}(url)
	}

	for {
		vacancies = append(vacancies, <-ch...)
		counter++
		if counter == len(urls) {
			break
		}
	}

	return vacancies
}

func getURLs(text string) []string {
	urls := []string{}

	for i := 0; i < 20; i++ {
		urls = append(
			urls,
			fmt.Sprintf(
				"https://api.hh.ru/vacancies?text=%s&per_page=100&page=%d",
				url.QueryEscape(text),
				i,
			),
		)
	}

	return urls
}

func getVacancies(r *http.Response, taskID int) []Vacancy {
	json := func() []byte {
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(r.Body); err != nil {
			log.Fatal(err)
		}
		return buf.Bytes()
	}()

	vacancies := []Vacancy{}
	p := fastjson.Parser{}
	value, err := p.ParseBytes(json)
	if err != nil {
		log.Fatal(err)
	}

	for _, o := range value.Get("items").GetArray() {
		if !o.GetBool("has_test") && o.GetStringBytes("response_url") == nil {
			v := &Vacancy{
				TaskID:      taskID,
				Number:      string(o.GetStringBytes("id")),
				Link:        string(o.GetStringBytes("alternate_url")),
				Name:        string(o.GetStringBytes("name")),
				SalaryFrom:  o.GetFloat64("salary", "from"),
				SalaryTo:    o.GetFloat64("salary", "to"),
				Area:        string(o.GetStringBytes("area", "name")),
				Company:     string(o.GetStringBytes("employer", "name")),
				Description: string(o.GetStringBytes("snippet", "responsibility")),
				AtPublished: string(o.GetStringBytes("published_at")),
			}
			vacancies = append(vacancies, *v)
		}
	}

	return vacancies
}
