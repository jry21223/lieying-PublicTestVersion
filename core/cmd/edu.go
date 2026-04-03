package cmd

import (
	"fmt"

	"github.com/kunlun-sec/lunying/pkg/edu"
	"github.com/spf13/cobra"
)

var (
	eduLevel     string
	eduProvince  string
	eduList      bool
	eduStats     bool
)

var eduCmd = &cobra.Command{
	Use:   "edu",
	Short: "教育SRC相关功能",
	Long: `教育行业安全测试相关功能：
  - 高校信息系统查询
  - 批量漏洞扫描
  - 统计数据`,
	Run: func(cmd *cobra.Command, args []string) {
		if eduList {
			listUniversities(eduLevel, eduProvince)
			return
		}

		if eduStats {
			showEduStats()
			return
		}

		// 默认显示帮助
		cmd.Help()
	},
}

var eduScanCmd = &cobra.Command{
	Use:   "scan [domain]",
	Short: "扫描单个高校域名",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]

		fmt.Println("=====================================")
		fmt.Println("  猎影教育SRC扫描模块")
		fmt.Println("  昆仑安全实验室(前逍遥安全实验室-逍遥)")
		fmt.Println("=====================================")
		fmt.Printf("  目标: %s\n", domain)
		fmt.Println("=====================================")

		scanner := edu.NewEduPOCScanner(domain)
		results, err := scanner.Scan()
		if err != nil {
			fmt.Printf("扫描失败: %v\n", err)
			return
		}

		if len(results) == 0 {
			fmt.Println("\n未发现漏洞")
		} else {
			fmt.Printf("\n发现 %d 个漏洞:\n", len(results))
			for _, r := range results {
				fmt.Printf("  [%s] %s: %s\n", r.SystemType, r.VulnType, r.URL)
			}
		}
	},
}

var eduBatchCmd = &cobra.Command{
	Use:   "batch",
	Short: "批量扫描高校",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=====================================")
		fmt.Println("  猎影教育SRC批量扫描模块")
		fmt.Println("  昆仑安全实验室(前逍遥安全实验室-逍遥)")
		fmt.Println("=====================================")

		scanner := edu.NewBatchScannerWithFilter(eduLevel, eduProvince, 10)
		scanner.ScanAll()
	},
}

func init() {
	eduCmd.Flags().BoolVarP(&eduList, "list", "l", false, "列出高校信息")
	eduCmd.Flags().BoolVarP(&eduStats, "stats", "s", false, "显示统计数据")
	eduCmd.Flags().StringVarP(&eduLevel, "level", "", "", "按级别过滤 (985/211/本科/专科)")
	eduCmd.Flags().StringVarP(&eduProvince, "province", "", "", "按省份过滤")

	eduBatchCmd.Flags().StringVarP(&eduLevel, "level", "", "", "按级别扫描 (985/211)")
	eduBatchCmd.Flags().StringVarP(&eduProvince, "province", "", "", "按省份扫描")

	eduCmd.AddCommand(eduScanCmd)
	eduCmd.AddCommand(eduBatchCmd)
}

func listUniversities(level, province string) {
	db := edu.NewUniversityDB()

	var universities []edu.University
	if level != "" {
		universities = db.GetByLevel(level)
	} else if province != "" {
		universities = db.GetByProvince(province)
	} else {
		universities = db.GetAll()
	}

	fmt.Println("=====================================")
	fmt.Println("  高校信息列表")
	fmt.Println("=====================================")
	fmt.Printf("  总数: %d 所高校\n\n", len(universities))

	for _, uni := range universities {
		fmt.Printf("🏫 %s\n", uni.Name)
		fmt.Printf("   级别: %s | 省份: %s\n", uni.Level, uni.Province)
		fmt.Printf("   域名: %v\n", uni.Domains)
		fmt.Println()
	}
}

func showEduStats() {
	db := edu.NewUniversityDB()
	universities := db.GetAll()

	stats := map[string]int{}
	provinceStats := map[string]int{}

	for _, uni := range universities {
		stats[uni.Level]++
		provinceStats[uni.Province]++
	}

	fmt.Println("=====================================")
	fmt.Println("  教育SRC统计信息")
	fmt.Println("=====================================")
	fmt.Println("\n按级别统计:")
	for level, count := range stats {
		fmt.Printf("  %s: %d 所\n", level, count)
	}

	fmt.Println("\n按省份统计:")
	for province, count := range provinceStats {
		fmt.Printf("  %s: %d 所\n", province, count)
	}

	fmt.Printf("\n总计: %d 所高校\n", len(universities))
	fmt.Println("=====================================")
}