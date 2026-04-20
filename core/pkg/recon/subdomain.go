package recon

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type SubdomainEnumerator struct {
	domain     string
	subdomains []string
	mu         sync.Mutex
}

func NewSubdomainEnumerator(domain string) *SubdomainEnumerator {
	return &SubdomainEnumerator{
		domain:     domain,
		subdomains: []string{},
	}
}

func (se *SubdomainEnumerator) addSubdomain(subdomain string) {
	se.mu.Lock()
	defer se.mu.Unlock()
	for _, s := range se.subdomains {
		if s == subdomain {
			return
		}
	}
	se.subdomains = append(se.subdomains, subdomain)
}

func (se *SubdomainEnumerator) Enumerate() ([]string, error) {
	fmt.Printf("🔍 开始枚举子域: %s\n", se.domain)

	// 优先使用外部工具（更高效）
	if se.useSubfinder() {
		se.verifyWithHttpx()
		return se.subdomains, nil
	}

	// 外部工具不可用，使用内置字典枚举
	fmt.Println("⚠️  subfinder 未安装，使用内置字典枚举（较慢）")
	se.detectWildcardAndEnumerate()

	return se.subdomains, nil
}

// useSubfinder 使用 subfinder 进行被动子域名枚举
func (se *SubdomainEnumerator) useSubfinder() bool {
	// 检查 subfinder 是否可用
	_, err := exec.LookPath("subfinder")
	if err != nil {
		return false
	}

	fmt.Println("📡 使用 subfinder 进行被动子域名枚举...")

	// 执行 subfinder
	cmd := exec.Command("subfinder", "-d", se.domain, "-silent", "-json")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("subfinder 执行错误: %v\n", err)
		return false
	}

	// 解析 JSON 输出
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var result struct {
			Host string `json:"host"`
		}
		if err := json.Unmarshal([]byte(line), &result); err == nil && result.Host != "" {
			se.addSubdomain(result.Host)
		}
	}

	fmt.Printf("✅ subfinder 发现 %d 个子域\n", len(se.subdomains))
	return len(se.subdomains) > 0
}

// verifyWithHttpx 使用 httpx 验证子域名存活状态
func (se *SubdomainEnumerator) verifyWithHttpx() bool {
	if len(se.subdomains) == 0 {
		return false
	}

	// 检查 httpx 是否可用
	_, err := exec.LookPath("httpx")
	if err != nil {
		fmt.Println("⚠️  httpx 未安装，跳过存活验证")
		return true
	}

	originalCount := len(se.subdomains)
	fmt.Println("🔍 使用 httpx 验证子域名存活状态...")

	// 写入临时文件
	tmpFile := fmt.Sprintf("/tmp/subdomains_%s.txt", se.domain)
	data := strings.Join(se.subdomains, "\n")
	os.WriteFile(tmpFile, []byte(data), 0644)

	// 执行 httpx
	cmd := exec.Command("httpx", "-l", tmpFile, "-silent", "-status-code", "-json")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("httpx 执行警告: %v\n", err)
		os.Remove(tmpFile)
		return true
	}

	os.Remove(tmpFile)

	// 解析存活子域名
	var aliveSubdomains []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var result struct {
			URL        string `json:"url"`
			StatusCode int    `json:"status_code"`
		}
		if err := json.Unmarshal([]byte(line), &result); err == nil && result.StatusCode > 0 {
			// 提取主机名
			host := strings.TrimPrefix(result.URL, "http://")
			host = strings.TrimPrefix(host, "https://")
			host = strings.Split(host, "/")[0]
			host = strings.Split(host, ":")[0]
			aliveSubdomains = append(aliveSubdomains, host)
		}
	}

	// 更新子域名列表（只保留存活的）
	se.mu.Lock()
	se.subdomains = aliveSubdomains
	se.mu.Unlock()

	fmt.Printf("✅ httpx 验证完成: %d/%d 子域名存活\n", len(aliveSubdomains), originalCount)
	for _, sub := range aliveSubdomains {
		fmt.Printf("  ✅ %s\n", sub)
	}

	return true
}

