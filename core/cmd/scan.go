package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kunlun-sec/lunying/pkg/scan"
	"github.com/kunlun-sec/lunying/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	scanOutput string
	scanInput  string
	scanNuclei bool
	scanSQLi   bool
	scanXSS    bool
	scanUpload bool
	scanUnauth bool
)

var scanCmd = &cobra.Command{
	Use:   "scan [target]",
	Short: "漏洞扫描",
	Long: `对目标进行漏洞扫描，包括：
  - Nuclei漏洞扫描
  - SQL注入检测
  - XSS检测
  - 文件上传漏洞检测
  - 未授权访问检测

支持单目标扫描或从JSON文件批量读取目标（使用 -i 参数）`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var targets []string

		// 从JSON文件读取目标
		if scanInput != "" {
			targets = loadTargetsFromJSON(scanInput)
			if len(targets) == 0 {
				fmt.Fprintf(os.Stderr, "未从文件中找到有效目标\n")
				os.Exit(1)
			}
		} else if len(args) > 0 {
			targets = []string{args[0]}
		} else {
			fmt.Fprintf(os.Stderr, "请指定目标或使用 -i 参数提供JSON文件\n")
			os.Exit(1)
		}

		fmt.Println("=====================================")
		fmt.Println("  猎影漏洞扫描模块")
		fmt.Println("  昆仑安全实验室(前逍遥安全实验室-逍遥)")
		fmt.Println("=====================================")
		fmt.Printf("  原始目标数量: %d\n", len(targets))

		// 用 httpx 过滤存活目标
		targets = filterAliveTargets(targets)
		if len(targets) == 0 {
			fmt.Fprintf(os.Stderr, "没有存活的目标可扫描\n")
			os.Exit(1)
		}

		fmt.Printf("  存活目标数量: %d\n", len(targets))
		fmt.Println("=====================================")
		fmt.Println()

		// Nuclei始终批量扫描（更高效）
		if scanNuclei {
			runBatchNucleiScan(targets)
		}

		// 其他模块逐个扫描
		if scanSQLi || scanXSS || scanUpload || scanUnauth {
			for i, target := range targets {
				fmt.Printf("[%d/%d] 扫描目标: %s\n", i+1, len(targets), target)

				manager := scan.NewScanManager(target)

				if scanSQLi {
					fmt.Println("[+] SQL注入检测...")
					manager.RunSQLiScan()
				}

				if scanXSS {
					fmt.Println("[+] XSS检测...")
					manager.RunXSSScan()
				}

				if scanUpload {
					fmt.Println("[+] 文件上传漏洞检测...")
					manager.RunUploadScan()
				}

				if scanUnauth {
					fmt.Println("[+] 未授权访问检测...")
					manager.RunUnauthScan()
				}

				result := manager.GetResult()
				outputResult(result, scanOutput)
			}
		}

		// 未指定任何模块时执行全面扫描
		if !scanNuclei && !scanSQLi && !scanXSS && !scanUpload && !scanUnauth {
			// nuclei批量扫描
			runBatchNucleiScan(targets)

			// 其他模块逐个扫描
			for i, target := range targets {
				fmt.Printf("[%d/%d] 扫描目标: %s\n", i+1, len(targets), target)

				manager := scan.NewScanManager(target)
				manager.RunSQLiScan()
				manager.RunXSSScan()
				manager.RunUploadScan()
				manager.RunUnauthScan()

				result := manager.GetResult()
				outputResult(result, scanOutput)
			}
		}
	},
}

func init() {
	scanCmd.Flags().StringVarP(&scanOutput, "output", "o", "", "输出文件路径 (JSON格式)")
	scanCmd.Flags().StringVarP(&scanInput, "input", "i", "", "从JSON文件读取目标 (支持recon输出)")
	scanCmd.Flags().BoolVarP(&scanNuclei, "nuclei", "n", false, "仅执行Nuclei扫描")
	scanCmd.Flags().BoolVarP(&scanSQLi, "sqli", "s", false, "仅执行SQL注入检测")
	scanCmd.Flags().BoolVarP(&scanXSS, "xss", "x", false, "仅执行XSS检测")
	scanCmd.Flags().BoolVarP(&scanUpload, "upload", "u", false, "仅执行文件上传检测")
	scanCmd.Flags().BoolVarP(&scanUnauth, "unauth", "a", false, "仅执行未授权访问检测")
}

