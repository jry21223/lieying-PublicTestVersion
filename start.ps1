# 猎影 (Lieying) 渗透测试平台启动脚本
# 昆仑安全实验室(前逍遥安全实验室-逍遥)

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "   猎影 (Lieying) 渗透测试平台启动脚本" -ForegroundColor Cyan
Write-Host "   昆仑安全实验室(前逍遥安全实验室-逍遥)" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# 获取脚本所在目录
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptPath

# 检查环境
Write-Host "[1/3] 检查环境..." -ForegroundColor Yellow

# 检查 Go
$goPath = Get-Command go -ErrorAction SilentlyContinue
if (-not $goPath) {
    Write-Host "[错误] 未找到 Go，请先安装 Go 1.21+" -ForegroundColor Red
    Write-Host "下载地址: https://golang.org/dl/" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "按任意键退出..."
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit 1
}
Write-Host "  [OK] Go 已安装: $(go version)" -ForegroundColor Green

# 检查 Node.js
$nodePath = Get-Command node -ErrorAction SilentlyContinue
if (-not $nodePath) {
    Write-Host "[错误] 未找到 Node.js，请先安装 Node.js 18+" -ForegroundColor Red
    Write-Host "下载地址: https://nodejs.org/" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "按任意键退出..."
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    exit 1
}
Write-Host "  [OK] Node.js 已安装: $(node --version)" -ForegroundColor Green

Write-Host ""

# 启动后端
Write-Host "[2/3] 启动后端API服务..." -ForegroundColor Yellow
$backendJob = Start-Job -ScriptBlock {
    param($path)
    Set-Location "$path\core"
    Write-Host "  正在安装依赖..." -ForegroundColor Gray
    go mod download
    Write-Host "  启动服务..." -ForegroundColor Gray
    go run src/main.go server :3000
} -ArgumentList $scriptPath

# 等待后端启动
Start-Sleep -Seconds 3

# 检查后端是否成功启动
$backendStatus = Receive-Job -Job $backendJob -Keep
if ($backendJob.State -eq "Failed") {
    Write-Host "[警告] 后端启动可能出现问题，继续尝试启动前端..." -ForegroundColor Yellow
}

# 启动前端
Write-Host "[3/3] 启动前端开发服务器..." -ForegroundColor Yellow
$frontendJob = Start-Job -ScriptBlock {
    param($path)
    Set-Location "$path\frontend"
    Write-Host "  正在安装依赖..." -ForegroundColor Gray
    npm install
    Write-Host "  启动服务..." -ForegroundColor Gray
    npm run dev
} -ArgumentList $scriptPath

Write-Host ""
Write-Host "==========================================" -ForegroundColor Green
Write-Host "   启动完成！" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green
Write-Host ""
Write-Host "请等待几秒钟让服务完全启动..." -ForegroundColor Yellow
Write-Host ""
Write-Host "访问地址:" -ForegroundColor Cyan
Write-Host "  - 前端界面: http://localhost:5173" -ForegroundColor White
Write-Host "  - 后端API:  http://localhost:3000" -ForegroundColor White
Write-Host ""
Write-Host "按任意键打开浏览器 (或按 Ctrl+C 退出)..." -ForegroundColor Yellow

# 等待按键
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

# 打开浏览器
Start-Process "http://localhost:5173"

# 保持脚本运行
Write-Host ""
Write-Host "服务正在后台运行，关闭此窗口将停止服务..." -ForegroundColor Yellow
Write-Host "按任意键停止服务并退出..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

# 停止任务
Stop-Job -Job $backendJob -ErrorAction SilentlyContinue
Stop-Job -Job $frontendJob -ErrorAction SilentlyContinue
Remove-Job -Job $backendJob -ErrorAction SilentlyContinue
Remove-Job -Job $frontendJob -ErrorAction SilentlyContinue

Write-Host "服务已停止" -ForegroundColor Green
