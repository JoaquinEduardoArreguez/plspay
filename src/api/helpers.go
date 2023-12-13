package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"github.com/justinas/nosurf"
)

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *Application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("template %s does not exist", name))
		return
	}

	buffer := new(bytes.Buffer)

	err := ts.Execute(buffer, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
	}

	buffer.WriteTo(w)
}

func (app *Application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.Flash = app.session.PopString(r, "flash")
	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CsrfToken = nosurf.Token(r)
	return td
}

func (app *Application) authenticatedUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}
