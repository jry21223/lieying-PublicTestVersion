package scan

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUploadScannerDoesNotFlagFormPresenceAlone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/upload":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<html><body><form action="/upload" method="post" enctype="multipart/form-data"><input type="file" name="file"></form></body></html>`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	scanner := NewUploadScanner(server.URL)
	scanner.httpClient = server.Client()
	results, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected no upload finding for form presence alone, got %+v", results)
	}
}

func TestUploadScannerUsesFormActionAndInputName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/upload":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<html><body><form action="/api/upload" method="post" enctype="multipart/form-data"><input type="file" name="avatar"></form></body></html>`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/upload":
			file, header, err := r.FormFile("avatar")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer file.Close()
			if strings.HasSuffix(header.Filename, ".php") {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(fmt.Sprintf(`{"path":"/uploads/%s","status":"success"}`, header.Filename)))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	scanner := NewUploadScanner(server.URL)
	scanner.httpClient = server.Client()
	results, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected dangerous upload acceptance to be detected")
	}
	if results[0].InputName != "avatar" {
		t.Fatalf("expected scanner to use form input name, got %+v", results[0])
	}
	if !strings.Contains(results[0].FormAction, "/api/upload") {
		t.Fatalf("expected scanner to use form action, got %+v", results[0])
	}
}
