# 快速启动指南

**昆仑安全实验室(前逍遥安全实验室-逍遥)** 出品

---

## 🚀 快速启动（推荐）

### Windows 用户

双击运行 `start.bat` 脚本，自动完成环境检查、依赖安装和服务启动。

### 手动启动

#### 1. 环境要求

- **Go** 1.21+ [下载](https://golang.org/dl/)
- **Node.js** 18+ [下载](https://nodejs.org/)
- **Git** (可选)

#### 2. 启动后端服务

```bash
cd core
go mod download
go run src/main.go server :3000
```

服务启动后，API 地址：`http://localhost:3000`

#### 3. 启动前端界面

在**新终端**中执行：

```bash
cd frontend
npm install
npm run dev
```

前端地址：`http://localhost:5173`

#### 4. 访问界面

打开浏览器访问：`http://localhost:5173`

---

## 📁 项目结构

```
lieying/
├── core/              # Go 后端
│   └── src/main.go   # 程序入口
├── frontend/          # React 前端
│   └── src/          # 源代码
├── docs/             # 文档
├── README.md         # 项目介绍
└── start.bat         # Windows 启动脚本
```

---

## 🛠️ 常用命令

### 后端命令

```bash
# 启动 API 服务器
go run src/main.go server :3000

# 运行测试
go test ./...

# 构建二进制文件
go build -o lieying.exe src/main.go
```

### 前端命令

```bash
# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览生产版本
npm run preview
```

### CLI 命令

```bash
# 信息收集
./lieying recon example.com

# 漏洞扫描
./lieying scan example.com

# AI 助手
./lieying ai

# 教育SRC
./lieying edu stats
```

---

## 🔧 配置说明

配置文件位置：`~/.lieying/config.yaml`

首次运行会自动创建默认配置。

---

## ❓ 常见问题

### Q: 端口被占用怎么办？

修改启动端口：

```bash
# 后端使用 8080 端口
go run src/main.go server :8080

# 前端修改 vite.config.ts 中的端口
```

### Q: 如何停止服务？

- 后端：在终端按 `Ctrl+C`
- 前端：在终端按 `Ctrl+C`

### Q: 如何更新项目？

```bash
git pull origin main
cd core && go mod download
cd ../frontend && npm install
```

---

## 📞 技术支持

- 📧 邮箱: 1978512375@qq.com
- 微信: XY5431008

---

**Made with ❤️ by 昆仑安全实验室(前逍遥安全实验室-逍遥)**
