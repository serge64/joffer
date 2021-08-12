package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	return func(w http.ResponseWriter, r *http.Request) {
		// u := r.Context().Value(ctxKeyUser).(*model.User)
	}
}

func (a *api) CreateGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// u := r.Context().Value(ctxKeyUser).(*model.User)
		request := &group{}
		json.NewDecoder(r.Body).Decode(request)

		group := &model.Group{
			UserID:    1, //u.ID,
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

		a.controller.respond(w, r, http.StatusCreated, id)
	}
}

func (a *api) UpdateGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		iD := vars["id"]
		id, _ := strconv.Atoi(iD)

		// u := r.Context().Value(ctxKeyUser).(*model.User)
		request := &group{}
		json.NewDecoder(r.Body).Decode(request)

		group := &model.Group{
			ID:        id,
			UserID:    1, //,
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

func (a *api) CreateLetter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// u := r.Context().Value(ctxKeyUser).(*model.User)
		request := &letter{}
		json.NewDecoder(r.Body).Decode(request)

		letter := &model.Letter{
			UserID: 1, //u.ID,
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

		// u := r.Context().Value(ctxKeyUser).(*model.User)
		request := &letter{}
		json.NewDecoder(r.Body).Decode(request)

		letter := &model.Letter{
			ID:     id,
			UserID: 1, //,
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
		// u := r.Context().Value(ctxKeyUser).(*model.User)
		vars := mux.Vars(r)
		iD := vars["id"]
		id, _ := strconv.Atoi(iD)

		if err := a.controller.store.Letter().Delete(1, id); err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		a.controller.respond(w, r, http.StatusOK, nil)
	}
}
