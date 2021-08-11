package server

import (
	"net/http"

	"guthub.com/serge64/joffer/internal/server/controller"
	"guthub.com/serge64/joffer/internal/storage"

	"github.com/gorilla/mux"
)

type handler struct {
	router     *mux.Router
	controller *controller.Controller
}

func newHandler(store storage.Store) *handler {
	h := &handler{
		router:     mux.NewRouter(),
		controller: controller.New(store),
	}

	h.configureRouting()

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *handler) configureRouting() {
	h.router.NotFoundHandler = h.controller.NotFoundHandler()
	h.router.MethodNotAllowedHandler = h.controller.MethodNotAllowedHandler()

	h.router.Use(
		h.controller.Middleware().ContentTypeHandler,
		h.controller.Middleware().LoggingHandler,
	)
}
