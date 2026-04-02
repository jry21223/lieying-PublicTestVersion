# 猎影渗透测试平台

一个现代化的渗透测试管理平台，包含资产发现、漏洞扫描、AI分析等功能。

## 项目结构

```
.
├── core/           # Go 后端
│   ├── internal/
│   │   ├── api/      # API 路由和处理
│   │   ├── models/   # 数据模型
│   │   ├── repository/ # 数据仓库
│   │   ├── service/ # 业务逻辑
│   │   └── config/  # 配置
│   └── src/
│       └── main.go
└── frontend/        # React 前端
    └── src/
        ├── components/  # React 组件
        ├── services/    # API 服务
        └── App.tsx
```

## 快速开始

### 后端

```bash
cd core
go build -o lieying.exe ./src
./lieying.exe server
```

### 前端

```bash
cd frontend
npm install
npm run dev
```

## 功能模块

- 资产概览
- 信息收集
- 漏洞扫描任务
- 漏洞结果
- 漏洞看板
- 漏洞管理
- AI对话
- AI配置
- 教育SRC
- 报告生成

## 许可证

MIT
