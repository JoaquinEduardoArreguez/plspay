package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	serverMux := pat.New()

	// Home
	serverMux.Get("/", dynamicMiddleware.ThenFunc(app.home))

	// Groups
	serverMux.Get("/groups/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createGroupForm))
	serverMux.Post("/groups/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createGroup))
	serverMux.Get("/groups/:id", dynamicMiddleware.ThenFunc(app.getGroupById))

	// Users
	serverMux.Get("/users/create", dynamicMiddleware.ThenFunc(app.createUserForm))
	serverMux.Get("/users/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	serverMux.Get("/users/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	serverMux.Post("/users/create", dynamicMiddleware.ThenFunc(app.createUser))
	serverMux.Post("/users/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	serverMux.Post("/users/login", dynamicMiddleware.ThenFunc(app.loginUser))
	serverMux.Post("/users/logout", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.logoutUser))
	serverMux.Get("/users/:id", dynamicMiddleware.ThenFunc(app.getUserById))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	serverMux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(serverMux)
}
