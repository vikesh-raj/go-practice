package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/vikesh-raj/go-practice/quiz/models"
)

func (a *Application) handleError(w http.ResponseWriter) {
	writeErrorResponse(w, a.lastError.Error(), http.StatusServiceUnavailable)
}

func (a *Application) handleReady(w http.ResponseWriter, r *http.Request) {
	if a.lastError != nil {
		a.handleError(w)
	} else {
		res := BasicResponse{
			Status: "OK",
		}
		writeResponse(w, &res, http.StatusOK)
	}
}

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <title>Quiz</title>
	<style>
    .err {
        background-color: #FFCCCC;
    }
	</style>
</head>
<body>
<h1>Quiz</h1>
<div>Version : {{.Version}}</div><br/>
{{if .ErrorMessage}}
<div class="err">{{.ErrorMessage}}</div><br/>
{{end}}
{{if .User}}
<ul>
<li><a href="/view?user={{.User}}">View Scores</a></li>
<li><a href="/">Logout</a></li>
</ul>

<form action="/quiz?user={{.User}}" method="post">
<p> Take Quiz :
	<select name = "quiz">
	{{range $index, $item := .Quizzes}}
   		<option value = "{{$item}}">{{$item}}</option>
	{{end}}
</select>
<input type="submit" value="Submit">
</form>

{{else}}
<form action="" method="get">
<label for="user">User:</label>
<input type="text" id="user" name="user"><br><br>
<input type="submit" value="Submit">
</form>
{{end}}
</body>
</html>
`

const viewHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <title>Quiz</title>
	<style>
    .err {
        background-color: #FFCCCC;
    }
	table,
	th,
	td {
	  padding: 10px;
	  border: 1px solid black;
	  border-collapse: collapse;
	}
	</style>
</head>
<body>
<div><a href="/?user={{.User}}">Home</a></div>
{{if .Error}}
<div class="err">{{.Error}}<br/></div><br/>
{{end}}
<h1>Quizzes</h1>
{{if .User}}
	{{if .Quizzes}}
		<table>
			<tr>
				<th>No</th>
	  			<th>Quiz</th>
	  			<th>Total</th>
	  			<th>Score</th>
				<th>Date</th>
			</tr>
			{{range $index, $item := .Quizzes}}
				<tr>
					<td>{{inc $index}}</td>
					<td>{{$item.Quiz}}</td>
					<td>{{$item.Total}}</td>
					<td>{{$item.Score}}</td>
					<td>{{printTime $item.Time}}</td>
				</tr>
			{{end}}
		</table> 
	{{else}}
	<p>You have not taken any quizzes.
	{{end}}
{{else}}
<div class="err"><a href="/">Please login</a></div>
{{end}}
</body>
</html>
`

const quizHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <title>Quiz</title>
	<style>
	.msg {
		background-color: #E0DA96;
	}
	.err {
        background-color: #FFCCCC;
    }
	</style>
</head>
<body>
<div><a href="/?user={{.User}}">Home</a></div>
<h1>Quiz {{.Quiz}} </h1>
{{if .User}}
	{{if .Message}}
		<div class="msg">{{.Message}}<br/></div><br/>
	{{end}}
	{{if .Error}}
		<div class="err">{{.Error}}<br/></div><br/>
	{{end}}

	{{if .QuestionsComplete}}
		<div> Quiz complete. Your score is {{.Score}} out of {{.Total}}</div>
	{{else}}
		<form action="" method="post">
		<input type="hidden" id="quiz" name="quiz" value="{{.Quiz}}">
		<input type="hidden" id="questionidx" name="questionidx" value="{{.QuestionIndex}}">
		<input type="hidden" id="score" name="score" value="{{.Score}}">

		<div id="question">
			<label for="question">{{.Question}}</label><br/>
			{{range $index, $item := .Answers}}
				<input type="radio" id="{{$index}}" name="{{$index}}" value="{{$index}}">
				<label for="{{$index}}">{{$item}}</label><br>
			{{end}}
		</div>
		<input type="submit" value="Submit">
		</form>
	
	{{end}}

