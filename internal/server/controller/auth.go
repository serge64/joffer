package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"guthub.com/serge64/joffer/internal/model"
	"guthub.com/serge64/joffer/internal/storage"
)

type auth struct {
	controller *Controller
}

type request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	errInvalidEmailOrPassword = errors.New("некорректная электронная почта или пароль")
	errInvalidEmail           = errors.New("некорректная электронная почта")
	errDontUniqueEmail        = errors.New("электронная почта уже используется")
)

func (a *auth) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		json.NewDecoder(r.Body).Decode(req)

		u, err := a.controller.store.User().FindByEmail(req.Email)
		if err != nil {
			a.controller.error(w, r, http.StatusUnauthorized, errInvalidEmailOrPassword)
			return
		} else if !u.ComparePassword(req.Password) {
			a.controller.error(w, r, http.StatusUnauthorized, errInvalidEmailOrPassword)
			return
		}

		if err := a.setSession(w, r, u.ID); err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (a *auth) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		json.NewDecoder(r.Body).Decode(req)

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		id, err := a.controller.store.User().Create(u)
		if err != nil {
			if err == storage.ErrUserNotValid {
				a.controller.error(w, r, http.StatusUnprocessableEntity, errInvalidEmail)
			} else {
				a.controller.error(w, r, http.StatusUnprocessableEntity, errDontUniqueEmail)
			}
			return
		}

		if err := a.setSession(w, r, id); err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (a *auth) LogOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := a.clearSession(w, r); err != nil {
			a.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func (a *auth) setSession(w http.ResponseWriter, r *http.Request, id int) error {
	session, _ := a.controller.session.Get(r, sessionName)
	session.Values["user_id"] = id
	if err := a.controller.session.Save(r, w, session); err != nil {
		return err
	}
	return nil
}

func (a *auth) clearSession(w http.ResponseWriter, r *http.Request) error {
	session, err := a.controller.session.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	if err := a.controller.session.Save(r, w, session); err != nil {
		return err
	}
	return nil
}
