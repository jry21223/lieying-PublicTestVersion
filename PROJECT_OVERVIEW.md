# 猎影 (Lieying) 项目概述

## 🎯 项目愿景

「猎影」是一个开源的渗透测试平台，目标是让白帽黑客能够一站式完成从信息收集到报告提交的全流程，效率提升10倍。项目采用 MIT 许可证，核心功能永久免费，支持本地优先运行，注重隐私和易用性。

---

## 🛠️ 技术栈

### 后端 (Backend)
- **语言**: Go 1.21+
- **Web框架**: 标准库 `net/http` (轻量级、高性能)
- **ORM**: GORM (SQLite/PostgreSQL 支持)
- **数据库**: SQLite (默认本地存储)，可选 PostgreSQL
- **并发**: Goroutines + Channels (任务调度)
- **WebSocket**: Gorilla WebSocket (实时推送)
- **工具集成**: `exec.Command` 调用外部二进制 (Nuclei、Amass、Subfinder 等)

### 前端 (Frontend)
- **框架**: React 18 + TypeScript
- **构建工具**: Vite
- **样式**: Tailwind CSS
- **图标**: Lucide React
- **状态管理**: Zustand (可选 React Context)
- **HTTP客户端**: Axios (REST API)
- **WebSocket**: 原生 WebSocket API

### 桌面应用
- **框架**: Tauri (跨平台打包，轻量级)
- **通信**: HTTP REST API (推荐，便于调试) 或 Tauri Commands

### AI 引擎
- **本地**: Ollama (无需联网，保护隐私)
- **云端**: DeepSeek、Claude (可选，需API Key)
- **功能**: 漏洞优先级排序、POC生成、误报过滤、报告生成

### 插件系统
- **协议**: gRPC (高性能、多语言支持)
- **支持语言**: Go、Python、JavaScript
- **功能**: 插件注册、调用、热加载

---

## 🏗️ 核心功能模块

### 1. 信息收集模块 (Reconnaissance)
- ✅ **子域名枚举**: 多源收集 (字典、API、DNS)
- ✅ **端口扫描**: 全端口快速扫描 (TCP/UDP)
- ✅ **指纹识别**: Web框架、CMS、中间件识别
- ✅ **目录爆破**: 敏感路径和文件发现
- 🔄 **集成工具**: Amass、Subfinder、Nmap、Wappalyzer

### 2. 漏洞扫描模块 (Vulnerability Scanning)
- ✅ **SQL注入检测**: 基于错误、时间、布尔的注入检测
- ✅ **XSS漏洞检测**: 反射型、存储型、DOM型XSS
- ✅ **文件上传漏洞**: 上传点检测、绕过测试
- ✅ **未授权访问**: 接口和页面权限检测
- ✅ **弱口令检测**: 常见口令爆破
- 🔄 **Nuclei集成**: Nuclei模板扫描
- 🔄 **Xray/Afrog**: 支持第三方扫描器
- 🔄 **自定义POC**: 用户可添加自定义POC

### 3. 逻辑测试模块 (Logic Testing) ⏳
- **越权检测**: 水平越权、垂直越权
- **支付绕过**: 支付逻辑漏洞检测
- **验证码复用**: 验证码安全性测试
- **会话劫持**: 会话管理安全测试

### 4. 教育SRC专属模块 (Education SRC) ✅
- **高校数据库**: 100+ 所985/211高校域名
- **教务系统指纹**: 正方、强智、金智、URP、青果
- **专项POC库**: 针对各教务系统的POC
- **批量扫描**: 并发扫描多所高校
- **漏洞统计**: 教育SRC漏洞趋势分析

### 5. AI辅助模块 (AI Assistant)
- **漏洞优先级排序**: 输入漏洞列表，返回排序结果
- **POC自动生成**: 根据漏洞描述生成测试Payload
- **误报过滤**: 判断扫描结果是否为真实漏洞
- **WAF绕过辅助**: 提供绕过建议和Payload
- **报告自动生成**: 根据漏洞信息生成规范报告

### 6. 报告生成模块 (Reporting)
- **一键导出**: Word、Markdown、PDF格式
- **SRC格式**: 符合主流SRC平台的报告格式
- **漏洞分类**: CVSS评分、严重级别
- **自定义模板**: 支持用户自定义报告模板

