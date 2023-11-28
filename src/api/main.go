package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	// Configs
	serverAddress := flag.String("serverAddress", ":3000", "HTTP network address, host:port .")
	postgresDsn := flag.String("postgresDsn", "", "Postgres database DSN.")
	flag.Parse()

	// Initialization
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	_, initDatabaseError := initDatabase(*postgresDsn)
	if initDatabaseError != nil {
		errorLog.Fatal(initDatabaseError)
	}

	app := &Application{errorLog: errorLog, infoLog: infoLog}

	infoLog.Printf("Starting server on '%v'", *serverAddress)
	server := &http.Server{Addr: *serverAddress, ErrorLog: errorLog, Handler: app.routes()}
	serveError := server.ListenAndServe()

	errorLog.Fatal(serveError)
}

func initDatabase(dsn string) (*gorm.DB, error) {
	database, openDatabaseError := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if openDatabaseError != nil {
		return nil, openDatabaseError
	}

	migrateDatabaseError := database.AutoMigrate(
		&models.Group{},
		&models.User{},
		&models.Expense{},
		&models.Transaction{},
	)
	if migrateDatabaseError != nil {
		return nil, migrateDatabaseError
	}

	return database, nil
}
