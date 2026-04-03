package ai

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// AIClient 统一的AI客户端接口
type AIClient interface {
	IsAvailable() bool
	Chat(messages []Message) (string, error)
	Generate(prompt string) (string, error)
	ListModels() ([]OllamaModel, error)
}

type Assistant struct {
	client    AIClient
	analyzer  *VulnerabilityAnalyzer
	pocGen    *POCGenerator
	context   []Message
	isRunning bool
}

func NewAssistant(client *OllamaClient) *Assistant {
	return &Assistant{
		client:    client,
		analyzer:  NewVulnerabilityAnalyzer(client),
		pocGen:    NewPOCGenerator(client),
		context:   []Message{},
		isRunning: false,
	}
}

// NewAssistantWithClient 使用任意客户端创建助手
func NewAssistantWithClient(client AIClient) *Assistant {
	return &Assistant{
		client:    client,
		analyzer:  nil, // 分析器和POC生成器需要特定客户端
		pocGen:    nil,
		context:   []Message{},
		isRunning: false,
	}
}

func (a *Assistant) StartInteractiveMode() {
	a.isRunning = true

	fmt.Println("=====================================")
	fmt.Println("  🤖 你好！我是小影，你的AI渗透测试助手")
	fmt.Println("  昆仑安全实验室(前逍遥安全实验室-逍遥)")
	fmt.Println("=====================================")
	fmt.Println()

	if !a.client.IsAvailable() {
		fmt.Println("⚠️  警告：AI服务未正确配置")
		fmt.Println()
		fmt.Println("使用 Ollama:")
		fmt.Println("  1. 安装Ollama: https://ollama.com")
		fmt.Println("  2. 启动Ollama服务: ollama serve")
		fmt.Println("  3. 拉取模型: ollama pull qwen2.5:14b")
		fmt.Println()
		fmt.Println("使用 DeepSeek/OpenAI:")
		fmt.Println("  1. 设置环境变量: export OPENAI_API_KEY=your_key")
		fmt.Println("  2. 可选设置: export OPENAI_BASE_URL=https://api.deepseek.com")
		fmt.Println()
		return
	} else {
		fmt.Println("✅ AI服务连接正常")
		fmt.Println()
	}

	fmt.Println("我可以帮你：")
	fmt.Println("  1. 分析漏洞优先级")
	fmt.Println("  2. 生成POC代码")
	fmt.Println("  3. 生成WAF绕过Payload")
	fmt.Println("  4. 检查误报")
	fmt.Println("  5. 规划攻击路径")
	fmt.Println("  6. 回答渗透测试相关问题")
	fmt.Println()
	fmt.Println("输入 'help' 查看帮助，输入 'exit' 退出")
	fmt.Println("=====================================")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for a.isRunning {
		fmt.Print("🤖 小影 > ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		a.handleInput(input)
	}
}

func (a *Assistant) handleInput(input string) {
	inputLower := strings.ToLower(input)

	switch inputLower {
	case "exit", "quit", "bye":
		a.stop()
	case "help", "?":
		a.showHelp()
	case "status":
		a.showStatus()
	case "clear":
		a.clearContext()
		fmt.Println("🧹 对话上下文已清除")
	default:
		a.chat(input)
	}
}

func (a *Assistant) chat(input string) {
	if !a.client.IsAvailable() {
		fmt.Println("❌ AI服务不可用")
		return
	}

	a.context = append(a.context, Message{
		Role:    "user",
		Content: input,
	})

	systemPrompt := `你是小影，一个专业的渗透测试AI助手。你由昆仑安全实验室(前逍遥安全实验室-逍遥)开发。

你的专长包括：
1. 漏洞分析和评估
2. POC生成和优化
3. 渗透测试技术指导
4. 安全工具使用建议
5. 攻击路径规划

请用中文回答，保持专业、友好、简洁的风格。如果用户问的是技术问题，请提供详细的步骤和示例。`

	messages := append([]Message{{
		Role:    "system",
		Content: systemPrompt,
	}}, a.context...)

	fmt.Println("🤔 小影正在思考...")
	response, err := a.client.Chat(messages)
	if err != nil {
		fmt.Printf("❌ AI响应错误: %v\n", err)
		return
	}

	a.context = append(a.context, Message{
		Role:    "assistant",
		Content: response,
	})

	fmt.Println()
	fmt.Println(response)
	fmt.Println()
}

func (a *Assistant) AnalyzeVulnerabilities(vulns []map[string]interface{}) {
	if !a.client.IsAvailable() {
		fmt.Println("❌ AI服务不可用")
		return
	}

	if a.analyzer == nil {
		fmt.Println("❌ 此功能需要 Ollama 客户端")
		return
	}

	priorities, err := a.analyzer.AnalyzeAndPrioritize(vulns)
	if err != nil {
		fmt.Printf("❌ 分析失败: %v\n", err)
		return
	}

	fmt.Println("\n📊 AI漏洞优先级分析结果：")
	fmt.Println("=====================================")
	for i, p := range priorities {
		fmt.Printf("%d. %s\n", i+1, p.Name)
		fmt.Printf("   URL: %s\n", p.URL)
		fmt.Printf("   类型: %s | 严重级别: %s\n", p.Type, p.Severity)
		fmt.Printf("   优先级: %d/10 | 奖金潜力: %.1f/10\n", p.Priority, p.RewardScore)
		fmt.Printf("   原因: %s\n", p.Reason)
		fmt.Printf("   置信度: %.0f%%\n", p.Confidence*100)
		fmt.Println()
	}
	fmt.Println("=====================================")
}

func (a *Assistant) GeneratePOC(vulnType, description, url string) {
	if !a.client.IsAvailable() {
		fmt.Println("❌ AI服务不可用")
		return
	}

	if a.pocGen == nil {
		fmt.Println("❌ 此功能需要 Ollama 客户端")
		return
	}

	poc, err := a.pocGen.GenerateFromDescription(vulnType, description, url)
	if err != nil {
		fmt.Printf("❌ POC生成失败: %v\n", err)
		return
	}

	fmt.Println("\n📋 生成的POC：")
	fmt.Println("=====================================")
	fmt.Printf("名称: %s\n", poc.Name)
	fmt.Printf("类型: %s | 严重级别: %s\n", poc.Type, poc.Severity)
	fmt.Printf("描述: %s\n", poc.Description)
	fmt.Println()
	fmt.Println("代码:")
	fmt.Println(poc.POCContent)
	fmt.Println()
	fmt.Printf("使用方法: %s\n", poc.Usage)
	if len(poc.Payloads) > 0 {
		fmt.Println("测试Payloads:")
		for _, p := range poc.Payloads {
			fmt.Printf("  - %s\n", p)
		}
	}
	fmt.Println("=====================================")
}

func (a *Assistant) GenerateWAFBypass(originalPayload, wafType string) {
	if !a.client.IsAvailable() {
		fmt.Println("❌ AI服务不可用")
		return
	}

	if a.pocGen == nil {
		fmt.Println("❌ 此功能需要 Ollama 客户端")
		return
	}

	payloads, err := a.pocGen.GenerateWAFBypassPayload(originalPayload, wafType)
	if err != nil {
		fmt.Printf("❌ 生成失败: %v\n", err)
		return
	}

	fmt.Println("\n🛡️ WAF绕过Payloads：")
	fmt.Println("=====================================")
	for i, p := range payloads {
		fmt.Printf("%d. %s\n", i+1, p)
	}
	fmt.Println("=====================================")
}

func (a *Assistant) CheckFalsePositive(vuln map[string]interface{}) {
	if !a.client.IsAvailable() {
		fmt.Println("❌ AI服务不可用")
		return
	}

	if a.analyzer == nil {
		fmt.Println("❌ 此功能需要 Ollama 客户端")
		return
	}

	result, err := a.analyzer.CheckFalsePositive(vuln)
	if err != nil {
		fmt.Printf("❌ 检查失败: %v\n", err)
		return
	}

	fmt.Println("\n🔍 误报检查结果：")
	fmt.Println("=====================================")
	if result.IsFalsePositive {
		fmt.Printf("⚠️  可能是误报 (置信度: %.0f%%)\n", result.Confidence*100)
	} else {
		fmt.Printf("✅ 可能是真实漏洞 (置信度: %.0f%%)\n", result.Confidence*100)
	}
	fmt.Printf("原因: %s\n", result.Reason)
	fmt.Println("=====================================")
}

func (a *Assistant) GenerateAttackPath(target string, vulns []map[string]interface{}) {
	if !a.client.IsAvailable() {
		fmt.Println("❌ AI服务不可用")
		return
	}

	if a.analyzer == nil {
		fmt.Println("❌ 此功能需要 Ollama 客户端")
		return
	}

	path, err := a.analyzer.GenerateAttackPath(target, vulns)
	if err != nil {
		fmt.Printf("❌ 生成失败: %v\n", err)
		return
	}

	fmt.Println("\n🎯 攻击路径规划：")
	fmt.Println("=====================================")
	fmt.Println(path)
	fmt.Println("=====================================")
}

func (a *Assistant) showHelp() {
	fmt.Println("\n📖 小影帮助：")
	fmt.Println("=====================================")
	fmt.Println("基本命令：")
	fmt.Println("  help      - 显示帮助")
	fmt.Println("  status    - 检查服务状态")
	fmt.Println("  clear     - 清除对话上下文")
	fmt.Println("  exit      - 退出")
	fmt.Println()
	fmt.Println("你可以问我：")
	fmt.Println("  - 如何测试SQL注入？")
	fmt.Println("  - 这个漏洞怎么利用？")
	fmt.Println("  - 生成一个XSS POC")
	fmt.Println("  - 如何绕过WAF？")
	fmt.Println("  - 等等...")
	fmt.Println("=====================================")
}

func (a *Assistant) showStatus() {
	fmt.Println("\n📊 系统状态：")
	fmt.Println("=====================================")
	if a.client.IsAvailable() {
		fmt.Println("✅ AI服务: 已连接")
		models, err := a.client.ListModels()
		if err == nil && len(models) > 0 {
			fmt.Printf("📦 可用模型: %d 个\n", len(models))
			for _, m := range models {
				fmt.Printf("   - %s\n", m.Name)
			}
		}
	} else {
		fmt.Println("❌ AI服务: 未配置")
	}
	fmt.Printf("💬 对话历史: %d 条消息\n", len(a.context))
	fmt.Println("=====================================")
}

func (a *Assistant) clearContext() {
	a.context = []Message{}
}

func (a *Assistant) stop() {
	a.isRunning = false
	fmt.Println("\n👋 再见！祝你挖洞顺利！")
	fmt.Println("   昆仑安全实验室(前逍遥安全实验室-逍遥)")
}

func (a *Assistant) IsRunning() bool {
	return a.isRunning
}