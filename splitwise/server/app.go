package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vikesh-raj/go-practice/splitwise/providers"
)

// Application holds the application state
type Application struct {
	Opts      Opts
	lastError error
	router    *mux.Router
	provider  providers.DBProvider
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

	a.provider = providers.NewInMemoryDB()
	return nil
}

// StartServer starts the server listening on the given port
func (a *Application) StartServer(port int) error {
	return http.ListenAndServe(":"+strconv.Itoa(port), a.router)
}

func (a *Application) addRoutes() {
	a.router = mux.NewRouter()
	a.router.HandleFunc("/", a.handleIndex).Methods("GET")
	a.router.HandleFunc("/view", a.handleView).Methods("GET")
	a.router.HandleFunc("/add", a.handleAdd).Methods("GET")
	a.router.HandleFunc("/add", a.handleAddPost).Methods("POST")
	a.router.HandleFunc("/settle", a.handleSettle).Methods("GET")
	a.router.HandleFunc("/settle", a.handleSettlePost).Methods("POST")
	a.router.HandleFunc("/readiness", a.handleReady).Methods("GET")
	a.router.HandleFunc("/liveness", a.handleReady).Methods("GET")
}
