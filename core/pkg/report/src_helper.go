package report

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SRCHelper struct {
	platforms []SRCPlatform
	httpClient *http.Client
}

type SRCPlatform struct {
	Name      string
	URL       string
	Type      string
	LoginURL  string
	SubmitURL string
	APIKey    string
}

type SubmissionData struct {
	Title       string
	VulnType    string
	Severity    string
	URL         string
	Description string
	Evidence    string
	POC         string
	Solution    string
}

func NewSRCHelper() *SRCHelper {
	return &SRCHelper{
		platforms: []SRCPlatform{
			{
				Name:      "补天",
				URL:       "https://www.butian.net",
				Type:      "butian",
				LoginURL:  "https://www.butian.net/login",
				SubmitURL: "https://www.butian.net/bug/submit",
			},
			{
				Name:      "阿里云先知",
				URL:       "https://xianzhicom.aliyun.com",
				Type:      "aliyun",
				LoginURL:  "https://account.aliyun.com/login",
				SubmitURL: "https://xianzhicom.aliyun.com/submit",
			},
			{
				Name:      "腾讯安全应急响应中心",
				URL:       "https://security.tencent.com",
				Type:      "tencent",
				LoginURL:  "https://security.tencent.com/login",
				SubmitURL: "https://security.tencent.com/submit",
			},
			{
				Name:      "360安全应急响应中心",
				URL:       "https://security.360.cn",
				Type:      "360",
				LoginURL:  "https://security.360.cn/login",
				SubmitURL: "https://security.360.cn/submit",
			},
			{
				Name:      "京东安全应急响应中心",
				URL:       "https://security.jd.com",
				Type:      "jd",
				LoginURL:  "https://security.jd.com/login",
				SubmitURL: "https://security.jd.com/submit",
			},
			{
				Name:      "百度安全应急响应中心",
				URL:       "https://security.baidu.com",
				Type:      "baidu",
				LoginURL:  "https://security.baidu.com/login",
				SubmitURL: "https://security.baidu.com/submit",
			},
			{
				Name:      "字节跳动安全中心",
				URL:       "https://security.bytedance.com",
				Type:      "bytedance",
				LoginURL:  "https://security.bytedance.com/login",
				SubmitURL: "https://security.bytedance.com/submit",
			},
			{
				Name:      "美团安全应急响应中心",
				URL:       "https://security.meituan.com",
				Type:      "meituan",
				LoginURL:  "https://security.meituan.com/login",
				SubmitURL: "https://security.meituan.com/submit",
			},
			{
				Name:      "滴滴出行安全应急响应中心",
				URL:       "https://security.didiglobal.com",
				Type:      "didi",
				LoginURL:  "https://security.didiglobal.com/login",
				SubmitURL: "https://security.didiglobal.com/submit",
			},
			{
				Name:      "网易安全中心",
				URL:       "https://security.163.com",
				Type:      "netease",
				LoginURL:  "https://security.163.com/login",
				SubmitURL: "https://security.163.com/submit",
			},
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (sh *SRCHelper) ListPlatforms() []SRCPlatform {
	return sh.platforms
}

func (sh *SRCHelper) GetPlatform(name string) *SRCPlatform {
	for _, p := range sh.platforms {
		if strings.EqualFold(p.Name, name) || strings.EqualFold(p.Type, name) {
			return &p
		}
	}
	return nil
}

func (sh *SRCHelper) GenerateSubmissionForm(platformType string, data SubmissionData) map[string]string {
	form := map[string]string{}

	switch platformType {
	case "butian":
		form["title"] = data.Title
		form["vul_type"] = sh.mapVulnType("butian", data.VulnType)
		form["vul_level"] = sh.mapSeverity("butian", data.Severity)
		form["vul_url"] = data.URL
		form["vul_description"] = data.Description
		form["vul_poc"] = data.POC
		form["vul_solution"] = data.Solution
	case "aliyun":
		form["title"] = data.Title
		form["type"] = sh.mapVulnType("aliyun", data.VulnType)
		form["severity"] = sh.mapSeverity("aliyun", data.Severity)
		form["url"] = data.URL
		form["description"] = data.Description
		form["poc"] = data.POC
		form["fix"] = data.Solution
	default:
		form["title"] = data.Title
		form["vuln_type"] = data.VulnType
		form["severity"] = data.Severity
		form["url"] = data.URL
		form["description"] = data.Description
		form["poc"] = data.POC
		form["solution"] = data.Solution
	}

	return form
}

func (sh *SRCHelper) mapVulnType(platform, vulnType string) string {
	mappings := map[string]map[string]string{
		"butian": {
			"SQL注入": "SQL注入",
			"XSS": "XSS跨站脚本",
			"文件上传": "文件上传",
			"未授权访问": "未授权访问/权限绕过",
			"信息泄露": "敏感信息泄露",
			"命令执行": "命令执行",
			"代码执行": "代码执行",
		},
		"aliyun": {
			"SQL注入": "SQL注入",
			"XSS": "XSS",
			"文件上传": "任意文件上传",
			"未授权访问": "未授权访问",
			"信息泄露": "敏感信息泄露",
			"命令执行": "命令执行",
			"代码执行": "代码执行",
		},
	}

	if platformMappings, ok := mappings[platform]; ok {
		if mapped, ok := platformMappings[vulnType]; ok {
			return mapped
		}
	}

	return vulnType
}

func (sh *SRCHelper) mapSeverity(platform, severity string) string {
	mappings := map[string]map[string]string{
		"butian": {
			"Critical": "严重",
			"High": "高危",
			"Medium": "中危",
			"Low": "低危",
		},
		"aliyun": {
			"Critical": "严重",
			"High": "高危",
			"Medium": "中危",
			"Low": "低危",
		},
	}

	if platformMappings, ok := mappings[platform]; ok {
		if mapped, ok := platformMappings[severity]; ok {
			return mapped
		}
	}

	return severity
}

func (sh *SRCHelper) SubmitToPlatform(platform *SRCPlatform, data SubmissionData) error {
	fmt.Printf("📤 正在提交到 %s...\n", platform.Name)

	formData := sh.GenerateSubmissionForm(platform.Type, data)

	values := url.Values{}
	for k, v := range formData {
		values.Set(k, v)
	}

	resp, err := sh.httpClient.PostForm(platform.SubmitURL, values)
	if err != nil {
		return fmt.Errorf("提交失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("✅ 成功提交到 %s\n", platform.Name)
		return nil
	}

	return fmt.Errorf("提交失败，状态码: %d", resp.StatusCode)
}

func (sh *SRCHelper) BatchSubmit(data SubmissionData, platforms []string) []error {
	var errors []error

	for _, platformName := range platforms {
		platform := sh.GetPlatform(platformName)
		if platform == nil {
			errors = append(errors, fmt.Errorf("未知平台: %s", platformName))
			continue
		}

		if err := sh.SubmitToPlatform(platform, data); err != nil {
			errors = append(errors, fmt.Errorf("%s: %v", platform.Name, err))
		}
	}

	return errors
}

func (sh *SRCHelper) GenerateMarkdownReport(data SubmissionData) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", data.Title))
	sb.WriteString("## 漏洞信息\n\n")
	sb.WriteString(fmt.Sprintf("- **漏洞类型**: %s\n", data.VulnType))
	sb.WriteString(fmt.Sprintf("- **严重级别**: %s\n", data.Severity))
	sb.WriteString(fmt.Sprintf("- **漏洞URL**: %s\n\n", data.URL))
	sb.WriteString("## 漏洞描述\n\n")
	sb.WriteString(data.Description)
	sb.WriteString("\n\n## 复现步骤\n\n")
	sb.WriteString(data.POC)
	sb.WriteString("\n\n## 修复建议\n\n")
	sb.WriteString(data.Solution)
	sb.WriteString("\n\n## 提交记录\n\n")
	sb.WriteString(fmt.Sprintf("- 提交时间: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString("- 提交者: 昆仑安全实验室\n")

	return sb.String()
}

func (sh *SRCHelper) ExportToJSON(data SubmissionData) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