// runBatchNucleiScan 批量Nuclei扫描，一次性扫描所有目标
func runBatchNucleiScan(targets []string) {
	fmt.Println("[+] Nuclei批量漏洞扫描...")
	scanner := scan.NewNucleiScanner()
	for _, target := range targets {
		scanner.AddTarget(target)
	}
	results, err := scanner.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Nuclei扫描错误: %v\n", err)
		return
	}

	fmt.Printf("[+] Nuclei批量扫描完成，发现 %d 个漏洞\n", len(results))

	// 按目标分组输出
	targetVulns := make(map[string][]scan.NucleiResult)
	for _, r := range results {
		targetVulns[r.Host] = append(targetVulns[r.Host], r)
	}

	for target, vulns := range targetVulns {
		fmt.Printf("\n目标: %s - 发现 %d 个漏洞\n", target, len(vulns))
		for _, v := range vulns {
			fmt.Printf("  [%s] %s\n", v.Info.Severity, v.Info.Name)
		}
	}

	// 输出到文件
	if scanOutput != "" {
		outputData := struct {
			TotalTargets int                 `json:"total_targets"`
			TotalVulns   int                 `json:"total_vulns"`
			Results      []scan.NucleiResult `json:"results"`
		}{
			TotalTargets: len(targets),
			TotalVulns:   len(results),
			Results:      results,
		}
		jsonData, err := json.MarshalIndent(outputData, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "JSON序列化失败: %v\n", err)
			return
		}
		err = os.WriteFile(scanOutput, jsonData, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "写入文件失败: %v\n", err)
			return
		}
		fmt.Printf("\n结果已保存到: %s\n", scanOutput)
	}
}

// loadTargetsFromJSON 从JSON文件中提取目标列表
func loadTargetsFromJSON(filePath string) []string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取文件失败: %v\n", err)
		return nil
	}

	var targets []string

	// 尝试解析为recon输出格式
	var reconResult struct {
		Target     string   `json:"target"`
		Subdomains []string `json:"subdomains"`
	}
	if err := json.Unmarshal(data, &reconResult); err == nil {
		// 添加主目标
		if reconResult.Target != "" {
			targets = append(targets, reconResult.Target)
		}
		// 添加子域名
		for _, sub := range reconResult.Subdomains {
			if sub != "" {
				if !utils.HasProtocol(sub) {
					sub = "https://" + sub
				}
				targets = append(targets, sub)
			}
		}
		if len(targets) > 0 {
			return targets
		}
	}

	// 尝试解析为简单字符串数组
	var simpleTargets []string
	if err := json.Unmarshal(data, &simpleTargets); err == nil {
		for _, t := range simpleTargets {
			if t != "" {
				if !utils.HasProtocol(t) {
					t = "https://" + t
				}
				targets = append(targets, t)
			}
		}
		return targets
	}

	// 尝试解析为每行一个URL的文本格式
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			if !utils.HasProtocol(line) {
				line = "https://" + line
			}
			targets = append(targets, line)
		}
	}

	return targets
}

// filterAliveTargets 用 httpx 过滤存活目标
func filterAliveTargets(targets []string) []string {
	// 检查 httpx 是否可用
	_, err := exec.LookPath("httpx")
	if err != nil {
		fmt.Println("⚠️  httpx 未安装，跳过存活验证")
		return targets
	}

	// 写入临时文件
	tmpFile := filepath.Join(os.TempDir(), "targets_tmp.txt")
	fileData := strings.Join(targets, "\n")
	os.WriteFile(tmpFile, []byte(fileData), 0644)

	// 执行 httpx
	cmd := exec.Command("httpx", "-l", tmpFile, "-silent", "-status-code", "-json")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("httpx 执行警告: %v\n", err)
		os.Remove(tmpFile)
		return targets
	}

	// 解析存活目标
	var alive []string
	outputLines := strings.Split(string(output), "\n")
	for _, line := range outputLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var result struct {
			URL        string `json:"url"`
			StatusCode int    `json:"status_code"`
		}
		if err := json.Unmarshal([]byte(line), &result); err == nil && result.StatusCode > 0 {
			alive = append(alive, result.URL)
		}
	}

	os.Remove(tmpFile)
	fmt.Printf("[+] httpx 验证完成: %d/%d 目标存活\n", len(alive), len(targets))
	return alive
}