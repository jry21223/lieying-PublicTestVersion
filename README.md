# 猎影渗透测试平台 - CLI 版本

[English](#english) | [中文](#中文)

---

## 中文

### 项目简介

猎影渗透测试平台 CLI 版本是一个专业的渗透测试工具集，提供信息收集、漏洞扫描、AI辅助、教育SRC等功能。

**原项目**: [xyz-1008/lieying-PublicTestVersion](https://github.com/xyz-1008/lieying-PublicTestVersion)

**本项目在原项目基础上进行了以下改进**:
- 重构为 CLI 工具，支持命令行操作
- 使用 Cobra 框架，提供更好的命令行体验
- 支持多种 LLM 后端 (Ollama/OpenAI/DeepSeek)
- 添加 Nuclei 漏洞扫描集成

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
| Nuclei | 漏洞扫描模板 | `go install github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest` |
| Ollama | 本地 LLM 运行 | [ollama.com](https://ollama.com) |

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
- Refactored as a CLI tool with command-line interface
- Built with Cobra framework for better CLI experience
- Support for multiple LLM backends (Ollama/OpenAI/DeepSeek)
- Integrated Nuclei vulnerability scanning

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