package main

import (
	"html/template"
	"net/http"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	templateFilePaths := []string{
		"./ui/html/home.page.template.html",
		"./ui/html/base.layout.template.html",
		"./ui/html/footer.partial.template.html",
	}

	templateSet, templateParseFilesError := template.ParseFiles(templateFilePaths...)

	if templateParseFilesError != nil {
		app.serverError(w, templateParseFilesError)
		return
	}

	templateSetExecuteError := templateSet.Execute(w, nil)

	if templateSetExecuteError != nil {
		app.serverError(w, templateSetExecuteError)
	}

	w.Write([]byte("PlsPay"))
}

func (app *Application) createGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new group"))
}

func (app *Application) listGroups(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List all groups"))
}
