package scan

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type UnauthScanner struct {
	target     string
	results    []UnauthResult
	httpClient *http.Client
}

type UnauthResult struct {
	URL         string
	Path        string
	Type        string
	StatusCode  int
	Evidence    string
	Severity    string
	Confirmed   bool
}

func NewUnauthScanner(target string) *UnauthScanner {
	return &UnauthScanner{
		target:  target,
		results: []UnauthResult{},
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (us *UnauthScanner) Scan() ([]UnauthResult, error) {
	fmt.Println("🔍 开始未授权访问检测...")

	us.testAdminPaths()
	us.testAPIEndpoints()
	us.testSensitivePaths()

	fmt.Printf("✅ 未授权访问检测完成，发现 %d 个漏洞\n", len(us.results))
	return us.results, nil
}

func (us *UnauthScanner) testAdminPaths() {
	adminPaths := []string{
		"/admin",
		"/admin/",
		"/admin/index.php",
		"/admin/login.php",
		"/admin/dashboard",
		"/administrator",
		"/administrator/",
		"/manage",
		"/manage/",
		"/manager",
		"/manager/",
		"/backend",
		"/backend/",
		"/console",
		"/console/",
		"/control",
		"/control/",
		"/panel",
		"/panel/",
		"/cms",
		"/cms/",
		"/wp-admin",
		"/wp-admin/",
		"/phpmyadmin",
		"/phpmyadmin/",
	}

	parsedURL, err := url.Parse(us.target)
	if err != nil {
		return
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host

	for _, path := range adminPaths {
		testURL := baseURL + path
		
		resp, err := us.httpClient.Get(testURL)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		bodyStr := string(body)
		
		if resp.StatusCode == 200 {
			if !strings.Contains(bodyStr, "login") &&
			   !strings.Contains(bodyStr, "登录") &&
			   !strings.Contains(bodyStr, "password") &&
			   !strings.Contains(bodyStr, "密码") &&
			   !strings.Contains(bodyStr, "username") &&
			   !strings.Contains(bodyStr, "用户名") {
				
				result := UnauthResult{
					URL:        testURL,
					Path:       path,
					Type:       "Unauthorized Admin Access",
					StatusCode: resp.StatusCode,
					Evidence:   "管理页面可直接访问，无需登录",
					Severity:   "Critical",
					Confirmed:  false,
				}
				us.results = append(us.results, result)
				fmt.Printf("⚠️  发现未授权访问: %s\n", testURL)
			}
		}
	}
}

func (us *UnauthScanner) testAPIEndpoints() {
	apiPaths := []string{
		"/api/users",
		"/api/user",
		"/api/admin",
		"/api/config",
		"/api/settings",
		"/api/system",
		"/api/database",
		"/api/backup",
		"/api/logs",
		"/api/debug",
		"/swagger-ui.html",
		"/swagger",
		"/api-docs",
		"/v2/api-docs",
		"/api/v1/users",
		"/api/v2/users",
		"/rest/users",
		"/graphql",
	}

	parsedURL, err := url.Parse(us.target)
	if err != nil {
		return
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host

	for _, path := range apiPaths {
		testURL := baseURL + path
		
		resp, err := us.httpClient.Get(testURL)
		if err != nil {
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == 200 {
			result := UnauthResult{
				URL:        testURL,
				Path:       path,
				Type:       "Unauthorized API Access",
				StatusCode: resp.StatusCode,
				Evidence:   "API接口可直接访问，无需认证",
				Severity:   "High",
				Confirmed:  false,
			}
			us.results = append(us.results, result)
			fmt.Printf("⚠️  发现未授权API: %s\n", testURL)
		}
	}
}

func (us *UnauthScanner) testSensitivePaths() {
	sensitivePaths := []string{
		"/.env",
		"/.git/config",
		"/.git/HEAD",
		"/.svn/entries",
		"/.htaccess",
		"/.htpasswd",
		"/config.php",
		"/config.inc.php",
		"/configuration.php",
		"/wp-config.php",
		"/database.yml",
		"/db.conf",
		"/application.yml",
		"/application.properties",
		"/settings.py",
		"/settings.php",
		"/robots.txt",
		"/sitemap.xml",
		"/crossdomain.xml",
		"/clientaccesspolicy.xml",
		"/.well-known/security.txt",
	}

	parsedURL, err := url.Parse(us.target)
	if err != nil {
		return
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host

	for _, path := range sensitivePaths {
		testURL := baseURL + path
		
		resp, err := us.httpClient.Get(testURL)
		if err != nil {
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == 200 {
			result := UnauthResult{
				URL:        testURL,
				Path:       path,
				Type:       "Sensitive File Exposure",
				StatusCode: resp.StatusCode,
				Evidence:   "敏感文件可直接访问",
				Severity:   "High",
				Confirmed:  false,
			}
			us.results = append(us.results, result)
			fmt.Printf("⚠️  发现敏感文件: %s\n", testURL)
		}
	}
}

func (us *UnauthScanner) GetResults() []UnauthResult {
	return us.results
}
