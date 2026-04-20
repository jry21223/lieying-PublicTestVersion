package scan

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type NucleiScanner struct {
	targets    []string
	templates  []string
	outputFile string
	severity   []string
}

type NucleiResult struct {
	Template      string `json:"template"`
	TemplateURL   string `json:"template-url"`
	TemplateID    string `json:"template-id"`
	TemplatePath  string `json:"template-path"`
	Info          NucleiInfo `json:"info"`
	Type          string `json:"type"`
	Host          string `json:"host"`
	Matched       string `json:"matched-at"`
	Extracted     string `json:"extracted-results"`
	IP            string `json:"ip"`
	Timestamp     string `json:"timestamp"`
	MatcherStatus bool   `json:"matcher-status"`
	MatcherName   string `json:"matcher-name"`
}

type NucleiInfo struct {
	Name           string   `json:"name"`
	Author         string   `json:"author"`
	Tags           []string `json:"tags"`
	Description    string   `json:"description"`
	Reference      []string `json:"reference"`
	Severity       string   `json:"severity"`
	Classification NucleiClassification `json:"classification"`
}

type NucleiClassification struct {
	CVEID       []string `json:"cve-id"`
	CWEID       []string `json:"cwe-id"`
	CVSSScore   float64  `json:"cvss-score"`
	CVSSMetrics string   `json:"cvss-metrics"`
}

func NewNucleiScanner() *NucleiScanner {
	return &NucleiScanner{
		targets:    []string{},
		templates:  []string{},
		severity:   []string{"critical", "high", "medium", "low"},
		outputFile: filepath.Join(os.TempDir(), "nuclei_output.json"),
	}
}

func (ns *NucleiScanner) AddTarget(target string) {
	ns.targets = append(ns.targets, target)
}

func (ns *NucleiScanner) SetTemplates(templates []string) {
	ns.templates = templates
}

func (ns *NucleiScanner) SetSeverity(severity []string) {
	ns.severity = severity
}

func (ns *NucleiScanner) Scan() ([]NucleiResult, error) {
	if len(ns.targets) == 0 {
		return []NucleiResult{}, fmt.Errorf("没有指定扫描目标")
	}

	fmt.Println("🔍 开始Nuclei漏洞扫描...")
	fmt.Printf("📊 扫描目标数量: %d\n", len(ns.targets))

	if !ns.isNucleiInstalled() {
		fmt.Println("⚠️  Nuclei未安装，尝试自动安装...")
		if err := ns.installNuclei(); err != nil {
			return []NucleiResult{}, fmt.Errorf("Nuclei安装失败: %v", err)
		}
	}

	args := []string{
		"-j",
		"-o", ns.outputFile,
		"-severity", strings.Join(ns.severity, ","),
		"-c", "25",
		"-rate-limit", "100",
		"-bulk-size", "25",
		"-timeout", "5",
		"-retries", "1",
		"-stats",
		"-stats-interval", "5",
	}

	if len(ns.templates) > 0 {
		for _, template := range ns.templates {
			args = append(args, "-t", template)
		}
	}

	for _, target := range ns.targets {
		args = append(args, "-u", target)
	}

	cmd := exec.Command("nuclei", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("🚀 执行Nuclei扫描...")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Nuclei扫描警告: %v\n", err)
	}

	results, err := ns.parseResults()
	if err != nil {
		return []NucleiResult{}, err
	}

	fmt.Printf("✅ Nuclei扫描完成，发现 %d 个漏洞\n", len(results))
	return results, nil
}

func (ns *NucleiScanner) isNucleiInstalled() bool {
	_, err := exec.LookPath("nuclei")
	return err == nil
}

func (ns *NucleiScanner) installNuclei() error {
	fmt.Println("📥 正在安装Nuclei...")
	
	cmd := exec.Command("go", "install", "-v", "github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return err
	}
	
	fmt.Println("✅ Nuclei安装完成")
	fmt.Println("📥 正在更新Nuclei模板...")
	
	updateCmd := exec.Command("nuclei", "-ut")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	
	if err := updateCmd.Run(); err != nil {
		fmt.Printf("模板更新警告: %v\n", err)
	}
	
	return nil
}

func (ns *NucleiScanner) parseResults() ([]NucleiResult, error) {
	data, err := os.ReadFile(ns.outputFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []NucleiResult{}, nil
		}
		return []NucleiResult{}, err
	}

	lines := strings.Split(string(data), "\n")
	var results []NucleiResult

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var result NucleiResult
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			fmt.Printf("解析结果错误: %v\n", err)
			continue
		}
		results = append(results, result)
	}

	os.Remove(ns.outputFile)
	return results, nil
}

func (ns *NucleiScanner) UpdateTemplates() error {
	if !ns.isNucleiInstalled() {
		return fmt.Errorf("Nuclei未安装")
	}

	fmt.Println("📥 正在更新Nuclei模板...")
	cmd := exec.Command("nuclei", "-ut")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}
