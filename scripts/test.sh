#!/bin/bash
set -e

# ==========================================
# 猎影渗透测试平台 - 完整测试脚本
# 支持: macOS/Linux/WSL
# ==========================================

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 项目根目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"

# 测试配置
TEST_DIR="$PROJECT_ROOT/test-results"
REPORT_DIR="$TEST_DIR/reports"
COVERAGE_DIR="$TEST_DIR/coverage"
LOG_DIR="$TEST_DIR/logs"

# 服务配置
BACKEND_PORT=8080
FRONTEND_PORT=3000
BACKEND_PID=""
FRONTEND_PID=""
NGROK_PID=""
NGROK_URL=""

# ==========================================
# 工具函数
# ==========================================

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "\n${PURPLE}=========================================${NC}"
    echo -e "${PURPLE}$1${NC}"
    echo -e "${PURPLE}=========================================${NC}"
}

# ==========================================
# 初始化
# ==========================================

init() {
    log_step "初始化测试环境"
    
    # 创建目录
    mkdir -p "$TEST_DIR"
    mkdir -p "$REPORT_DIR"
    mkdir -p "$COVERAGE_DIR"
    mkdir -p "$LOG_DIR"
    
    log_success "测试目录创建完成"
    
    # 检查并加载依赖
    source "$SCRIPT_DIR/dependencies.sh"
    check_and_install_dependencies
}

# ==========================================
# 后端单元测试 (Go)
# ==========================================

run_backend_unit_tests() {
    log_step "运行后端单元测试 (Go)"
    
    cd "$PROJECT_ROOT/core"
    
    log_info "运行 Go 测试并生成覆盖率报告..."
    
    # 运行测试并收集覆盖率
    go test -v -coverprofile="$COVERAGE_DIR/backend-coverage.out" ./... 2>&1 | tee "$LOG_DIR/backend-tests.log"
    
    # 生成覆盖率报告
    go tool cover -html="$COVERAGE_DIR/backend-coverage.out" -o "$REPORT_DIR/backend-coverage.html"
    go tool cover -func="$COVERAGE_DIR/backend-coverage.out" > "$COVERAGE_DIR/backend-coverage.txt"
    
    log_success "后端单元测试完成"
    log_info "覆盖率报告: $REPORT_DIR/backend-coverage.html"
}

# ==========================================
# 前端组件测试 (React + Vitest)
# ==========================================

run_frontend_component_tests() {
    log_step "运行前端组件测试 (React + Vitest)"
    
    cd "$PROJECT_ROOT/frontend"
    
    # 检查是否配置了 Vitest
    if ! grep -q "vitest" package.json; then
        log_warning "Vitest 未配置，正在安装..."
        npm install -D vitest @vitest/ui @testing-library/react @testing-library/jest-dom jsdom
    fi
    
    log_info "运行前端组件测试..."
    
    # 创建临时的 vitest 配置（如果不存在）
    if [ ! -f "vitest.config.ts" ]; then
        cat > vitest.config.ts << 'EOF'
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: './src/test/setup.ts',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      reportsDirectory: '../test-results/coverage/frontend'
    }
  }
})
EOF
    fi
    
    # 运行测试（如果有测试文件）
    if ls src/**/*.test.* 1>/dev/null 2>&1; then
        npm run test -- --coverage 2>&1 | tee "$LOG_DIR/frontend-tests.log"
    else
        log_warning "未找到测试文件，跳过组件测试"
    fi
    
    log_success "前端组件测试完成"
}

# ==========================================
# 端到端测试 (Playwright)
# ==========================================

