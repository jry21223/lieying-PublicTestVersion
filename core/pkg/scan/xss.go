package scan

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type XSSScanner struct {
	target     string
	results    []XSSResult
	httpClient *http.Client
}

type XSSResult struct {
	URL       string
	Parameter string
	Type      string
	Payload   string
	Evidence  string
	Severity  string
	Confirmed bool
}

func NewXSSScanner(target string) *XSSScanner {
	return &XSSScanner{
		target:  target,
		results: []XSSResult{},
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (xs *XSSScanner) Scan() ([]XSSResult, error) {
	fmt.Println("🔍 开始XSS检测...")

	xs.testReflectedXSS()
	xs.testStoredXSS()

	fmt.Printf("✅ XSS检测完成，发现 %d 个漏洞\n", len(xs.results))
	return xs.results, nil
}

func (xs *XSSScanner) testReflectedXSS() {
	xssPayloads := []string{
		"<script>alert('XSS')</script>",
		"<img src=x onerror=alert('XSS')>",
		"<svg onload=alert('XSS')>",
		"<body onload=alert('XSS')>",
		"<iframe src=javascript:alert('XSS')>",
		"<a href=javascript:alert('XSS')>click</a>",
	}

	parsedURL, err := url.Parse(xs.target)
	if err != nil {
		return
	}

	query := parsedURL.Query()
	for param := range query {
		for _, payload := range xssPayloads {
			testURL := replaceQueryParam(xs.target, param, payload)

			resp, err := xs.httpClient.Get(testURL)
			if err != nil {
				continue
			}

			body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
			resp.Body.Close()
			if err != nil {
				continue
			}

			contentType := strings.ToLower(resp.Header.Get("Content-Type"))
			bodyStr := string(body)
			if !strings.Contains(contentType, "html") {
				continue
			}
			if isEscapedReflection(bodyStr, payload) {
				continue
			}
			if !strings.Contains(bodyStr, payload) {
				continue
			}
			if !containsDangerousHTMLContext(bodyStr, payload) {
				continue
			}

			result := XSSResult{
				URL:       testURL,
				Parameter: param,
				Type:      "Reflected XSS",
				Payload:   payload,
				Evidence:  "Payload 以未转义 HTML/JS 上下文反射",
				Severity:  "High",
				Confirmed: true,
			}
			xs.results = append(xs.results, result)
			fmt.Printf("⚠️  发现XSS: %s (参数: %s)\n", testURL, param)
			return
		}
	}
}

func (xs *XSSScanner) testStoredXSS() {
	storedPayloads := []string{
		"<script>alert('StoredXSS')</script>",
		"<img src=x onerror=alert('StoredXSS')>",
		"<svg onload=alert('StoredXSS')>",
	}

	commonForms := []string{
		"/search",
		"/comment",
		"/feedback",
		"/contact",
		"/register",
		"/signup",
		"/post",
		"/reply",
	}

	parsedURL, err := url.Parse(xs.target)
	if err != nil {
		return
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host

	for _, form := range commonForms {
		for _, payload := range storedPayloads {
			testURL := baseURL + form
			resp, err := xs.httpClient.PostForm(testURL, url.Values{
				"q":       {payload},
				"search":  {payload},
				"name":    {payload},
				"email":   {payload},
				"message": {payload},
			})
			if err != nil {
				continue
			}
			resp.Body.Close()

			readResp, err := xs.httpClient.Get(testURL)
			if err != nil {
				continue
			}
			body, err := io.ReadAll(io.LimitReader(readResp.Body, 1024*1024))
			readResp.Body.Close()
			if err != nil {
				continue
			}

			bodyStr := string(body)
			contentType := strings.ToLower(readResp.Header.Get("Content-Type"))
			if !strings.Contains(contentType, "html") {
				continue
			}
			if isEscapedReflection(bodyStr, payload) {
				continue
			}
			if !strings.Contains(bodyStr, payload) {
				continue
			}
			if !containsDangerousHTMLContext(bodyStr, payload) {
				continue
			}

			result := XSSResult{
				URL:       testURL,
				Parameter: "message",
				Type:      "Stored XSS",
				Payload:   payload,
				Evidence:  "Payload 提交后再次访问页面时以未转义 HTML/JS 上下文出现",
				Severity:  "High",
				Confirmed: true,
			}
			xs.results = append(xs.results, result)
			fmt.Printf("⚠️  发现存储型XSS: %s\n", testURL)
			return
		}
	}
}

func isEscapedReflection(body, payload string) bool {
	escapedForms := []string{
		url.QueryEscape(payload),
		strings.ReplaceAll(strings.ReplaceAll(payload, "<", "&lt;"), ">", "&gt;"),
	}
	for _, escaped := range escapedForms {
		if escaped != payload && strings.Contains(body, escaped) && !strings.Contains(body, payload) {
			return true
		}
	}
	return false
}

func containsDangerousHTMLContext(body, payload string) bool {
	bodyLower := strings.ToLower(body)
	payloadLower := strings.ToLower(payload)
	dangerousFragments := []string{
		"<script",
		"onerror=",
		"onerror =",
		"onload=",
		"onload =",
		"javascript:",
	}
	for _, fragment := range dangerousFragments {
		if strings.Contains(payloadLower, fragment) && strings.Contains(bodyLower, fragment) {
			return true
		}
	}
	return false
}

func (xs *XSSScanner) GetResults() []XSSResult {
	return xs.results
}
