package controller

import (
	"fmt"
	"net/http"
)

type api struct {
	controller *Controller
}

func (a *api) AddProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := a.controller.config.ClientID
		url := fmt.Sprintf("https://hh.ru/oauth/authorize?response_type=code&client_id=%s", id)
		a.controller.respond(w, r, http.StatusOK, url)
	}
}
