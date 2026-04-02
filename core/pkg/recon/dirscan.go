package recon

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type DirScanner struct {
	baseURL     string
	paths       []string
	foundDirs   []*DirResult
	mu          sync.Mutex
	timeout     time.Duration
	workers     int
	httpClient  *http.Client
}

type DirResult struct {
	URL           string
	Path          string
	StatusCode    int
	ContentLength int64
	Title         string
	IsSensitive   bool
}

func NewDirScanner(baseURL string) *DirScanner {
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &DirScanner{
		baseURL: baseURL,
		paths:   []string{},
		foundDirs: []*DirResult{},
		timeout: 3 * time.Second,
		workers: 50,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (ds *DirScanner) SetTimeout(timeout time.Duration) {
	ds.timeout = timeout
	ds.httpClient.Timeout = timeout
}

func (ds *DirScanner) SetWorkers(workers int) {
	ds.workers = workers
}

func (ds *DirScanner) AddPath(path string) {
	ds.paths = append(ds.paths, path)
}

func (ds *DirScanner) AddCommonPaths() {
	commonPaths := []string{
		"/",
		"/admin",
		"/admin.php",
		"/admin/login",
		"/admin/index.php",
		"/administrator",
		"/api",
		"/api/v1",
		"/api/v2",
		"/auth",
		"/backup",
		"/backup.zip",
		"/backup.sql",
		"/.backup",
		"/bak",
		"/bak.zip",
		"/config",
		"/config.php",
		"/config.inc.php",
		"/.config",
		"/conf",
		"/database",
		"/database.php",
		"/db",
		"/db.php",
		"/.env",
		"/.env.example",
		"/.git",
		"/.git/config",
		"/.git/HEAD",
		"/.htaccess",
		"/.htpasswd",
		"/index.php",
		"/install",
		"/install.php",
		"/install/index.php",
		"/login",
		"/login.php",
		"/logout",
		"/manage",
		"/manager",
		"/old",
		"/phpinfo.php",
		"/.php_cs.cache",
		"/private",
		"/protected",
		"/public",
		"/robots.txt",
		"/sitemap.xml",
		"/sitemap.xml.gz",
		"/sql",
		"/sql.zip",
		"/temp",
		"/tmp",
		"/test",
		"/testing",
		"/user",
		"/users",
		"/uploads",
		"/upload",
		"/web.config",
		"/.well-known",
		"/.well-known/security.txt",
		"/favicon.ico",
		"/crossdomain.xml",
		"/clientaccesspolicy.xml",
		"/.idea",
		"/.vscode",
		"/.DS_Store",
		"/README.md",
		"/README.txt",
		"/CHANGELOG.md",
		"/LICENSE",
		"/LICENSE.txt",
		"/composer.json",
		"/composer.lock",
		"/package.json",
		"/package-lock.json",
		"/yarn.lock",
	}

	for _, path := range commonPaths {
		ds.paths = append(ds.paths, path)
	}
}

func (ds *DirScanner) Scan() ([]*DirResult, error) {
	fmt.Printf("🔍 开始目录扫描: %s\n", ds.baseURL)
	fmt.Printf("📊 扫描路径数量: %d\n", len(ds.paths))

	if len(ds.paths) == 0 {
		return []*DirResult{}, nil
	}

	pathChan := make(chan string, len(ds.paths))
	resultChan := make(chan *DirResult, len(ds.paths))
	var wg sync.WaitGroup

	for i := 0; i < ds.workers; i++ {
		wg.Add(1)
		go ds.worker(pathChan, resultChan, &wg)
	}

	for _, path := range ds.paths {
		pathChan <- path
	}
	close(pathChan)

	wg.Wait()
	close(resultChan)

	for result := range resultChan {
		ds.foundDirs = append(ds.foundDirs, result)
	}

	fmt.Printf("✅ 发现目录/文件: %d\n", len(ds.foundDirs))
	return ds.foundDirs, nil
}

func (ds *DirScanner) worker(pathChan <-chan string, resultChan chan<- *DirResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for path := range pathChan {
		result := ds.scanPath(path)
		if result != nil {
			ds.mu.Lock()
			resultChan <- result
			ds.mu.Unlock()
		}
	}
}

func (ds *DirScanner) scanPath(path string) *DirResult {
	url := ds.baseURL + path

	resp, err := ds.httpClient.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	title := ds.extractTitle(string(bodyBytes))
	isSensitive := ds.isSensitivePath(path, resp.StatusCode, title)

	result := &DirResult{
		URL:           url,
		Path:          path,
		StatusCode:    resp.StatusCode,
		ContentLength: resp.ContentLength,
		Title:         title,
		IsSensitive:   isSensitive,
	}

	if isSensitive {
		fmt.Printf("⚠️  发现敏感路径: %s (%d)\n", url, resp.StatusCode)
	} else if resp.StatusCode == 200 || resp.StatusCode == 301 || resp.StatusCode == 302 {
		fmt.Printf("✅ 发现: %s (%d)\n", url, resp.StatusCode)
	}

	return result
}

func (ds *DirScanner) extractTitle(body string) string {
	re := regexp.MustCompile(`<title[^>]*>([^<]+)</title>`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		title := strings.TrimSpace(matches[1])
		if len(title) > 100 {
			title = title[:97] + "..."
		}
		return title
	}
	return ""
}

func (ds *DirScanner) isSensitivePath(path string, statusCode int, title string) bool {
	sensitiveKeywords := []string{
		"backup", "bak", "config", "conf", "env", "git", "sql",
		"database", "db", "private", "protected", "secret", "key",
		"password", "passwd", "htaccess", "htpasswd",
	}

	pathLower := strings.ToLower(path)
	titleLower := strings.ToLower(title)

	for _, keyword := range sensitiveKeywords {
		if strings.Contains(pathLower, keyword) {
			return true
		}
		if strings.Contains(titleLower, keyword) {
			return true
		}
	}

	if statusCode == 200 && (strings.HasSuffix(path, ".zip") || 
		strings.HasSuffix(path, ".sql") || 
		strings.HasSuffix(path, ".tar") ||
		strings.HasSuffix(path, ".tar.gz") ||
		strings.HasSuffix(path, ".rar") ||
		strings.HasSuffix(path, ".7z")) {
		return true
	}

	return false
}

func (ds *DirScanner) GetFoundDirs() []*DirResult {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	result := make([]*DirResult, len(ds.foundDirs))
	copy(result, ds.foundDirs)
	return result
}