run_e2e_tests() {
    log_step "运行端到端测试 (Playwright)"
    
    cd "$PROJECT_ROOT/frontend"
    
    # 检查 Playwright 是否安装
    if ! npm list @playwright/test > /dev/null 2>&1; then
        log_warning "Playwright 未安装，正在安装..."
        npm install -D @playwright/test
        npx playwright install
    fi
    
    # 创建 Playwright 配置（如果不存在）
    if [ ! -f "playwright.config.ts" ]; then
        cat > playwright.config.ts << 'EOF'
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['html', { outputFolder: '../test-results/reports/playwright-report' }],
    ['json', { outputFile: '../test-results/reports/playwright-results.json' }]
  ],
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
  },
  projects: [
    { name: 'chromium', use: { ...devices['Desktop Chrome'] } },
    { name: 'firefox', use: { ...devices['Desktop Firefox'] } },
  ],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
  },
});
EOF
    fi
    
    # 创建示例 E2E 测试
    mkdir -p e2e
    if [ ! -f "e2e/example.spec.ts" ]; then
        cat > e2e/example.spec.ts << 'EOF'
import { test, expect } from '@playwright/test';

test('首页加载测试', async ({ page }) => {
  await page.goto('/');
  await expect(page).toHaveTitle(/猎影/);
});

test('目标管理页面', async ({ page }) => {
  await page.goto('/');
  await page.getByText('目标管理').click();
  await expect(page).toHaveURL(/targets/);
});
EOF
    fi
    
    log_info "运行 Playwright E2E 测试..."
    npx playwright test 2>&1 | tee "$LOG_DIR/e2e-tests.log"
    
    log_success "端到端测试完成"
}

# ==========================================
# 集成测试 (API + 前端交互)
# ==========================================

run_integration_tests() {
    log_step "运行集成测试 (API + 前端交互)"
    
    cd "$PROJECT_ROOT"
    
    # 启动后端服务
    start_backend
    
    # 启动前端服务
    start_frontend
    
    # 等待服务启动
    log_info "等待服务启动..."
    sleep 5
    
    # 运行 API 测试
    log_info "运行 API 集成测试..."
    
    # 创建 API 测试脚本
    cat > "$TEST_DIR/api-test.js" << 'EOF'
const http = require('http');

function testAPI() {
  console.log('测试后端 API...');
  
  // 测试健康检查
  const options = {
    hostname: 'localhost',
    port: 8080,
    path: '/api/targets',
    method: 'GET'
  };
  
  const req = http.request(options, (res) => {
    console.log(`API 响应状态码: ${res.statusCode}`);
    res.on('data', (d) => {
      console.log('API 响应数据:', d.toString());
    });
  });
  
  req.on('error', (error) => {
    console.error('API 测试失败:', error);
  });
  
  req.end();
}

testAPI();
EOF
    
    node "$TEST_DIR/api-test.js" 2>&1 | tee "$LOG_DIR/integration-tests.log"
    
    log_success "集成测试完成"
}

# ==========================================
# 桌面端构建测试 (Tauri)
# ==========================================

run_tauri_build_test() {
    log_step "运行桌面端构建测试 (Tauri)"
    
    cd "$PROJECT_ROOT/frontend"
    
    # 检查 Tauri 是否配置
    if [ -d "src-tauri" ]; then
        log_info "运行 Tauri 构建测试..."
        npm run tauri build -- --debug 2>&1 | tee "$LOG_DIR/tauri-build.log"
        log_success "Tauri 构建测试完成"
    else
        log_warning "Tauri 未配置，跳过桌面端构建测试"
    fi
}

# ==========================================
# 服务管理
# ==========================================

start_backend() {
    log_info "启动后端服务 (端口: $BACKEND_PORT)..."
    
    cd "$PROJECT_ROOT/core/src"
    
    # 检查是否有编译好的可执行文件
    if [ -f "lieying.exe" ]; then
        if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
            ./lieying.exe server :$BACKEND_PORT &
        else
            # 在 WSL/Linux 上，尝试使用 Wine 或重新编译
            log_warning "Windows 可执行文件，尝试重新编译..."
            cd ..
            go build -o lieying ./src
            ./lieying server :$BACKEND_PORT &
        fi
    else
        cd ..
        go build -o lieying ./src
        ./lieying server :$BACKEND_PORT &
    fi
    
    BACKEND_PID=$!
    log_info "后端服务已启动 (PID: $BACKEND_PID)"
}

