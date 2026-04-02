@echo off
chcp 65001 >nul
title 猎影渗透测试平台 - 便携版

REM ==========================================
REM 猎影 (Lieying) 便携版启动脚本
REM 昆仑安全实验室(前逍遥安全实验室-逍遥)
REM ==========================================

echo.
echo ==========================================
echo    猎影 (Lieying) 渗透测试平台
echo    便携版 - 零配置，解压即用
echo    昆仑安全实验室(前逍遥安全实验室-逍遥)
echo ==========================================
echo.

REM 设置运行时路径
set "BASE_DIR=%~dp0"
set "GO_PATH=%BASE_DIR%runtime\go\bin"
set "NODE_PATH=%BASE_DIR%runtime\node"
set "BACKEND_PORT=3000"
set "FRONTEND_PORT=5173"

REM 添加运行时到PATH
set "PATH=%GO_PATH%;%NODE_PATH%;%PATH%"

REM 检查运行时是否存在
if not exist "%GO_PATH%\go.exe" (
    echo [X] 错误: 未找到 Go 运行时
    echo.
    echo 请确保 runtime\go 目录存在
    echo 如果缺失，请重新下载便携版
    echo.
    pause
    exit /b 1
)

if not exist "%NODE_PATH%\node.exe" (
    echo [X] 错误: 未找到 Node.js 运行时
    echo.
    echo 请确保 runtime\node 目录存在
    echo 如果缺失，请重新下载便携版
    echo.
    pause
    exit /b 1
)

echo [OK] Go 运行时: %GO_PATH%
echo [OK] Node.js 运行时: %NODE_PATH%
echo.

REM 验证版本
echo [1/4] 验证运行时版本...
for /f "tokens=*" %%a in ('"%GO_PATH%\go.exe" version') do (set GO_VERSION=%%a)
for /f "tokens=*" %%a in ('"%NODE_PATH%\node.exe" --version') do (set NODE_VERSION=%%a)
echo   Go: %GO_VERSION%
echo   Node.js: %NODE_VERSION%
echo.

REM 启动后端
echo [2/4] 启动后端API服务...
echo   端口: %BACKEND_PORT%
start "猎影后端" cmd /c "cd /d "%BASE_DIR%core" && "%GO_PATH%\go.exe" mod download && "%GO_PATH%\go.exe" run src/main.go server :%BACKEND_PORT%"
timeout /t 3 /nobreak >nul
echo   [OK] 后端服务已启动
echo.

REM 启动前端
echo [3/4] 启动前端开发服务器...
echo   端口: %FRONTEND_PORT%
start "猎影前端" cmd /c "cd /d "%BASE_DIR%frontend" && "%NODE_PATH%\npm.cmd" install && "%NODE_PATH%\npm.cmd" run dev -- --port %FRONTEND_PORT%"
timeout /t 5 /nobreak >nul
echo   [OK] 前端服务已启动
echo.

REM 完成
echo ==========================================
echo    启动完成！
echo ==========================================
echo.
echo 访问地址:
echo   - 前端界面: http://localhost:%FRONTEND_PORT%
echo   - 后端API:  http://localhost:%BACKEND_PORT%
echo.
echo 正在打开浏览器...
timeout /t 3 /nobreak >nul
start http://localhost:%FRONTEND_PORT%
echo.
echo 提示:
echo   - 关闭此窗口不会停止服务
echo   - 请手动关闭后端和前端的命令行窗口
echo   - 数据保存在: %USERPROFILE%\.lieying\
echo.
pause
