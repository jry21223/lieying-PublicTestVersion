package main

import (
	"fmt"
	"os"

	"github.com/kunlun-sec/lunying/internal/app"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		app.RunServer()
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("猎影渗透测试平台 - Core")
	fmt.Println("Lieying Penetration Testing Platform")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  lieying server   - 启动API服务器")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lieying server")
}
