@echo off
chcp 65001 >nul
title 猎影渗透测试平台启动器

echo ==========================================
echo    猎影 (Lieying) 渗透测试平台
echo    昆仑安全实验室(前逍遥安全实验室-逍遥)
echo ==========================================
echo.

REM 检查Go
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo [X] 错误: 未找到 Go
    echo.
    echo 请先安装 Go 1.21+
    echo 下载地址: https://golang.org/dl/
    echo.
    pause
    exit /b 1
)
echo [OK] Go 已安装

REM 检查Node.js
where node >nul 2>nul
if %errorlevel% neq 0 (
    echo [X] 错误: 未找到 Node.js
    echo.
    echo 请先安装 Node.js 18+
    echo 下载地址: https://nodejs.org/
    echo.
    pause
    exit /b 1
)
echo [OK] Node.js 已安装

echo.
echo ==========================================
echo    正在启动服务...
echo ==========================================
echo.

REM 获取当前目录
set "BASE_DIR=%~dp0"

REM 启动后端
echo [1/2] 启动后端服务...
start "猎影后端" cmd /c "cd /d "%BASE_DIR%core" && go mod download && go run src/main.go server :3000"

REM 等待3秒
timeout /t 3 /nobreak >nul

REM 启动前端
echo [2/2] 启动前端服务...
start "猎影前端" cmd /c "cd /d "%BASE_DIR%frontend" && npm install && npm run dev"

echo.
echo ==========================================
echo    服务启动中，请稍候...
echo ==========================================
echo.
echo 访问地址:
echo   - 前端: http://localhost:5173
echo   - 后端: http://localhost:3000
echo.

timeout /t 5 /nobreak >nul

echo 正在打开浏览器...
start http://localhost:5173

echo.
echo 提示: 关闭此窗口不会停止服务
echo      请手动关闭后端和前端的命令行窗口
echo.
pause
