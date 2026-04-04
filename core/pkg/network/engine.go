package network

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type NetworkEngine struct {
	httpClient  *http.Client
	proxyURL    *url.URL
	proxyAuth   *ProxyAuth
	rateLimiter *RateLimiter
	userAgents  []string
	userAgent   string
	timeout     time.Duration
	mu          sync.RWMutex
	config      *EngineConfig
}

type EngineConfig struct {
	Proxy           string            `json:"proxy"`
	ProxyAuth       *ProxyAuth        `json:"proxy_auth"`
	Concurrency     int               `json:"concurrency"`
	QPS             float64           `json:"qps"`
	Timeout         int               `json:"timeout"`
	UserAgent       string            `json:"user_agent"`
	Insecure        bool              `json:"insecure"`
	FollowRedirects bool              `json:"follow_redirects"`
	Headers         map[string]string `json:"headers"`
}

type ProxyAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Request struct {
	Method          string            `json:"method"`
	URL             string            `json:"url"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	Timeout         int               `json:"timeout"`
	FollowRedirects bool              `json:"follow_redirects"`
}

type Response struct {
	ID       string            `json:"id"`
	URL      string            `json:"url,omitempty"`
	Status   int               `json:"status"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
	Duration int64             `json:"duration"`
	Error    string            `json:"error,omitempty"`
}

type PortScanResult struct {
	Port    int    `json:"port"`
	Open    bool   `json:"open"`
	Service string `json:"service"`
	Banner  string `json:"banner"`
}

func NewNetworkEngine(config *EngineConfig) (*NetworkEngine, error) {
	if config == nil {
		config = &EngineConfig{
			Concurrency:     10,
			QPS:             50,
			Timeout:         30,
			UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			FollowRedirects: true,
			Insecure:        false,
			Headers:         make(map[string]string),
		}
	}

	var proxyURL *url.URL
	if config.Proxy != "" {
		var err error
		proxyURL, err = url.Parse(config.Proxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}
	}

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			if proxyURL != nil {
				return proxyURL, nil
			}
			return http.ProxyFromEnvironment(req)
		},
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.Insecure,
		},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout: time.Duration(config.Timeout) * time.Second,
		}).DialContext,
	}

	if !config.FollowRedirects {
		transport.DisableKeepAlives = false
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(config.Timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !config.FollowRedirects {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
	}

	engine := &NetworkEngine{
		httpClient:  httpClient,
		proxyURL:    proxyURL,
		proxyAuth:   config.ProxyAuth,
		rateLimiter: NewRateLimiter(config.QPS),
		userAgents:  userAgents,
		userAgent:   config.UserAgent,
		timeout:     time.Duration(config.Timeout) * time.Second,
		config:      config,
	}

	return engine, nil
}

func (e *NetworkEngine) DoRequest(ctx context.Context, req *Request) (*Response, error) {
	startTime := time.Now()

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, strings.NewReader(req.Body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range e.config.Headers {
		httpReq.Header.Set(key, value)
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	if httpReq.Header.Get("User-Agent") == "" {
		if e.userAgent != "" {
			httpReq.Header.Set("User-Agent", e.userAgent)
		}
	}

	if httpReq.Header.Get("Accept") == "" {
		httpReq.Header.Set("Accept", "*/*")
	}

	if httpReq.Header.Get("Accept-Language") == "" {
		httpReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	}

	timeout := e.timeout
	if req.Timeout > 0 {
		timeout = time.Duration(req.Timeout) * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := e.rateLimiter.Wait(ctx, req.URL); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	resp, err := e.httpClient.Do(httpReq.WithContext(ctx))
	duration := time.Since(startTime).Milliseconds()

	response := &Response{
		ID:       uuid.New().String(),
		Duration: duration,
	}

	if err != nil {
		response.Error = err.Error()
		return response, err
	}
	defer resp.Body.Close()

	response.Status = resp.StatusCode
	response.Headers = make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			response.Headers[key] = values[0]
		}
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		response.Error = fmt.Sprintf("failed to read response body: %v", err)
		return response, err
	}
	response.Body = string(body)

	return response, nil
}

func (e *NetworkEngine) Get(ctx context.Context, targetURL string) (*Response, error) {
	return e.DoRequest(ctx, &Request{
		Method: "GET",
		URL:    targetURL,
	})
}

func (e *NetworkEngine) Post(ctx context.Context, targetURL, body string) (*Response, error) {
	return e.DoRequest(ctx, &Request{
		Method: "POST",
		URL:    targetURL,
		Body:   body,
	})
}

func (e *NetworkEngine) Head(ctx context.Context, targetURL string) (*Response, error) {
	return e.DoRequest(ctx, &Request{
		Method: "HEAD",
		URL:    targetURL,
	})
}

func (e *NetworkEngine) Options(ctx context.Context, targetURL string) (*Response, error) {
	return e.DoRequest(ctx, &Request{
		Method: "OPTIONS",
		URL:    targetURL,
	})
}

func (e *NetworkEngine) ScanPort(ip string, port int, timeout time.Duration) *PortScanResult {
	result := &PortScanResult{
		Port: port,
		Open: false,
	}

	addr := net.JoinHostPort(ip, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return result
	}
	defer conn.Close()

	result.Open = true
	result.Service = guessService(port)

	if timeout > 100*time.Millisecond {
		conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		buf := make([]byte, 1024)
		if n, _ := conn.Read(buf); n > 0 {
			result.Banner = strings.TrimSpace(string(buf[:n]))
		}
	}

	return result
}

func (e *NetworkEngine) ScanPorts(ip string, ports []int, concurrency int) []*PortScanResult {
	results := make([]*PortScanResult, len(ports))
	semaphore := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, port := range ports {
		wg.Add(1)
		go func(idx int, p int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			results[idx] = e.ScanPort(ip, p, 1*time.Second)
		}(i, port)
	}

	wg.Wait()
	return results
}

func (e *NetworkEngine) Resolve(host string) ([]net.IP, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, fmt.Errorf("DNS resolution failed: %w", err)
	}
	return ips, nil
}

func (e *NetworkEngine) Download(ctx context.Context, downloadURL, destPath string) error {
	resp, err := e.Get(ctx, downloadURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	if resp.Status < 200 || resp.Status >= 300 {
		return fmt.Errorf("download returned status %d", resp.Status)
	}

	return nil
}

func (e *NetworkEngine) SetProxy(proxyURL string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if proxyURL == "" {
		e.proxyURL = nil
		return nil
	}

	parsed, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %w", err)
	}

	e.proxyURL = parsed
	return nil
}

func (e *NetworkEngine) GetProxy() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.proxyURL == nil {
		return ""
	}
	return e.proxyURL.String()
}

