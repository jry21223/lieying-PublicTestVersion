package scan

import (
	"encoding/json"
	"fmt"
	"sync"
)

type ScanManager struct {
	target  string
	results *ScanResult
	mu      sync.Mutex
}

type ScanResult struct {
	Target          string                 `json:"target"`
	NucleiResults   []NucleiResult         `json:"nuclei_results,omitempty"`
	SQLiResults     []SQLiResult           `json:"sqli_results,omitempty"`
	XSSResults      []XSSResult            `json:"xss_results,omitempty"`
	UploadResults   []UploadResult         `json:"upload_results,omitempty"`
	UnauthResults   []UnauthResult         `json:"unauth_results,omitempty"`
	TotalVulns      int                    `json:"total_vulns"`
	CriticalCount   int                    `json:"critical_count"`
	HighCount       int                    `json:"high_count"`
	MediumCount     int                    `json:"medium_count"`
	LowCount        int                    `json:"low_count"`
}

func NewScanManager(target string) *ScanManager {
	return &ScanManager{
		target: target,
		results: &ScanResult{
			Target:        target,
			NucleiResults: []NucleiResult{},
			SQLiResults:   []SQLiResult{},
			XSSResults:    []XSSResult{},
			UploadResults: []UploadResult{},
			UnauthResults: []UnauthResult{},
		},
	}
}

func (sm *ScanManager) RunFullScan() (*ScanResult, error) {
	fmt.Println("=====================================")
	fmt.Println("  开始全面漏洞扫描")
	fmt.Println("  目标:", sm.target)
	fmt.Println("  昆仑安全实验室(前逍遥安全实验室-逍遥)")
	fmt.Println("=====================================")
	fmt.Println()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		sm.RunNucleiScan()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sm.RunSQLiScan()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sm.RunXSSScan()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sm.RunUploadScan()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sm.RunUnauthScan()
	}()

	wg.Wait()

	sm.calculateStats()

	fmt.Println()
	fmt.Println("=====================================")
	fmt.Println("  漏洞扫描完成！")
	fmt.Println("=====================================")
	fmt.Printf("  总漏洞数: %d\n", sm.results.TotalVulns)
	fmt.Printf("  严重: %d | 高危: %d | 中危: %d | 低危: %d\n",
		sm.results.CriticalCount, sm.results.HighCount,
		sm.results.MediumCount, sm.results.LowCount)
	fmt.Println("=====================================")

	return sm.results, nil
}

func (sm *ScanManager) RunNucleiScan() {
	fmt.Println("[1/5] 开始Nuclei漏洞扫描...")
	scanner := NewNucleiScanner()
	scanner.AddTarget(sm.target)
	results, err := scanner.Scan()
	if err != nil {
		fmt.Printf("Nuclei扫描错误: %v\n", err)
		return
	}

	sm.mu.Lock()
	sm.results.NucleiResults = results
	sm.mu.Unlock()

	fmt.Printf("[1/5] Nuclei扫描完成，发现 %d 个漏洞\n", len(results))
}

func (sm *ScanManager) RunSQLiScan() {
	fmt.Println("[2/5] 开始SQL注入检测...")
	scanner := NewSQLiScanner(sm.target)
	results, err := scanner.Scan()
	if err != nil {
		fmt.Printf("SQL注入检测错误: %v\n", err)
		return
	}

	sm.mu.Lock()
	sm.results.SQLiResults = results
	sm.mu.Unlock()

	fmt.Printf("[2/5] SQL注入检测完成，发现 %d 个漏洞\n", len(results))
}

func (sm *ScanManager) RunXSSScan() {
	fmt.Println("[3/5] 开始XSS检测...")
	scanner := NewXSSScanner(sm.target)
	results, err := scanner.Scan()
	if err != nil {
		fmt.Printf("XSS检测错误: %v\n", err)
		return
	}

	sm.mu.Lock()
	sm.results.XSSResults = results
	sm.mu.Unlock()

	fmt.Printf("[3/5] XSS检测完成，发现 %d 个漏洞\n", len(results))
}

func (sm *ScanManager) RunUploadScan() {
	fmt.Println("[4/5] 开始文件上传漏洞检测...")
	scanner := NewUploadScanner(sm.target)
	results, err := scanner.Scan()
	if err != nil {
		fmt.Printf("文件上传检测错误: %v\n", err)
		return
	}

	sm.mu.Lock()
	sm.results.UploadResults = results
	sm.mu.Unlock()

	fmt.Printf("[4/5] 文件上传漏洞检测完成，发现 %d 个漏洞\n", len(results))
}

func (sm *ScanManager) RunUnauthScan() {
	fmt.Println("[5/5] 开始未授权访问检测...")
	scanner := NewUnauthScanner(sm.target)
	results, err := scanner.Scan()
	if err != nil {
		fmt.Printf("未授权访问检测错误: %v\n", err)
		return
	}

	sm.mu.Lock()
	sm.results.UnauthResults = results
	sm.mu.Unlock()

	fmt.Printf("[5/5] 未授权访问检测完成，发现 %d 个漏洞\n", len(results))
}

func (sm *ScanManager) calculateStats() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, result := range sm.results.NucleiResults {
		sm.results.TotalVulns++
		switch result.Info.Severity {
		case "critical":
			sm.results.CriticalCount++
		case "high":
			sm.results.HighCount++
		case "medium":
			sm.results.MediumCount++
		case "low":
			sm.results.LowCount++
		}
	}

	for _, result := range sm.results.SQLiResults {
		sm.results.TotalVulns++
		switch result.Severity {
		case "Critical":
			sm.results.CriticalCount++
		case "High":
			sm.results.HighCount++
		}
	}

	for _, result := range sm.results.XSSResults {
		sm.results.TotalVulns++
		switch result.Severity {
		case "Critical":
			sm.results.CriticalCount++
		case "High":
			sm.results.HighCount++
		}
	}

	for _, result := range sm.results.UploadResults {
		sm.results.TotalVulns++
		switch result.Severity {
		case "Critical":
			sm.results.CriticalCount++
		case "High":
			sm.results.HighCount++
		}
	}

	for _, result := range sm.results.UnauthResults {
		sm.results.TotalVulns++
		switch result.Severity {
		case "Critical":
			sm.results.CriticalCount++
		case "High":
			sm.results.HighCount++
		}
	}
}

func (sm *ScanManager) GetResult() *ScanResult {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.results
}

func (sm *ScanManager) ToJSON() (string, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	data, err := json.MarshalIndent(sm.results, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
