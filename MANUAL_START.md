# 手动启动指南

如果自动脚本无法运行，请按照以下步骤手动启动项目。

---

## 前置要求

确保已安装：
- **Go** 1.21+ : https://golang.org/dl/
- **Node.js** 18+ : https://nodejs.org/

---

## 步骤1: 打开终端

按 `Win + R`，输入 `cmd`，回车打开命令提示符。

---

## 步骤2: 启动后端服务

在终端中依次执行：

```cmd
cd /d D:\ai项目\测试项目\开发\渗透测试工具自动化\core
go mod download
go run src/main.go server :3000
```

看到类似以下输出表示启动成功：
```
🚀 猎影API服务器启动中...
   监听地址: :3000
   数据库: C:\Users\你的用户名\.lieying\lieying.db
```

**保持此窗口运行，不要关闭！**

---

## 步骤3: 打开新终端启动前端

按 `Win + R`，再次输入 `cmd`，回车打开**另一个**命令提示符。

执行：

```cmd
cd /d D:\ai项目\测试项目\开发\渗透测试工具自动化\frontend
npm install
npm run dev
```

看到类似以下输出表示启动成功：
```
  VITE v5.x.x  ready in xxx ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
  ➜  press h + enter to show help
```

**保持此窗口运行，不要关闭！**

---

## 步骤4: 访问界面

打开浏览器，访问：

```
http://localhost:5173
```

---

## 常见问题

### Q1: 提示 "go" 不是内部或外部命令

**解决**: Go 未正确安装或未添加到环境变量
1. 重新安装 Go
2. 重启电脑
3. 在终端输入 `go version` 验证

### Q2: 提示 "npm" 不是内部或外部命令

**解决**: Node.js 未正确安装
1. 重新安装 Node.js
2. 重启电脑
3. 在终端输入 `node -v` 验证

### Q3: 端口被占用

**解决**: 修改端口号

后端使用其他端口（如 8080）：
```cmd
go run src/main.go server :8080
```

前端修改 `frontend/vite.config.ts`：
```typescript
export default defineConfig({
  server: {
    port: 8081,  // 修改这里
  },
})
```

### Q4: npm install 很慢或失败

**解决**: 使用淘宝镜像
```cmd
npm config set registry https://registry.npmmirror.com
npm install
```

### Q5: 如何停止服务

- **后端**: 在第一个终端窗口按 `Ctrl + C`
- **前端**: 在第二个终端窗口按 `Ctrl + C`

---

## 联系方式

- 📧 邮箱: 1978512375@qq.com
- 微信: XY5431008

---

**Made with ❤️ by 昆仑安全实验室(前逍遥安全实验室-逍遥)**
