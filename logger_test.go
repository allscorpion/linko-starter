package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func printErrorMessage(expected any, got any) string {
	return fmt.Sprintf("\nexpected log message: \n\t%sGot:\n\t%s", expected, got)
}

func Test_requestLogger(t *testing.T) {
	logBuffer := &bytes.Buffer{}

	logger := slog.New(slog.NewTextHandler(logBuffer, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Time(slog.TimeKey, time.Date(2023, 10, 1, 12, 34, 57, 0, time.UTC))
			}
			return a
		},
	}))

	requestLoggerMiddleware := requestLogger(logger)
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	loggedHandler := requestLoggerMiddleware(dummyHandler)

	req := httptest.NewRequest("GET", "http://localhost:8080/test?foo=bar", nil)

	rr := httptest.NewRecorder()

	loggedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	expectedLog := `time=2023-10-01T12:34:57.000Z level=INFO msg="Served request" method=GET path=/test client_ip=192.0.2.x
`
	if logBuffer.String() != expectedLog {
		t.Errorf("Expected log output:\n%s\nGot:\n%s", expectedLog, logBuffer.String())
	}
}
