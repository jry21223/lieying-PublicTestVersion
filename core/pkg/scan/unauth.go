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
	URL        string
	Path       string
	Type       string
	StatusCode int
	Evidence   string
	Severity   string
	Confirmed  bool
}

type sensitivePathRule struct {
	path       string
	indicators []string
}

func NewUnauthScanner(target string) *UnauthScanner {
	return &UnauthScanner{
		target:  target,
		results: []UnauthResult{},
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
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

	adminIndicators := []string{
		"admin dashboard",
		"dashboard",
		"后台管理",
		"管理后台",
		"控制台",
		"系统管理",
		"用户管理",
		"权限管理",
		"management console",
		"site administration",
	}

	baseURL, err := getBaseURL(us.target)
	if err != nil {
		return
	}

	for _, path := range adminPaths {
		testURL := baseURL + path
		resp, body, err := us.fetchResponse(testURL)
		if err != nil {
			continue
		}

		bodyLower := strings.ToLower(body)
		contentType := strings.ToLower(resp.Header.Get("Content-Type"))
		if resp.StatusCode != http.StatusOK || !strings.Contains(contentType, "html") {
			continue
		}
		if looksLikeLoginPage(bodyLower) || looksLikePublicDoc(bodyLower) {
			continue
		}
		if !containsAny(bodyLower, adminIndicators) {
			continue
		}

		result := UnauthResult{
			URL:        testURL,
			Path:       path,
			Type:       "Unauthorized Admin Access",
			StatusCode: resp.StatusCode,
			Evidence:   "后台页面存在管理功能特征且无需登录即可访问",
			Severity:   "Critical",
			Confirmed:  true,
		}
		us.results = append(us.results, result)
		fmt.Printf("⚠️  发现未授权访问: %s\n", testURL)
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

	apiIndicators := []string{
		"\"users\"",
		"\"email\"",
		"\"username\"",
		"\"role\"",
		"\"password\"",
		"\"token\"",
		"\"config\"",
		"\"settings\"",
		"\"database\"",
		"\"debug\"",
		"\"trace\"",
		"\"stack\"",
		"\"admin\"",
		"connectionstring",
	}

	baseURL, err := getBaseURL(us.target)
	if err != nil {
		return
	}

	for _, path := range apiPaths {
		testURL := baseURL + path
		resp, body, err := us.fetchResponse(testURL)
		if err != nil {
			continue
		}

		bodyLower := strings.ToLower(body)
		contentType := strings.ToLower(resp.Header.Get("Content-Type"))
		if resp.StatusCode != http.StatusOK || looksLikeLoginPage(bodyLower) {
			continue
		}
		if isPublicDocumentationPath(path) {
			continue
		}
		if !strings.Contains(contentType, "json") && !containsAny(bodyLower, apiIndicators) {
			continue
		}
		if !containsAny(bodyLower, apiIndicators) {
			continue
		}

		result := UnauthResult{
			URL:        testURL,
			Path:       path,
			Type:       "Unauthorized API Access",
			StatusCode: resp.StatusCode,
			Evidence:   "接口响应包含敏感数据或调试/配置特征，且无需认证",
			Severity:   "High",
			Confirmed:  true,
		}
		us.results = append(us.results, result)
		fmt.Printf("⚠️  发现未授权API: %s\n", testURL)
	}
}

func (us *UnauthScanner) testSensitivePaths() {
	rules := []sensitivePathRule{
		{path: "/.env", indicators: []string{"app_key", "db_password", "secret", "api_key", "password", "token"}},
		{path: "/.git/config", indicators: []string{"[core]", "[remote ", "repositoryformatversion", "url = "}},
		{path: "/.git/HEAD", indicators: []string{"ref:"}},
		{path: "/.svn/entries", indicators: []string{"svn", "dir"}},
		{path: "/.htaccess", indicators: []string{"authuserfile", "rewriteengine", "deny from all"}},
		{path: "/.htpasswd", indicators: []string{":$apr1$", ":$2y$", ":$2a$", ":$1$"}},
		{path: "/config.php", indicators: []string{"database", "password", "username", "host"}},
		{path: "/config.inc.php", indicators: []string{"database", "password", "username", "host"}},
		{path: "/configuration.php", indicators: []string{"database", "password", "username", "host"}},
		{path: "/wp-config.php", indicators: []string{"db_name", "db_user", "db_password", "db_host"}},
		{path: "/database.yml", indicators: []string{"adapter:", "database:", "password:", "username:"}},
		{path: "/db.conf", indicators: []string{"database", "password", "username", "host"}},
		{path: "/application.yml", indicators: []string{"spring:", "datasource:", "password:", "username:"}},
		{path: "/application.properties", indicators: []string{"spring.datasource", "password=", "username="}},
		{path: "/settings.py", indicators: []string{"secret_key", "database", "password", "allowed_hosts"}},
		{path: "/settings.php", indicators: []string{"database", "password", "username", "host"}},
	}

	baseURL, err := getBaseURL(us.target)
	if err != nil {
		return
	}

	for _, rule := range rules {
		testURL := baseURL + rule.path
		resp, body, err := us.fetchResponse(testURL)
		if err != nil {
			continue
		}
		if resp.StatusCode != http.StatusOK {
			continue
		}

		bodyLower := strings.ToLower(body)
		if !containsAny(bodyLower, rule.indicators) {
			continue
		}

		result := UnauthResult{
			URL:        testURL,
			Path:       rule.path,
			Type:       "Sensitive File Exposure",
			StatusCode: resp.StatusCode,
			Evidence:   "敏感文件可直接访问且内容包含敏感配置特征",
			Severity:   "High",
			Confirmed:  true,
		}
		us.results = append(us.results, result)
		fmt.Printf("⚠️  发现敏感文件: %s\n", testURL)
	}
}

func (us *UnauthScanner) fetchResponse(target string) (*http.Response, string, error) {
	resp, err := us.httpClient.Get(target)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return resp, string(body), nil
}

func getBaseURL(target string) (string, error) {
	parsedURL, err := url.Parse(target)
	if err != nil {
		return "", err
	}
	return parsedURL.Scheme + "://" + parsedURL.Host, nil
}

func containsAny(body string, indicators []string) bool {
	for _, indicator := range indicators {
		if strings.Contains(body, strings.ToLower(indicator)) {
			return true
		}
	}
	return false
}

func looksLikeLoginPage(body string) bool {
	loginIndicators := []string{
		"login",
		"sign in",
		"signin",
		"password",
		"username",
		"登录",
		"密码",
		"用户名",
		"统一身份认证",
		"cas",
		"sso",
	}
	return containsAny(body, loginIndicators)
}

func looksLikePublicDoc(body string) bool {
	docIndicators := []string{"documentation", "api doc", "swagger", "openapi", "redoc"}
	return containsAny(body, docIndicators)
}

func isPublicDocumentationPath(path string) bool {
	return strings.Contains(path, "swagger") || strings.Contains(path, "api-docs")
}

func (us *UnauthScanner) GetResults() []UnauthResult {
	return us.results
}
