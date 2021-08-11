package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type middleware struct {
	controller *Controller
}

func (m *middleware) LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		format := fmt.Sprintln(r.Method, r.RequestURI, rw.code, time.Since(start))

		switch {
		case rw.code >= 500:
			logrus.Error(format)
		case rw.code >= 400:
			logrus.Warn(format)
		default:
			logrus.Info(format)
		}
	})
}
