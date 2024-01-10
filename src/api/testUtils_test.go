package main

import (
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"regexp"
	"testing"
)

var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value="{{.CsrfToken}}" />`)

type testServer struct {
	*httptest.Server
}

func newTestApplication(t *testing.T) *Application {
	return &Application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
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
