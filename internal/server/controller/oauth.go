package controller

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"guthub.com/serge64/joffer/internal/model"
)

type oauth struct {
	controller *Controller
}

func (o *oauth) HeadHunter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ctxKeyUser).(*model.User)
		code := r.URL.Query().Get("code")
		if code != "" {
			p := &model.Profile{
				UserID:     u.ID,
				PlatformID: 1,
			}
			if err := o.controller.store.Profile().Create(p, code, o.controller.config); err != nil {
				logrus.Error(err)
			}
		}
		http.Redirect(w, r, "/profiles", http.StatusFound)
	}
}
