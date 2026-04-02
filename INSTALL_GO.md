# Go 安装指南

**昆仑安全实验室(前逍遥安全实验室-逍遥)** 出品

---

## 🚀 快速安装 Go

### 方式1：自动下载（推荐）

双击运行 `download-go.bat`，自动下载并安装 Go。

### 方式2：手动安装

#### 步骤1：下载 Go

浏览器访问以下任一链接：

- 官方：https://golang.google.cn/dl/go1.21.5.windows-amd64.zip
- 阿里云：https://mirrors.aliyun.com/golang/go1.21.5.windows-amd64.zip
- 清华：https://mirrors.tuna.tsinghua.edu.cn/golang/go1.21.5.windows-amd64.zip

#### 步骤2：解压安装

1. 下载完成后，解压 `go1.21.5.windows-amd64.zip`
2. 将解压出的 `go` 文件夹复制到项目目录：
   ```
   D:\ai项目\测试项目\开发\渗透测试工具自动化\runtime\go\
   ```

3. 验证安装：
   ```
   D:\ai项目\测试项目\开发\渗透测试工具自动化\runtime\go\bin\go.exe version
   ```
   应显示：`go version go1.21.5 windows/amd64`

#### 步骤3：运行项目

双击 `start-with-go.bat` 启动完整项目。

---

## 📋 目录结构

安装完成后，项目结构如下：

```
lieying/
├── runtime/
│   └── go/                    # Go 运行时
│       ├── bin/
│       │   └── go.exe        # Go 可执行文件
│       ├── lib/
│       └── ...
├── core/                      # 猎影后端
├── frontend/                  # 猎影前端
├── start-with-go.bat          # 完整版启动脚本 ⭐
└── download-go.bat            # 自动下载脚本
```

---

## ✅ 验证安装

打开命令提示符，执行：

```cmd
cd /d D:\ai项目\测试项目\开发\渗透测试工具自动化
runtime\go\bin\go.exe version
```

如果显示版本信息，说明安装成功！

---

## 🎯 运行完整项目

安装 Go 后，双击 `start-with-go.bat`：

1. 自动启动后端API服务 (端口3000)
2. 自动启动前端开发服务器 (端口5173)
3. 自动打开浏览器访问界面

访问地址：
- 前端：http://localhost:5173
- 后端：http://localhost:3000

---

## ❓ 常见问题

### Q1: 下载速度慢

**解决**：使用国内镜像
- 阿里云：https://mirrors.aliyun.com/golang/
- 清华：https://mirrors.tuna.tsinghua.edu.cn/golang/

### Q2: 解压失败

**解决**：使用 7-Zip 或 WinRAR 解压
- 7-Zip：https://www.7-zip.org/

### Q3: 提示找不到 go.exe

**解决**：检查路径是否正确
- 正确路径：`runtime\go\bin\go.exe`
- 不要多一层 `go` 目录

---

## 📞 技术支持

- 📧 邮箱: 1978512375@qq.com
- 微信: XY5431008

---

**Made with ❤️ by 昆仑安全实验室(前逍遥安全实验室-逍遥)**
