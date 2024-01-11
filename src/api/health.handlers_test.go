package main

import (
	"net/http"
	"testing"
)

func Test_ping(t *testing.T) {
	appMock := newTestApplication(t)

	tlsServerMock := newTestServer(t, appMock.routes())
	defer tlsServerMock.Close()

	code, _, body := tlsServerMock.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