{{else}}
<a href="/">Please login</a>
{{end}}
</body>
</html>
`

func printTime(t time.Time) string {
	return t.Format("Mon Jan 2 2006 03:04:05 pm")
}

func inc(i int) int {
	return i + 1
}

var funcMap = template.FuncMap{
	"printTime": printTime,
	"inc":       inc,
}

func parseTemplate(page string) *template.Template {
	return template.Must(template.New("").Funcs(funcMap).Parse(page))
}

var indexTemplate = parseTemplate(indexHTML)

type indexPageParams struct {
	Version      string
	ErrorMessage string
	Opts         Opts
	User         string
	Quizzes      []string
}

func (a *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	indexPageParams := indexPageParams{
		Version: version,
		Opts:    a.Opts,
		User:    r.URL.Query().Get("user"),
		Quizzes: a.getQuizzes(),
	}
	if a.lastError != nil {
		indexPageParams.ErrorMessage = a.lastError.Error()
	}

	indexTemplate.Execute(w, &indexPageParams)
}

var viewTemplate = parseTemplate(viewHTML)

type viewPageParams struct {
	User    string
	Error   string
	Quizzes []models.UserQuizzes
}

func (a *Application) handleView(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	viewPageParams := viewPageParams{
		User: user,
	}
	quizzes, err := a.provider.GetQuizzes(user)
	if err != nil {
		viewPageParams.Error = "Error : " + err.Error()
	} else {
		viewPageParams.Quizzes = quizzes
	}

	viewTemplate.Execute(w, &viewPageParams)
}

var quizTemplate = parseTemplate(quizHTML)

type quizPageParams struct {
	User              string
	Message           string
	Error             string
	Quiz              string
	Score             int
	Total             int
	QuestionIndex     int
	Question          string
	Answers           []string
	QuestionsComplete bool
}

func (a *Application) handleQuiz(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	quizPageParams := quizPageParams{
		User: user,
	}

	msg, err := a.handleQuizPostRequest(user, &quizPageParams, r)
	if err != nil {
		quizPageParams.Error = "Error : " + err.Error()
	}
	quizPageParams.Message = msg
	quizTemplate.Execute(w, &quizPageParams)
}

func (a *Application) handleQuizPostRequest(user string, quizPageParams *quizPageParams, r *http.Request) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", err
	}

	// fmt.Println(r.PostForm)

	quiz := r.PostForm.Get("quiz")
	if quiz == "" {
		return "", fmt.Errorf("no quiz name found")
	}
	quizPageParams.Quiz = quiz

	questions := a.getQuestions(quiz)
	if len(questions) == 0 {
		return "", fmt.Errorf("no questions found")
	}

	quizPageParams.Total = len(questions)
	quizPageParams.QuestionIndex, err = parseInt(r.PostForm.Get("questionidx"))
	if err != nil {
		return "", err
	}

	quizPageParams.Score, err = parseInt(r.PostForm.Get("score"))
	if err != nil {
		return "", err
	}

	complete := r.PostForm.Get("questioncomplete")
	if complete == "true" {
		quizPageParams.QuestionsComplete = true
		return "", nil
	}

	fmt.Println("quesion idx = ", quizPageParams.QuestionIndex)
	updateScore(questions, quizPageParams, r)

	if quizPageParams.QuestionIndex >= len(questions) {
		quizPageParams.QuestionsComplete = true
		qz := models.UserQuizzes{
			User:  user,
			Quiz:  quiz,
			Total: len(questions),
			Score: quizPageParams.Score,
			Time:  time.Now(),
		}
		err := a.provider.AddScore(qz)
		if err != nil {
			return "", err
		}
		return "", nil
	} else {
		q := questions[quizPageParams.QuestionIndex]
		quizPageParams.Question = q.Question
		quizPageParams.Answers = q.Answers
		quizPageParams.QuestionIndex++
	}

	fmt.Println(quizPageParams)
	return "", nil
}

func updateScore(questions []models.Question, quizPageParams *quizPageParams, r *http.Request) {
	if quizPageParams.QuestionIndex == 0 || quizPageParams.QuestionIndex > len(questions) {
		return
	}
	q := questions[quizPageParams.QuestionIndex-1]
	answer := findAnswer(r.PostForm, len(q.Answers))
	fmt.Println("answer = ", answer)
	if answer == q.CorrectAnswer {
		quizPageParams.Score++
	}
}

func findAnswer(values url.Values, max int) int {
	for i := 0; i < max; i++ {
		s := fmt.Sprintf("%d", i)
		if values.Get(s) == s {
			return i
		}
	}
	return -1
}

func parseInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}

	idx, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse idx to number : %s", value)
	}
	return int(idx), nil
}
