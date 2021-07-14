package models

type Question struct {
	Question      string
	Answers       []string
	Marks         int
	NegativeMark  int
	CorrectAnswer int
}

type Quiz struct {
	Name      string
	Questions []Question
}

var Quizzes = []Quiz{
	{
		Name: "basic",
		Questions: []Question{
			{
				Question: "What is the color of sky ?",
				Answers: []string{
					"red", "blue", "yellow", "green",
				},
				CorrectAnswer: 1,
			},
			{
				Question: "What is the favourite color ?",
				Answers: []string{
					"red", "blue", "yellow", "green",
				},
				CorrectAnswer: 0,
			},
		},
	},
	{
		Name: "color",
		Questions: []Question{
			{
				Question: "What is the color of sky ?",
				Answers: []string{
					"red", "blue", "yellow", "green",
				},
				CorrectAnswer: 1,
			},
			{
				Question: "What is the favourite color ?",
				Answers: []string{
					"red", "blue", "yellow", "green",
				},
				CorrectAnswer: 0,
			},
		},
	},
}
