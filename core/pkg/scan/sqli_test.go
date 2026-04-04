package scan

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSQLiScannerIgnoresExistingSQLText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Guide: common SQL syntax mistakes"))
	}))
	defer server.Close()

	scanner := NewSQLiScanner(server.URL + "/?id=1")
	scanner.httpClient = server.Client()
	results, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	for _, result := range results {
		if result.Type == "Error-based SQL Injection" {
			t.Fatalf("expected no SQLi finding for existing SQL text, got %+v", result)
		}
	}
}

func TestSQLiScannerIgnoresUniformlySlowEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5200 * time.Millisecond)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	scanner := NewSQLiScanner(server.URL + "/?id=1")
	scanner.httpClient = server.Client()
	results, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	for _, result := range results {
		if result.Type == "Time-based SQL Injection" {
			t.Fatalf("expected no time-based SQLi on uniformly slow endpoint, got %+v", result)
		}
	}
}

func TestSQLiScannerDetectsInjectedDatabaseError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if strings.Contains(id, "'") {
			_, _ = w.Write([]byte("SQL syntax error near ''"))
			return
		}
		_, _ = w.Write([]byte(fmt.Sprintf("user=%s", id)))
	}))
	defer server.Close()

	scanner := NewSQLiScanner(server.URL + "/?id=1")
	scanner.httpClient = server.Client()
	results, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected SQLi finding for injected database error")
	}
}
