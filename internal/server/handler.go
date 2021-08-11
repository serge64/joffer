package server

import (
	"net/http"

	"guthub.com/serge64/joffer/internal/server/controller"
	"guthub.com/serge64/joffer/internal/storage"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type handler struct {
	router     *mux.Router
	controller *controller.Controller
}

func newHandler(store storage.Store, session storage.Session) *handler {
	h := &handler{
		router:     mux.NewRouter(),
		controller: controller.New(store, session),
	}

	h.configureRouting()

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *handler) configureRouting() {
	h.router.Use(
		h.controller.Middleware().LoggingHandler,
		handlers.CORS(handlers.AllowedOrigins([]string{"*"})),
	)

	h.router.HandleFunc("/login", h.controller.Pages().Login())
	h.router.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))
	h.router.HandleFunc("/auth/login", h.controller.Auth().SignIn()).Methods("POST")
	h.router.HandleFunc("/auth/signup", h.controller.Auth().SignUp()).Methods("POST")

	auth := h.router.NewRoute().Subrouter()

	auth.Use(h.controller.Middleware().AuthenticateHandler)
	auth.HandleFunc("/auth/logout", h.controller.Auth().LogOut()).Methods("POST")
	auth.HandleFunc("/", h.controller.Pages().Index())
	auth.HandleFunc("/profiles", h.controller.Pages().Profiles())
}
