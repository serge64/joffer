package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"guthub.com/serge64/joffer/internal/model"
)

type api struct {
	controller *Controller
}

type group struct {
	Name      string
	Resume    string
	Letter    string
	Positions []string
}

type letter struct {
	Name string
	Body string
}

func (a *api) AddProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := a.controller.config.ClientID
		url := fmt.Sprintf("https://hh.ru/oauth/authorize?response_type=code&client_id=%s", id)
		a.controller.respond(w, r, http.StatusOK, url)
	}
}

func (a *api) Profile() http.HandlerFunc {
	type profileResponse struct {
		Name    string         `json:"name"`
		Email   string         `json:"email"`
		Resumes []string       `json:"resumes"`
		Groups  []model.Group  `json:"groups"`
		Letters []model.Letter `json:"letters"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)

		profile, err := a.controller.store.Profile().Find(u.ID, 1)
		if err != nil {
			logrus.Error(err)
			a.controller.respond(w, r, http.StatusOK, nil)
			return
		}

		resumes, err := a.controller.store.Resume().Find(profile.ID)
		if err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		letters, err := a.controller.store.Letter().Find(u.ID)
		if err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		groups, err := a.controller.store.Group().Find(profile.ID)
		if err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		toArray := func(m []model.Resume) []string {
			array := []string{}
			for _, v := range m {
				array = append(array, v.Name)
			}
			return array
		}

		res := &profileResponse{
			Name:    profile.Name,
			Email:   profile.Email,
			Resumes: toArray(resumes),
			Letters: letters,
			Groups:  groups,
		}

		a.controller.respond(w, r, http.StatusOK, res)
	}
}

func (a *api) DeleteProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)
		err := a.controller.store.Profile().Delete(u.ID, 1)
		if err != nil {
			logrus.Error(err)
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}
		a.controller.respond(w, r, http.StatusOK, nil)
	}
}

func (a *api) CreateGroup() http.HandlerFunc {
	type response struct {
		ID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)
		request := &group{}
		json.NewDecoder(r.Body).Decode(request)

		profile, err := a.controller.store.Profile().Find(u.ID, 1)
		if err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		group := &model.Group{
			ProfileID: profile.ID,
			Name:      request.Name,
			Resume:    request.Resume,
			Letter:    request.Letter,
			Positions: request.Positions,
		}

		id, err := a.controller.store.Group().Create(group)
		if err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		res := &response{
			ID: strconv.Itoa(id),
		}

		a.controller.respond(w, r, http.StatusCreated, res)
	}
}

func (a *api) UpdateGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		iD := vars["id"]
		id, _ := strconv.Atoi(iD)

		u := r.Context().Value(ctxKeyUser).(*model.User)
		request := &group{}
		json.NewDecoder(r.Body).Decode(request)

		profile, err := a.controller.store.Profile().Find(u.ID, 1)
		if err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		group := &model.Group{
			ID:        id,
			ProfileID: profile.ID,
			Name:      request.Name,
			Resume:    request.Resume,
			Letter:    request.Letter,
			Positions: request.Positions,
		}

		if err := a.controller.store.Group().Update(group); err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		a.controller.respond(w, r, http.StatusOK, nil)
	}
}

func (a *api) DeleteGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		iD := vars["id"]
		id, _ := strconv.Atoi(iD)

		if err := a.controller.store.Group().Delete(id); err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		a.controller.respond(w, r, http.StatusOK, nil)
	}
}

func (a *api) GroupResponse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)
		vars := mux.Vars(r)
		iD := vars["id"]
		id, _ := strconv.Atoi(iD)

		if err := a.controller.store.Group().Response(u.ID, id); err != nil {
			logrus.Error(err)
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		logrus.WithFields(logrus.Fields{
			"group_id": id,
		}).Warnf("successful application for a vacancy by group")

		a.controller.respond(w, r, http.StatusOK, nil)
	}
}

func (a *api) CreateLetter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)
		request := &letter{}
		json.NewDecoder(r.Body).Decode(request)

		letter := &model.Letter{
			UserID: u.ID,
			Name:   request.Name,
			Body:   request.Body,
		}

		id, err := a.controller.store.Letter().Create(letter)
		if err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		a.controller.respond(w, r, http.StatusCreated, id)
	}
}

func (a *api) UpdateLetter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		iD := vars["id"]
		id, _ := strconv.Atoi(iD)

		u := r.Context().Value(ctxKeyUser).(*model.User)
		request := &letter{}
		json.NewDecoder(r.Body).Decode(request)

		letter := &model.Letter{
			ID:     id,
			UserID: u.ID,
			Name:   request.Name,
			Body:   request.Body,
		}

		if err := a.controller.store.Letter().Update(letter); err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		a.controller.respond(w, r, http.StatusOK, nil)
	}
}

func (a *api) DeleteLetter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)
		vars := mux.Vars(r)
		iD := vars["id"]
		id, _ := strconv.Atoi(iD)

		if err := a.controller.store.Letter().Delete(u.ID, id); err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		a.controller.respond(w, r, http.StatusOK, nil)
	}
}

func (a *api) Filters() http.HandlerFunc {
	type filter struct {
		Sites  []string `json:"sites"`
		Groups []string `json:"groups"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)

		profile, err := a.controller.store.Profile().Find(u.ID, 1)
		if err != nil {
			logrus.Error(err)
			a.controller.error(w, r, http.StatusNotFound, errors.New("-"))
			return
		}

		sites, err := a.controller.store.Platform().Find()
		if err != nil {
			logrus.Error(err)
			sites = []string{}
		}

		groups, err := a.controller.store.Group().FindList(profile.ID)
		if err != nil {
			logrus.Error(err)
			groups = []string{}
		}

		res := &filter{
			Sites:  sites,
			Groups: groups,
		}

		a.controller.respond(w, r, http.StatusOK, res)
	}
}

