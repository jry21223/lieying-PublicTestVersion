package scan

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
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

func TestParseUploadFormRejectsCrossOriginAction(t *testing.T) {
	_, ok := parseUploadForm("https://example.com/upload", `<form action="https://evil.example/upload" method="post" enctype="multipart/form-data"><input type="file" name="file"></form>`)
	if ok {
		t.Fatal("expected cross-origin form action to be skipped")
	}
}

func TestParseUploadFormAcceptsFlexibleFileInputAttributes(t *testing.T) {
	cases := []string{
		`<form action="/upload" method="post" enctype="multipart/form-data"><input type="file" name="f"></form>`,
		`<form action="/upload" method="post" enctype="multipart/form-data"><INPUT TYPE = "FILE" NAME = "f"></form>`,
		`<form action="/upload" method="post" enctype="multipart/form-data"><input name='f'   type = file></form>`,
	}

	for _, body := range cases {
		form, ok := parseUploadForm("https://example.com/upload", body)
		if !ok {
			t.Fatalf("expected upload form to be parsed for %q", body)
		}
		if form.inputName != "f" {
			t.Fatalf("expected input name f, got %+v", form)
		}
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
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/uploads/shell-"):
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(executionMarkerFromPath(r.URL.Path)))
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
	if !results[0].Confirmed {
		t.Fatalf("expected reachable uploaded file to be confirmed, got %+v", results[0])
	}
}

func TestIsSameOriginTreatsExplicitDefaultPortAsSameOrigin(t *testing.T) {
	if !isSameOrigin("https://example.com/upload", "https://example.com:443/file.php") {
		t.Fatal("expected explicit default https port to match same origin")
	}
	if !isSameOrigin("http://example.com/upload", "http://example.com:80/file.php") {
		t.Fatal("expected explicit default http port to match same origin")
	}
}

func executionMarkerFromPath(uploadedPath string) string {
	sum := md5.Sum([]byte(filenameMarkerFromPath(uploadedPath)))
	return fmt.Sprintf("%x", sum)
}

func filenameMarkerFromPath(uploadedPath string) string {
	base := path.Base(uploadedPath)
	base = strings.TrimSuffix(base, path.Ext(base))
	return strings.Replace(base, "shell-", "lieying-upload-marker-", 1)
}

func TestUploadScannerDoesNotConfirmWhenUploadedFileIsNotReachable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/upload":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<html><body><form action="/api/upload" method="post" enctype="multipart/form-data"><input type="file" name="avatar"></form></body></html>`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/upload":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"path":"/uploads/shell.php","status":"success"}`))
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
		t.Fatalf("expected no confirmed upload finding when uploaded file is unreachable, got %+v", results)
	}
}
