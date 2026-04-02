package recon

import (
	"encoding/json"
	"fmt"
	"sync"
)

type ReconManager struct {
	target string
	result *ReconResult
	mu     sync.Mutex
}

type ReconResult struct {
	Target      string          `json:"target"`
	Subdomains  []string        `json:"subdomains"`
	Ports       []PortInfo      `json:"ports"`
	Fingerprint *FingerprintData `json:"fingerprint"`
	Directories []DirInfo       `json:"directories"`
}

type PortInfo struct {
	Port     int    `json:"port"`
	Service  string `json:"service"`
	Protocol string `json:"protocol"`
	Status   string `json:"status"`
}

type FingerprintData struct {
	URL        string   `json:"url"`
	CMS        string   `json:"cms"`
	CMSVersion string   `json:"cms_version"`
	WebServer  string   `json:"web_server"`
	Framework  string   `json:"framework"`
	WAF        string   `json:"waf"`
	TechStack  []string `json:"tech_stack"`
}

type DirInfo struct {
	URL           string `json:"url"`
	Path          string `json:"path"`
	StatusCode    int    `json:"status_code"`
	ContentLength int64  `json:"content_length"`
	Title         string `json:"title"`
	IsSensitive   bool   `json:"is_sensitive"`
}

func NewReconManager(target string) *ReconManager {
	return &ReconManager{
		target: target,
		result: &ReconResult{
			Target:      target,
			Subdomains:  []string{},
			Ports:       []PortInfo{},
			Fingerprint: nil,
			Directories: []DirInfo{},
		},
	}
}

func (rm *ReconManager) RunFullRecon() (*ReconResult, error) {
	fmt.Println("=====================================")
	fmt.Println("  开始全面信息收集")
	fmt.Println("  目标:", rm.target)
	fmt.Println("  昆仑安全实验室(前逍遥安全实验室-逍遥)")
	fmt.Println("=====================================")
	fmt.Println()

	rootDomain := ExtractRootDomain(rm.target)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		rm.runSubdomainEnum(rootDomain)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rm.runPortScan(rm.target)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rm.runFingerprint(rm.target)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rm.runDirScan(rm.target)
	}()

	wg.Wait()

	fmt.Println()
	fmt.Println("=====================================")
	fmt.Println("  信息收集完成！")
	fmt.Println("=====================================")
	fmt.Printf("  子域数量: %d\n", len(rm.result.Subdomains))
	fmt.Printf("  开放端口: %d\n", len(rm.result.Ports))
	fmt.Printf("  目录/文件: %d\n", len(rm.result.Directories))
	if rm.result.Fingerprint != nil {
		fmt.Printf("  CMS: %s\n", rm.result.Fingerprint.CMS)
		fmt.Printf("  Web服务器: %s\n", rm.result.Fingerprint.WebServer)
		fmt.Printf("  WAF: %s\n", rm.result.Fingerprint.WAF)
	}
	fmt.Println("=====================================")

	return rm.result, nil
}

func (rm *ReconManager) runSubdomainEnum(domain string) {
	fmt.Println("[1/4] 开始子域枚举...")
	enum := NewSubdomainEnumerator(domain)
	subdomains, err := enum.Enumerate()
	if err != nil {
		fmt.Printf("子域枚举错误: %v\n", err)
		return
	}

	rm.mu.Lock()
	rm.result.Subdomains = subdomains
	rm.mu.Unlock()

	fmt.Printf("[1/4] 子域枚举完成，发现 %d 个子域\n", len(subdomains))
}

func (rm *ReconManager) runPortScan(host string) {
	fmt.Println("[2/4] 开始端口扫描...")
	scanner := NewPortScanner(host)
	scanner.AddCommonPorts()
	openPorts, err := scanner.Scan()
	if err != nil {
		fmt.Printf("端口扫描错误: %v\n", err)
		return
	}

	rm.mu.Lock()
	for _, port := range openPorts {
		rm.result.Ports = append(rm.result.Ports, PortInfo{
			Port:     port,
			Service:  scanner.GetServiceName(port),
			Protocol: "tcp",
			Status:   "open",
		})
	}
	rm.mu.Unlock()

	fmt.Printf("[2/4] 端口扫描完成，发现 %d 个开放端口\n", len(openPorts))
}

func (rm *ReconManager) runFingerprint(url string) {
	fmt.Println("[3/4] 开始Web指纹识别...")
	fingerprinter := NewFingerprinter(url)
	result, err := fingerprinter.Fingerprint()
	if err != nil {
		fmt.Printf("指纹识别错误: %v\n", err)
		return
	}

	rm.mu.Lock()
	rm.result.Fingerprint = &FingerprintData{
		URL:        result.URL,
		CMS:        result.CMS,
		CMSVersion: result.CMSVersion,
		WebServer:  result.WebServer,
		Framework:  result.Framework,
		WAF:        result.WAF,
		TechStack:  result.TechStack,
	}
	rm.mu.Unlock()

	fmt.Println("[3/4] Web指纹识别完成")
}

func (rm *ReconManager) runDirScan(url string) {
	fmt.Println("[4/4] 开始目录扫描...")
	scanner := NewDirScanner(url)
	scanner.AddCommonPaths()
	foundDirs, err := scanner.Scan()
	if err != nil {
		fmt.Printf("目录扫描错误: %v\n", err)
		return
	}

	rm.mu.Lock()
	for _, dir := range foundDirs {
		rm.result.Directories = append(rm.result.Directories, DirInfo{
			URL:           dir.URL,
			Path:          dir.Path,
			StatusCode:    dir.StatusCode,
			ContentLength: dir.ContentLength,
			Title:         dir.Title,
			IsSensitive:   dir.IsSensitive,
		})
	}
	rm.mu.Unlock()

	fmt.Printf("[4/4] 目录扫描完成，发现 %d 个目录/文件\n", len(foundDirs))
}

func (rm *ReconManager) GetResult() *ReconResult {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	return rm.result
}

func (rm *ReconManager) ToJSON() (string, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	data, err := json.MarshalIndent(rm.result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
