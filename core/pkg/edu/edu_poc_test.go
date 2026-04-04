package edu

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEduPOCScannerDoesNotFlagGenericReachablePages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/jwgl":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("校园统一服务 jwgl portal"))
		case "/jwgl/admin", "/jwgl/ueditor/":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("统一身份认证登录"))
		default:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("normal page"))
		}
	}))
	defer server.Close()

	scanner := NewEduPOCScanner(server.URL)
	scanner.httpClient = server.Client()
	results, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected no edu vuln findings for generic reachable pages, got %+v", results)
	}
}

func TestEduPOCScannerDetectsConfirmedExposure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/jwgl":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("青果教务管理系统"))
		case "/jwgl/admin":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("<html><title>Admin Dashboard</title><body>系统管理 用户管理 权限配置</body></html>"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	scanner := NewEduPOCScanner(server.URL)
	scanner.httpClient = server.Client()
	results, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected confirmed edu exposure to be detected")
	}
	if results[0].URL != fmt.Sprintf("%s/jwgl/admin", server.URL) {
		t.Fatalf("unexpected result URL: %+v", results[0])
	}
}
