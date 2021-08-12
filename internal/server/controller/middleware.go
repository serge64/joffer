package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type middleware struct {
	controller *Controller
}

type ctxKey int8

var (
	sessionName        = "session"
	ctxKeyUser  ctxKey = 1
)

func (m *middleware) AuthenticateHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.controller.session.Get(r, sessionName)
		if err != nil {
			m.controller.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		u, err := m.controller.store.User().Find(id.(int))
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (m *middleware) LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		if !strings.HasPrefix(r.RequestURI, "/public") {
			format := fmt.Sprintln(r.Method, r.RequestURI, rw.code, time.Since(start))

			switch {
			case rw.code >= 500:
				logrus.Error(format)
			case rw.code >= 400:
				logrus.Warn(format)
			default:
				logrus.Info(format)
			}
		}
	})
}