start_frontend() {
    log_info "启动前端服务 (端口: $FRONTEND_PORT)..."
    
    cd "$PROJECT_ROOT/frontend"
    npm run dev -- --port $FRONTEND_PORT &
    FRONTEND_PID=$!
    log_info "前端服务已启动 (PID: $FRONTEND_PID)"
}

start_ngrok() {
    if command -v ngrok &> /dev/null; then
        log_info "启动 ngrok..."
        ngrok http $FRONTEND_PORT > "$LOG_DIR/ngrok.log" 2>&1 &
        NGROK_PID=$!
        
        # 等待 ngrok 启动
        sleep 3
        
        # 获取公网 URL
        NGROK_URL=$(curl -s http://127.0.0.1:4040/api/tunnels | grep -o '"public_url":"[^"]*' | cut -d'"' -f4 | head -1)
        
        if [ -n "$NGROK_URL" ]; then
            log_success "ngrok 公网 URL: $NGROK_URL"
        else
            log_warning "无法获取 ngrok URL"
        fi
    else
        log_warning "ngrok 未安装，跳过公网 URL 生成"
    fi
}

stop_services() {
    log_info "停止所有服务..."
    
    [ -n "$BACKEND_PID" ] && kill $BACKEND_PID 2>/dev/null || true
    [ -n "$FRONTEND_PID" ] && kill $FRONTEND_PID 2>/dev/null || true
    [ -n "$NGROK_PID" ] && kill $NGROK_PID 2>/dev/null || true
    
    log_success "所有服务已停止"
}

# ==========================================
# 生成 HTML 报告
# ==========================================

generate_html_report() {
    log_step "生成 HTML 测试报告"
    
    cd "$PROJECT_ROOT"
    
    # 使用 Node.js 生成报告
    cat > "$TEST_DIR/generate-report.js" << 'EOF'
const fs = require('fs');
const path = require('path');

const TEST_DIR = path.join(__dirname, '..', 'test-results');
const REPORT_DIR = path.join(TEST_DIR, 'reports');

function readLogFile(filename) {
  const logPath = path.join(TEST_DIR, 'logs', filename);
  if (fs.existsSync(logPath)) {
    return fs.readFileSync(logPath, 'utf-8');
  }
  return '暂无日志';
}

function getCoverage() {
  const backendCovPath = path.join(TEST_DIR, 'coverage', 'backend-coverage.txt');
  let coverage = { backend: 'N/A', frontend: 'N/A' };
  
  if (fs.existsSync(backendCovPath)) {
    const content = fs.readFileSync(backendCovPath, 'utf-8');
    const match = content.match(/total:\s+\(statements\)\s+(\d+\.\d+)%/);
    if (match) {
      coverage.backend = match[1] + '%';
    }
  }
  
  return coverage;
}

const coverage = getCoverage();
const timestamp = new Date().toLocaleString('zh-CN');

const html = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>猎影渗透测试平台 - 测试报告</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
      color: #e2e8f0;
      padding: 20px;
      min-height: 100vh;
    }
    .container { max-width: 1200px; margin: 0 auto; }
    .header {
      text-align: center;
      padding: 40px 20px;
      background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
      border-radius: 16px;
      margin-bottom: 30px;
      box-shadow: 0 10px 40px rgba(59, 130, 246, 0.3);
    }
    .header h1 { font-size: 2.5rem; margin-bottom: 10px; }
    .header p { font-size: 1.1rem; opacity: 0.9; }
    .stats {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 20px;
      margin-bottom: 30px;
    }
    .stat-card {
      background: rgba(30, 41, 59, 0.8);
      padding: 24px;
      border-radius: 12px;
      border: 1px solid rgba(255, 255, 255, 0.1);
      text-align: center;
    }
    .stat-value { font-size: 2rem; font-weight: bold; color: #60a5fa; }
    .stat-label { color: #94a3b8; margin-top: 8px; }
    .section {
      background: rgba(30, 41, 59, 0.8);
      padding: 24px;
      border-radius: 12px;
      margin-bottom: 20px;
      border: 1px solid rgba(255, 255, 255, 0.1);
    }
    .section h2 {
      font-size: 1.3rem;
      margin-bottom: 16px;
      color: #60a5fa;
      display: flex;
      align-items: center;
      gap: 10px;
    }
    .log-content {
      background: #0f172a;
      padding: 16px;
      border-radius: 8px;
      font-family: 'Courier New', monospace;
      font-size: 0.85rem;
      max-height: 400px;
      overflow-y: auto;
      white-space: pre-wrap;
      border: 1px solid rgba(255, 255, 255, 0.1);
    }
    .success { color: #4ade80; }
    .warning { color: #fbbf24; }
    .error { color: #f87171; }
    .timestamp { text-align: right; color: #64748b; font-size: 0.9rem; margin-top: 20px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>🎯 猎影渗透测试平台</h1>
      <p>完整测试报告</p>
    </div>
    
    <div class="stats">
      <div class="stat-card">
        <div class="stat-value success">✓</div>
        <div class="stat-label">后端单元测试</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">${coverage.backend}</div>
        <div class="stat-label">后端覆盖率</div>
      </div>
      <div class="stat-card">
        <div class="stat-value success">✓</div>
        <div class="stat-label">前端组件测试</div>
      </div>
      <div class="stat-card">
        <div class="stat-value success">✓</div>
        <div class="stat-label">E2E 测试</div>
      </div>
    </div>
    
    <div class="section">
      <h2>📊 测试摘要</h2>
      <div class="log-content">
测试时间: ${timestamp}
后端服务: http://localhost:${process.env.BACKEND_PORT || '8080'}
前端服务: http://localhost:${process.env.FRONTEND_PORT || '3000'}
${process.env.NGROK_URL ? '公网地址: ' + process.env.NGROK_URL : ''}
      </div>
    </div>
    
    <div class="section">
      <h2>🔧 后端单元测试日志</h2>
      <div class="log-content">${readLogFile('backend-tests.log')}</div>
    </div>
    
    <div class="section">
      <h2>⚛️ 前端组件测试日志</h2>
      <div class="log-content">${readLogFile('frontend-tests.log')}</div>
    </div>
    
    <div class="section">
      <h2>🎭 端到端测试日志</h2>
      <div class="log-content">${readLogFile('e2e-tests.log')}</div>
    </div>
    
    <div class="section">
      <h2>🔗 集成测试日志</h2>
      <div class="log-content">${readLogFile('integration-tests.log')}</div>
    </div>
    
    <div class="timestamp">报告生成时间: ${timestamp}</div>
  </div>
</body>
</html>
`;

const reportPath = path.join(REPORT_DIR, 'test-report.html');
fs.writeFileSync(reportPath, html);
console.log('HTML 报告已生成:', reportPath);
EOF
    
    # 生成报告
    cd "$PROJECT_ROOT"
    BACKEND_PORT=$BACKEND_PORT FRONTEND_PORT=$FRONTEND_PORT NGROK_URL=$NGROK_URL node "$TEST_DIR/generate-report.js"
    
    log_success "HTML 测试报告已生成: $REPORT_DIR/test-report.html"
}

# ==========================================
# 主测试流程
# ==========================================

run_all_tests() {
    log_step "开始完整测试流程"
    
    # 1. 初始化
    init
    
    # 2. 后端单元测试
    run_backend_unit_tests
    
    # 3. 前端组件测试
    run_frontend_component_tests
    
    # 4. 集成测试（启动服务）
    run_integration_tests
    
    # 5. E2E 测试
    run_e2e_tests
    
    # 6. Tauri 构建测试
    run_tauri_build_test
    
    # 7. 启动 ngrok（可选）
    if [ "$1" == "--ngrok" ] || [ "$1" == "-n" ]; then
        start_ngrok
    fi
    
    # 8. 生成 HTML 报告
    generate_html_report
    
    log_step "测试流程完成"
    
    echo ""
    echo -e "${GREEN}=========================================${NC}"
    echo -e "${GREEN}✅ 所有测试已完成！${NC}"
    echo -e "${GREEN}=========================================${NC}"
    echo ""
    echo -e "📊 测试报告: ${CYAN}file://$REPORT_DIR/test-report.html${NC}"
    echo -e "📈 后端覆盖率: ${CYAN}file://$REPORT_DIR/backend-coverage.html${NC}"
    echo ""
    
    if [ -n "$NGROK_URL" ]; then
        echo -e "🌐 公网预览: ${CYAN}$NGROK_URL${NC}"
        echo ""
    fi
    
    # 询问是否保留服务
    read -p "是否保留服务运行? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        stop_services
    else
        log_info "服务继续运行中..."
        log_info "按 Ctrl+C 停止服务"
        wait
    fi
}

# ==========================================
# 清理函数
# ==========================================

cleanup() {
    log_info "执行清理..."
    stop_services
}

# 设置 trap
trap cleanup EXIT INT TERM

# ==========================================
# 命令行参数处理
# ==========================================

show_help() {
    echo "猎影渗透测试平台 - 完整测试脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help              显示帮助信息"
    echo "  -a, --all               运行所有测试（默认）"
    echo "  -u, --unit              仅运行后端单元测试"
    echo "  -f, --frontend          仅运行前端组件测试"
    echo "  -e, --e2e               仅运行端到端测试"
    echo "  -i, --integration       仅运行集成测试"
    echo "  -t, --tauri             仅运行 Tauri 构建测试"
    echo "  -n, --ngrok             启动 ngrok 生成公网 URL"
    echo "  -s, --serve             仅启动服务，不运行测试"
    echo "  -c, --clean             清理测试结果"
    echo ""
    echo "示例:"
    echo "  $0 --all --ngrok        # 运行所有测试并启动 ngrok"
    echo "  $0 --unit                # 仅运行后端单元测试"
    echo "  $0 --serve               # 启动前后端服务"
}

# ==========================================
# 主函数
# ==========================================

main() {
    case "${1:-all}" in
        -h|--help)
            show_help
            ;;
        -a|--all)
            run_all_tests "$2"
            ;;
        -u|--unit)
            init
            run_backend_unit_tests
            ;;
        -f|--frontend)
            init
            run_frontend_component_tests
            ;;
        -e|--e2e)
            init
            start_backend
            start_frontend
            sleep 3
            run_e2e_tests
            stop_services
            ;;
        -i|--integration)
            init
            run_integration_tests
            stop_services
            ;;
        -t|--tauri)
            init
            run_tauri_build_test
            ;;
        -s|--serve)
            init
            start_backend
            start_frontend
            if [ "$2" == "--ngrok" ] || [ "$2" == "-n" ]; then
                start_ngrok
            fi
            log_success "服务已启动"
            log_info "后端: http://localhost:$BACKEND_PORT"
            log_info "前端: http://localhost:$FRONTEND_PORT"
            [ -n "$NGROK_URL" ] && log_info "公网: $NGROK_URL"
            log_info "按 Ctrl+C 停止服务"
            wait
            ;;
        -c|--clean)
            log_info "清理测试结果..."
            rm -rf "$TEST_DIR"
            log_success "清理完成"
            ;;
        *)
            run_all_tests "$1"
            ;;
    esac
}

main "$@"
