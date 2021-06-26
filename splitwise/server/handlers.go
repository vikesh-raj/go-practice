package server

import (
	"net/http"
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
    <title>Geocode Proxy</title>
</head>
<body>
<h1>Geocode Proxy</h1>
<div>Version : {{.Version}}</div><br/>
{{if .ErrorMessage}}
<div>{{.ErrorMessage}}</div><br/>
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
}

func (a *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	indexPageParams := indexPageParams{
		Version: version,
		Opts:    a.Opts,
	}
	if a.lastError != nil {
		indexPageParams.ErrorMessage = a.lastError.Error()
	}

	indexTemplate.Execute(w, &indexPageParams)
}
