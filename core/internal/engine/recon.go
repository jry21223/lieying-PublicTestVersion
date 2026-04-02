package engine

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ReconResult struct {
	Target       string           `json:"target"`
	Subdomains   []string         `json:"subdomains"`
	Ports        []PortInfo       `json:"ports"`
	Fingerprints []FingerprintInfo `json:"fingerprints"`
	Directories  []DirectoryInfo   `json:"directories"`
	Timestamp    time.Time        `json:"timestamp"`
}

type PortInfo struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Service  string `json:"service"`
	Version  string `json:"version"`
	State    string `json:"state"`
}

type FingerprintInfo struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Confidence int    `json:"confidence"`
}

type DirectoryInfo struct {
	Path    string `json:"path"`
	Status  int    `json:"status"`
	Size    int64  `json:"size"`
	DirType string `json:"type"`
}

type ReconEngine struct {
	client *http.Client
}

func NewReconEngine() *ReconEngine {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:   100,
		MaxConnsPerHost: 10,
	}

	return &ReconEngine{
		client: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

type FOFAClient struct {
	email string
	key   string
}

type HunterClient struct {
	apiKey string
}

type ZoomEyeClient struct {
	apiKey string
}

func NewFOFAClient(email, key string) *FOFAClient {
	return &FOFAClient{email: email, key: key}
}

func NewHunterClient(apiKey string) *HunterClient {
	return &HunterClient{apiKey: apiKey}
}

func NewZoomEyeClient(apiKey string) *ZoomEyeClient {
	return &ZoomEyeClient{apiKey: apiKey}
}

func (f *FOFAClient) Search(query string, page, size int) ([]string, error) {
	baseURL := "https://fofa.info/api/v1/search/all"
	q := url.QueryEscape(query)
	apiURL := fmt.Sprintf("%s?email=%s&key=%s&qbase64=%s&page=%d&size=%d",
		baseURL, f.email, f.key, q, page, size)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Error   bool     `json:"error"`
		Results []string `json:"results"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

func (h *HunterClient) Search(query string, page, limit int) ([]string, error) {
	apiURL := "https://hunter.qianxin.com/api/home/root"
	data := url.Values{}
	data.Set("api_key", h.apiKey)
	data.Set("search", query)
	data.Set("page", fmt.Sprintf("%d", page))
	data.Set("limit", fmt.Sprintf("%d", limit))

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Code int `json:"code"`
		Data struct {
			Arr []string `json:"arr"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Data.Arr, nil
}

func (z *ZoomEyeClient) Search(query string, page int) ([]string, error) {
	apiURL := "https://api.zoomeye.org/api/search"
	data := url.Values{}
	data.Set("query", query)
	data.Set("page", fmt.Sprintf("%d", page))

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("API-KEY", z.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Matches []string `json:"matches"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Matches, nil
}

func (e *ReconEngine) EnumerateSubdomains(domain string) ([]string, error) {
	var subdomains []string

	commonPrefixes := []string{
		"www", "mail", "ftp", "localhost", "webmail", "smtp",
		"pop", "ns1", "webdisk", "ns2", "dev", "www2",
		"admin", "forum", "news", "vpn", "ns", "mail2",
		"mysql", "old", "lists", "support", "mobile", "mx",
		"static", "docs", "beta", "shop", "test", "api",
		"cdn", "blog", "file", "media", "login", "git", "m",
	}

	baseDomain := domain
	if strings.HasPrefix(domain, "http") {
		if u, err := url.Parse(domain); err == nil {
			baseDomain = u.Hostname()
		}
	}

	parts := strings.Split(baseDomain, ".")
	if len(parts) < 2 {
		baseDomain = domain
	} else {
		baseDomain = strings.Join(parts[len(parts)-2:], ".")
	}

	for _, prefix := range commonPrefixes {
		subdomain := fmt.Sprintf("%s.%s", prefix, baseDomain)
		resp, err := e.client.Get("http://" + subdomain)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode < 400 {
				subdomains = append(subdomains, subdomain)
			}
		}
	}

	if len(subdomains) == 0 {
		subdomains = append(subdomains, domain)
	}

	return subdomains, nil
}

func (e *ReconEngine) ScanPorts(host string, ports []int) []PortInfo {
	var results []PortInfo

	serviceMap := map[int]string{
		21:   "FTP",
		22:   "SSH",
		23:   "Telnet",
		25:   "SMTP",
		53:   "DNS",
		80:   "HTTP",
		110:  "POP3",
		143:  "IMAP",
		443:  "HTTPS",
		465:  "SMTPS",
		587:  "SMTP",
		993:  "IMAPS",
		995:  "POP3S",
		1433: "MSSQL",
		1521: "Oracle",
		3306: "MySQL",
		3389: "RDP",
		5432: "PostgreSQL",
		5900: "VNC",
		6379: "Redis",
		8080: "HTTP-Proxy",
		8443: "HTTPS-Alt",
		9200: "Elasticsearch",
		27017: "MongoDB",
	}

	for _, port := range ports {
		addr := fmt.Sprintf("%s:%d", host, port)
		conn, err := tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true})
		if err == nil {
			conn.Close()
			results = append(results, PortInfo{
				Port:     port,
				Protocol: "tcp",
				Service:  serviceMap[port],
				State:    "open",
			})
		}
	}

	return results
}

func (e *ReconEngine) Fingerprint(target string) []FingerprintInfo {
	var results []FingerprintInfo

	resp, err := e.client.Get(target)
	if err != nil {
		return results
	}
	defer resp.Body.Close()

	server := resp.Header.Get("Server")
	if server != "" {
		results = append(results, FingerprintInfo{
			Name:       server,
			Version:    "",
			Confidence: 90,
		})
	}

	xPoweredBy := resp.Header.Get("X-Powered-By")
	if xPoweredBy != "" {
		results = append(results, FingerprintInfo{
			Name:       xPoweredBy,
			Version:    "",
			Confidence: 80,
		})
	}

	return results
}

func (e *ReconEngine) ScanDirectories(target string, dirs []string) []DirectoryInfo {
	var results []DirectoryInfo

	for _, dir := range dirs {
		resp, err := e.client.Get(target + dir)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 || resp.StatusCode == 301 || resp.StatusCode == 302 {
			body, _ := io.ReadAll(resp.Body)
			contentType := resp.Header.Get("Content-Type")

			dirType := "unknown"
			if strings.Contains(contentType, "text/html") {
				dirType = "directory"
			} else if strings.Contains(contentType, "application") {
				dirType = "file"
			}

			results = append(results, DirectoryInfo{
				Path:    dir,
				Status:  resp.StatusCode,
				Size:    int64(len(body)),
				DirType: dirType,
			})
		}
	}

	return results
}

func (e *ReconEngine) FullRecon(target string, modules []string, platform string, platformAPIKey string) (*ReconResult, error) {
	result := &ReconResult{
		Target:       target,
		Subdomains:   []string{},
		Ports:        []PortInfo{},
		Fingerprints: []FingerprintInfo{},
		Directories:  []DirectoryInfo{},
		Timestamp:    time.Now(),
	}

	host := target
	if strings.HasPrefix(target, "http") {
		if u, err := url.Parse(target); err == nil {
			host = u.Hostname()
		}
	}

	for _, module := range modules {
		switch module {
		case "subdomain":
			if platform == "fofa" && platformAPIKey != "" {
				email := "user@example.com"
				key := platformAPIKey
				client := NewFOFAClient(email, key)
				if subs, err := client.Search(fmt.Sprintf("domain=\"%s\"", host), 1, 100); err == nil {
					result.Subdomains = subs
				}
			} else if platform == "hunter" && platformAPIKey != "" {
				client := NewHunterClient(platformAPIKey)
				if subs, err := client.Search(fmt.Sprintf("domain=\"%s\"", host), 1, 100); err == nil {
					result.Subdomains = subs
				}
			} else if platform == "zoomeye" && platformAPIKey != "" {
				client := NewZoomEyeClient(platformAPIKey)
				if subs, err := client.Search(fmt.Sprintf("domain:\"%s\"", host), 1); err == nil {
					result.Subdomains = subs
				}
			} else {
				if subs, err := e.EnumerateSubdomains(host); err == nil {
					result.Subdomains = subs
				}
			}
		case "port":
			defaultPorts := []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 465, 587, 993, 995, 1433, 1521, 3306, 3389, 5432, 5900, 6379, 8080, 8443, 9200, 27017}
			result.Ports = e.ScanPorts(host, defaultPorts)
		case "fingerprint":
			result.Fingerprints = e.Fingerprint(target)
		case "directory":
			defaultDirs := []string{"/admin", "/api", "/backup", "/config", "/login", "/test", "/debug", "/.git"}
			result.Directories = e.ScanDirectories(target, defaultDirs)
		}
	}

	return result, nil
}