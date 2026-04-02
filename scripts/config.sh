#!/bin/bash
# ==========================================
# 猎影渗透测试平台 - 测试配置文件
# ==========================================

# 服务端口配置
export BACKEND_PORT=8080
export FRONTEND_PORT=3000

# 测试目录配置
export TEST_DIR="$PROJECT_ROOT/test-results"
export REPORT_DIR="$TEST_DIR/reports"
export COVERAGE_DIR="$TEST_DIR/coverage"
export LOG_DIR="$TEST_DIR/logs"

# 测试配置
export RUN_UNIT_TESTS=true
export RUN_FRONTEND_TESTS=true
export RUN_E2E_TESTS=true
export RUN_INTEGRATION_TESTS=true
export RUN_TAURI_TESTS=true

# 覆盖率阈值（百分比）
export BACKEND_COVERAGE_THRESHOLD=70
export FRONTEND_COVERAGE_THRESHOLD=60

# 超时配置（秒）
export TEST_TIMEOUT=300
export SERVICE_START_TIMEOUT=30

# 数据库配置（测试用）
export TEST_DB_PATH=":memory:"

# API 配置
export API_BASE_URL="http://localhost:$BACKEND_PORT"
export FRONTEND_URL="http://localhost:$FRONTEND_PORT"

# 浏览器配置（Playwright）
export PLAYWRIGHT_BROWSERS="chromium,firefox"
export PLAYWRIGHT_HEADLESS=true

# ngrok 配置
export NGROK_REGION="cn"  # us, eu, ap, au, sa, jp, in
