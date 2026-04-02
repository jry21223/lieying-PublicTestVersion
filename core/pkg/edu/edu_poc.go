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

	eps.scanZFSoft()
	eps.scanQZSoft()
	eps.scanKingo()
	eps.scanURP()
	eps.scanJWGL()

	fmt.Printf("✅ 教务系统扫描完成，发现 %d 个漏洞\n", len(eps.results))
	return eps.results, nil
}

// ScanAll 扫描所有（API用）
func (eps *EduPOCScanner) ScanAll() ([]EduPOCResult, error) {
	return eps.Scan()
}

// 正方教务系统
func (eps *EduPOCScanner) scanZFSoft() {
	paths := []string{
		"/jwglxt",
		"/jsxsd",
		"/xsxk",
		"/",
	}

	vulnPaths := []struct {
		Path        string
		VulnType    string
		Description string
	}{
		{"/jwglxt/xtgl/login_slogin.html", "信息泄露", "登录页面泄露"},
		{"/jsxsd/", "信息泄露", "正方系统默认页面"},
		{"/jwglxt/kjgl/kjgl_queryBySql?sql=", "SQL注入", "SQL注入漏洞"},
		{"/jsxsd/xsxk/xsxk_index?jx0502zbid=", "SQL注入", "选课SQL注入"},
		{"/jwglxt/xtgl/init_cxAreaOrProvince.html", "未授权访问", "地区信息未授权访问"},
		{"/jsxsd/framework/jslib/", "目录遍历", "JS库目录遍历"},
	}

	for _, path := range paths {
		baseURL := eps.target + path
		resp, err := eps.httpClient.Get(baseURL)
		if err != nil {
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		bodyStr := string(body)

		if strings.Contains(bodyStr, "正方") || strings.Contains(bodyStr, "zfsoft") {
			fmt.Printf("✅ 发现正方教务系统: %s\n", baseURL)

			for _, vuln := range vulnPaths {
				vulnURL := eps.target + vuln.Path
				resp, err := eps.httpClient.Get(vulnURL)
				if err != nil {
					continue
				}
				resp.Body.Close()

				if resp.StatusCode == 200 {
					result := EduPOCResult{
						SystemType:  "正方教务系统",
						VulnType:    vuln.VulnType,
						Severity:    "High",
						URL:         vulnURL,
						Evidence:    "页面可访问",
						Description: vuln.Description,
					}
					eps.results = append(eps.results, result)
					fmt.Printf("⚠️  发现漏洞: %s - %s\n", vuln.VulnType, vulnURL)
				}
			}
			return
		}
	}
}

// 强智教务系统
func (eps *EduPOCScanner) scanQZSoft() {
	paths := []string{
		"/qzsoft",
		"/qz",
		"/",
	}

	vulnPaths := []struct {
		Path        string
		VulnType    string
		Description string
	}{
		{"/qzsoft/", "信息泄露", "强智系统默认页面"},
		{"/qzsoft/proxy", "SSRF", "代理SSRF漏洞"},
		{"/qzsoft/ueditor/", "文件上传", "UEditor上传漏洞"},
		{"/qzsoft/admin", "未授权访问", "后台未授权访问"},
	}

	for _, path := range paths {
		baseURL := eps.target + path
		resp, err := eps.httpClient.Get(baseURL)
		if err != nil {
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		bodyStr := string(body)

		if strings.Contains(bodyStr, "强智") || strings.Contains(bodyStr, "qzsoft") {
			fmt.Printf("✅ 发现强智教务系统: %s\n", baseURL)

			for _, vuln := range vulnPaths {
				vulnURL := eps.target + vuln.Path
				resp, err := eps.httpClient.Get(vulnURL)
				if err != nil {
					continue
				}
				resp.Body.Close()

				if resp.StatusCode == 200 {
					result := EduPOCResult{
						SystemType:  "强智教务系统",
						VulnType:    vuln.VulnType,
						Severity:    "High",
						URL:         vulnURL,
						Evidence:    "页面可访问",
						Description: vuln.Description,
					}
					eps.results = append(eps.results, result)
					fmt.Printf("⚠️  发现漏洞: %s - %s\n", vuln.VulnType, vulnURL)
				}
			}
			return
		}
}
}

// 金智教育
func (eps *EduPOCScanner) scanKingo() {
	paths := []string{
		"/kingo",
		"/",
	}

	vulnPaths := []struct {
		Path        string
		VulnType    string
		Description string
	}{
		{"/kingo/", "信息泄露", "金智系统默认页面"},
		{"/kingo/admin", "未授权访问", "后台未授权访问"},
		{"/kingo/upload", "文件上传", "文件上传漏洞"},
	}

	for _, path := range paths {
		baseURL := eps.target + path
		resp, err := eps.httpClient.Get(baseURL)
		if err != nil {
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		bodyStr := string(body)

		if strings.Contains(bodyStr, "金智") || strings.Contains(bodyStr, "kingo") {
			fmt.Printf("✅ 发现金智教务系统: %s\n", baseURL)

			for _, vuln := range vulnPaths {
				vulnURL := eps.target + vuln.Path
				resp, err := eps.httpClient.Get(vulnURL)
				if err != nil {
					continue
				}
				resp.Body.Close()

				if resp.StatusCode == 200 {
					result := EduPOCResult{
						SystemType:  "金智教务系统",
						VulnType:    vuln.VulnType,
						Severity:    "High",
						URL:         vulnURL,
						Evidence:    "页面可访问",
						Description: vuln.Description,
					}
					eps.results = append(eps.results, result)
					fmt.Printf("⚠️  发现漏洞: %s - %s\n", vuln.VulnType, vulnURL)
				}
			}
			return
		}
	}
}

// URP教务系统
func (eps *EduPOCScanner) scanURP() {
	paths := []string{
		"/urp",
		"/",
	}

	vulnPaths := []struct {
		Path        string
		VulnType    string
		Description string
	}{
		{"/urp/", "信息泄露", "URP系统默认页面"},
		{"/urp/servlet/", "未授权访问", "Servlet未授权访问"},
		{"/urp/servlet/com.runqian.report.view.Servlet", "SQL注入", "报表SQL注入"},
	}

	for _, path := range paths {
		baseURL := eps.target + path
		resp, err := eps.httpClient.Get(baseURL)
		if err != nil {
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		bodyStr := string(body)

		if strings.Contains(bodyStr, "URP") || strings.Contains(bodyStr, "urp") {
			fmt.Printf("✅ 发现URP教务系统: %s\n", baseURL)

			for _, vuln := range vulnPaths {
				vulnURL := eps.target + vuln.Path
				resp, err := eps.httpClient.Get(vulnURL)
				if err != nil {
					continue
				}
				resp.Body.Close()

				if resp.StatusCode == 200 {
					result := EduPOCResult{
						SystemType:  "URP教务系统",
						VulnType:    vuln.VulnType,
						Severity:    "High",
						URL:         vulnURL,
						Evidence:    "页面可访问",
						Description: vuln.Description,
					}
					eps.results = append(eps.results, result)
					fmt.Printf("⚠️  发现漏洞: %s - %s\n", vuln.VulnType, vulnURL)
				}
			}
			return
		}
	}
}

// 青果教务系统
func (eps *EduPOCScanner) scanJWGL() {
	paths := []string{
		"/jwgl",
		"/",
	}

	vulnPaths := []struct {
		Path        string
		VulnType    string
		Description string
	}{
		{"/jwgl/", "信息泄露", "青果系统默认页面"},
		{"/jwgl/admin", "未授权访问", "后台未授权访问"},
		{"/jwgl/ueditor/", "文件上传", "UEditor上传漏洞"},
	}

	for _, path := range paths {
		baseURL := eps.target + path
		resp, err := eps.httpClient.Get(baseURL)
		if err != nil {
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		bodyStr := string(body)

		if strings.Contains(bodyStr, "青果") || strings.Contains(bodyStr, "jwgl") {
			fmt.Printf("✅ 发现青果教务系统: %s\n", baseURL)

			for _, vuln := range vulnPaths {
				vulnURL := eps.target + vuln.Path
				resp, err := eps.httpClient.Get(vulnURL)
				if err != nil {
					continue
				}
				resp.Body.Close()

				if resp.StatusCode == 200 {
					result := EduPOCResult{
						SystemType:  "青果教务系统",
						VulnType:    vuln.VulnType,
						Severity:    "High",
						URL:         vulnURL,
						Evidence:    "页面可访问",
						Description: vuln.Description,
					}
					eps.results = append(eps.results, result)
					fmt.Printf("⚠️  发现漏洞: %s - %s\n", vuln.VulnType, vulnURL)
				}
			}
			return
		}
	}
}

func (eps *EduPOCScanner) GetResults() []EduPOCResult {
	return eps.results
}
