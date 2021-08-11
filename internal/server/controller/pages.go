package controller

import (
	"html/template"
	"net/http"
)

type pages struct {
	template *template.Template
}

type data struct {
	IsFirst bool
}

func (p *pages) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.template.ExecuteTemplate(w, "login.html", nil)
	}
}

func (p *pages) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := &data{
			IsFirst: true,
		}
		p.template.ExecuteTemplate(w, "index.html", d)
	}
}

func (p *pages) Profiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := &data{}
		p.template.ExecuteTemplate(w, "profiles.html", d)
	}
}
