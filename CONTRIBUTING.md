# 贡献指南

感谢您对 **猎影 (Lieying)** 项目的关注！我们欢迎所有形式的贡献，包括但不限于：

- 🐛 提交 Bug 报告
- 💡 提出新功能建议
- 📝 改进文档
- 🔧 提交代码修复
- 🌍 翻译项目
- 🎨 设计改进

***

## 📋 目录

1. [行为准则](#行为准则)
2. [如何贡献](#如何贡献)
3. [开发环境搭建](#开发环境搭建)
4. [代码规范](#代码规范)
5. [提交规范](#提交规范)
6. [POC 贡献指南](#poc-贡献指南)
7. [文档贡献](#文档贡献)
8. [社区交流](#社区交流)

***

## 行为准则

参与本项目即表示您同意遵守以下准则：

- 🤝 尊重他人，保持友善
- 💬 欢迎建设性的讨论
- 🎯 专注于项目目标
- 🙏 接受不同的观点和经验

***

## 如何贡献

### 报告 Bug

如果您发现了 Bug，请通过 [GitHub Issues](https://github.com/kunlun-sec/lieying/issues) 提交报告。

**提交前请检查：**

- [ ] 搜索现有 Issues，确认问题未被报告
- [ ] 使用最新的代码版本测试
- [ ] 提供详细的复现步骤

**Bug 报告模板：**

```markdown
**描述**
清晰简洁地描述 Bug

**复现步骤**
1. 执行 '...'
2. 点击 '...'
3. 看到错误

**期望行为**
描述期望的正确行为

**截图**
如果适用，添加截图

**环境信息**
- OS: [例如 Ubuntu 20.04]
- Go 版本: [例如 1.21]
- 项目版本: [例如 v1.0.0]

**附加信息**
其他相关信息
```

### 提出新功能

有新功能建议？请通过 GitHub Issues 提交。

**功能请求模板：**

```markdown
**功能描述**
清晰简洁地描述功能

**使用场景**
描述这个功能的使用场景

**期望解决方案**
描述您期望的实现方式

**替代方案**
描述您考虑过的替代方案

**附加信息**
其他相关信息或截图
```

### 提交代码

1. **Fork 项目**
   ```bash
   git clone https://github.com/YOUR_USERNAME/lieying.git
   cd lieying
   ```
2. **创建分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或
   git checkout -b fix/bug-description
   ```
3. **提交更改**
   ```bash
   git add .
   git commit -m "feat: 添加新功能"
   git push origin feature/your-feature-name
   ```
4. **创建 Pull Request**
   - 访问原项目仓库
   - 点击 "New Pull Request"
   - 选择您的分支
   - 填写 PR 描述

***

## 开发环境搭建

### 后端开发 (Go)

```bash
# 1. 克隆仓库
git clone https://github.com/kunlun-sec/lieying.git
cd lieying/core

# 2. 安装依赖
go mod download

# 3. 运行测试
go test ./...

# 4. 启动服务
go run src/main.go server :3000
```

### 前端开发 (React)

```bash
# 1. 进入前端目录
cd lieying/frontend

# 2. 安装依赖
npm install

# 3. 启动开发服务器
npm run dev

# 4. 构建
npm run build
```

### 完整开发环境

```bash
# 终端 1 - 启动后端
cd core && go run src/main.go server :3000

# 终端 2 - 启动前端
cd frontend && npm run dev

# 终端 3 - 运行测试
cd core && go test ./... -v
```

***

## 代码规范

### Go 代码规范

- 使用 `gofmt` 格式化代码
- 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
- 使用有意义的变量名
- 添加必要的注释

**示例：**

```go
// ScanTarget 执行目标扫描
// 参数:
//   - target: 目标地址
//   - options: 扫描选项
// 返回:
//   - *ScanResult: 扫描结果
//   - error: 错误信息
func ScanTarget(target string, options ScanOptions) (*ScanResult, error) {
    // 实现代码
}
```

### 前端代码规范

- 使用 TypeScript
- 遵循 ESLint 规则
- 使用函数式组件
- 添加组件注释

**示例：**

```typescript
/**
 * 漏洞卡片组件
 * @param title - 漏洞标题
 * @param severity - 严重级别
 * @param onVerify - 验证回调
 */
interface VulnerabilityCardProps {
  title: string;
  severity: 'critical' | 'high' | 'medium' | 'low';
  onVerify: () => void;
}

export default function VulnerabilityCard({ 
  title, 
  severity, 
  onVerify 
}: VulnerabilityCardProps) {
  // 组件实现
}
```

***

## 提交规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范。

### 提交类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 代码重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

### 提交示例

```bash
# 新功能
feat: 添加教育SRC批量扫描功能

# Bug 修复
fix: 修复 SQL 注入检测误报问题

# 文档
docs: 更新 API 文档

# 性能优化
perf: 优化端口扫描速度
```

***

## POC 贡献指南

### POC 格式

POC 使用 YAML 格式，参考 Nuclei 模板规范：

```yaml
id: example-sqli

info:
  name: Example SQL Injection
  author: 您的名字
  severity: critical
  description: |
    描述漏洞详情
  reference:
    - https://example.com/vuln
  tags: sqli,example

requests:
  - method: GET
    path:
      - "{{BaseURL}}/api/user?id=1' AND 1=1--"
    
    matchers:
      - type: word
        words:
          - "user_id"
        part: body
```

### POC 目录结构

```
core/pkg/scan/pocs/
├── cve/
│   └── CVE-2024-XXXX.yaml
├── sqli/
│   └── generic-sqli.yaml
├── xss/
│   └── reflected-xss.yaml
└── upload/
    └── arbitrary-upload.yaml
```

### 提交 POC

1. 在对应目录创建 POC 文件
2. 确保 POC 经过测试
3. 提交 PR 时说明测试目标

***

## 文档贡献

### 文档位置

- `README.md` - 项目主页
- `docs/USAGE.md` - 使用指南
- `docs/API.md` - API 文档
- `docs/DEVELOPMENT.md` - 开发文档

### 文档规范

- 使用 Markdown 格式
- 添加必要的代码示例
- 保持中英文标点一致
- 添加目录结构

***

## 社区交流

### 联系方式

- 📧 邮箱: 1978512375\@qq.com
- 微信：XY5431008

### 贡献者荣誉

所有贡献者将被记录在 [CONTRIBUTORS.md](CONTRIBUTORS.md) 中。

***

## 许可证

通过贡献代码，您同意将其授权给项目，遵循 [MIT 许可证](LICENSE)。

***

**再次感谢您的贡献！**

*Made with ❤️ by 昆仑安全实验室(前逍遥安全实验室-逍遥)*
