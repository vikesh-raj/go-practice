package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
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
    <title>Splitwise</title>
	<style>
    .err {
        background-color: #FFCCCC;
    }
	</style>
</head>
<body>
<h1>Splitwise</h1>
<div>Version : {{.Version}}</div><br/>
{{if .ErrorMessage}}
<div class="err">{{.ErrorMessage}}</div><br/>
{{end}}
{{if .User}}
<ul>
<li><a href="/view?user={{.User}}">View Statement</a></li>
<li><a href="/add?user={{.User}}">Add Transaction</a></li>
<li><a href="/settle?user={{.User}}">Settle</a></li>
<li><a href="/">Logout</a></li>
</ul>
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
    <title>Splitwise</title>
	<style>
    .err {
        background-color: #FFCCCC;
    }
	</style>
</head>
<body>
<div><a href="/?user={{.User}}">Home</a></div>
<h1>Statement</h1>
{{if .User}}
Fetching Statements
{{else}}
<div class="err"><a href="/">Please login</a></div>
{{end}}
</body>
</html>
`

const addHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <title>Splitwise</title>
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
<h1>Add Transaction</h1>
{{if .User}}
	{{if .Message}}
		<div class="msg">{{.Message}}<br/></div><br/>
	{{end}}
	{{if .Error}}
		<div class="err">{{.Error}}<br/></div><br/>
	{{end}}

	<form action="" method="post">
	<input type="hidden" id="user" value="{{.User}}">
	<label for="other_user">Other User:</label>
	<input type="text" id="other_user" name="other_user" value="{{.OtherUser}}"><br><br>
	<ul>
	<li>
		<label for="total">Total:</label>
		<input type="number" id="total" name="total">
		<label for="percentage">Your Percentage:</label>
		<input type="number" id="percentage" name="percentage", value="50"><br><br>
	</li>
	<li>
		<label for="amount">Amount:</label>
		<input type="number" id="amount" name="amount"><br><br>
	</li>
	</ul>
	<label for="remarks">Remarks:</label>
	<input type="text" id="remarks" name="remarks" value="{{.Remarks}}"><br><br>
	<input type="submit" value="Submit">
	</form>

{{else}}
<a href="/">Please login</a>
{{end}}
</body>
</html>
`

const settleHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <title>Splitwise</title>
</head>
<style>
	.msg {
		background-color: #E0DA96;
	}
	.err {
        background-color: #FFCCCC;
    }
