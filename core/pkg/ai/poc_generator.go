package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

type POCGenerator struct {
	client *OllamaClient
}

type GeneratedPOC struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Severity    string   `json:"severity"`
	Description string   `json:"description"`
	POCContent  string   `json:"poc_content"`
	Usage       string   `json:"usage"`
	Payloads    []string `json:"payloads,omitempty"`
}

func NewPOCGenerator(client *OllamaClient) *POCGenerator {
	return &POCGenerator{
		client: client,
	}
}

func (pg *POCGenerator) GenerateFromDescription(vulnType, description, url string) (*GeneratedPOC, error) {
	if !pg.client.IsAvailable() {
		return nil, fmt.Errorf("Ollama服务不可用")
	}

	fmt.Println("🤖 AI正在生成POC...")

	prompt := fmt.Sprintf(`你是一个专业的渗透测试工程师，请根据以下漏洞描述生成一个POC（概念验证）脚本：

漏洞类型：%s
漏洞描述：%s
目标URL：%s

请生成一个完整的POC，包括：
1. POC名称
2. 漏洞类型和严重级别
3. 详细描述
4. POC代码（Python或Go语言）
5. 使用方法
6. 测试用的Payload列表

请返回JSON格式：
{
  "name": "POC名称",
  "type": "漏洞类型",
  "severity": "严重级别",
  "description": "详细描述",
  "poc_content": "POC代码",
  "usage": "使用方法",
  "payloads": ["payload1", "payload2"]
}

只返回JSON，不要其他内容。`, vulnType, description, url)

	response, err := pg.client.Generate(prompt)
	if err != nil {
		return nil, err
	}

	response = pg.extractJSON(response)

	var poc GeneratedPOC
	if err := json.Unmarshal([]byte(response), &poc); err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %v", err)
	}

	fmt.Println("✅ POC生成完成")
	return &poc, nil
}

func (pg *POCGenerator) GenerateWAFBypassPayload(originalPayload, wafType string) ([]string, error) {
	if !pg.client.IsAvailable() {
		return nil, fmt.Errorf("Ollama服务不可用")
	}

	fmt.Println("🤖 AI正在生成WAF绕过Payload...")

	prompt := fmt.Sprintf(`你是一个专业的渗透测试工程师，请为以下Payload生成WAF绕过变种：

原始Payload：%s
WAF类型：%s

请生成5-10个绕过Payload，使用以下技术：
1. URL编码
2. Unicode编码
3. 大小写混淆
4. 注释插入
5. 空格替换
6. 其他绕过技术

请返回JSON格式：
{
  "payloads": ["payload1", "payload2", "payload3", "payload4", "payload5"]
}

只返回JSON，不要其他内容。`, originalPayload, wafType)

	response, err := pg.client.Generate(prompt)
	if err != nil {
		return nil, err
	}

	response = pg.extractJSON(response)

	var result struct {
		Payloads []string `json:"payloads"`
	}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %v", err)
	}

	fmt.Printf("✅ 生成 %d 个WAF绕过Payload\n", len(result.Payloads))
	return result.Payloads, nil
}

func (pg *POCGenerator) OptimizeExistingPOC(pocContent, optimizationGoal string) (string, error) {
	if !pg.client.IsAvailable() {
		return "", fmt.Errorf("Ollama服务不可用")
	}

	fmt.Println("🤖 AI正在优化POC...")

	prompt := fmt.Sprintf(`你是一个专业的渗透测试工程师，请优化以下POC代码：

优化目标：%s

原始POC代码：
%s

请优化后的POC代码，要求：
1. 提高成功率
2. 减少误报
3. 提高执行效率
4. 增强稳定性

请直接返回优化后的代码，不要其他内容。`, optimizationGoal, pocContent)

	response, err := pg.client.Generate(prompt)
	if err != nil {
		return "", err
	}

	fmt.Println("✅ POC优化完成")
	return response, nil
}

func (pg *POCGenerator) GenerateNucleiTemplate(vulnName, vulnType, severity, description string) (string, error) {
	if !pg.client.IsAvailable() {
		return "", fmt.Errorf("Ollama服务不可用")
	}

	fmt.Println("🤖 AI正在生成Nuclei模板...")

	prompt := fmt.Sprintf(`你是一个专业的渗透测试工程师，请生成一个Nuclei POC模板：

漏洞名称：%s
漏洞类型：%s
严重级别：%s
漏洞描述：%s

请生成一个标准的Nuclei YAML模板，包括：
1. id
2. info（名称、作者、严重级别、描述、参考）
3. http请求
4. 匹配器

请直接返回YAML格式的模板代码。`, vulnName, vulnType, severity, description)

	response, err := pg.client.Generate(prompt)
	if err != nil {
		return "", err
	}

	fmt.Println("✅ Nuclei模板生成完成")
	return response, nil
}

func (pg *POCGenerator) extractJSON(response string) string {
	start := strings.Index(response, "[")
	end := strings.LastIndex(response, "]")

	if start != -1 && end != -1 && end > start {
		return response[start : end+1]
	}

	start = strings.Index(response, "{")
	end = strings.LastIndex(response, "}")

	if start != -1 && end != -1 && end > start {
		return response[start : end+1]
	}

	return response
}
