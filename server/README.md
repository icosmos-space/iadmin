## server项目结构

```shell
├── api
│   └── v1
├── config
├── core
├── docs
├── global
├── initialize
│   └── internal
├── middleware
├── model
│   ├── request
│   └── response
├── packfile
├── resource
│   ├── excel
│   ├── page
│   └── template
├── router
├── service
├── source
└── utils
    ├── timer
    └── upload
```

| 文件夹       | 说明                    | 描述                        |
| ------------ | ----------------------- | --------------------------- |
| `api`        | api层                   | api层 |
| `--v1`       | v1版本接口              | v1版本接口                  |
| `config`     | 配置包                  | config.yaml对应的配置结构体 |
| `core`       | 核心文件                | 核心组件(zap, viper, server)的初始化 |
| `docs`       | swagger文档目录         | swagger文档目录 |
| `global`     | 全局对象                | 全局对象 |
| `initialize` | 初始化 | router,redis,gorm,validator, timer的初始化 |
| `--internal` | 初始化内部函数 | gorm 的 longger 自定义,在此文件夹的函数只能由 `initialize` 层进行调用 |
| `middleware` | 中间件层 | 用于存放 `gin` 中间件代码 |
| `model`      | 模型层                  | 模型对应数据表              |
| `--request`  | 入参结构体              | 接收前端发送到后端的数据。  |
| `--response` | 出参结构体              | 返回给前端的数据结构体      |
| `packfile`   | 静态文件打包            | 静态文件打包 |
| `resource`   | 静态资源文件夹          | 负责存放静态文件                |
| `--excel` | excel导入导出默认路径 | excel导入导出默认路径 |
| `--page` | 表单生成器 | 表单生成器 打包后的dist |
| `--template` | 模板 | 模板文件夹,存放的是代码生成器的模板 |
| `router`     | 路由层                  | 路由层 |
| `service`    | service层               | 存放业务逻辑问题 |
| `source` | source层 | 存放初始化数据的函数 |
| `utils`      | 工具包                  | 工具函数封装            |
| `--timer` | timer | 定时器接口封装 |
| `--upload`      | oss                  | oss接口封装        |

## 构建说明

- **生产（默认）**：不含 Swagger UI 内嵌页、`swaggo/files`、`mcp` 整包及 MCP HTTP 路由；体积更小。  
  `go build -ldflags="-s -w" -o server .`
- **开发（文档 + MCP）**：  
  `go run -tags=dev .` 或 `go build -tags=dev -ldflags="-s -w" -o server .`  
  需 API 文档时请先执行 `swag init`（见仓库根 `Makefile` 的 `doc` 目标）。

### 按需编译 SQL 驱动（减小体积）

默认 **不** 加 `driver_custom` 时，与原先一致：五种 SQL 驱动（MySQL / PostgreSQL / Oracle / SQL Server / SQLite）都会编进二进制。

若需只包含部分驱动，加上 **`driver_custom`**，并列出需要的 **`driver_*`** 标签（至少一个）：

| 标签 | 含义 |
|------|------|
| `driver_mysql` | MySQL |
| `driver_pgsql` | PostgreSQL |
| `driver_oracle` | Oracle |
| `driver_mssql` | SQL Server |
| `driver_sqlite` | SQLite（glebarez / modernc） |

示例（仅 SQLite，常与生产默认配合）：

```bash
go build -tags "driver_custom,driver_sqlite" -ldflags "-s -w" -o server .
```

示例（MySQL + SQLite）：

```bash
go build -tags "driver_custom,driver_mysql,driver_sqlite" -ldflags "-s -w" -o server .
```

**注意**：启用 `driver_custom` 后，运行时的 `system.db-type` / `db-list` 只能使用已编入的驱动，否则会在连接时 panic 并提示重新编译标签。

本地用 gopls 编辑时，若需要分析某组合驱动，可在 `go.buildTags` / `gopls.build.buildFlags` 中写入相同的 `-tags=...`。

### 初始化页与后端驱动一致

`internal/initdbcaps` 与 `driver_custom` / `driver_*` 标签同步。接口 **`GET /init/db-types`** 返回 `{ dbTypes: ["mysql", ...] }`，前端初始化页只展示当前二进制支持的类型；提交 **`POST /init/initdb`** 时也会校验，避免选到未编译的驱动。

### 前端嵌入二进制（embedweb）

占位目录：**`initialize/resource/web/dist/`**（与 `initialize/embedded_web.go` 中 `//go:embed` 相对路径一致）。

1. 在仓库根构建前端：`cd web && npm run build`（或 `yarn build`）。
2. 将 **`web/dist/` 下全部文件** 复制到 **`server/initialize/resource/web/dist/`**（覆盖占位 `index.html`）。
3. 编译：`go build -tags=embedweb -ldflags "-s -w" -o server .`

启用后，未匹配 API 的页面请求会回退到 `index.html`（SPA）。请保证 **`system.router-prefix`** 与前端环境变量 **`VITE_BASE_API`** 一致（例如均为 `/api`），否则接口路径与静态页会错位。

不加 `-tags=embedweb` 时不嵌入前端，行为与原先一致（可继续用 Nginx 反代静态资源）。

