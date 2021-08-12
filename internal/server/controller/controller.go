package controller

import (
	"encoding/json"
	"html/template"
	"net/http"

	"guthub.com/serge64/joffer/internal/config"
	"guthub.com/serge64/joffer/internal/storage"
)

type Controller struct {
	middleware *middleware
	pages      *pages
	auth       *auth
	oauth      *oauth
	api        *api
	store      storage.Store
	session    storage.Session
	config     *config.Config
}

type errorResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Path    string `json:"path"`
}

func New(store storage.Store, session storage.Session, config *config.Config) *Controller {
	return &Controller{
		store:   store,
		session: session,
		config:  config,
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

func (c *Controller) Pages() *pages {
	if c.pages == nil {
		c.pages = &pages{
			template: template.Must(template.ParseGlob("templates/html/*html")),
		}
	}
	return c.pages
}

func (c *Controller) Auth() *auth {
	if c.auth == nil {
		c.auth = &auth{
			controller: c,
		}
	}
	return c.auth
}

func (c *Controller) OAuth() *oauth {
	if c.oauth == nil {
		c.oauth = &oauth{
			controller: c,
		}
	}
	return c.oauth
}

func (c *Controller) API() *api {
	if c.api == nil {
		c.api = &api{
			controller: c,
		}
	}
	return c.api
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
