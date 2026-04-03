package cmd

import (
	"fmt"
	"os"

	"github.com/kunlun-sec/lunying/internal/app"
	"github.com/spf13/cobra"
)

var (
	serverPort    string
	serverDBPath  string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动API服务器",
	Long: `启动猎影平台API服务器，提供RESTful API接口。
支持目标管理、漏洞管理、资产管理等功能。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=====================================")
		fmt.Println("  猎影API服务器")
		fmt.Println("  昆仑安全实验室(前逍遥安全实验室-逍遥)")
		fmt.Println("=====================================")
		fmt.Printf("  端口: %s\n", serverPort)
		fmt.Printf("  数据库: %s\n", serverDBPath)
		fmt.Println("=====================================")

		// 确保数据目录存在
		if err := os.MkdirAll("data", 0755); err != nil {
			fmt.Printf("创建数据目录失败: %v\n", err)
			os.Exit(1)
		}

		app.RunServer()
	},
}

func init() {
	serverCmd.Flags().StringVarP(&serverPort, "port", "p", ":8081", "服务器监听端口")
	serverCmd.Flags().StringVarP(&serverDBPath, "db", "d", "data/lieying.db", "数据库路径")
}