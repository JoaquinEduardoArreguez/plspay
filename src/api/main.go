package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"github.com/JoaquinEduardoArreguez/plspay/package/models/repositories"
	"github.com/golangcollege/sessions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Application struct {
	errorLog        *log.Logger
	infoLog         *log.Logger
	session         *sessions.Session
	groupRepository *repositories.GroupRepository
	userRepository  *repositories.UserRepository
	templateCache   map[string]*template.Template
}

func main() {

	// Configs
	serverAddress := flag.String("serverAddress", ":3000", "HTTP network address, host:port .")
	postgresDsn := flag.String("postgresDsn", "", "Postgres database DSN.")
	sessionsKey := flag.String("sessionsKey", "jaodh+pPbnzHbS*+9Pk8qGWhTzbpa@ps", "Session secret key.")
	flag.Parse()

	// Initialization
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	database, initDatabaseError := initDatabase(*postgresDsn)
	if initDatabaseError != nil {
		errorLog.Fatal(initDatabaseError)
	}

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*sessionsKey))
	session.Lifetime = 12 * time.Hour

	app := &Application{
		errorLog:        errorLog,
		infoLog:         infoLog,
		session:         session,
		groupRepository: repositories.NewGroupRepository(database),
		userRepository:  repositories.NewUserRepository(database),
		templateCache:   templateCache,
	}

	infoLog.Printf("Starting server on '%v'", *serverAddress)
	server := &http.Server{Addr: *serverAddress, ErrorLog: errorLog, Handler: app.routes()}
	serveError := server.ListenAndServe()

	errorLog.Fatal(serveError)
}

func initDatabase(dsn string) (*gorm.DB, error) {
	database, openDatabaseError := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
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
