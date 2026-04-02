@echo off
chcp 65001 >nul
title 猎影渗透测试平台 - 完整版

echo ==========================================
echo    猎影 (Lieying) 渗透测试平台
echo    完整版 - 内置Go运行时
echo    昆仑安全实验室(前逍遥安全实验室-逍遥)
echo ==========================================
echo.

REM 设置路径
set "BASE_DIR=%~dp0"
set "GO_PATH=%BASE_DIR%runtime\go\bin"
set "BACKEND_PORT=3000"
set "FRONTEND_PORT=5173"

REM 检查Go是否存在
if not exist "%GO_PATH%\go.exe" (
    echo [X] 未找到 Go 运行时
echo.
    echo 请先运行 download-go.bat 安装 Go
    echo 或手动下载 Go 并解压到 runtime\go\
    echo.
    echo 下载地址:
    echo   https://golang.google.cn/dl/go1.21.5.windows-amd64.zip
    echo.
    pause
    exit /b 1
)

echo [OK] 找到 Go: %GO_PATH%
echo.

REM 添加Go到PATH
set "PATH=%GO_PATH%;%PATH%"

REM 验证Go
echo [1/4] 验证 Go 版本...
for /f "tokens=*" %%a in ('"%GO_PATH%\go.exe" version') do (set GO_VERSION=%%a)
echo   %GO_VERSION%
echo.

REM 启动后端
echo [2/4] 启动后端API服务...
echo   端口: %BACKEND_PORT%
start "猎影后端" cmd /k "cd /d "%BASE_DIR%core" && echo 安装依赖... && "%GO_PATH%\go.exe" mod download && echo 启动服务... && "%GO_PATH%\go.exe" run src/main.go server :%BACKEND_PORT%"
timeout /t 5 /nobreak >nul
echo   [OK] 后端服务已启动
echo.

REM 启动前端
echo [3/4] 启动前端开发服务器...
echo   端口: %FRONTEND_PORT%
start "猎影前端" cmd /k "cd /d "%BASE_DIR%frontend" && echo 安装依赖... && npm install && echo 启动服务... && npm run dev -- --port %FRONTEND_PORT%"
timeout /t 8 /nobreak >nul
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
