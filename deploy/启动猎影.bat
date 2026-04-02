@echo off
chcp 65001 >nul
title 猎影渗透测试平台

echo ==========================================
echo   猎影渗透测试平台 - 一键启动
echo   Lieying Penetration Testing Platform
echo ==========================================
echo.

cd /d "%~dp0"

echo [1/4] 检查端口占用...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :8081 ^| findstr LISTENING') do (
    echo  关闭占用端口 8081 的进程...
    taskkill /F /PID %%a >nul 2>&1
)
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :5173 ^| findstr LISTENING') do (
    echo  关闭占用端口 5173 的进程...
    taskkill /F /PID %%a >nul 2>&1
)

echo.
echo [2/4] 启动后端服务 (端口 8081)...
if not exist "data\lieying.db" (
    echo  创建数据库...
    if not exist "data" mkdir data
)
start "猎影后端" cmd /k "lieying_server.exe server"
echo  等待后端初始化...
timeout /t 2 /nobreak >nul

echo.
echo [3/4] 启动前端服务 (端口 5173)...
if not exist "frontend" (
    echo  错误: 未找到 frontend 目录
    pause
    exit /b 1
)
start "猎影前端" cmd /k "cd frontend && npx serve -l 5173 -s ."
echo  等待前端启动...
timeout /t 3 /nobreak >nul

echo.
echo ==========================================
echo   启动完成！
echo.
echo   前端地址: http://localhost:5173
echo   后端API:  http://localhost:8081
echo.
echo   按任意键打开浏览器...
pause >nul
start http://localhost:5173
echo.
echo ==========================================
