package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Application holds the application state
type Application struct {
	Opts      Opts
	lastError error
	router    *mux.Router
}

// CreateApplication creates the application given the options
func CreateApplication(opts Opts) (*Application, error) {
	a := Application{
		Opts: opts,
	}
	err := a.Init()
	a.lastError = err
	return &a, err
}

// Init initializes the application
func (a *Application) Init() error {
	a.addRoutes()
	return nil
}

// StartServer starts the server listening on the given port
func (a *Application) StartServer(port int) error {
	return http.ListenAndServe(":"+strconv.Itoa(port), a.router)
}

func (a *Application) addRoutes() {
	a.router = mux.NewRouter()
	a.router.HandleFunc("/", a.handleIndex).Methods("GET")
	a.router.HandleFunc("/readiness", a.handleReady).Methods("GET")
	a.router.HandleFunc("/liveness", a.handleReady).Methods("GET")
}
