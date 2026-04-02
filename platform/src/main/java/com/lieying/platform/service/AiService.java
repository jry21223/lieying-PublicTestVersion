package com.lieying.platform.service;

import lombok.RequiredArgsConstructor;
import org.springframework.ai.chat.ChatClient;
import org.springframework.ai.chat.prompt.Prompt;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class AiService {
    private final ChatClient chatClient;

    public String generatePoc(String vulnDescription) {
        String prompt = String.format("""
            你是一个专业的安全研究人员。请根据以下漏洞描述生成一个POC（概念验证）代码。
            漏洞描述：%s
            
            要求：
            1. 提供完整可运行的代码
            2. 代码要有详细注释
            3. 仅用于授权的安全测试
            
            返回格式：
            - 语言：[Python/Go]
            - 代码：
            ```[语言]
            代码内容
            ```
            """, vulnDescription);

        return chatClient.call(new Prompt(prompt)).getResult().getOutput().getContent();
    }

    public String generateReport(String vulnDetails) {
        String prompt = String.format("""
            请根据以下漏洞详情生成一份专业的SRC漏洞报告。
            漏洞详情：%s
            
            报告应包含：
            1. 漏洞标题
            2. 漏洞描述
            3. 复现步骤
            4. 漏洞危害
            5. 修复建议
            """, vulnDetails);

        return chatClient.call(new Prompt(prompt)).getResult().getOutput().getContent();
    }

    public String analyzeAttackPath(String reconData) {
        String prompt = String.format("""
            你是一个红队专家。请根据以下信息收集结果，分析可能的攻击路径。
            信息收集结果：%s
            
            请提供：
            1. 高价值目标识别
            2. 推荐的攻击路径
            3. 优先级建议
            """, reconData);

        return chatClient.call(new Prompt(prompt)).getResult().getOutput().getContent();
    }
}