### 7. POC市场 (POC Marketplace) ⏳
- **用户上传**: 用户可上传自己的POC
- **自动验证**: 系统自动验证POC有效性
- **积分系统**: 上传高质量POC获得积分
- **下载POC**: 使用积分下载他人POC

### 8. 插件系统 (Plugin System) ⏳
- **多语言支持**: Go、Python、JavaScript插件
- **gRPC协议**: 高性能插件通信
- **插件管理器**: 加载、卸载、调用插件
- **插件市场**: 用户分享和下载插件

### 9. 团队协作 (Team Collaboration) ⏳
- **共享工作区**: 团队成员共享目标和漏洞
- **实时同步**: WebSocket实时同步数据
- **评论讨论**: 漏洞评论和讨论
- **权限管理**: 不同角色不同权限

---

## 🏗️ 架构设计

```
┌─────────────────────────────────────────────────────────────────┐
│                        猎影平台架构                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   CLI工具    │  │   Web界面    │  │  桌面应用    │      │
│  │   (Go)       │  │  (React)     │  │  (Tauri)     │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                   │                   │               │
│         └───────────────────┼───────────────────┘               │
│                             │                                   │
│  ┌──────────────────────────┴──────────────────────────┐       │
│  │              API Gateway (net/http)                   │       │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │       │
│  │  │ REST API   │  │ WebSocket   │  │    CORS     │ │       │
│  │  └─────────────┘  └─────────────┘  └─────────────┘ │       │
│  └──────────────────────────┬──────────────────────────┘       │
│                             │                                   │
│  ┌──────────────────────────┴──────────────────────────┐       │
│  │                   Core Engine                        │       │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ │       │
│  │  │  Recon  │ │  Scan   │ │   AI    │ │ Report  │ │       │
│  │  │  Module │ │ Module  │ │ Module  │ │ Module  │ │       │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘ │       │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐             │       │
│  │  │  Task   │ │  Plugin │ │  Edu    │             │       │
│  │  │ Scheduler│ │ Manager │ │  SRC    │             │       │
│  │  └─────────┘ └─────────┘ └─────────┘             │       │
│  └──────────────────────────┬──────────────────────────┘       │
│                             │                                   │
│  ┌──────────────────────────┴──────────────────────────┐       │
│  │                   Data Layer                          │       │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │       │
│  │  │   SQLite    │  │  File Sys   │  │   Config    │ │       │
│  │  │  (Local)    │  │   (POC)     │  │   (YAML)    │ │       │
│  │  └─────────────┘  └─────────────┘  └─────────────┘ │       │
│  └─────────────────────────────────────────────────────┘       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📊 数据模型

### Target (目标资产)
```go
type Target struct {
    ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
    Name        string         `gorm:"size:255;not null" json:"name"`
    Type        string         `gorm:"size:50;not null" json:"type"` // domain/ip/url
    Value       string         `gorm:"size:500;not null;index" json:"value"`
    Description string         `gorm:"type:text" json:"description"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
}
```

### Vulnerability (漏洞)
```go
type Vulnerability struct {
    ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
    Title       string         `gorm:"size:500;not null" json:"title"`
    Description string         `gorm:"type:text" json:"description"`
    Severity    string         `gorm:"size:20;not null" json:"severity"` // critical/high/medium/low/info
    URL         string         `gorm:"size:1000" json:"url"`
    Confirmed   bool           `gorm:"default:false" json:"confirmed"`
    CreatedAt   time.Time      `json:"created_at"`
}
```

### ScanTask (扫描任务)
```go
type ScanTask struct {
    ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
    TargetID    uuid.UUID      `gorm:"type:uuid;index" json:"target_id"`
    Name        string         `gorm:"size:255;not null" json:"name"`
    Type        string         `gorm:"size:50;not null" json:"type"` // recon/scan
    Status      string         `gorm:"size:50;default:pending" json:"status"`
    Progress    int            `gorm:"default:0" json:"progress"`
    Result      string         `gorm:"type:text" json:"result"`
    Error       string         `gorm:"type:text" json:"error"`
    StartedAt   *time.Time     `json:"started_at"`
    FinishedAt  *time.Time     `json:"finished_at"`
}
```

---

## 🚀 API 接口设计

### 目标管理
- `GET /api/targets` - 列出所有目标
- `POST /api/targets` - 创建目标
- `DELETE /api/targets/:id` - 删除目标

### 漏洞管理
- `GET /api/vulnerabilities` - 列出所有漏洞
- `POST /api/vulnerabilities/:id/confirm` - 确认漏洞

### 信息收集
- `POST /api/recon` - 执行信息收集
- `GET /api/recon/:id/status` - 获取收集状态

### 漏洞扫描
- `POST /api/scan` - 执行漏洞扫描
- `GET /api/scan/:id/status` - 获取扫描状态

### 教育SRC
- `GET /api/edu/universities` - 获取高校列表
- `POST /api/edu/scan` - 扫描教务系统
- `POST /api/edu/batch` - 批量扫描

### AI功能
- `GET /api/ai/status` - 检查AI服务状态
- `POST /api/ai/poc` - AI生成POC
- `POST /api/ai/analyze` - AI分析漏洞

### WebSocket
- `WS /api/ws` - 实时推送扫描进度、日志、AI分析结果

---

## 📦 目录结构

```
lieying/
├── core/                    # 后端核心
│   ├── internal/
│   │   ├── api/            # API 路由和处理器
│   │   ├── app/            # 应用初始化
│   │   ├── models/         # 数据模型
│   │   ├── repository/     # 数据访问层
│   │   ├── service/        # 业务逻辑层
│   │   ├── engine/         # 核心引擎
│   │   └── scanner/        # 扫描器
│   ├── pkg/
│   │   ├── recon/          # 信息收集包
│   │   ├── scan/           # 漏洞扫描包
│   │   ├── ai/             # AI 集成包
│   │   ├── edu/            # 教育SRC包
│   │   ├── report/         # 报告生成包
│   │   └── logger/         # 日志包
│   └── src/
│       └── main.go         # 主入口
├── frontend/               # 前端
│   ├── src/
│   │   ├── components/     # React 组件
│   │   ├── services/       # API 服务
│   │   ├── stores/         # 状态管理
│   │   └── App.tsx         # 主应用
│   └── package.json
├── engine/                 # 独立引擎 (可选)
├── platform/               # 平台相关
├── portable/               # 便携版
├── docs/                   # 文档
├── README.md               # 项目说明
└── start.bat              # 启动脚本
```

---

## 🎯 开发优先级

### Phase 1: 核心功能完善 (当前)
- ✅ 修复假数据问题
- ✅ 目标资产详情显示
- ✅ 漏洞数据持久化
- 🔄 完善漏洞扫描模块
- 🔄 完善信息收集模块

### Phase 2: AI 集成
- Ollama 本地AI集成
- 漏洞优先级排序
- POC 自动生成
- 误报过滤

### Phase 3: 任务调度与WebSocket
- goroutine + channel 任务调度
- WebSocket 实时推送
- 任务取消和超时控制

### Phase 4: Tauri 桌面应用
- Tauri 配置和集成
- 文件系统访问
- 桌面应用打包

### Phase 5: 插件系统
- gRPC 协议定义
- 插件管理器
- 插件开发示例

### Phase 6: 高级功能
- POC 市场
- 团队协作
- 逻辑测试模块

---

## 📝 开发规范

### 后端开发规范
- 使用 Go 标准库 `net/http` 构建 REST API
- 使用 GORM 操作数据库
- 所有函数添加错误处理
- 使用 `context` 进行超时和取消控制
- 代码注释清晰，说明功能和参数

### 前端开发规范
- 使用 React + TypeScript 组件化开发
- 使用 Tailwind CSS 进行样式
- 状态管理使用 Zustand 或 React Context
- API 调用封装在 services 层
- 组件命名使用 PascalCase

### Git 提交规范
- `feat: 新功能`
- `fix: 修复bug`
- `docs: 文档更新`
- `style: 代码格式`
- `refactor: 重构`
- `test: 测试`
- `chore: 构建/工具`

---

## 🤝 贡献指南

我们欢迎所有形式的贡献！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

---

## 🙏 致谢

感谢以下开源项目：
- [Nuclei](https://github.com/projectdiscovery/nuclei)
- [Ollama](https://github.com/ollama/ollama)
- [Tauri](https://github.com/tauri-apps/tauri)
- [GORM](https://github.com/go-gorm/gorm)
- [React](https://github.com/facebook/react)
- [Tailwind CSS](https://github.com/tailwindlabs/tailwindcss)

---

**Made with ❤️ by 昆仑安全实验室(前逍遥安全实验室-逍遥)**

**颠覆行业，让安全测试触手可及！**
