package scan

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestXSSScannerIgnoresEscapedReflection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("received=" + url.QueryEscape(r.URL.Query().Get("q"))))
	}))
	defer server.Close()

	scanner := NewXSSScanner(server.URL + "/?q=hello")
	scanner.httpClient = server.Client()
	results, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected no XSS finding for escaped reflection, got %+v", results)
	}
}

func TestXSSScannerDetectsRawReflectedPayloadInHTML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(fmt.Sprintf("<html><body>%s</body></html>", r.URL.Query().Get("q"))))
	}))
	defer server.Close()

	scanner := NewXSSScanner(server.URL + "/?q=hello")
	scanner.httpClient = server.Client()
	results, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected reflected XSS finding for raw HTML reflection")
	}
	if results[0].Type != "Reflected XSS" {
		t.Fatalf("expected reflected XSS result, got %+v", results[0])
	}
}

func TestXSSScannerDetectsStoredXSSAfterReadback(t *testing.T) {
	var stored string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/comment":
			stored = r.FormValue("message")
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodGet && r.URL.Path == "/comment":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte("<html><body>" + stored + "</body></html>"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	scanner := NewXSSScanner(server.URL)
	scanner.httpClient = server.Client()
	results, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	found := false
	for _, result := range results {
		if result.Type == "Stored XSS" && strings.Contains(result.URL, "/comment") {
			found = true
		}
	}
	if !found {
		t.Fatal("expected stored XSS finding after payload readback")
	}
}
