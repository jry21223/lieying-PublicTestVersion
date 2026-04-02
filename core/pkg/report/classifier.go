package report

import (
	"fmt"
	"strings"
)

type VulnClassifier struct {
	classifications map[string]VulnClassification
}

type VulnClassification struct {
	Type        string
	Category    string
	CWEID       string
	Severity    string
	AvgReward   float64
	Description string
}

func NewVulnClassifier() *VulnClassifier {
	return &VulnClassifier{
		classifications: map[string]VulnClassification{
			"SQL注入": {
				Type:        "SQL注入",
				Category:    "注入类",
				CWEID:       "CWE-89",
				Severity:    "Critical",
				AvgReward:   5000,
				Description: "SQL注入漏洞允许攻击者执行任意SQL命令",
			},
			"XSS": {
				Type:        "XSS",
				Category:    "注入类",
				CWEID:       "CWE-79",
				Severity:    "High",
				AvgReward:   2000,
				Description: "跨站脚本漏洞允许攻击者注入恶意脚本",
			},
			"文件上传": {
				Type:        "文件上传",
				Category:    "文件操作",
				CWEID:       "CWE-434",
				Severity:    "High",
				AvgReward:   3000,
				Description: "不安全的文件上传允许上传恶意文件",
			},
			"未授权访问": {
				Type:        "未授权访问",
				Category:    "访问控制",
				CWEID:       "CWE-306",
				Severity:    "High",
				AvgReward:   2500,
				Description: "缺少身份验证或访问控制",
			},
			"信息泄露": {
				Type:        "信息泄露",
				Category:    "信息泄露",
				CWEID:       "CWE-200",
				Severity:    "Medium",
				AvgReward:   1000,
				Description: "敏感信息意外暴露给未授权用户",
			},
			"命令执行": {
				Type:        "命令执行",
				Category:    "代码执行",
				CWEID:       "CWE-78",
				Severity:    "Critical",
				AvgReward:   8000,
				Description: "允许攻击者在服务器上执行任意命令",
			},
			"代码执行": {
				Type:        "代码执行",
				Category:    "代码执行",
				CWEID:       "CWE-94",
				Severity:    "Critical",
				AvgReward:   10000,
				Description: "允许攻击者在服务器上执行任意代码",
			},
			"目录遍历": {
				Type:        "目录遍历",
				Category:    "文件操作",
				CWEID:       "CWE-22",
				Severity:    "High",
				AvgReward:   2000,
				Description: "允许访问受限目录之外的文件",
			},
			"SSRF": {
				Type:        "SSRF",
				Category:    "服务端请求伪造",
				CWEID:       "CWE-918",
				Severity:    "High",
				AvgReward:   4000,
				Description: "服务器端请求伪造漏洞",
			},
			"CSRF": {
				Type:        "CSRF",
				Category:    "访问控制",
				CWEID:       "CWE-352",
				Severity:    "Medium",
				AvgReward:   1500,
				Description: "跨站请求伪造漏洞",
			},
			"XXE": {
				Type:        "XXE",
				Category:    "XML处理",
				CWEID:       "CWE-611",
				Severity:    "Critical",
				AvgReward:   6000,
				Description: "XML外部实体注入漏洞",
			},
			"反序列化": {
				Type:        "反序列化",
				Category:    "代码执行",
				CWEID:       "CWE-502",
				Severity:    "Critical",
				AvgReward:   7000,
				Description: "不安全的反序列化漏洞",
			},
			"逻辑漏洞": {
				Type:        "逻辑漏洞",
				Category:    "业务逻辑",
				CWEID:       "CWE-840",
				Severity:    "Medium",
				AvgReward:   2000,
				Description: "业务逻辑实现缺陷",
			},
			"越权访问": {
				Type:        "越权访问",
				Category:    "访问控制",
				CWEID:       "CWE-639",
				Severity:    "High",
				AvgReward:   3000,
				Description: "水平或垂直权限绕过",
			},
			"弱口令": {
				Type:        "弱口令",
				Category:    "认证",
				CWEID:       "CWE-521",
				Severity:    "Medium",
				AvgReward:   1000,
				Description: "使用弱密码或默认密码",
			},
		},
	}
}

func (vc *VulnClassifier) Classify(vulnType string) *VulnClassification {
	vulnType = strings.TrimSpace(vulnType)

	for key, classification := range vc.classifications {
		if strings.Contains(strings.ToLower(vulnType), strings.ToLower(key)) {
			return &classification
		}
	}

	return &VulnClassification{
		Type:        vulnType,
		Category:    "其他",
		CWEID:       "CWE-Unknown",
		Severity:    "Medium",
		AvgReward:   1000,
		Description: "未分类漏洞",
	}
}

func (vc *VulnClassifier) AutoClassify(vulnName, vulnDesc string) *VulnClassification {
	combined := strings.ToLower(vulnName + " " + vulnDesc)

	keywords := map[string]string{
		"sql":           "SQL注入",
		"xss":           "XSS",
		"cross-site":    "XSS",
		"upload":        "文件上传",
		"unauthorized":  "未授权访问",
		"unauth":        "未授权访问",
		"information":   "信息泄露",
		"disclosure":    "信息泄露",
		"rce":           "代码执行",
		"execution":     "代码执行",
		"command":       "命令执行",
		"traversal":     "目录遍历",
		"ssrf":          "SSRF",
		"csrf":          "CSRF",
		"xxe":           "XXE",
		"deserialize":   "反序列化",
		"logic":         "逻辑漏洞",
		"privilege":     "越权访问",
		"weak":          "弱口令",
		"password":      "弱口令",
	}

	for keyword, vulnType := range keywords {
		if strings.Contains(combined, keyword) {
			return vc.Classify(vulnType)
		}
	}

	return vc.Classify("其他")
}

func (vc *VulnClassifier) GetAllClassifications() map[string]VulnClassification {
	return vc.classifications
}

func (vc *VulnClassifier) GetCategoryStats(vulns []map[string]interface{}) map[string]int {
	stats := map[string]int{}

	for _, vuln := range vulns {
		vulnType := ""
		if t, ok := vuln["type"].(string); ok {
			vulnType = t
		}

		classification := vc.Classify(vulnType)
		stats[classification.Category]++
	}

	return stats
}

func (vc *VulnClassifier) CalculateTotalReward(vulns []map[string]interface{}) float64 {
	total := 0.0

	for _, vuln := range vulns {
		vulnType := ""
		if t, ok := vuln["type"].(string); ok {
			vulnType = t
		}

		classification := vc.Classify(vulnType)
		total += classification.AvgReward
	}

	return total
}

func (vc *VulnClassifier) PrintClassificationReport(vulns []map[string]interface{}) {
	fmt.Println("\n📊 漏洞分类统计：")
	fmt.Println("=====================================")

	categoryStats := vc.GetCategoryStats(vulns)
	for category, count := range categoryStats {
		fmt.Printf("  %s: %d\n", category, count)
	}

	totalReward := vc.CalculateTotalReward(vulns)
	fmt.Printf("\n💰 预估总奖金: ¥%.0f\n", totalReward)

	fmt.Println("=====================================")
}
