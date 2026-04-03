package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "1.0.0"
	BuildDate = "2026-04-02"
	Author    = "昆仑安全实验室(前逍遥安全实验室-逍遥)"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=====================================")
		fmt.Println("  猎影渗透测试平台 CLI")
		fmt.Println("=====================================")
		fmt.Printf("  版本:     %s\n", Version)
		fmt.Printf("  构建日期: %s\n", BuildDate)
		fmt.Printf("  作者:     %s\n", Author)
		fmt.Println("=====================================")
		fmt.Println()
		fmt.Println("功能模块:")
		fmt.Println("  ✅ 信息收集 (recon)")
		fmt.Println("  ✅ 漏洞扫描 (scan)")
		fmt.Println("  ✅ AI助手 (ai)")
		fmt.Println("  ✅ 教育SRC (edu)")
		fmt.Println("  ✅ 配置管理 (config)")
		fmt.Println("  ✅ API服务 (server)")
		fmt.Println("=====================================")
	},
}

func init() {
	// 无额外flags
}