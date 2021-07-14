package models

import "time"

type UserQuizzes struct {
	ID    string
	User  string
	Quiz  string
	Total int
	Score int
	Time  time.Time
}

type QuizSession struct {
	ID          string
	Quiz        string
	QuestionIdx int
	Score       int
}
