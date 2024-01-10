package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_secureHeaders(t *testing.T) {
	t.Parallel()

	responseRecorder := httptest.NewRecorder()
	requestMock, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	nextHttpHandlerMock := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })

	secureHeaders(nextHttpHandlerMock).ServeHTTP(responseRecorder, requestMock)

	recorderResult := responseRecorder.Result()

	frameOptions := recorderResult.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}

	xssProtection := recorderResult.Header.Get("X-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf("want %q; got %q", "1; mode=block", xssProtection)
	}

	if recorderResult.StatusCode != http.StatusOK {
		t.Errorf("want %q; got %q", http.StatusOK, recorderResult.StatusCode)
	}

	defer recorderResult.Body.Close()
	body, err := io.ReadAll(responseRecorder.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
