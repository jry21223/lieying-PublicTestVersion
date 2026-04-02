#!/bin/bash
# ==========================================
# 猎影渗透测试平台 - 依赖检测和安装脚本
# ==========================================

# 颜色定义（重复定义以确保独立运行）
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

# 检测操作系统
detect_os() {
    case "$OSTYPE" in
        darwin*)
            echo "macOS"
            ;;
        linux-gnu*)
            if grep -qi microsoft /proc/version; then
                echo "WSL"
            else
                echo "Linux"
            fi
            ;;
        msys*|cygwin*|win32*)
            echo "Windows"
            ;;
        *)
            echo "Unknown"
            ;;
    esac
}

# 检测 Go
check_go() {
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}')
        log_success "Go 已安装: $GO_VERSION"
        return 0
    else
        log_warning "Go 未安装"
        return 1
    fi
}

# 检测 Node.js
check_node() {
    if command -v node &> /dev/null; then
        NODE_VERSION=$(node -v)
        log_success "Node.js 已安装: $NODE_VERSION"
        return 0
    else
        log_warning "Node.js 未安装"
        return 1
    fi
}

# 检测 npm
check_npm() {
    if command -v npm &> /dev/null; then
        NPM_VERSION=$(npm -v)
        log_success "npm 已安装: v$NPM_VERSION"
        return 0
    else
        log_warning "npm 未安装"
        return 1
    fi
}

# 检测 ngrok
check_ngrok() {
    if command -v ngrok &> /dev/null; then
        NGROK_VERSION=$(ngrok version 2>&1 | head -1)
        log_success "ngrok 已安装: $NGROK_VERSION"
        return 0
    else
        log_warning "ngrok 未安装（可选）"
        return 1
    fi
}

# 检测 Playwright 浏览器
check_playwright() {
    if command -v npx &> /dev/null; then
        if npx playwright --version &> /dev/null; then
            PW_VERSION=$(npx playwright --version)
            log_success "Playwright 已安装: $PW_VERSION"
            return 0
        fi
    fi
    log_warning "Playwright 未安装（可选）"
    return 1
}

# 安装 Go（macOS）
install_go_macos() {
    log_info "正在使用 Homebrew 安装 Go..."
    if command -v brew &> /dev/null; then
        brew install go
    else
        log_error "Homebrew 未安装，请先安装 Homebrew: https://brew.sh/"
        return 1
    fi
}

# 安装 Go（Linux）
install_go_linux() {
    log_info "正在安装 Go..."
    local GO_VERSION="1.21.5"
    local GO_TAR="go$GO_VERSION.linux-amd64.tar.gz"
    local GO_URL="https://go.dev/dl/$GO_TAR"
    
    wget "$GO_URL" -O "/tmp/$GO_TAR"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "/tmp/$GO_TAR"
    
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    echo 'export GOPATH=$HOME/go' >> ~/.bashrc
    echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
    
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    
    rm "/tmp/$GO_TAR"
}

# 安装 Node.js（macOS）
install_node_macos() {
    log_info "正在使用 Homebrew 安装 Node.js..."
    if command -v brew &> /dev/null; then
        brew install node
    else
        log_error "Homebrew 未安装，请先安装 Homebrew: https://brew.sh/"
        return 1
    fi
}

# 安装 Node.js（Linux）
install_node_linux() {
    log_info "正在使用 NVM 安装 Node.js..."
    if [ ! -d ~/.nvm ]; then
        curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
    fi
    export NVM_DIR="$HOME/.nvm"
    [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
    nvm install --lts
    nvm use --lts
}

# 安装 ngrok
install_ngrok() {
    log_info "正在安装 ngrok..."
    OS=$(detect_os)
    
    case "$OS" in
        macOS)
            if command -v brew &> /dev/null; then
                brew install ngrok
            else
                log_warning "请手动安装 ngrok: https://ngrok.com/download"
            fi
            ;;
        Linux|WSL)
            curl -s https://ngrok-agent.s3.amazonaws.com/ngrok.asc | sudo tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null
            echo "deb https://ngrok-agent.s3.amazonaws.com buster main" | sudo tee /etc/apt/sources.list.d/ngrok.list
            sudo apt update && sudo apt install ngrok
            ;;
        *)
            log_warning "请手动安装 ngrok: https://ngrok.com/download"
            ;;
    esac
}

# 安装项目依赖
install_project_dependencies() {
    log_info "正在安装项目依赖..."
    
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
    
    # 安装后端依赖
    log_info "安装后端 Go 依赖..."
    cd "$PROJECT_ROOT/core"
    go mod download
    go mod tidy
    
    # 安装前端依赖
    log_info "安装前端 Node 依赖..."
    cd "$PROJECT_ROOT/frontend"
    npm install
    
    log_success "项目依赖安装完成"
}

# 主检测函数
check_and_install_dependencies() {
    echo ""
    echo "========================================="
    echo "  猎影平台 - 依赖检测"
    echo "========================================="
    echo ""
    
    OS=$(detect_os)
    log_info "检测到操作系统: $OS"
    echo ""
    
    # 检查所有依赖
    GO_OK=true
    NODE_OK=true
    NPM_OK=true
    
    check_go || GO_OK=false
    check_node || NODE_OK=false
    check_npm || NPM_OK=false
    check_ngrok
    check_playwright
    
    echo ""
    
    # 检查是否有缺失的必需依赖
    if [ "$GO_OK" = false ] || [ "$NODE_OK" = false ] || [ "$NPM_OK" = false ]; then
        log_warning "部分必需依赖未安装"
        
        read -p "是否自动安装缺失的依赖? (y/N): " -n 1 -r
        echo
        
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            case "$OS" in
                macOS)
                    [ "$GO_OK" = false ] && install_go_macos
                    [ "$NODE_OK" = false ] && install_node_macos
                    ;;
                Linux|WSL)
                    [ "$GO_OK" = false ] && install_go_linux
                    [ "$NODE_OK" = false ] && install_node_linux
                    ;;
                Windows)
                    log_warning "Windows 系统请手动安装依赖："
                    log_warning "  Go: https://go.dev/dl/"
                    log_warning "  Node.js: https://nodejs.org/"
                    ;;
            esac
            
            # 再次检查
            echo ""
            log_info "重新检查依赖..."
            check_go
            check_node
            check_npm
        fi
    fi
    
    # 安装项目依赖
    echo ""
    read -p "是否安装项目依赖? (Y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Nn]$ ]]; then
        install_project_dependencies
    fi
    
    echo ""
    log_success "依赖检查完成"
    echo ""
}

# 如果直接运行此脚本
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    check_and_install_dependencies
fi
