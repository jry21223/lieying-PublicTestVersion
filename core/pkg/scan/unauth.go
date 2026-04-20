package scan

import (
	"fmt"
	"net/url"

	"github.com/kunlun-sec/lunying/pkg/utils"
)

type UnauthScanner struct {
	target  string
	results []UnauthResult
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

func NewUnauthScanner(target string) *UnauthScanner {
	return &UnauthScanner{
		target:  target,
		results: []UnauthResult{},
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
		"/admin", "/admin/", "/admin/index.php", "/admin/dashboard",
		"/administrator", "/administrator/", "/manage", "/manage/",
		"/backend", "/backend/", "/console", "/console/",
		"/panel", "/panel/", "/wp-admin", "/wp-admin/",
		"/phpmyadmin", "/phpmyadmin/",
	}

	parsedURL, err := url.Parse(us.target)
	if err != nil {
		return
	}
	baseURL := parsedURL.Scheme + "://" + parsedURL.Host

	loginWords := []string{
		"login", "登录", "password", "密码", "username", "用户名",
		"sign in", "signin", "用户登录", "请登录",
	}
	adminWords := []string{
		"dashboard", "管理", "logout", "退出", "欢迎", "welcome",
		"用户管理", "系统设置", "数据统计", "control panel",
	}

	for _, path := range adminPaths {
		testURL := baseURL + path
		body, code := utils.FetchBody(testURL, 256*1024)
		if code != 200 || body == "" {
			continue
		}
		hasLogin := utils.ContainsAnyLower(body, loginWords)
		hasAdmin := utils.ContainsAnyLower(body, adminWords)

		if !hasLogin && hasAdmin {
			us.results = append(us.results, UnauthResult{
				URL:        testURL,
				Path:       path,
				Type:       "Unauthorized Admin Access",
				StatusCode: code,
				Evidence:   "管理页面无需认证即可访问，页面包含管理功能内容",
				Severity:   utils.SeverityCritical,
				Confirmed:  true,
			})
			fmt.Printf("⚠️  [已确认] 管理后台未授权: %s\n", testURL)
		}
	}
}

func (us *UnauthScanner) testAPIEndpoints() {
	apiPaths := []struct {
		path       string
		indicators []string
	}{
		{"/api/users", []string{`"id"`, `"username"`, `"email"`, `"user"`}},
		{"/api/user", []string{`"id"`, `"name"`, `"email"`}},
		{"/api/admin", []string{`"id"`, `"role"`, `"admin"`}},
		{"/api/config", []string{`"key"`, `"value"`, `"config"`, `"setting"`}},
		{"/api/settings", []string{`"setting"`, `"config"`, `"value"`}},
		{"/swagger-ui.html", []string{"swagger", "Swagger UI", "openapi"}},
		{"/swagger", []string{"swagger", "openapi", "paths"}},
		{"/api-docs", []string{"swagger", "openapi", `"paths"`}},
		{"/v2/api-docs", []string{"swagger", "openapi", `"paths"`}},
		{"/api/v1/users", []string{`"id"`, `"username"`, `"email"`}},
		{"/graphql", []string{`"data"`, `"errors"`, `"__schema"`}},
	}

	parsedURL, err := url.Parse(us.target)
	if err != nil {
		return
	}
	baseURL := parsedURL.Scheme + "://" + parsedURL.Host

	for _, ep := range apiPaths {
		testURL := baseURL + ep.path
		body, code := utils.FetchBody(testURL, 256*1024)
		if code != 200 || body == "" {
			continue
		}
		if utils.ContainsAny(body, ep.indicators) {
			us.results = append(us.results, UnauthResult{
				URL:        testURL,
				Path:       ep.path,
				Type:       "Unauthorized API Access",
				StatusCode: code,
				Evidence:   utils.ExtractSnippet(body, ep.indicators, 20, 60),
				Severity:   utils.SeverityHigh,
				Confirmed:  true,
			})
			fmt.Printf("⚠️  [已确认] 未授权API数据暴露: %s\n", testURL)
		}
	}
}

func (us *UnauthScanner) testSensitivePaths() {
	sensitivePaths := []struct {
		path       string
		indicators []string
		severity   string
	}{
		{"/.env", []string{"DB_PASSWORD", "APP_KEY", "SECRET", "DATABASE_URL", "REDIS_"}, utils.SeverityCritical},
		{"/.git/config", []string{"[core]", "[remote", "repositoryformatversion"}, utils.SeverityHigh},
		{"/.git/HEAD", []string{"ref: refs/", "branch"}, utils.SeverityMedium},
		{"/wp-config.php", []string{"DB_NAME", "DB_USER", "DB_PASSWORD", "table_prefix"}, utils.SeverityCritical},
		{"/config.php", []string{"password", "database", "host", "DB_"}, utils.SeverityHigh},
		{"/database.yml", []string{"adapter:", "database:", "password:", "username:"}, utils.SeverityHigh},
		{"/application.yml", []string{"spring:", "datasource:", "password:", "secret"}, utils.SeverityHigh},
		{"/application.properties", []string{"spring.datasource", "password=", "secret="}, utils.SeverityHigh},
		{"/settings.py", []string{"SECRET_KEY", "DATABASES", "PASSWORD"}, utils.SeverityHigh},
		{"/.htpasswd", []string{"$apr1$", "$2y$", ":{SHA}", ":$"}, utils.SeverityHigh},
		{"/crossdomain.xml", []string{"<cross-domain-policy", "allow-access-from"}, utils.SeverityLow},
	}

	parsedURL, err := url.Parse(us.target)
	if err != nil {
		return
	}
	baseURL := parsedURL.Scheme + "://" + parsedURL.Host

	for _, sp := range sensitivePaths {
		testURL := baseURL + sp.path
		body, code := utils.FetchBody(testURL, 256*1024)
		if code != 200 || body == "" {
			continue
		}
		if utils.ContainsAny(body, sp.indicators) {
			us.results = append(us.results, UnauthResult{
				URL:        testURL,
				Path:       sp.path,
				Type:       "Sensitive File Exposure",
				StatusCode: code,
				Evidence:   utils.ExtractSnippet(body, sp.indicators, 20, 60),
				Severity:   sp.severity,
				Confirmed:  true,
			})
			fmt.Printf("⚠️  [已确认] 敏感文件泄露: %s\n", testURL)
		}
	}
}

func (us *UnauthScanner) GetResults() []UnauthResult {
	return us.results
}