# 更新日志

所有 notable 的变更都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

---

## [1.0.0] - 2024-12-27

### 🎉 初始发布

猎影 (Lieying) v1.0.0 正式发布！

### ✨ 新增功能

#### 信息收集 (Reconnaissance)
- 🔍 子域名枚举 - 多源子域名收集
- 🌐 端口扫描 - 全端口快速扫描
- 🎯 指纹识别 - Web框架和组件识别
- 📁 目录扫描 - 敏感路径和文件发现

#### 漏洞扫描 (Vulnerability Scanning)
- 💉 SQL注入检测 - 基于错误和时间的注入检测
- 🎭 XSS漏洞检测 - 反射型和存储型XSS
- 📤 文件上传漏洞 - 上传点检测和绕过
- 🔓 未授权访问 - 接口和页面权限检测
- 🔥 Nuclei模板 - 集成Nuclei漏洞模板

#### AI助手 (AI Assistant)
- 🤖 智能分析 - 本地Ollama集成，无需联网
- 📝 POC生成 - 自动生成漏洞验证代码
- 📊 报告生成 - 智能SRC报告撰写
- 🛡️ WAF绕过 - 提供绕过建议和Payload

#### 教育SRC专项 (Education SRC)
- 🏫 高校数据库 - 100+所985/211高校域名
- 📚 教务系统POC - 正方、强智、金智、URP、青果
- ⚡ 批量扫描 - 并发扫描多所高校
- 📈 漏洞统计 - 教育SRC漏洞趋势分析

#### 报告与提交 (Reporting)
- 📄 自动报告 - 生成专业渗透测试报告
- 🎯 SRC提交 - 自动化SRC平台提交
- 📋 漏洞分类 - CVSS评分和严重级别
- 📤 多格式导出 - PDF、Word、Markdown

### 🛠️ 技术特性

#### 后端 (Go)
- 本地优先架构 - SQLite零配置
- RESTful API - 标准库 net/http
- 并发处理 - Goroutines + Channels
- 日志系统 - 分级日志支持
- 配置管理 - YAML配置文件

#### 前端 (React + TypeScript)
- 现代化UI - Tailwind CSS
- 响应式设计 - 支持多设备
- 桌面应用 - Tauri集成
- 图标系统 - Lucide React

### 📦 安装方式

- 源码安装
- 二进制下载
- Docker部署

### 📚 文档

- README.md - 项目介绍
- docs/USAGE.md - 使用指南
- CONTRIBUTING.md - 贡献指南
- LICENSE - MIT许可证

### 🙏 致谢

感谢所有为猎影项目做出贡献的开发者！

---

## 版本说明

### 版本号格式

版本号格式：主版本号.次版本号.修订号

- **主版本号**：重大更新，可能不兼容旧版本
- **次版本号**：新增功能，向下兼容
- **修订号**：问题修复，向下兼容

### 版本标签

- 🎉 `Added` - 新功能
- 🔧 `Changed` - 变更
- 🗑️ `Deprecated` - 废弃
- 🐛 `Fixed` - 修复
- 🛡️ `Security` - 安全

---

**Made with ❤️ by 昆仑安全实验室(前逍遥安全实验室-逍遥)**
