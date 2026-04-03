package cmd

import (
	"fmt"
	"os"

	"github.com/kunlun-sec/lunying/pkg/scan"
	"github.com/spf13/cobra"
)

var (
	scanOutput string
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
  - 未授权访问检测`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		manager := scan.NewScanManager(target)

		fmt.Println("=====================================")
		fmt.Println("  猎影漏洞扫描模块")
		fmt.Println("  昆仑安全实验室(前逍遥安全实验室-逍遥)")
		fmt.Println("=====================================")
		fmt.Printf("  目标: %s\n", target)
		fmt.Println("=====================================")
		fmt.Println()

		// 如果没有指定具体模块，执行全面扫描
		if !scanNuclei && !scanSQLi && !scanXSS && !scanUpload && !scanUnauth {
			result, err := manager.RunFullScan()
			if err != nil {
				fmt.Fprintf(os.Stderr, "扫描失败: %v\n", err)
				os.Exit(1)
			}

			outputResult(result, scanOutput)
			return
		}

		// 执行指定的模块
		if scanNuclei {
			fmt.Println("[+] Nuclei漏洞扫描...")
			manager.RunNucleiScan()
		}

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
	},
}

func init() {
	scanCmd.Flags().StringVarP(&scanOutput, "output", "o", "", "输出文件路径 (JSON格式)")
	scanCmd.Flags().BoolVarP(&scanNuclei, "nuclei", "n", false, "仅执行Nuclei扫描")
	scanCmd.Flags().BoolVarP(&scanSQLi, "sqli", "s", false, "仅执行SQL注入检测")
	scanCmd.Flags().BoolVarP(&scanXSS, "xss", "x", false, "仅执行XSS检测")
	scanCmd.Flags().BoolVarP(&scanUpload, "upload", "u", false, "仅执行文件上传检测")
	scanCmd.Flags().BoolVarP(&scanUnauth, "unauth", "a", false, "仅执行未授权访问检测")
}