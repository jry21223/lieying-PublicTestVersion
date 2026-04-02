@echo off
chcp 65001 >nul
echo ==========================================
echo    猎影 (Lieying) 渗透测试平台启动脚本
echo    昆仑安全实验室(前逍遥安全实验室-逍遥)
echo ==========================================
echo.

REM 检查环境
echo [1/3] 检查环境...

where go >nul 2>nul
if %errorlevel% neq 0 (
    echo [错误] 未找到 Go，请先安装 Go 1.21+
    echo 下载地址: https://golang.org/dl/
    pause
    exit /b 1
)

where node >nul 2>nul
if %errorlevel% neq 0 (
    echo [错误] 未找到 Node.js，请先安装 Node.js 18+
    echo 下载地址: https://nodejs.org/
    pause
    exit /b 1
)

echo [OK] 环境检查通过
echo.

REM 启动后端
echo [2/3] 启动后端API服务...
start "猎影后端" cmd /k "cd /d %~dp0core && echo 正在安装依赖... && go mod download && echo 启动服务... && go run src/main.go server :3000"

REM 等待后端启动
timeout /t 3 /nobreak >nul

REM 启动前端
echo [3/3] 启动前端开发服务器...
start "猎影前端" cmd /k "cd /d %~dp0frontend && echo 正在安装依赖... && npm install && echo 启动服务... && npm run dev"

echo.
echo ==========================================
echo    启动完成！
echo ==========================================
echo.
echo 请等待几秒钟让服务完全启动...
echo.
echo 访问地址:
echo   - 前端界面: http://localhost:5173
echo   - 后端API:  http://localhost:3000
echo.
echo 按任意键打开浏览器...
pause >nul

start http://localhost:5173
