@echo off
chcp 65001 >nul
echo ==========================================
echo    下载 Go 运行时
echo    昆仑安全实验室(前逍遥安全实验室-逍遥)
echo ==========================================
echo.

set "GO_VERSION=1.21.5"
set "GO_ZIP=go%GO_VERSION%.windows-amd64.zip"
set "GO_URL=https://golang.google.cn/dl/%GO_ZIP%"
set "INSTALL_DIR=%~dp0runtime"

echo [1/3] 准备下载 Go %GO_VERSION%...
echo   下载地址: %GO_URL%
echo   安装目录: %INSTALL_DIR%
echo.

REM 创建目录
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

echo [2/3] 正在下载...
echo   这可能需要几分钟，请耐心等待...
echo.

REM 使用 PowerShell 下载
powershell -Command "& {Invoke-WebRequest -Uri '%GO_URL%' -OutFile '%TEMP%\%GO_ZIP%' -UseBasicParsing}"

if errorlevel 1 (
    echo [X] 下载失败，尝试使用备用镜像...
    set "GO_URL=https://mirrors.aliyun.com/golang/%GO_ZIP%"
    powershell -Command "& {Invoke-WebRequest -Uri '%GO_URL%' -OutFile '%TEMP%\%GO_ZIP%' -UseBasicParsing}"
)

if not exist "%TEMP%\%GO_ZIP%" (
    echo [X] 下载失败，请手动下载:
    echo   1. 访问 https://golang.google.cn/dl/
    echo   2. 下载 %GO_ZIP%
    echo   3. 解压到 %INSTALL_DIR%\go\
    pause
    exit /b 1
)

echo   [OK] 下载完成
echo.

echo [3/3] 解压安装...
powershell -Command "& {Expand-Archive -Path '%TEMP%\%GO_ZIP%' -DestinationPath '%INSTALL_DIR%' -Force}"

if errorlevel 1 (
    echo [X] 解压失败
    pause
    exit /b 1
)

echo   [OK] 解压完成
echo.

REM 验证安装
if exist "%INSTALL_DIR%\go\bin\go.exe" (
    echo ==========================================
    echo    安装成功！
    echo ==========================================
    echo.
    "%INSTALL_DIR%\go\bin\go.exe" version
    echo.
    echo Go 已安装到: %INSTALL_DIR%\go\
    echo.
    echo 现在可以运行项目了！
    echo 请运行 start-with-go.bat
) else (
    echo [X] 安装失败，请检查目录
echo   %INSTALL_DIR%\go\
)

echo.
pause
