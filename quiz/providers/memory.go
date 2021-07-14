package providers

import (
	"sync"

	"github.com/vikesh-raj/go-practice/quiz/models"
)

type inMemoryDB struct {
	mutex   *sync.Mutex
	quizzes []models.UserQuizzes
}

func NewInMemoryStore() UserQuiz {
	return &inMemoryDB{
		mutex:   &sync.Mutex{},
		quizzes: make([]models.UserQuizzes, 0),
	}
}

func (db *inMemoryDB) AddScore(o models.UserQuizzes) error {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.quizzes = append(db.quizzes, o)
	return nil
}

func (db *inMemoryDB) GetQuizzes(user string) ([]models.UserQuizzes, error) {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	quizzes := make([]models.UserQuizzes, 0)
	for i := len(db.quizzes) - 1; i >= 0; i-- {
		quiz := db.quizzes[i]
		if quiz.User == user {
			quizzes = append(quizzes, quiz)
		}
	}
	return quizzes, nil
}
