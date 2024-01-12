package main

import (
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/JoaquinEduardoArreguez/plspay/package/models"
	"github.com/JoaquinEduardoArreguez/plspay/package/services"
	"github.com/golangcollege/sessions"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value="{{.CsrfToken}}" />`)

type testServer struct {
	*httptest.Server
}

func newTestApplication(t *testing.T) *Application {
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	session := sessions.New([]byte("jaodh+pPbnzHbS*+9Pk8qGWhTzbpa@ps"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	database, initDatabaseError := initTestDatabase(t)
	if initDatabaseError != nil {
		t.Fatal(initDatabaseError)
	}

	return &Application{
		errorLog:       log.New(io.Discard, "", 0),
		infoLog:        log.New(io.Discard, "", 0),
		session:        session,
		templateCache:  templateCache,
		groupService:   services.NewGroupService(database),
		expenseService: services.NewExpenseService(database),
		userService:    services.NewUserService(database),
	}
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	tlsServer := httptest.NewTLSServer(h)

	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	tlsServer.Client().Jar = cookieJar

	// This function is called after a 3XX response is received by the client,
	// returning http.ErrUseLastResponse error forces it to immediately return
	// the received response.
	tlsServer.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{tlsServer}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal("Error on request", err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal("Error reading body", err)
	}
	return rs.StatusCode, rs.Header, body
}

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, []byte) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}

func extractCSRFToken(t *testing.T, body []byte) string {
	matches := csrfTokenRX.FindSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
	return html.UnescapeString(string(matches[1]))
}

func initTestDatabase(t *testing.T) (*gorm.DB, error) {
	testDsn := "host=localhost user=postgresTest password=adminTest dbname=plspayTest port=5433 sslmode=disable"
	database, openDatabaseError := gorm.Open(postgres.Open(testDsn), &gorm.Config{
		TranslateError: true,
	})
	if openDatabaseError != nil {
		return nil, openDatabaseError
	}

	database.Migrator().DropTable()

	migrateDatabaseError := database.AutoMigrate(
		&models.Group{},
		&models.User{},
		&models.Expense{},
		&models.Transaction{},
		&models.Balance{},
	)
	if migrateDatabaseError != nil {
		return nil, migrateDatabaseError
	}

	return database, nil
}
