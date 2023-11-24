package main

import "net/http"

func (app *Application) routes() *http.ServeMux {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", app.home)
	serverMux.HandleFunc("/groups/create", app.createGroup)
	serverMux.HandleFunc("/groups/list", app.listGroups)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	serverMux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return serverMux
}
