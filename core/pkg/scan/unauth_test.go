package scan

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUnauthScannerIgnoresPublicFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/robots.txt":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("User-agent: *\nDisallow: /admin\n"))
		case "/sitemap.xml":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("<urlset></urlset>"))
		case "/.well-known/security.txt":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Contact: mailto:security@example.com\n"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	scanner := NewUnauthScanner(server.URL)
	results, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	for _, result := range results {
		if result.Type == "Sensitive File Exposure" {
			t.Fatalf("expected no sensitive file exposure for public files, got %+v", result)
		}
	}
}

func TestUnauthScannerDetectsSensitiveFileWithIndicators(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/.env":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("APP_KEY=test\nDB_PASSWORD=secret\n"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	scanner := NewUnauthScanner(server.URL)
	results, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	found := false
	for _, result := range results {
		if result.Path == "/.env" && result.Type == "Sensitive File Exposure" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected .env exposure to be detected")
	}
}
