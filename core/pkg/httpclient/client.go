package httpclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	httpClient *http.Client
	limiter    *rate.Limiter
	retryCount int
	baseURL    string
	headers    http.Header
	cookies    []*http.Cookie
}

type ClientConfig struct {
	ProxyURL     string
	QPS          int
	Timeout      time.Duration
	FollowRedirects bool
	SkipTLSVerify bool
	RetryCount   int
}

func NewClient(config ClientConfig) (*Client, error) {
	transport := &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	if config.ProxyURL != "" {
		proxyURL, err := url.Parse(config.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	if config.SkipTLSVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	if !config.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	qps := config.QPS
	if qps <= 0 {
		qps = 10
	}
	limiter := rate.NewLimiter(rate.Limit(qps), qps)

	retryCount := config.RetryCount
	if retryCount <= 0 {
		retryCount = 3
	}

	return &Client{
		httpClient: client,
		limiter:    limiter,
		retryCount: retryCount,
		headers:    make(http.Header),
		cookies:    make([]*http.Cookie, 0),
	}, nil
}

func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = strings.TrimSuffix(baseURL, "/")
}

func (c *Client) SetHeader(key, value string) {
	c.headers.Set(key, value)
}

func (c *Client) AddHeader(key, value string) {
	c.headers.Add(key, value)
}

func (c *Client) SetCookie(cookie *http.Cookie) {
	c.cookies = append(c.cookies, cookie)
}

func (c *Client) SetUserAgent(ua string) {
	c.SetHeader("User-Agent", ua)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	if err := c.limiter.Wait(ctx); err != nil {
		return nil, err
	}

	for key, values := range c.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	var resp *http.Response
	var err error

	for i := 0; i < c.retryCount; i++ {
		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode < 500 {
			break
		}
		if i < c.retryCount-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}

	return resp, err
}

func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	url := c.buildURL(path)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Post(ctx context.Context, path string, contentType string, body []byte) (*http.Response, error) {
	url := c.buildURL(path)
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

func (c *Client) Put(ctx context.Context, path string, contentType string, body []byte) (*http.Response, error) {
	url := c.buildURL(path)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	url := c.buildURL(path)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) buildURL(path string) string {
	if c.baseURL == "" {
		return path
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return c.baseURL + path
}

func RandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
	return userAgents[time.Now().UnixNano()%int64(len(userAgents))]
}
