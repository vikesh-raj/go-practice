package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vikesh-raj/go-practice/quiz/models"
	"github.com/vikesh-raj/go-practice/quiz/providers"
)

// Application holds the application state
type Application struct {
	Opts      Opts
	lastError error
	router    *mux.Router
	provider  providers.UserQuiz
	quizzes   []models.Quiz
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

	a.quizzes = models.Quizzes
	a.provider = providers.NewInMemoryStore()

	return nil
}

func (a *Application) getQuizzes() []string {
	quizzes := make([]string, 0)
	for _, q := range a.quizzes {
		quizzes = append(quizzes, q.Name)
	}
	return quizzes
}

func (a *Application) getQuestions(quiz string) []models.Question {
	for _, q := range a.quizzes {
		if q.Name == quiz {
			return q.Questions
		}
	}
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
	a.router.HandleFunc("/quiz", a.handleQuiz).Methods("POST")
	a.router.HandleFunc("/readiness", a.handleReady).Methods("GET")
	a.router.HandleFunc("/liveness", a.handleReady).Methods("GET")
}
