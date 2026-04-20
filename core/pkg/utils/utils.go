package utils

import (
	"io"
	"net/http"
	"strings"
	"time"
)

// Severity levels as constants
const (
	SeverityCritical = "Critical"
	SeverityHigh     = "High"
	SeverityMedium   = "Medium"
	SeverityLow      = "Low"
)

// Shared HTTP client with connection pooling
var sharedClient = &http.Client{
	Timeout: 15 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 3 {
			return http.ErrUseLastResponse
		}
		return nil
	},
}

// FetchBody fetches URL content with size limit, returns body and status code
func FetchBody(rawURL string, maxSize int64) (string, int) {
	resp, err := sharedClient.Get(rawURL)
	if err != nil {
		return "", 0
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxSize))
	return string(body), resp.StatusCode
}

// ContainsAny checks if body contains any of the keywords
func ContainsAny(body string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(body, kw) {
			return true
		}
	}
	return false
}

// ContainsAnyLower checks if body contains any keyword (case-insensitive)
func ContainsAnyLower(body string, keywords []string) bool {
	bodyLower := strings.ToLower(body)
	for _, kw := range keywords {
		if strings.Contains(bodyLower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

// ExtractSnippet extracts text around first matched pattern
func ExtractSnippet(body string, patterns []string, before, after int) string {
	for _, p := range patterns {
		idx := strings.Index(body, p)
		if idx >= 0 {
			start := idx - before
			if start < 0 {
				start = 0
			}
			end := idx + len(p) + after
			if end > len(body) {
				end = len(body)
			}
			return "..." + strings.TrimSpace(body[start:end]) + "..."
		}
	}
	return ""
}

// HasProtocol checks if target has http/https prefix
func HasProtocol(target string) bool {
	return strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://")
}