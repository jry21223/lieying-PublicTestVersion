package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

type VulnerabilityAnalyzer struct {
	client *OllamaClient
}

type VulnPriority struct {
	Name        string  `json:"name"`
	URL         string  `json:"url"`
	Type        string  `json:"type"`
	Severity    string  `json:"severity"`
	Priority    int     `json:"priority"`
	Reason      string  `json:"reason"`
	Confidence  float64 `json:"confidence"`
	RewardScore float64 `json:"reward_score"`
}

type FalsePositiveResult struct {
	IsFalsePositive bool    `json:"is_false_positive"`
	Confidence      float64 `json:"confidence"`
	Reason          string  `json:"reason"`
}

func NewVulnerabilityAnalyzer(client *OllamaClient) *VulnerabilityAnalyzer {
	return &VulnerabilityAnalyzer{
		client: client,
	}
}

func (va *VulnerabilityAnalyzer) AnalyzeAndPrioritize(vulns []map[string]interface{}) ([]VulnPriority, error) {
	if !va.client.IsAvailable() {
		return nil, fmt.Errorf("Ollama服务不可用")
	}

	fmt.Println("🤖 AI正在分析漏洞优先级...")

	vulnData, err := json.Marshal(vulns)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(`你是一个专业的渗透测试工程师，请分析以下漏洞列表，并按以下标准排序（1-10分，10分最高）：

1. 漏洞严重性（Critical > High > Medium > Low）
2. 利用难度（越容易利用分数越高）
3. 潜在奖金（SRC平台历史数据，SQL注入、RCE等通常奖金更高）
4. 影响范围（影响用户越多分数越高）

请返回JSON格式：
[
  {
    "name": "漏洞名称",
    "url": "漏洞URL",
    "type": "漏洞类型",
    "severity": "严重级别",
    "priority": 优先级数字1-10,
    "reason": "优先级原因",
    "confidence": 置信度0-1,
    "reward_score": 预期奖金分数1-10
  }
]

漏洞列表：
%s

只返回JSON数组，不要其他内容。`, string(vulnData))

	response, err := va.client.Generate(prompt)
	if err != nil {
		return nil, err
	}

	response = va.extractJSON(response)

	var priorities []VulnPriority
	if err := json.Unmarshal([]byte(response), &priorities); err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %v", err)
	}

	fmt.Printf("✅ AI分析完成，已排序 %d 个漏洞\n", len(priorities))
	return priorities, nil
}

func (va *VulnerabilityAnalyzer) CheckFalsePositive(vuln map[string]interface{}) (*FalsePositiveResult, error) {
	if !va.client.IsAvailable() {
		return nil, fmt.Errorf("Ollama服务不可用")
	}

	vulnData, err := json.Marshal(vuln)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(`你是一个专业的渗透测试工程师，请判断以下漏洞是否可能是误报：

判断标准：
1. 响应内容是否确实证明漏洞存在
2. 是否有其他解释（如自定义错误页面）
3. 触发条件是否过于苛刻
4. 是否需要特殊环境才能复现

请返回JSON格式：
{
  "is_false_positive": true/false,
  "confidence": 置信度0-1,
  "reason": "判断原因"
}

漏洞信息：
%s

只返回JSON，不要其他内容。`, string(vulnData))

	response, err := va.client.Generate(prompt)
	if err != nil {
		return nil, err
	}

	response = va.extractJSON(response)

	var result FalsePositiveResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %v", err)
	}

	return &result, nil
}

func (va *VulnerabilityAnalyzer) GenerateAttackPath(target string, vulns []map[string]interface{}) (string, error) {
	if !va.client.IsAvailable() {
		return "", fmt.Errorf("Ollama服务不可用")
	}

	fmt.Println("🤖 AI正在生成攻击路径...")

	vulnData, err := json.Marshal(vulns)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`你是一个专业的渗透测试工程师，请为以下目标生成最优攻击路径：

目标：%s
已发现漏洞：
%s

请生成一个详细的攻击路径，包括：
1. 第一步应该攻击哪个漏洞
2. 如何利用该漏洞获取更多信息
3. 后续攻击步骤
4. 最终目标（如获取shell、读取敏感数据等）

请以中文返回详细的攻击路径建议。`, target, string(vulnData))

	response, err := va.client.Generate(prompt)
	if err != nil {
		return "", err
	}

	fmt.Println("✅ 攻击路径生成完成")
	return response, nil
}

func (va *VulnerabilityAnalyzer) RecommendNextStep(scanResults map[string]interface{}) (string, error) {
	if !va.client.IsAvailable() {
		return "", fmt.Errorf("Ollama服务不可用")
	}

	resultsData, err := json.Marshal(scanResults)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`你是一个专业的渗透测试工程师，基于以下扫描结果，请推荐下一步行动：

扫描结果：
%s

请推荐下一步应该做什么，选项包括：
1. 继续深入扫描特定漏洞
2. 尝试利用已发现的漏洞
3. 扩大扫描范围
4. 生成报告
5. 其他建议

请给出具体建议和理由。`, string(resultsData))

	return va.client.Generate(prompt)
}

func (va *VulnerabilityAnalyzer) extractJSON(response string) string {
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
