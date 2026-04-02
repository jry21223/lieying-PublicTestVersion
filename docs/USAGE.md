# 猎影平台使用指南

**昆仑安全实验室(前逍遥安全实验室-逍遥)** 出品

---

## 📑 目录

1. [快速开始](#快速开始)
2. [CLI命令行使用](#cli命令行使用)
3. [Web界面使用](#web界面使用)
4. [功能详解](#功能详解)
5. [配置说明](#配置说明)
6. [常见问题](#常见问题)

---

## 快速开始

### 1. 安装

#### 从源码安装

```bash
# 克隆仓库
git clone https://github.com/kunlun-sec/lieying.git
cd lieying

# 安装后端依赖
cd core
go mod download

# 安装前端依赖
cd ../frontend
npm install
```

#### 从Release下载

访问 [Releases页面](https://github.com/kunlun-sec/lieying/releases) 下载对应平台的二进制文件。

### 2. 启动服务

```bash
# 启动后端API服务
cd core
go run src/main.go server :3000

# 启动前端界面 (新终端)
cd frontend
npm run dev
```

### 3. 访问界面

打开浏览器访问: `http://localhost:5173`

---

## CLI命令行使用

### 全局命令

```bash
lieying [command] [flags]
```

### 可用命令

#### 1. 信息收集 (recon)

```bash
# 全面信息收集
lieying recon example.com

# 只收集子域名
lieying recon example.com --subdomain-only

# 只扫描端口
lieying recon example.com --port-only
```

**参数说明:**
- `--subdomain-only`: 仅子域名枚举
- `--port-only`: 仅端口扫描
- `--fingerprint-only`: 仅指纹识别
- `--dir-only`: 仅目录扫描
- `-o, --output`: 输出文件路径

#### 2. 漏洞扫描 (scan)

```bash
# 全面漏洞扫描
lieying scan example.com

# 只扫描SQL注入
lieying scan example.com --sqli-only

# 使用特定POC
lieying scan example.com --poc /path/to/poc.yaml
```

**参数说明:**
- `--sqli-only`: 仅SQL注入检测
- `--xss-only`: 仅XSS检测
- `--upload-only`: 仅文件上传检测
- `--unauth-only`: 仅未授权访问检测
- `--poc`: 指定POC文件路径
- `-o, --output`: 输出文件路径

#### 3. AI助手 (ai)

```bash
# 启动交互式AI助手
lieying ai

# 直接提问
lieying ai "如何绕过WAF检测SQL注入？"

# 生成POC
lieying ai --generate-poc --vuln-type sqli --url "http://example.com/login"
```

#### 4. 报告生成 (report)

```bash
# 生成HTML报告
lieying report example.com --format html

# 生成PDF报告
lieying report example.com --format pdf

# 生成Markdown报告
lieying report example.com --format md
```

#### 5. 教育SRC (edu)

```bash
# 查看高校统计
lieying edu stats

# 扫描单个教务系统
lieying edu scan jwgl.example.edu.cn

# 批量扫描985高校
lieying edu batch 985

# 批量扫描211高校
lieying edu batch 211

# 按省份扫描
lieying edu batch --province 北京
```

#### 6. 服务管理 (server)

```bash
# 启动API服务器
lieying server :3000

# 后台运行
lieying server :3000 --daemon
```

---

## Web界面使用

### 1. 资产概览

主界面显示：
- 发现目标数量
- 高危漏洞统计
- 已确认漏洞数
- 最新发现的漏洞列表

### 2. 信息收集

**功能模块:**
- 子域名枚举
- 端口扫描
- 指纹识别
- 目录扫描

**操作步骤:**
1. 点击"信息收集"菜单
2. 输入目标域名
3. 选择扫描模块
4. 点击"开始收集"
5. 查看扫描结果

### 3. 漏洞扫描

**支持的漏洞类型:**
- SQL注入
- XSS跨站脚本
- 文件上传
- 未授权访问

**操作步骤:**
1. 点击"漏洞扫描"菜单
2. 输入目标URL
3. 点击"开始扫描"
4. 查看漏洞列表
5. 使用AI验证漏洞

### 4. AI助手

**功能:**
- 智能问答
- POC生成
- 报告撰写
- WAF绕过建议

**使用方式:**
1. 点击"AI助手"菜单
2. 在输入框中提问
3. 获取AI回复

### 5. 教育SRC

**功能模块:**
- 统计概览
- 高校列表
- 批量扫描
- 扫描结果

**操作步骤:**
1. 点击"教育SRC"菜单
2. 选择"批量扫描"
3. 配置扫描参数
4. 开始扫描
5. 导出扫描结果

---

## 功能详解

### 信息收集模块

#### 子域名枚举

**原理:**
- 字典爆破
- 证书透明度日志
- 搜索引擎收集

**输出字段:**
- 子域名
- IP地址
- 状态码
- 标题

#### 端口扫描

**扫描方式:**
- SYN半开扫描
- Connect全连接扫描
- 服务指纹识别

**常见端口:**
- 21 (FTP)
- 22 (SSH)
- 80 (HTTP)
- 443 (HTTPS)
- 3306 (MySQL)
- 3389 (RDP)

#### 指纹识别

**识别内容:**
- Web服务器 (Nginx/Apache/IIS)
- 开发框架 (Spring/Django/ThinkPHP)
- CMS系统 (WordPress/Drupal)
- JavaScript库 (jQuery/Vue/React)

### 漏洞扫描模块

#### SQL注入检测

**检测方法:**
- 错误回显检测
- 时间盲注检测
- 布尔盲注检测
- 联合查询检测

**支持的注入点:**
- GET参数
- POST参数
- HTTP头
- Cookie

#### XSS检测

**检测类型:**
- 反射型XSS
- 存储型XSS
- DOM型XSS

**Payload示例:**
```
<script>alert(1)</script>
<img src=x onerror=alert(1)>
```

#### 文件上传漏洞

**检测点:**
- 后缀名绕过
- Content-Type绕过
- 内容检测绕过
- 条件竞争

### AI助手模块

#### POC生成

**支持的漏洞类型:**
- SQL注入
- 命令执行
- 文件读取
- 未授权访问
- 反序列化

**生成示例:**
```yaml
id: example-sqli
info:
  name: Example SQL Injection
  severity: critical
requests:
  - method: GET
    path:
      - "{{BaseURL}}/api/user?id=1' AND 1=1--"
    matchers:
      - type: word
        words:
          - "user_id"
```

#### 报告生成

**报告内容:**
- 漏洞概述
- 详细描述
- 复现步骤
- 修复建议
- CVSS评分

---

## 配置说明

### 配置文件位置

```
~/.lieying/config.yaml
```

### 配置示例

```yaml
version: "1.0.0"

database:
  path: ~/.lieying/lieying.db

network:
  mode: offline  # offline/proxy/direct
  proxy:
    enabled: false
    type: socks5
    host: 127.0.0.1
    port: 1080
  allowed_hosts: []

ai:
  mode: local  # local/remote
  local:
    provider: ollama
    endpoint: http://localhost:11434
    model: qwen2.5:7b
  remote:
    provider: openai
    api_key: ""
    base_url: https://api.deepseek.com
    model: deepseek-chat

scan:
  concurrency: 10
  timeout: 30
  nuclei_template_path: ~/.lieying/nuclei-templates
  max_depth: 3

recon:
  subdomain_wordlist: ~/.lieying/wordlists/subdomains.txt
  dir_wordlist: ~/.lieying/wordlists/dirs.txt
  port_range: "1-65535"
  timeout: 5

edu:
  default_concurrency: 10
  batch_size: 50

log:
  level: info  # debug/info/warn/error/fatal
  path: ~/.lieying/logs
  format: json  # json/text
```

### 环境变量

```bash
# 配置文件路径
export LIEYING_CONFIG=/path/to/config.yaml

# 数据库路径
export LIEYING_DB=/path/to/database.db

# 日志级别
export LIEYING_LOG_LEVEL=debug

# AI端点
export LIEYING_AI_ENDPOINT=http://localhost:11434
```

---

## 常见问题

### Q1: 如何更新POC库?

```bash
# 自动更新
lieying update --poc

# 手动下载
wget -O ~/.lieying/nuclei-templates.zip https://github.com/projectdiscovery/nuclei-templates/archive/refs/heads/main.zip
unzip ~/.lieying/nuclei-templates.zip -d ~/.lieying/
```

### Q2: 如何添加自定义POC?

```bash
# 创建POC目录
mkdir -p ~/.lieying/custom-pocs

# 添加POC文件
cp your-poc.yaml ~/.lieying/custom-pocs/

# 扫描时使用
cd ~/.lieying/custom-pocs
lieying scan example.com --poc ./your-poc.yaml
```

### Q3: 如何配置代理?

编辑配置文件 `~/.lieying/config.yaml`:

```yaml
network:
  mode: proxy
  proxy:
    enabled: true
    type: socks5  # socks5/http/https
    host: 127.0.0.1
    port: 1080
```

### Q4: 如何离线使用?

```yaml
network:
  mode: offline
```

离线模式下，所有功能均可正常使用，但不会连接外部网络。

### Q5: 如何导出扫描结果?

```bash
# JSON格式
lieying scan example.com -o result.json

# CSV格式
lieying scan example.com -o result.csv

# HTML报告
lieying report example.com --format html -o report.html
```

### Q6: 如何贡献POC?

1. Fork 项目仓库
2. 在 `core/pkg/scan/pocs/` 目录下添加POC文件
3. 提交 Pull Request
4. 等待审核合并

### Q7: 扫描速度太慢怎么办?

调整并发数:

```yaml
scan:
  concurrency: 50  # 增加并发数
  timeout: 10      # 减少超时时间
```

或使用命令行参数:

```bash
lieying scan example.com --concurrency 50 --timeout 10
```

### Q8: 如何查看日志?

```bash
# 实时查看日志
tail -f ~/.lieying/logs/lieying_$(date +%Y-%m-%d).log

# 查看特定级别日志
grep "ERROR" ~/.lieying/logs/lieying_$(date +%Y-%m-%d).log
```

---

## 📞 技术支持

- 📧 邮箱: contact@kunlun-sec.com
- 💬 Discord: https://discord.gg/lieying
- 🐛 Issue: https://github.com/kunlun-sec/lieying/issues

---

**Made with ❤️ by 昆仑安全实验室(前逍遥安全实验室-逍遥)**
