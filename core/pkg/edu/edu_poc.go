package edu

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type EduPOCScanner struct {
	target     string
	results    []EduPOCResult
	httpClient *http.Client
}

type EduPOCResult struct {
	SystemType  string
	VulnType    string
	Severity    string
	URL         string
	Evidence    string
	Payload     string
	Description string
}

type eduSystemRule struct {
	SystemType   string
	DetectPaths  []string
	Fingerprints []string
	VulnRules    []eduVulnRule
}

type eduVulnRule struct {
	Path        string
	VulnType    string
	Severity    string
	Description string
	Indicators  []string
}

func NewEduPOCScanner(target string) *EduPOCScanner {
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "http://" + target
	}
	return &EduPOCScanner{
		target:  target,
		results: []EduPOCResult{},
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (eps *EduPOCScanner) Scan() ([]EduPOCResult, error) {
	fmt.Println("🔍 开始教务系统专项扫描...")

	for _, rule := range eps.systemRules() {
		eps.scanSystem(rule)
	}

	fmt.Printf("✅ 教务系统扫描完成，发现 %d 个漏洞\n", len(eps.results))
	return eps.results, nil
}

// ScanAll 扫描所有（API用）
func (eps *EduPOCScanner) ScanAll() ([]EduPOCResult, error) {
	return eps.Scan()
}

func (eps *EduPOCScanner) scanSystem(rule eduSystemRule) {
	matchedBaseURL := ""
	for _, path := range rule.DetectPaths {
		baseURL := eps.target + path
		resp, body, err := eps.fetch(baseURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			continue
		}

		bodyLower := strings.ToLower(body)
		if !containsAny(bodyLower, rule.Fingerprints) {
			continue
		}
		matchedBaseURL = baseURL
		break
	}

	if matchedBaseURL == "" {
		return
	}

	fmt.Printf("✅ 发现%s: %s\n", rule.SystemType, matchedBaseURL)
	for _, vuln := range rule.VulnRules {
		vulnURL := eps.target + vuln.Path
		resp, body, err := eps.fetch(vulnURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			continue
		}

		bodyLower := strings.ToLower(body)
		if looksLikeEduLoginPage(bodyLower) {
			continue
		}
		if !containsAny(bodyLower, vuln.Indicators) {
			continue
		}

		result := EduPOCResult{
			SystemType:  rule.SystemType,
			VulnType:    vuln.VulnType,
			Severity:    vuln.Severity,
			URL:         vulnURL,
			Evidence:    "页面包含与漏洞类型匹配的暴露特征",
			Description: vuln.Description,
		}
		eps.results = append(eps.results, result)
		fmt.Printf("⚠️  发现漏洞: %s - %s\n", vuln.VulnType, vulnURL)
	}
}

func (eps *EduPOCScanner) fetch(target string) (*http.Response, string, error) {
	resp, err := eps.httpClient.Get(target)
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

func (eps *EduPOCScanner) systemRules() []eduSystemRule {
	return []eduSystemRule{
		{
			SystemType:   "正方教务系统",
			DetectPaths:  []string{"/jwglxt", "/jsxsd", "/xsxk", "/"},
			Fingerprints: []string{"正方", "zfsoft", "正方教务"},
			VulnRules: []eduVulnRule{
				{Path: "/jwglxt/kjgl/kjgl_queryBySql?sql=", VulnType: "SQL注入", Severity: "High", Description: "SQL注入漏洞", Indicators: []string{"sql syntax", "mysql", "ora-", "syntax error"}},
				{Path: "/jwglxt/xtgl/init_cxAreaOrProvince.html", VulnType: "未授权访问", Severity: "High", Description: "地区信息未授权访问", Indicators: []string{"地区", "province", "city", "json"}},
				{Path: "/jsxsd/framework/jslib/", VulnType: "目录遍历", Severity: "High", Description: "JS库目录遍历", Indicators: []string{"index of", "directory listing", "parent directory", "jslib"}},
			},
		},
		{
			SystemType:   "强智教务系统",
			DetectPaths:  []string{"/qzsoft", "/qz", "/"},
			Fingerprints: []string{"强智", "qzsoft", "强智教务"},
			VulnRules: []eduVulnRule{
				{Path: "/qzsoft/proxy", VulnType: "SSRF", Severity: "High", Description: "代理SSRF漏洞", Indicators: []string{"proxy", "http://", "https://", "target url"}},
				{Path: "/qzsoft/ueditor/", VulnType: "文件上传", Severity: "High", Description: "UEditor上传漏洞", Indicators: []string{"upload", "ueditor", "success", "file url"}},
				{Path: "/qzsoft/admin", VulnType: "未授权访问", Severity: "High", Description: "后台未授权访问", Indicators: []string{"admin dashboard", "系统管理", "用户管理", "权限配置"}},
			},
		},
		{
			SystemType:   "金智教务系统",
			DetectPaths:  []string{"/kingo", "/"},
			Fingerprints: []string{"金智", "kingo", "金智教育"},
			VulnRules: []eduVulnRule{
				{Path: "/kingo/admin", VulnType: "未授权访问", Severity: "High", Description: "后台未授权访问", Indicators: []string{"admin dashboard", "系统管理", "用户管理", "权限配置"}},
				{Path: "/kingo/upload", VulnType: "文件上传", Severity: "High", Description: "文件上传漏洞", Indicators: []string{"upload", "success", "file url", "uploaded"}},
			},
		},
		{
			SystemType:   "URP教务系统",
			DetectPaths:  []string{"/urp", "/"},
			Fingerprints: []string{"urp教务", "urp系统", "runqian"},
			VulnRules: []eduVulnRule{
				{Path: "/urp/servlet/", VulnType: "未授权访问", Severity: "High", Description: "Servlet未授权访问", Indicators: []string{"servlet", "admin dashboard", "系统管理", "用户管理"}},
				{Path: "/urp/servlet/com.runqian.report.view.Servlet", VulnType: "SQL注入", Severity: "High", Description: "报表SQL注入", Indicators: []string{"sql syntax", "mysql", "ora-", "runqian"}},
			},
		},
		{
			SystemType:   "青果教务系统",
			DetectPaths:  []string{"/jwgl", "/"},
			Fingerprints: []string{"青果", "青果教务", "教务管理系统"},
			VulnRules: []eduVulnRule{
				{Path: "/jwgl/admin", VulnType: "未授权访问", Severity: "High", Description: "后台未授权访问", Indicators: []string{"admin dashboard", "系统管理", "用户管理", "权限配置"}},
				{Path: "/jwgl/ueditor/", VulnType: "文件上传", Severity: "High", Description: "UEditor上传漏洞", Indicators: []string{"upload", "ueditor", "success", "file url"}},
			},
		},
	}
}

func containsAny(body string, indicators []string) bool {
	for _, indicator := range indicators {
		if strings.Contains(body, strings.ToLower(indicator)) {
			return true
		}
	}
	return false
}

func looksLikeEduLoginPage(body string) bool {
	return containsAny(body, []string{"登录", "统一身份认证", "cas", "sso", "password", "username", "用户名", "密码"})
}

func (eps *EduPOCScanner) GetResults() []EduPOCResult {
	return eps.results
}
