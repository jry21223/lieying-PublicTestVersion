package scan

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestReplaceQueryParamPreservesRepeatedParameters(t *testing.T) {
	updated := replaceQueryParam("http://example.com/?id=1&id=2&note=a%2Bb", "id", "'")
	parsed, err := url.Parse(updated)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	values := parsed.Query()["id"]
	if len(values) != 2 {
		t.Fatalf("expected duplicate id parameters to be preserved, got %q", parsed.RawQuery)
	}
	if values[0] != "'" || values[1] != "'" {
		t.Fatalf("expected all id parameters to change, got %q", parsed.RawQuery)
	}
}

func TestReplaceQueryParamHandlesEncodedValues(t *testing.T) {
	updated := replaceQueryParam("http://example.com/?q=a%2Bb&name=John+Doe", "q", "<script>alert(1)</script>")
	parsed, err := url.Parse(updated)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if got := parsed.Query().Get("q"); got != "<script>alert(1)</script>" {
		t.Fatalf("expected encoded parameter to be replaced, got %q", got)
	}
	if got := parsed.Query().Get("name"); got != "John Doe" {
		t.Fatalf("expected unrelated parameter to stay unchanged, got %q", got)
	}
}
