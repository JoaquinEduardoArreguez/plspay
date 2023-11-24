package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	// Configs
	serverAddress := flag.String("serverAddress", ":3000", "HTTP network address, host:port .")
	flag.Parse()

	// Initialization
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := &Application{errorLog: errorLog, infoLog: infoLog}

	infoLog.Printf("Starting server on '%v'", *serverAddress)
	server := &http.Server{Addr: *serverAddress, ErrorLog: errorLog, Handler: app.routes()}
	serveError := server.ListenAndServe()

	errorLog.Fatal(serveError)
}
