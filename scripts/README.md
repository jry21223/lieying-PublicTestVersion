# 猎影渗透测试平台 - 测试脚本

这是猎影渗透测试平台的完整测试套件，支持多种测试类型和自动化流程。

## 📋 功能特性

- **后端单元测试** (Go) - 测试后端业务逻辑
- **前端组件测试** (React + Vitest) - 测试前端组件
- **端到端测试** (Playwright) - 模拟真实用户操作
- **集成测试** - 测试后端 API 和前端交互
- **桌面端构建测试** (Tauri) - 测试桌面应用构建
- **自动依赖检测和安装**
- **统一 HTML 测试报告**
- **代码覆盖率报告**
- **ngrok 公网预览**
- **服务自动管理**

## 🚀 快速开始

### 前置要求

- macOS/Linux/WSL 环境
- Bash shell
- Git（可选）

### 1. 给脚本添加执行权限

```bash
chmod +x scripts/test.sh
chmod +x scripts/dependencies.sh
```

### 2. 检测并安装依赖

```bash
scripts/dependencies.sh
```

### 3. 运行所有测试

```bash
scripts/test.sh --all
```

或者直接运行（默认就是运行所有测试）：

```bash
scripts/test.sh
```

## 📖 使用方法

### 命令行选项

```
用法: scripts/test.sh [选项]

选项:
  -h, --help              显示帮助信息
  -a, --all               运行所有测试（默认）
  -u, --unit              仅运行后端单元测试
  -f, --frontend          仅运行前端组件测试
  -e, --e2e               仅运行端到端测试
  -i, --integration       仅运行集成测试
  -t, --tauri             仅运行 Tauri 构建测试
  -n, --ngrok             启动 ngrok 生成公网 URL
  -s, --serve             仅启动服务，不运行测试
  -c, --clean             清理测试结果
```

### 使用示例

#### 1. 运行所有测试并启动 ngrok

```bash
scripts/test.sh --all --ngrok
```

#### 2. 仅运行后端单元测试

```bash
scripts/test.sh --unit
```

#### 3. 仅运行前端组件测试

```bash
scripts/test.sh --frontend
```

#### 4. 运行端到端测试

```bash
scripts/test.sh --e2e
```

#### 5. 运行集成测试

```bash
scripts/test.sh --integration
```

#### 6. 仅启动服务，不运行测试

```bash
scripts/test.sh --serve
```

带 ngrok：

```bash
scripts/test.sh --serve --ngrok
```

#### 7. 清理测试结果

```bash
scripts/test.sh --clean
```

## 📊 测试报告

测试完成后，报告会生成在 `test-results/reports/` 目录下：

- `test-report.html` - 统一的 HTML 测试报告
- `backend-coverage.html` - 后端代码覆盖率报告
- `playwright-report/` - Playwright E2E 测试报告

### 打开测试报告

```bash
# 在浏览器中打开主报告
open test-results/reports/test-report.html  # macOS
xdg-open test-results/reports/test-report.html  # Linux
```

## 🔧 配置

测试配置文件位于 `scripts/config.sh`，你可以修改以下配置：

```bash
# 服务端口
BACKEND_PORT=8080
FRONTEND_PORT=3000

# 测试配置
RUN_UNIT_TESTS=true
RUN_FRONTEND_TESTS=true
RUN_E2E_TESTS=true
RUN_INTEGRATION_TESTS=true
RUN_TAURI_TESTS=true

# 覆盖率阈值
BACKEND_COVERAGE_THRESHOLD=70
FRONTEND_COVERAGE_THRESHOLD=60
```

## 📁 目录结构

```
scripts/
├── test.sh              # 主测试脚本
├── dependencies.sh      # 依赖检测和安装脚本
├── config.sh            # 测试配置文件
└── README.md            # 本文档

test-results/
├── reports/             # 测试报告
│   ├── test-report.html
│   ├── backend-coverage.html
│   └── playwright-report/
├── coverage/            # 覆盖率数据
└── logs/                # 测试日志
```

## 🐛 故障排除

### 问题：权限被拒绝

```bash
chmod +x scripts/*.sh
```

### 问题：端口被占用

修改 `scripts/config.sh` 中的端口配置，或者停止占用端口的进程：

```bash
# 查找占用 8080 端口的进程
lsof -ti:8080 | xargs kill -9  # macOS/Linux
netstat -ano | findstr :8080  # Windows
```

### 问题：依赖安装失败

- **macOS**: 确保已安装 Homebrew
- **Linux**: 确保有 sudo 权限
- **Windows**: 建议使用 WSL，或手动安装依赖

### 问题：ngrok 无法连接

确保已安装 ngrok 并登录：

```bash
ngrok config add-authtoken YOUR_AUTH_TOKEN
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
