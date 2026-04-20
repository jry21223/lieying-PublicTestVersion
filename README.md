# 猎影渗透测试平台 - CLI 版本

[English](#english) | [中文](#中文)

---

## 中文

### 项目简介

猎影渗透测试平台 CLI 版本是一个专业的渗透测试工具集，提供信息收集、漏洞扫描、AI辅助、教育SRC等功能。

**原项目**: [xyz-1008/lieying-PublicTestVersion](https://github.com/xyz-1008/lieying-PublicTestVersion)

**本项目在原项目基础上进行了以下改进**:

#### CLI 工具重构
- 重构为纯 CLI 工具，移除前端依赖，更轻量高效
- 使用 Cobra 框架，提供更好的命令行体验和参数解析
- 支持多种 LLM 后端 (Ollama/OpenAI/DeepSeek)

#### 扫描功能增强
- **批量扫描**: 支持从 JSON 文件批量读取目标 (`-i` 参数)，可无缝对接 recon 输出
- **存活验证**: 集成 httpx 在扫描前验证目标存活，大幅降低无效请求
- **Nuclei v3**: 兼容 Nuclei v3 参数 (`-j`, `-stats`, `-bulk-size`)，支持批量高效扫描
- **共享工具库**: 统一的 HTTP 客户端、严重等级常量、工具函数，消除代码重复

#### 子域名枚举优化
- **专业工具集成**: 优先使用 subfinder 进行被动子域名枚举（多源聚合）
- **存活验证**: 使用 httpx 验证子域名存活状态，自动过滤无效域名
- **泛解析检测**: 自动检测 DNS 泛解析，避免产生大量假阳性
- **中国域名支持**: 正确识别 `.edu.cn`, `.gov.cn` 等特殊顶级域

#### 工作流优化
- **一键流程**: `recon` → `scan -i` 实现信息收集到漏洞扫描的无缝衔接
- **进度显示**: Nuclei 扫描实时显示进度统计
- **结果分组**: 按目标分组展示扫描结果，更易阅读

### 功能特性

| 功能模块 | 说明 |
|----------|------|
| `recon` | 信息收集：子域名枚举、端口扫描、指纹识别、目录扫描 |
| `scan` | 漏洞扫描：Nuclei、SQL注入、XSS、文件上传、未授权访问 |
| `ai` | AI渗透测试助手，支持漏洞分析、POC生成 |
| `edu` | 教育SRC：高校信息查询、批量漏洞扫描 |
| `config` | 配置管理 |
| `server` | API服务器模式 |

### 安装

#### 方式一：从源码编译

```bash
# 克隆仓库
git clone https://github.com/jry21223/lieying-PublicTestVersion.git
cd lieying-PublicTestVersion/core

# 安装依赖
go mod download

# 编译
go build -o lieying ./src/

# 安装到系统
cp lieying /usr/local/bin/
```

### 快速开始

```bash
# 查看帮助
lieying --help

# 查看版本
lieying version

# 信息收集
lieying recon example.com

# 漏洞扫描
lieying scan https://example.com

# AI助手
lieying ai -p deepseek -k YOUR_API_KEY

# 教育SRC统计
lieying edu --stats
```

### 使用示例

#### 信息收集

```bash
# 全面信息收集
lieying recon example.com

# 仅子域名枚举
lieying recon example.com -s

# 仅端口扫描
lieying recon example.com -p

# 仅指纹识别
lieying recon https://example.com -f

# 仅目录扫描
lieying recon https://example.com -d

# 输出到文件
lieying recon example.com -o result.json
```

#### 漏洞扫描

```bash
# 全面漏洞扫描
lieying scan https://example.com

# 仅 Nuclei 扫描
lieying scan https://example.com -n

# 仅 SQL 注入检测
lieying scan https://example.com -s

# 仅 XSS 检测
lieying scan https://example.com -x

# 输出到文件
lieying scan https://example.com -o result.json
```

#### AI 助手

```bash
# 使用 Ollama (本地免费)
ollama serve
ollama pull qwen2.5:14b
lieying ai

# 使用 DeepSeek
export OPENAI_API_KEY=your-api-key
lieying ai -p deepseek -m deepseek-chat

# 使用 OpenAI
export OPENAI_API_KEY=your-api-key
lieying ai -p openai -m gpt-4
```

#### 教育SRC

```bash
# 统计信息
lieying edu --stats

# 列出高校
lieying edu --list

# 按级别过滤
lieying edu --list --level 985

# 按省份过滤
lieying edu --list --province 北京

# 扫描单个高校
lieying edu scan example.edu.cn

# 批量扫描
lieying edu batch --level 985
```

### 配置

```bash
# 初始化配置
lieying config init

# 查看配置
lieying config --show

# 设置配置项
lieying config --set api.key --value YOUR_API_KEY
lieying config --set network.timeout --value 60
```

配置文件位置: `~/.lieying/config.yaml`

### 依赖工具

| 工具 | 说明 | 安装 |
|------|------|------|
| Nuclei | 漏洞扫描模板引擎 | `go install github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest` |
| subfinder | 被动子域名枚举 | `go install github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest` |
| httpx | HTTP 存活验证 | `go install github.com/projectdiscovery/httpx/cmd/httpx@latest` |
| Ollama | 本地 LLM 运行 | [ollama.com](https://ollama.com) |

**推荐安装**:
```bash
# 一键安装 ProjectDiscovery 工具链
go install github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest
go install github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
go install github.com/projectdiscovery/httpx/cmd/httpx@latest

# 更新 Nuclei 模板
nuclei -ut
```

### 一键渗透测试流程

```bash
# 1. 信息收集（使用 subfinder+httpx 获取真实存活子域名）
lieying recon henu.edu.cn -o recon_result.json

# 2. 批量漏洞扫描（使用 Nuclei 对 recon 输出的目标进行扫描）
lieying scan -i recon_result.json -o scan_result.json

# 3. 查看结果
cat scan_result.json | jq '.results[] | {host: .host, severity: .info.severity, name: .info.name}'
```

> **注意**: Nuclei 现用于漏洞扫描（scan 命令），子域名枚举改用更专业的 subfinder + httpx 组合

### 项目结构

```
lieying/
├── core/                   # Go 后端
│   ├── cmd/               # CLI 命令
│   │   ├── root.go       # 根命令
│   │   ├── recon.go      # 信息收集
│   │   ├── scan.go       # 漏洞扫描
│   │   ├── ai.go         # AI助手
│   │   ├── edu.go        # 教育SRC
│   │   ├── config.go     # 配置管理
│   │   ├── server.go     # API服务
│   │   └── version.go    # 版本信息
│   ├── pkg/               # 核心包
│   │   ├── recon/        # 信息收集模块
│   │   ├── scan/         # 漏洞扫描模块
│   │   ├── ai/           # AI助手模块
│   │   └── edu/          # 教育SRC模块
│   ├── internal/          # 内部模块
│   └── src/               # 入口
├── frontend/              # React 前端
└── README.md
```

### 贡献

欢迎提交 Issue 和 Pull Request！

### 许可证

本项目基于原项目进行二次开发，遵循原项目许可证。

### 致谢

- 感谢 [昆仑安全实验室](https://github.com/xyz-1008) 提供原始项目
- 感谢所有开源安全工具的贡献者

---

## English

### Introduction

Lieying Penetration Testing Platform CLI Version is a professional penetration testing toolkit providing information gathering, vulnerability scanning, AI assistance, and EduSRC features.

**Original Project**: [xyz-1008/lieying-PublicTestVersion](https://github.com/xyz-1008/lieying-PublicTestVersion)

**Improvements made in this version**:

#### CLI Tool Refactoring
- Refactored as pure CLI tool, removed frontend dependencies for lighter footprint
- Built with Cobra framework for better CLI experience and parameter parsing
- Support for multiple LLM backends (Ollama/OpenAI/DeepSeek)

#### Scanning Enhancements
- **Batch Scanning**: Support reading targets from JSON files (`-i` parameter), seamless integration with recon output
- **Alive Verification**: Integrated httpx to verify target availability before scanning, reducing invalid requests
- **Nuclei v3**: Compatible with Nuclei v3 parameters (`-j`, `-stats`, `-bulk-size`), efficient batch scanning
- **Shared Utilities**: Unified HTTP client, severity constants, and utility functions, eliminating code duplication

#### Subdomain Enumeration Optimization
- **Professional Tool Integration**: Prioritize subfinder for passive subdomain enumeration (multi-source aggregation)
- **Alive Verification**: Use httpx to verify subdomain availability, automatically filter invalid domains
- **Wildcard Detection**: Automatic DNS wildcard detection to avoid massive false positives
- **Chinese Domain Support**: Correctly identify special TLDs like `.edu.cn`, `.gov.cn`

#### Workflow Optimization
- **One-click Workflow**: `recon` → `scan -i` for seamless information gathering to vulnerability scanning
- **Progress Display**: Real-time progress statistics during Nuclei scanning
- **Result Grouping**: Display scan results grouped by target for better readability

### Features

| Module | Description |
|--------|-------------|
| `recon` | Information gathering: subdomain enumeration, port scanning, fingerprinting, directory scanning |
| `scan` | Vulnerability scanning: Nuclei, SQL injection, XSS, file upload, unauthorized access |
| `ai` | AI penetration testing assistant with vulnerability analysis and POC generation |
| `edu` | EduSRC: university information query, batch vulnerability scanning |
| `config` | Configuration management |
| `server` | API server mode |

### Installation

```bash
# Clone repository
git clone https://github.com/jry21223/lieying-PublicTestVersion.git
cd lieying-PublicTestVersion/core

# Install dependencies
go mod download

# Build
go build -o lieying ./src/

# Install to system
cp lieying /usr/local/bin/
```

### Quick Start

```bash
# Show help
lieying --help

# Version
lieying version

# Information gathering
lieying recon example.com

# Vulnerability scanning
lieying scan https://example.com

# AI assistant
lieying ai -p deepseek -k YOUR_API_KEY

# EduSRC statistics
lieying edu --stats
```

### License

This project is based on the original project for secondary development, following the original project license.

### Acknowledgments

- Thanks to [Kunlun Security Lab](https://github.com/xyz-1008) for the original project
- Thanks to all contributors of open-source security tools

---

**Made with ❤️ by Security Researchers**