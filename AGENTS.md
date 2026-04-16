# CLAUDE.md — Goploy Project Guide

## Project Overview

Goploy (go + deploy) 是一个基于 Go + Vue 3 的 Web 部署系统，支持 Git/SVN/FTP/SFTP 代码发布与回滚、RBAC 权限管理、服务器监控、定时任务、终端模拟等功能。

- **版本**: 1.17.5
- **License**: GPLv3
- **文档**: https://docs.goploy.icu
- **API 文档**: https://api-doc.goploy.icu

## Tech Stack

| 层级 | 技术 |
|------|------|
| Backend | Go 1.21+ (toolchain go1.22.5) |
| Frontend | Vue 3 + TypeScript + Vite + Element Plus |
| Database | MySQL (utf8mb4, InnoDB) |
| Query Builder | Squirrel |
| Config | TOML (Koanf v2), 支持热重载 |
| Auth | JWT + Cookie + API Key + LDAP (可选) |
| WebSocket | Gorilla WebSocket |
| Logging | Logrus |
| Docker | Alpine Linux base image |

## Repository Structure

```
cmd/server/              # 应用入口和 API 处理
  main.go                # 入口, 版本号在此维护 (appVersion)
  api/                   # 按领域划分的 API 模块 (deploy/, project/, user/, server/ ...)
    middleware/           # 中间件 (登录日志、操作日志、权限检查、Webhook过滤)
  task/                  # 后台任务 (监控、部署队列)
  ws/                    # WebSocket (实时部署进度、Xterm终端、SFTP)
config/                  # 配置管理 (TOML), 事件驱动的配置重载
database/                # SQL 迁移文件 (按版本号命名, 如 1.17.2.sql)
  goploy.sql             # 完整初始化 schema
  embed.go               # Go embed 嵌入 SQL 文件
internal/
  cache/                 # 缓存层 (memory 实现, 可扩展 Redis)
  media/                 # 外部集成 (钉钉、飞书)
  model/                 # 数据访问层, 每个实体一个文件
  monitor/               # HTTP/TCP/Process/Script 健康检查
  notify/                # 通知系统 (微信、钉钉、飞书、Webhook)
  pipeline/docker/       # Docker 构建与执行
  pkg/                   # 工具包 (git, svn, ssh, cmd, copy, recorder)
  repo/                  # 仓库抽象 (Git/SVN/FTP/SFTP, Factory 模式)
  server/                # HTTP 服务框架 (路由、中间件、请求解码、响应封装)
  transmitter/           # 文件传输抽象 (SFTP/Rsync/Custom, Factory 模式)
  validator/             # 输入校验
web/                     # Vue 3 前端
  src/api/               # Axios API 客户端
  src/views/             # 页面组件
  src/components/        # 可复用组件
  src/composables/       # Composition API 工具
  src/store/             # Vuex 状态管理
  src/lang/              # 国际化 (中/英)
  src/router/            # 路由配置
```

## Build & Run

### Backend

```bash
# 开发运行
go run cmd/server/main.go

# 构建 (Linux amd64)
env GOOS=linux GOARCH=amd64 go build -o goploy cmd/server/main.go

# 跨平台构建 (使用 build.sh)
./build.sh
```

### Frontend

```bash
cd web
npm install
npm run dev          # 开发模式 (端口 8000, API 代理到后端)
npm run build        # 生产构建
```

### Configuration

配置文件: `goploy.toml` (参考 `goploy.example.toml`)

关键配置项:
- `[db]` — MySQL 连接 (必需)
- `[web]` — HTTP 端口 (默认 80)
- `[jwt]` — JWT 密钥
- `[app]` — deployLimit, shutdownTimeout, repositoryPath
- `env` — `production` | `development`

## Database

- 引擎: MySQL
- 初始化: `database/goploy.sql`
- 迁移: 启动时自动执行, 按版本号文件 (如 `database/1.17.2.sql`)
- 迁移逻辑: `model.Update(version)` 在 `internal/model/model.go`

## Code Conventions

### API 路由注册模式

```go
// 每个 API 模块实现 Handler() 返回 []server.Route
func (u User) Handler() []server.Route {
    return []server.Route{
        server.NewRoute("/user/getList", http.MethodGet, u.GetList).
            Permissions(config.ShowMemberPage),
        server.NewWhiteRoute("/user/login", http.MethodPost, u.Login),
    }
}

// 在 main.go 中注册
srv.Router.Register(user.User{})
```

### 请求/响应

- 请求: URL query params + JSON body (body 优先)
- 验证: Validator tag (`validate:"required,min=1,max=25"`)
- 响应: `response.JSON` 封装 `{code, message, data}`

### 工厂模式

仓库 (`internal/repo/factory.go`) 和传输器 (`internal/transmitter/factory.go`) 使用工厂模式创建具体实现。

### 错误处理

```go
if err != nil {
    return fmt.Errorf("meaningful context: %w", err)
}
```

### 配置热重载

通过 `config.GetEventBus()` 订阅配置变更事件。

## Testing

```bash
# 运行所有测试
go test ./...

# 带竞态检测
go test -race ./...

# 带覆盖率
go test -cover ./...
```

现有测试文件:
- `internal/server/decode_test.go`
- `internal/model/model_test.go`
- `cmd/server/api/user/ldap_test.go`

## Key API Endpoints

| 模块 | 端点 | 用途 |
|------|------|------|
| Deploy | `POST /deploy/publish` | 触发部署 |
| Deploy | `POST /deploy/rebuild` | 回滚 |
| Deploy | `GET /deploy/getList` | 部署列表 |
| Deploy | `POST /deploy/webhook` | Webhook 触发 |
| Project | CRUD | 项目管理 |
| User | `/user/login` | 登录 |
| Server | CRUD | 服务器管理 |
| Monitor | CRUD | 监控管理 |
| Cron | CRUD | 定时任务 |

## Environment Variables

| 变量 | 说明 |
|------|------|
| `GOPLOY_URL` | API 地址 (默认 http://localhost:3001) |
| `GOPLOY_API_KEY` | API 认证密钥 |
| `GOPLOY_NAMESPACE_ID` | 命名空间 ID (默认 1) |

## Version Management

版本号在以下位置同步维护:
- `cmd/server/main.go` — `appVersion` 常量 + `@version` 注解
- `database/goploy.sql` — Schema 版本
- `docker/Dockerfile` — `GOPLOY_VER` ARG
- `web/package.json` — `version` 字段

`build.sh` 提供自动化版本号更新。
