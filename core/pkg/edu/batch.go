package edu

import (
	"fmt"
	"sync"
)

type BatchScanner struct {
	universities []University
	results      map[string][]EduPOCResult
	mu           sync.Mutex
}

func NewBatchScanner(universities []University) *BatchScanner {
	return &BatchScanner{
		universities: universities,
		results:      make(map[string][]EduPOCResult),
	}
}

// NewBatchScannerWithFilter 创建带过滤条件的批量扫描器（API用）
func NewBatchScannerWithFilter(level, province string, concurrency int) *BatchScanner {
	db := NewUniversityDB()
	var universities []University

	if level != "" {
		universities = db.GetByLevel(level)
	} else if province != "" {
		universities = db.GetByProvince(province)
	} else {
		universities = db.GetAll()
	}

	return &BatchScanner{
		universities: universities,
		results:      make(map[string][]EduPOCResult),
	}
}

func (bs *BatchScanner) ScanAll() error {
	fmt.Println("=====================================")
	fmt.Println("  开始批量扫描高校教务系统")
	fmt.Println("  昆仑安全实验室(前逍遥安全实验室-逍遥)")
	fmt.Println("=====================================")
	fmt.Printf("📊 扫描目标数量: %d 所高校\n\n", len(bs.universities))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // 限制并发数

	for _, uni := range bs.universities {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(u University) {
			defer wg.Done()
			defer func() { <-semaphore }()

			bs.scanUniversity(u)
		}(uni)
	}

	wg.Wait()

	bs.printSummary()
	return nil
}

func (bs *BatchScanner) scanUniversity(uni University) {
	fmt.Printf("🔍 正在扫描: %s\n", uni.Name)

	for _, domain := range uni.Domains {
		scanner := NewEduPOCScanner(domain)
		results, err := scanner.Scan()
		if err != nil {
			continue
		}

		if len(results) > 0 {
			bs.mu.Lock()
			bs.results[uni.Name] = results
			bs.mu.Unlock()
			fmt.Printf("✅ %s 发现 %d 个漏洞\n", uni.Name, len(results))
		}
	}
}

func (bs *BatchScanner) printSummary() {
	fmt.Println("\n=====================================")
	fmt.Println("  批量扫描完成！")
	fmt.Println("=====================================")

	totalVulns := 0
	for uniName, results := range bs.results {
		fmt.Printf("\n🏫 %s: %d 个漏洞\n", uniName, len(results))
		for _, r := range results {
			fmt.Printf("   - [%s] %s: %s\n", r.SystemType, r.VulnType, r.URL)
		}
		totalVulns += len(results)
	}

	fmt.Printf("\n📊 总计: %d 所高校，%d 个漏洞\n", len(bs.results), totalVulns)
	fmt.Println("=====================================")
}

func (bs *BatchScanner) GetResults() map[string][]EduPOCResult {
	return bs.results
}

func (bs *BatchScanner) ScanByLevel(level string) error {
	db := NewUniversityDB()
	universities := db.GetByLevel(level)
	
	if len(universities) == 0 {
		return fmt.Errorf("没有找到 %s 级别的高校", level)
	}

	bs.universities = universities
	return bs.ScanAll()
}

func (bs *BatchScanner) ScanByProvince(province string) error {
	db := NewUniversityDB()
	universities := db.GetByProvince(province)
	
	if len(universities) == 0 {
		return fmt.Errorf("没有找到 %s 省份的高校", province)
	}

	bs.universities = universities
	return bs.ScanAll()
}