</style>
<body>
<div><a href="/?user={{.User}}">Home</a></div>
<h1>Settle Amount</h1>
{{if .User}}
	{{if .Message}}
		<div class="msg">{{.Message}}<br/></div><br/>
	{{end}}
	{{if .Error}}
		<div class="err">{{.Error}}<br/></div><br/>
	{{end}}

	<form action="" method="post">
	<input type="hidden" id="user" value="{{.User}}">
	<label for="other_user">Other User:</label>
	<input type="text" id="other_user" name="other_user" value="{{.OtherUser}}"><br><br>
	{{if .SettleAmount}}
		<input type="hidden" id="settle" name="settle" value="true">
		<label for="remarks">Remarks:</label>
		<input type="text" id="remarks" name="remarks" value="{{.Remarks}}"><br><br>
		<div>Settle Amount : {{.SettleAmount}} </div>
		<input type="submit" value="Settle">
	{{else}}
		<input type="submit" value="Fetch">
	{{end}}
	</form>
{{else}}
<a href="/">Please login</a>
{{end}}
</body>
</html>
`

func parseTemplate(page string) *template.Template {
	return template.Must(template.New("").Parse(page))
}

var indexTemplate = parseTemplate(indexHTML)

type indexPageParams struct {
	Version      string
	ErrorMessage string
	Opts         Opts
	User         string
}

func (a *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	indexPageParams := indexPageParams{
		Version: version,
		Opts:    a.Opts,
		User:    r.URL.Query().Get("user"),
	}
	if a.lastError != nil {
		indexPageParams.ErrorMessage = a.lastError.Error()
	}

	indexTemplate.Execute(w, &indexPageParams)
}

var viewTemplate = parseTemplate(viewHTML)

type viewPageParams struct {
	User string
}

func (a *Application) handleView(w http.ResponseWriter, r *http.Request) {
	viewPageParams := viewPageParams{
		User: r.URL.Query().Get("user"),
	}

	viewTemplate.Execute(w, &viewPageParams)
}

var addTemplate = parseTemplate(addHTML)

type addPageParams struct {
	User      string
	Message   string
	Error     string
	OtherUser string
	Remarks   string
}

func (a *Application) handleAdd(w http.ResponseWriter, r *http.Request) {
	addPageParams := addPageParams{
		User: r.URL.Query().Get("user"),
	}

	addTemplate.Execute(w, &addPageParams)
}

func (a *Application) handleAddPost(w http.ResponseWriter, r *http.Request) {

	addPageParams := addPageParams{
		User: r.URL.Query().Get("user"),
	}

	msg, err := a.handleAddPostRequest(r)
	if err != nil {
		addPageParams.Error = "Error : " + err.Error()
	}
	addPageParams.Message = msg
	addTemplate.Execute(w, &addPageParams)
}

func (a *Application) handleAddPostRequest(r *http.Request) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", err
	}

	otherUser := r.PostForm.Get("other_user")

	if otherUser == "" {
		return "", fmt.Errorf("specify other user")
	}

	amountStr := r.PostForm.Get("amount")
	totalStr := r.PostForm.Get("total")
	percentageStr := r.PostForm.Get("percentage")
	amount := 0.0

	if amountStr == "" {
		if totalStr == "" {
			return "", fmt.Errorf("specify amount or total")
		}
		if percentageStr == "" {
			return "", fmt.Errorf("specify percentage")
		}

		pct, err := strconv.ParseFloat(strings.TrimSpace(percentageStr), 64)
		if err != nil {
			return "", fmt.Errorf("unable to parse percentage to number : %s", percentageStr)
		}

		if pct > 100 {
			return "", fmt.Errorf("percentage greater than 100 : %.2f", pct)
		}

		if pct < 0 {
			return "", fmt.Errorf("percentage cannot be negative : %.2f", pct)
		}

		total, err := strconv.ParseFloat(strings.TrimSpace(totalStr), 64)
		if err != nil {
			return "", fmt.Errorf("unable to parse total to number : %s", totalStr)
		}

		amount = (total * pct) / 100.0

	} else {
		amount, err = strconv.ParseFloat(strings.TrimSpace(amountStr), 64)
		if err != nil {
			return "", fmt.Errorf("unable to parse percentage to number : %s", amountStr)
		}
	}

	if amount > 0 {

	}
	return "Added successfully", nil
}

var settleTemplate = parseTemplate(settleHTML)

type settlePageParams struct {
	User         string
	Message      string
	Error        string
	OtherUser    string
	Remarks      string
	SettleAmount float64
}

func (a *Application) handleSettle(w http.ResponseWriter, r *http.Request) {
	settlePageParams := settlePageParams{
		User: r.URL.Query().Get("user"),
	}

	settleTemplate.Execute(w, &settlePageParams)
}

func (a *Application) handleSettlePost(w http.ResponseWriter, r *http.Request) {

	settlePageParams := settlePageParams{
		User: r.URL.Query().Get("user"),
	}

	msg, err := a.handleSettlePostRequest(&settlePageParams, r)
	if err != nil {
		settlePageParams.Error = "Error : " + err.Error()
	}
	settlePageParams.Message = msg
	settleTemplate.Execute(w, &settlePageParams)
}

func (a *Application) handleSettlePostRequest(sp *settlePageParams, r *http.Request) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", err
	}

	otherUser := r.PostForm.Get("other_user")
	if otherUser == "" {
		return "", fmt.Errorf("specify other user")
	}
	sp.OtherUser = otherUser

	settle := r.PostForm.Get("settle")
	if settle != "true" {
		sp.SettleAmount = 25.0
	} else {
		return "Settle Successful", nil
	}
	return "", nil
}
