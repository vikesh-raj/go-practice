package providers

import (
	"github.com/vikesh-raj/go-practice/quiz/models"
)

type UserQuiz interface {
	AddScore(models.UserQuizzes) error
	GetQuizzes(user string) ([]models.UserQuizzes, error)
}
