package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brainly/olowek/stats"
)

func TestStatsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/stats", nil)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	s := stats.NewStats()

	rr := httptest.NewRecorder()
	handler := StatsHandler(s)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected HTTP 200 OK, got: '%d'", rr.Code)
	}

	contentType := rr.HeaderMap.Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("API response shuld have content type application/json, got: '%s'", contentType)
	}
}