func (a *api) Vacancies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := &model.Filter{}
		json.NewDecoder(r.Body).Decode(&filter)
		u := r.Context().Value(ctxKeyUser).(*model.User)

		filter.UserID = u.ID

		vacancies, err := a.controller.store.Vacancy().Find(filter)
		if err != nil {
			logrus.Error(err)
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		a.controller.respond(w, r, http.StatusOK, vacancies)
	}
}

func (a *api) Toggle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iD := mux.Vars(r)["id"]
		id, err := strconv.Atoi(iD)

		if err != nil {
			a.controller.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if err := a.controller.store.Vacancy().Toggle(id); err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		a.controller.respond(w, r, http.StatusOK, nil)
	}
}

func (a *api) Response() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)
		iD := mux.Vars(r)["id"]
		id, err := strconv.Atoi(iD)
		if err != nil {
			a.controller.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if err := a.controller.store.Vacancy().Response(u.ID, id); err != nil {
			logrus.Error(err)
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		logrus.WithFields(logrus.Fields{
			"vacancy_id": id,
		}).Warnf("successful application for a vacancy")

		a.controller.respond(w, r, http.StatusOK, nil)
	}
}

func (a *api) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)
		vars := mux.Vars(r)
		value := vars["value"]

		switch category := vars["category"]; category {
		case "position":
			found, err := a.controller.store.Search().Position(u.ID, value)
			if err != nil {
				logrus.Error(err)
				a.controller.error(w, r, http.StatusInternalServerError, err)
				return
			}
			a.controller.respond(w, r, http.StatusOK, found)
		case "company":
			found, err := a.controller.store.Search().Company(u.ID, value)
			if err != nil {
				logrus.Error(err)
				a.controller.error(w, r, http.StatusInternalServerError, err)
				return
			}
			a.controller.respond(w, r, http.StatusOK, found)
		case "area":
			found, err := a.controller.store.Search().Area(u.ID, value)
			if err != nil {
				logrus.Error(err)
				a.controller.error(w, r, http.StatusInternalServerError, err)
				return
			}
			a.controller.respond(w, r, http.StatusOK, found)
		}
	}
}
