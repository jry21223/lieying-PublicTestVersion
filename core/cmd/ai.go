package cmd

import (
	"fmt"
	"os"

	"github.com/kunlun-sec/lunying/pkg/ai"
	"github.com/spf13/cobra"
)

var (
	aiModel     string
	aiOllamaURL string
	aiProvider  string
	aiAPIKey    string
	aiAPIURL    string
)

var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "AI渗透测试助手",
	Long: `启动AI交互式助手模式，提供：
  - 漏洞分析和评估
  - POC生成和优化
  - 渗透测试技术指导
  - 安全工具使用建议
  - 攻击路径规划

支持多种 LLM 后端：
  - ollama (默认，本地免费)
  - openai (OpenAI API)
  - deepseek (DeepSeek API)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=====================================")
		fmt.Println("  猎影AI助手模块")
		fmt.Println("  昆仑安全实验室(前逍遥安全实验室-逍遥)")
		fmt.Println("=====================================")
		fmt.Println()

		// 优先使用环境变量
		apiKey := aiAPIKey
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}

		apiURL := aiAPIURL
		if apiURL == "" {
			apiURL = os.Getenv("OPENAI_BASE_URL")
		}

		var assistant *ai.Assistant

		switch aiProvider {
		case "openai", "deepseek":
			if apiKey == "" {
				fmt.Println("❌ 使用 OpenAI/DeepSeek 需要设置 API Key:")
				fmt.Println("   方式1: 设置环境变量 OPENAI_API_KEY")
				fmt.Println("   方式2: 使用 --api-key 参数")
				fmt.Println("   方式3: 运行 'lieying config init' 并编辑配置文件")
				os.Exit(1)
			}

			defaultURL := "https://api.openai.com"
			if aiProvider == "deepseek" {
				defaultURL = "https://api.deepseek.com"
			}
			if apiURL == "" {
				apiURL = defaultURL
			}

			fmt.Printf("📡 使用 %s API\n", aiProvider)
			fmt.Printf("   URL: %s\n", apiURL)
			fmt.Printf("   Model: %s\n\n", aiModel)

			client := ai.NewOpenAIClient(apiURL, apiKey, aiModel)
			assistant = ai.NewAssistantWithClient(client)

		default: // ollama
			fmt.Println("📡 使用 Ollama (本地)")
			fmt.Printf("   URL: %s\n", aiOllamaURL)
			fmt.Printf("   Model: %s\n\n", aiModel)

			client := ai.NewOllamaClient(aiOllamaURL, aiModel)
			assistant = ai.NewAssistant(client)
		}

		assistant.StartInteractiveMode()
	},
}

func init() {
	aiCmd.Flags().StringVarP(&aiModel, "model", "m", "", "模型名称 (ollama默认: qwen2.5:14b, deepseek默认: deepseek-chat)")
	aiCmd.Flags().StringVarP(&aiProvider, "provider", "p", "ollama", "LLM提供商 (ollama/openai/deepseek)")
	aiCmd.Flags().StringVarP(&aiOllamaURL, "ollama-url", "", "http://localhost:11434", "Ollama服务地址")
	aiCmd.Flags().StringVarP(&aiAPIKey, "api-key", "k", "", "OpenAI/DeepSeek API Key")
	aiCmd.Flags().StringVarP(&aiAPIURL, "api-url", "", "", "OpenAI/DeepSeek API地址")
}