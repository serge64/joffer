package controller

import (
	"net/http"
)

type oauth struct {
	controller *Controller
}

func (o *oauth) HeadHunter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// u := r.Context().Value(ctxKeyUser).(*model.User)
		// code := r.URL.Query().Get("code")]
		// if code != "" {
		// 	p := &model.Profile{
		// 		UserID:     u.ID,
		// 		PlatformID: 1,
		// 	}
		// }
		http.Redirect(w, r, "/profiles", http.StatusFound)
	}
}
