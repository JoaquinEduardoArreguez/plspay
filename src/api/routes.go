package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	serverMux := pat.New()

	serverMux.Get("/", http.HandlerFunc(app.home))

	serverMux.Get("/groups/create", http.HandlerFunc(app.createGroupForm))
	serverMux.Post("/groups/create", http.HandlerFunc(app.createGroup))
	serverMux.Get("/groups/:id", http.HandlerFunc(app.getGroupById))

	serverMux.Get("/users/create", http.HandlerFunc(app.createUserForm))
	serverMux.Post("/users/create", http.HandlerFunc(app.createUser))
	serverMux.Get("/users/:id", http.HandlerFunc(app.getUserById))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	serverMux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(serverMux)
}
