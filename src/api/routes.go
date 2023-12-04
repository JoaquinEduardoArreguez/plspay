package main

import "net/http"

func (app *Application) routes() *http.ServeMux {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", app.home)
	serverMux.HandleFunc("/groups/create", app.createGroup)
	serverMux.HandleFunc("/groups", app.getGroupById)
	serverMux.HandleFunc("/users/create", app.createUser)
	serverMux.HandleFunc("/users", app.getUserById)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	serverMux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return serverMux
}
