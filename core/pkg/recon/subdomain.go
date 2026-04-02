package recon

import (
	"fmt"
	"net"
	"strings"
	"sync"
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

	se.quickCheck()
	se.dictionaryBrute()

	return se.subdomains, nil
}

func (se *SubdomainEnumerator) quickCheck() {
	commonSubdomains := []string{
		"www", "mail", "ftp", "admin", "api", "dev", "test",
		"staging", "blog", "shop", "store", "webmail", "cpanel",
		"whm", "webdisk", "webstats", "m", "mobile", "wap",
		"img", "images", "video", "cdn", "static", "assets",
		"portal", "member", "members", "user", "users", "account",
		"login", "signup", "register", "auth", "oauth", "sso",
		"api", "api2", "api3", "rest", "graphql", "soap",
		"old", "new", "beta", "alpha", "demo", "docs",
		"help", "support", "contact", "about", "career", "jobs",
		"news", "press", "media", "status", "monitor", "stats",
		"cache", "proxy", "vpn", "ns1", "ns2", "dns",
	}

	for _, sub := range commonSubdomains {
		fullDomain := fmt.Sprintf("%s.%s", sub, se.domain)
		if se.checkDomain(fullDomain) {
			se.addSubdomain(fullDomain)
			fmt.Printf("✅ 发现子域: %s\n", fullDomain)
		}
	}
}

func (se *SubdomainEnumerator) dictionaryBrute() {
	commonWords := []string{
		"web", "data", "db", "database", "backup", "bak", "old",
		"new", "temp", "tmp", "dev", "stage", "prod", "production",
		"test", "testing", "uat", "sit", "qa", "quality",
		"admin", "adm", "manager", "manage", "cms", "wp", "blog",
		"forum", "bbs", "community", "social", "chat", "message",
		"shop", "store", "ecommerce", "cart", "checkout", "pay",
		"api", "service", "services", "soap", "rest", "graphql",
		"file", "files", "upload", "download", "assets", "static",
		"img", "image", "images", "pic", "photo", "video", "media",
		"doc", "docs", "document", "documents", "wiki", "help",
		"support", "contact", "about", "info", "information",
		"login", "signin", "auth", "oauth", "sso", "account",
		"user", "users", "member", "members", "profile", "setting",
		"config", "conf", "cfg", "setup", "install",
	}

	for _, word := range commonWords {
		fullDomain := fmt.Sprintf("%s.%s", word, se.domain)
		if se.checkDomain(fullDomain) {
			se.addSubdomain(fullDomain)
			fmt.Printf("✅ 发现子域: %s\n", fullDomain)
		}
	}
}

func (se *SubdomainEnumerator) checkDomain(domain string) bool {
	_, err := net.LookupHost(domain)
	return err == nil
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
	if len(parts) >= 2 {
		return strings.Join(parts[len(parts)-2:], ".")
	}
	return domain
}
