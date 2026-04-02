# 猎影便携版构建指南

**昆仑安全实验室(前逍遥安全实验室-逍遥)** 出品

---

## 📦 构建便携版

便携版包含：
- ✅ 猎影后端 (Go)
- ✅ 猎影前端 (React)
- ✅ Go 1.21+ 运行时
- ✅ Node.js 18+ 运行时
- ✅ 一键启动脚本

---

## 🚀 快速构建步骤

### 1. 准备环境

确保已安装：
- 7-Zip 或 WinRAR（用于打包）
- 约 2GB 磁盘空间

### 2. 下载运行时

#### 下载 Go

1. 访问 https://golang.org/dl/
2. 下载 `go1.21.x.windows-amd64.zip`
3. 解压到 `portable/runtime/go/`

#### 下载 Node.js

1. 访问 https://nodejs.org/dist/v18.x.x/
2. 下载 `node-v18.x.x-win-x64.zip`
3. 解压到 `portable/runtime/node/`

### 3. 复制项目文件

```batch
REM 复制后端
cd portable
xcopy /E /I ..\core core\

REM 复制前端
xcopy /E /I ..\frontend frontend\

REM 复制文档
copy ..\README.md .
copy ..\LICENSE .
```

### 4. 目录结构

```
portable/
├── runtime/
│   ├── go/                    # Go 运行时
│   │   ├── bin/
│   │   ├── lib/
│   │   └── ...
│   └── node/                  # Node.js 运行时
│       ├── node.exe
│       ├── npm.cmd
│       └── ...
├── core/                      # 猎影后端
│   ├── src/
│   ├── pkg/
│   ├── internal/
│   └── go.mod
├── frontend/                  # 猎影前端
│   ├── src/
│   ├── public/
│   └── package.json
├── start.bat                  # 启动脚本
├── README_PORTABLE.md         # 使用说明
└── LICENSE                    # 许可证
```

### 5. 打包发布

使用 7-Zip 打包：

```batch
REM 进入项目根目录
cd d:\ai项目\测试项目\开发\渗透测试工具自动化

REM 使用 7-Zip 打包
"C:\Program Files\7-Zip\7z.exe" a -tzip -r lieying-portable-v1.0.0.zip portable\
```

或使用 WinRAR：

```batch
"C:\Program Files\WinRAR\WinRAR.exe" a -r -afzip lieying-portable-v1.0.0.zip portable\
```

---

## 📋 自动化构建脚本

创建 `build-portable.bat`：

```batch
@echo off
chcp 65001 >nul
echo ==========================================
echo    猎影便携版构建脚本
echo    昆仑安全实验室(前逍遥安全实验室-逍遥)
echo ==========================================
echo.

set "VERSION=1.0.0"
set "OUTPUT=lieying-portable-v%VERSION%.zip"

echo [1/4] 清理旧文件...
if exist portable\core rmdir /S /Q portable\core
if exist portable\frontend rmdir /S /Q portable\frontend
if exist %OUTPUT% del %OUTPUT%
echo   [OK]
echo.

echo [2/4] 复制项目文件...
xcopy /E /I /Y /Q core portable\core\
xcopy /E /I /Y /Q frontend portable\frontend\>
echo   [OK]
echo.

echo [3/4] 检查运行时...
if not exist "portable\runtime\go\bin\go.exe" (
    echo   [警告] 未找到 Go 运行时
    echo   请下载并解压到 portable\runtime\go\
)
if not exist "portable\runtime\node\node.exe" (
    echo   [警告] 未找到 Node.js 运行时
    echo   请下载并解压到 portable\runtime\node\
)
echo.

echo [4/4] 打包文件...
if exist "C:\Program Files\7-Zip\7z.exe" (
    "C:\Program Files\7-Zip\7z.exe" a -tzip -r %OUTPUT% portable\
) else (
    echo   [错误] 未找到 7-Zip，请手动打包
    pause
    exit /b 1
)
echo   [OK]
echo.

echo ==========================================
echo    构建完成！
echo ==========================================
echo.
echo 输出文件: %OUTPUT%
echo.
pause
```

---

## 📥 运行时下载地址

### Go 1.21.5

- 官方：https://go.dev/dl/go1.21.5.windows-amd64.zip
- 镜像：https://golang.google.cn/dl/go1.21.5.windows-amd64.zip

### Node.js 18.19.0

- 官方：https://nodejs.org/dist/v18.19.0/node-v18.19.0-win-x64.zip
- 镜像：https://npmmirror.com/mirrors/node/v18.19.0/node-v18.19.0-win-x64.zip

---

## ✅ 验证便携版

构建完成后，测试以下功能：

1. **解压测试**
   - 将 zip 解压到干净目录
   - 确认所有文件存在

2. **启动测试**
   - 双击 `start.bat`
   - 确认服务正常启动
   - 访问 http://localhost:5173

3. **功能测试**
   - 信息收集模块
   - 漏洞扫描模块
   - 教育SRC模块
   - AI助手模块

---

## 📦 发布清单

发布前确认：

- [ ] 包含 Go 运行时
- [ ] 包含 Node.js 运行时
- [ ] 包含完整的 core 目录
- [ ] 包含完整的 frontend 目录
- [ ] 包含 start.bat 启动脚本
- [ ] 包含 README_PORTABLE.md
- [ ] 包含 LICENSE 文件
- [ ] 已在干净环境测试通过

---

## 📞 技术支持

- 📧 邮箱: 1978512375@qq.com
- 微信: XY5431008

---

**Made with ❤️ by 昆仑安全实验室(前逍遥安全实验室-逍遥)**