func (e *NetworkEngine) SetQPS(qps float64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rateLimiter.SetQPS(qps)
}

func (e *NetworkEngine) SetConcurrency(limit int) {
}

func (e *NetworkEngine) SetTimeout(timeout int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.timeout = time.Duration(timeout) * time.Second
}

func (e *NetworkEngine) SetUserAgent(ua string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.userAgent = ua
}

func (e *NetworkEngine) RandomizeUserAgent() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.userAgents) > 0 {
		e.userAgent = e.userAgents[time.Now().UnixNano()%int64(len(e.userAgents))]
	}
}

func (e *NetworkEngine) SetHeader(key, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.config.Headers[key] = value
}

func (e *NetworkEngine) SetInsecure(insecure bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.config.Insecure = insecure
}

func (e *NetworkEngine) GetConfig() *EngineConfig {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.config
}

func (e *NetworkEngine) BatchRequest(ctx context.Context, requests []*Request, concurrency int) []*Response {
	results := make([]*Response, len(requests))
	semaphore := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, req := range requests {
		wg.Add(1)
		go func(idx int, r *Request) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			resp, err := e.DoRequest(ctx, r)
			if err != nil {
				results[idx] = &Response{
					ID:    uuid.New().String(),
					Error: err.Error(),
					URL:   r.URL,
				}
			} else {
				results[idx] = resp
			}
		}(i, req)
	}

	wg.Wait()
	return results
}

func guessService(port int) string {
	commonPorts := map[int]string{
		20:    "ftp-data",
		21:    "ftp",
		22:    "ssh",
		23:    "telnet",
		25:    "smtp",
		53:    "dns",
		80:    "http",
		110:   "pop3",
		143:   "imap",
		443:   "https",
		465:   "smtps",
		587:   "smtp-submission",
		993:   "imaps",
		995:   "pop3s",
		1433:  "mssql",
		1521:  "oracle",
		3306:  "mysql",
		3389:  "rdp",
		5432:  "postgresql",
		5900:  "vnc",
		6379:  "redis",
		8080:  "http-proxy",
		8443:  "https-alt",
		27017: "mongodb",
	}

	if service, ok := commonPorts[port]; ok {
		return service
	}
	return "unknown"
}

func ExtractVars(raw string) []string {
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	matches := re.FindAllStringSubmatch(raw, -1)

	vars := make([]string, 0, len(matches))
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 && !seen[match[1]] {
			vars = append(vars, match[1])
			seen[match[1]] = true
		}
	}

	return vars
}

func PayloadGeneration(raw string, variables map[string][]string) []string {
	var results []string

	for varName, varValues := range variables {
		var newResults []string
		if len(results) == 0 {
			for _, value := range varValues {
				newResults = append(newResults, strings.ReplaceAll(raw, "{{"+varName+"}}", value))
			}
		} else {
			for _, existing := range results {
				for _, value := range varValues {
					newResults = append(newResults, strings.ReplaceAll(existing, "{{"+varName+"}}", value))
				}
			}
		}
		results = newResults
	}

	if len(results) == 0 {
		return []string{raw}
	}

	return results
}

type IntruderPayload struct {
	Position int    `json:"position"`
	Value    string `json:"value"`
}

type IntruderResult struct {
	Request  *Request  `json:"request"`
	Response *Response `json:"response"`
	Payload  string    `json:"payload"`
}

func (e *NetworkEngine) Intruder(ctx context.Context, baseReq *Request, positions []int, payloads [][]string) []*IntruderResult {
	results := make([]*IntruderResult, 0)

	for _, payloadList := range payloads {
		for _, payload := range payloadList {
			reqCopy := *baseReq
			body := baseReq.Body

			for _, pos := range positions {
				if pos < len(body) {
					body = body[:pos] + payload + body[pos:]
				}
			}

			reqCopy.Body = body
			resp, err := e.DoRequest(ctx, &reqCopy)
			if err != nil {
				results = append(results, &IntruderResult{
					Request:  &reqCopy,
					Response: &Response{Error: err.Error()},
					Payload:  payload,
				})
			} else {
				results = append(results, &IntruderResult{
					Request:  &reqCopy,
					Response: resp,
					Payload:  payload,
				})
			}
		}
	}

	return results
}

func ParseJSONResponse(body string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (e *NetworkEngine) Clone() *NetworkEngine {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return &NetworkEngine{
		httpClient:  e.httpClient,
		proxyURL:    e.proxyURL,
		proxyAuth:   e.proxyAuth,
		rateLimiter: e.rateLimiter,
		userAgents:  e.userAgents,
		userAgent:   e.userAgent,
		timeout:     e.timeout,
		config:      e.config,
	}
}
