package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kunlun-sec/lunying/pkg/recon"
	"github.com/spf13/cobra"
)

var (
	reconOutput     string
	reconSubdomain  bool
	reconPort       bool
	reconFingerprint bool
	reconDir        bool
)

var reconCmd = &cobra.Command{
	Use:   "recon [target]",
	Short: "信息收集",
	Long: `对目标进行全面的信息收集，包括：
  - 子域名枚举
  - 端口扫描
  - Web指纹识别
  - 目录扫描`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]

		manager := recon.NewReconManager(target)

		fmt.Println("=====================================")
		fmt.Println("  猎影信息收集模块")
		fmt.Println("  昆仑安全实验室(前逍遥安全实验室-逍遥)")
		fmt.Println("=====================================")
		fmt.Printf("  目标: %s\n", target)
		fmt.Println("=====================================")
		fmt.Println()

		// 如果没有指定具体模块，执行全面扫描
		if !reconSubdomain && !reconPort && !reconFingerprint && !reconDir {
			result, err := manager.RunFullRecon()
			if err != nil {
				fmt.Fprintf(os.Stderr, "信息收集失败: %v\n", err)
				os.Exit(1)
			}

			outputResult(result, reconOutput)
			return
		}

		// 执行指定的模块
		rootDomain := recon.ExtractRootDomain(target)

		if reconSubdomain {
			fmt.Println("[+] 子域名枚举...")
			manager.RunSubdomainEnum(rootDomain)
		}

		if reconPort {
			fmt.Println("[+] 端口扫描...")
			manager.RunPortScan(target)
		}

		if reconFingerprint {
			fmt.Println("[+] Web指纹识别...")
			manager.RunFingerprint(target)
		}

		if reconDir {
			fmt.Println("[+] 目录扫描...")
			manager.RunDirScan(target)
		}

		result := manager.GetResult()
		outputResult(result, reconOutput)
	},
}

func init() {
	reconCmd.Flags().StringVarP(&reconOutput, "output", "o", "", "输出文件路径 (JSON格式)")
	reconCmd.Flags().BoolVarP(&reconSubdomain, "subdomain", "s", false, "仅执行子域名枚举")
	reconCmd.Flags().BoolVarP(&reconPort, "port", "p", false, "仅执行端口扫描")
	reconCmd.Flags().BoolVarP(&reconFingerprint, "fingerprint", "f", false, "仅执行指纹识别")
	reconCmd.Flags().BoolVarP(&reconDir, "dir", "d", false, "仅执行目录扫描")
}

func outputResult(data interface{}, outputPath string) {
	if outputPath != "" {
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "JSON序列化失败: %v\n", err)
			return
		}
		err = os.WriteFile(outputPath, jsonData, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "写入文件失败: %v\n", err)
			return
		}
		fmt.Printf("\n结果已保存到: %s\n", outputPath)
	}
}