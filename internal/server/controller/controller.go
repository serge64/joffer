package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"guthub.com/serge64/joffer/internal/storage"
)

type Controller struct {
	middleware *middleware
	store      storage.Store
}

type errorResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Path    string `json:"path"`
}

func New(store storage.Store) *Controller {
	return &Controller{
		store: store,
	}
}

func (c *Controller) NotFoundHandler() http.HandlerFunc {
	return c.errorHandler(http.StatusNotFound)
}

func (c *Controller) MethodNotAllowedHandler() http.HandlerFunc {
	return c.errorHandler(http.StatusMethodNotAllowed)
}

func (c *Controller) errorHandler(code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.error(w, r, code, fmt.Errorf("-"))
	}
}

func (c *Controller) Middleware() *middleware {
	if c.middleware == nil {
		c.middleware = &middleware{
			controller: c,
		}
	}
	return c.middleware
}

func (c *Controller) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	res := &errorResponse{
		Code:    code,
		Status:  http.StatusText(code),
		Message: err.Error(),
		Path:    r.RequestURI,
	}
	c.respond(w, r, code, res)
}

func (c *Controller) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
