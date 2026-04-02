package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lieying/engine/internal/config"
	"github.com/lieying/engine/pkg/ai"
	"github.com/lieying/engine/pkg/network"
)

func main() {
	fmt.Println("=====================================")
	fmt.Println("  猎影渗透测试工具 - CLI模式")
	fmt.Println("  Lieying Penetration Testing Tool")
	fmt.Println("=====================================")
	fmt.Println()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	nm, err := network.NewNetworkManager(cfg.Network)
	if err != nil {
		log.Fatalf("Failed to init network manager: %v", err)
	}

	fmt.Println("📡 网络状态:")
	fmt.Printf("   模式: %s\n", cfg.Network.Mode)
	fmt.Printf("   离线模式: %v\n", nm.IsOffline())
	fmt.Println()

	aiManager := ai.NewAIClientManager(cfg.AI, nm.GetClient())
	fmt.Println("🤖 AI状态:")
	fmt.Printf("   模式: %s\n", cfg.AI.Mode)
	fmt.Printf("   可用: %v\n", aiManager.IsAvailable())
	fmt.Println()

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "scan":
		runScan()
	case "recon":
		runRecon()
	case "ai":
		runAI(aiManager)
	case "status":
		fmt.Println("✅ 所有系统状态检查完成")
	default:
		fmt.Printf("❌ 未知命令: %s\n", os.Args[1])
		printHelp()
	}
}

func printHelp() {
	fmt.Println("使用方法:")
	fmt.Println("  lieying scan    - 漏洞扫描")
	fmt.Println("  lieying recon   - 信息收集")
	fmt.Println("  lieying ai      - AI助手")
	fmt.Println("  lieying status  - 状态检查")
	fmt.Println()
}

func runScan() {
	fmt.Println("🔍 启动漏洞扫描模块...")
	fmt.Println("   (此功能正在开发中)")
}

func runRecon() {
	fmt.Println("🕵️ 启动信息收集模块...")
	fmt.Println("   (此功能正在开发中)")
}

func runAI(aiManager *ai.AIClientManager) {
	fmt.Println("💬 AI助手模式")
	fmt.Println("   (此功能正在开发中)")
}
