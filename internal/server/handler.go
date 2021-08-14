package server

import (
	"net/http"

	"guthub.com/serge64/joffer/internal/config"
	"guthub.com/serge64/joffer/internal/server/controller"
	"guthub.com/serge64/joffer/internal/storage"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type handler struct {
	router     *mux.Router
	controller *controller.Controller
}

func newHandler(store storage.Store, session storage.Session, config *config.Config) *handler {
	h := &handler{
		router:     mux.NewRouter(),
		controller: controller.New(store, session, config),
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
	h.router.HandleFunc("/auth/signin", h.controller.Auth().SignIn()).Methods("POST")
	h.router.HandleFunc("/auth/signup", h.controller.Auth().SignUp()).Methods("POST")

	auth := h.router.NewRoute().Subrouter()
	auth.Use(h.controller.Middleware().AuthenticateHandler)

	auth.HandleFunc("/auth/logout", h.controller.Auth().LogOut()).Methods("POST")
	auth.HandleFunc("/", h.controller.Pages().Index())
	auth.HandleFunc("/profiles", h.controller.Pages().Profiles())

	oauth := auth.PathPrefix("/oauth").Subrouter()
	oauth.HandleFunc("/hh", h.controller.OAuth().HeadHunter())

	api := auth.PathPrefix("/api").Subrouter()
	api.HandleFunc("/profile", h.controller.API().AddProfile()).Methods("POST")
	api.HandleFunc("/profile", h.controller.API().Profile()).Methods("GET")
	api.HandleFunc("/profile", h.controller.API().DeleteProfile()).Methods("DELETE")

	api.HandleFunc("/groups", h.controller.API().CreateGroup()).Methods("POST")
	api.HandleFunc("/groups/{id:[0-9]+}", h.controller.API().UpdateGroup()).Methods("PATCH")
	api.HandleFunc("/groups/{id:[0-9]+}", h.controller.API().DeleteGroup()).Methods("DELETE")
	api.HandleFunc("/groups/{id:[0-9]+}", h.controller.API().GroupResponse()).Methods("POST")

	api.HandleFunc("/letters", h.controller.API().CreateLetter()).Methods("POST")
	api.HandleFunc("/letters/{id:[0-9]+}", h.controller.API().UpdateLetter()).Methods("PATCH")
	api.HandleFunc("/letters/{id:[0-9]+}", h.controller.API().DeleteLetter()).Methods("DELETE")

	api.HandleFunc("/vacancies", h.controller.API().Vacancies()).Methods("POST")
	api.HandleFunc("/vacancies/{id:[0-9]+}", h.controller.API().Toggle()).Methods("PATCH")
	api.HandleFunc("/vacancies/{id:[0-9]+}", h.controller.API().Response()).Methods("POST")
	api.HandleFunc("/filters", h.controller.API().Filters()).Methods("GET")
	api.HandleFunc("/search/{category}", h.controller.API().Search()).Queries("value", "{value}").Methods("GET")
}