// detectWildcardAndEnumerate 内置枚举（带泛解析检测）
func (se *SubdomainEnumerator) detectWildcardAndEnumerate() {
	// 检测泛解析
	hasWildcard := se.detectWildcard()
	if hasWildcard {
		fmt.Println("⚠️  检测到 DNS 泛解析，将进行 HTTP 验证（较慢）")
	}

	se.quickCheck(hasWildcard)
	se.dictionaryBrute(hasWildcard)
}

// detectWildcard 检测泛解析
func (se *SubdomainEnumerator) detectWildcard() bool {
	randomDomain := fmt.Sprintf("random%snonexistent.%s", time.Now().Unix(), se.domain)
	_, err := exec.Command("host", randomDomain).Output()
	return err == nil
}

// quickCheck 快速检查常见子域名
func (se *SubdomainEnumerator) quickCheck(hasWildcard bool) {
	commonSubdomains := []string{
		"www", "mail", "ftp", "admin", "api", "dev", "test",
		"staging", "blog", "shop", "store", "webmail", "cpanel",
		"m", "mobile", "img", "cdn", "static", "assets",
		"portal", "bbs", "forum", "news", "wiki", "docs",
	}

	for _, sub := range commonSubdomains {
		fullDomain := fmt.Sprintf("%s.%s", sub, se.domain)
		if se.verifySubdomain(fullDomain, hasWildcard) {
			se.addSubdomain(fullDomain)
			fmt.Printf("✅ 发现子域: %s\n", fullDomain)
		}
	}
}

// dictionaryBrute 字典枚举
func (se *SubdomainEnumerator) dictionaryBrute(hasWildcard bool) {
	commonWords := []string{
		"web", "data", "db", "backup", "old", "new", "temp",
		"dev", "stage", "prod", "test", "uat", "sit", "qa",
		"admin", "manager", "cms", "wp", "blog", "forum",
		"shop", "store", "api", "service", "file", "upload",
		"img", "image", "video", "media", "doc", "help",
		"login", "auth", "account", "user", "member", "config",
	}

	for _, word := range commonWords {
		fullDomain := fmt.Sprintf("%s.%s", word, se.domain)
		if se.verifySubdomain(fullDomain, hasWildcard) {
			se.addSubdomain(fullDomain)
			fmt.Printf("✅ 发现子域: %s\n", fullDomain)
		}
	}
}

// verifySubdomain 验证子域名
func (se *SubdomainEnumerator) verifySubdomain(domain string, hasWildcard bool) bool {
	// DNS 解析
	_, err := exec.Command("host", domain).Output()
	if err != nil {
		return false
	}

	// 如果有泛解析，需要 HTTP 验证
	if hasWildcard {
		cmd := exec.Command("curl", "-s", "-o", "/dev/null", "-w", "%{http_code}",
			"--max-time", "3", fmt.Sprintf("http://%s", domain))
		output, err := cmd.Output()
		if err != nil {
			return false
		}
		statusCode := strings.TrimSpace(string(output))
		return statusCode != "000" && statusCode != ""
	}

	return true
}

func (se *SubdomainEnumerator) GetSubdomains() []string {
	se.mu.Lock()
	defer se.mu.Unlock()
	result := make([]string, len(se.subdomains))
	copy(result, se.subdomains)
	return result
}

func ExtractRootDomain(domain string) string {
	parts := strings.Split(domain, ".")

	// 处理特殊顶级域（中国教育、政府等）
	if len(parts) >= 3 {
		tld := strings.Join(parts[len(parts)-2:], ".")
		specialTLDs := []string{"edu.cn", "gov.cn", "ac.cn", "com.cn", "net.cn", "org.cn"}
		for _, special := range specialTLDs {
			if tld == special {
				return strings.Join(parts[len(parts)-3:], ".")
			}
		}
	}

	// 处理其他国家的特殊顶级域（如 .co.uk, .ac.uk）
	if len(parts) >= 3 {
		tld := parts[len(parts)-1]
		secondLevel := parts[len(parts)-2]
		if tld == "uk" && (secondLevel == "co" || secondLevel == "ac" || secondLevel == "gov") {
			return strings.Join(parts[len(parts)-3:], ".")
		}
	}

	// 普通域名，返回最后两部分
	if len(parts) >= 2 {
		return strings.Join(parts[len(parts)-2:], ".")
	}
	return domain
}
