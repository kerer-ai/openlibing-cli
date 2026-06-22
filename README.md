# openlibing-cli

AI 原生命令行工具，用于查询 [OpenLibing](https://www.openlibing.com/) CI/CD 可观测性平台数据。

## 设计哲学：SPC-First

**Everything is a Super Power。** 所有查询能力都由 SPC（Skill Pipeline Configuration）文件定义。添加新的查询能力 = 写一个 `.spc.yaml` 文件，零代码。

```
openlibing run pipeline-list --project-id 123 --status failed
```

上面这条命令的实际执行流程：

```
SPC 定义 (pipeline-list.spc.yaml)
       │
       ▼
┌─────────────────────────────────────────┐
│  Engine 7 步执行管线                     │
│                                         │
│  Parse  → Validate → Resolve → Call     │
│  查 SPC   参数校验    模板渲染   HTTP请求 │
│                                         │
│  Extract → Transform → Output           │
│  JSON提取  字段变换     终端表格          │
└─────────────────────────────────────────┘
```

## 快速开始

### 安装

```bash
# 构建
cd openlibing-cli
make build

# 二进制位于 bin/openlibing，可直接使用
./bin/openlibing --help
```

### 配置

首次使用，创建认证文件（可选，仅在使用需要认证的 API 时需要）：

```bash
mkdir -p ~/.openlibing

cat > ~/.openlibing/auth.yaml << EOF
openlibing:
  token: "your-openlibing-token"
  token_type: Bearer
EOF

chmod 600 ~/.openlibing/auth.yaml
```

自定义配置（可选）：

```bash
cat > ~/.openlibing/config.yaml << EOF
endpoint: https://www.openlibing.com
defaults:
  project_id: ""      # 默认项目 ID
  limit: 10
output:
  format: table
  color: true
EOF
```

## 命令

```
openlibing
├── run      <spc-name> [flags]   执行一个 Super Power
├── list     [--category]         列出可用的 Super Powers
├── inspect  <spc-name>           查看 SPC 完整定义
└── chat                          进入 AI 对话模式 (即将推出)
```

### `list` — 发现可用能力

```bash
# 表格展示所有可用的 Super Power
$ openlibing list

NAME             TYPE   CATEGORY   DESCRIPTION
────────         ─────  ─────────  ─────────────────
pipeline-detail  query  pipeline   Get detailed information about a specific...
pipeline-list    query  pipeline   Query pipeline runs for a GitCode project...
pipeline-logs    query  pipeline   Fetch execution logs for a specific job step...

# JSON 格式输出（适合脚本消费）
$ openlibing list --format json

# 按分类过滤
$ openlibing list --category pipeline
```

### `inspect` — 查看 SPC 详情

```bash
$ openlibing inspect pipeline-list

Name:        pipeline-list
Version:     1.0
Type:        query
Category:    pipeline
Origin:      builtin
Tags:        ci, status, list

Description:
  Query pipeline runs for a GitCode project.

Parameters:
  --project-id  string (required)
        GitCode project identifier
  --limit  integer [default: 10]
        Max number of results to return
  --status  string
        Filter by pipeline status
        Values: running, success, failed, pending

Source:
  GET gateway/openlibing-cicd/project/pipeline/pipeline-run/detail

Output: table (5 fields)
  id           → ID
  status       → Status [upper]
  branch       → Branch
  duration_ms  → Duration [duration]
  created      → Created

Examples:
  openlibing run pipeline-list --project-id 123
  # List last 10 pipelines for project 123

  openlibing run pipeline-list --project-id 123 --status failed --limit 5
  # Show 5 most recent failures
```

### `run` — 执行查询

```bash
# 查询项目的最近流水线
openlibing run pipeline-list --project-id 123

# 只看失败的，限制 5 条
openlibing run pipeline-list --project-id 123 --status failed --limit 5

# 查看某次流水线的详细信息
openlibing run pipeline-detail --run-id abc-def-123

# 获取执行日志（POST 请求）
openlibing run pipeline-logs \
  --project-id 123 \
  --pipeline-run-id run-abc \
  --job-run-id job-1 \
  --step-run-id step-a

# 输出为 JSON
openlibing run pipeline-list --project-id 123 --output json

# 输出为 YAML
openlibing run pipeline-detail --run-id abc --output yaml
```

### `chat` — AI 对话模式（即将推出）

```bash
$ openlibing chat

Chat mode is coming in a future release.
```

## 内置 Super Powers

### pipeline-list

| 属性 | 值 |
|------|-----|
| 端点 | `GET /gateway/openlibing-cicd/project/pipeline/pipeline-run/detail` |
| 输出 | table (ID, Status, Branch, Duration, Created) |
| 必填参数 | `--project-id` |
| 可选参数 | `--limit` (默认 10, 最大 100), `--status` (running/success/failed/pending) |

### pipeline-detail

| 属性 | 值 |
|------|-----|
| 端点 | `GET /gateway/openlibing-cicd/project/pipeline/pipeline-run/detail` |
| 输出 | json (Pipeline ID, Status, Branch, Duration, Stages 数量) |
| 必填参数 | `--run-id` |

### pipeline-logs

| 属性 | 值 |
|------|-----|
| 端点 | `POST /gateway/openlibing-cicd/project/pipeline/exec-log` |
| 输出 | raw (原始日志文本) |
| 必填参数 | `--project-id`, `--pipeline-run-id`, `--job-run-id`, `--step-run-id` |

## 自定义 SPC

在 `~/.openlibing/spc/` 目录下放置 `.spc.yaml` 文件，CLI 启动时自动发现并可用。

同名 SPC 覆盖规则：**项目 (./.openlibing/spc/) > 用户 (~/.openlibing/spc/) > 内置 (embedded)**

### SPC 文件示例

```yaml
name: my-custom-query
version: "1.0"
description: 我的自定义查询
type: query
category: pipeline

parameters:
  - name: project_id
    type: string
    required: true
    description: 项目 ID

source:
  method: GET
  endpoint: gateway/openlibing-cicd/project/pipeline/pipeline-run/detail
  query_params:
    projectId: "{{.project_id}}"

output:
  format: table
  fields:
    - name: id
      header: "Pipeline ID"
      path: "pipelineRunId"
    - name: status
      header: "Status"
      path: "status"
      transform: upper
```

然后直接使用：

```bash
openlibing run my-custom-query --project-id 456
```

## 输出格式

| 格式 | 说明 | 适用场景 |
|------|------|---------|
| `table` | 对齐的终端表格（默认） | 人类阅读 |
| `json` | 缩进 JSON | 脚本消费、管道给 `jq` |
| `yaml` | YAML 格式 | 配置文件、可读性 |
| `raw` | 原始 HTTP 响应 | 日志查看、调试 |

通过 `--output` / `-o` 覆盖 SPC 默认格式：

```bash
openlibing run pipeline-list --project-id 123 --output json | jq '.[] | select(.status=="FAILED")'
```

## 架构

```
┌─────────────────────────────────────────────────────────┐
│                    openlibing CLI                        │
│                                                         │
│  用户界面     run │ list │ inspect │ chat                │
│                     │                                   │
│  ┌──────────────────┴───────────────────────────────┐  │
│  │              SPC Engine (7-step pipeline)         │  │
│  │  Parse → Validate → Resolve → Call →             │  │
│  │  Extract → Transform → Output                    │  │
│  └──────┬──────────────────────────────┬───────────┘  │
│         │                              │               │
│  ┌──────┴──────┐              ┌────────┴───────────┐  │
│  │  Registry   │              │   API Client        │  │
│  │  3层加载     │              │   HTTP + Auth       │  │
│  │  SPC 索引   │              │   Retry + Timeout   │  │
│  └─────────────┘              └────────────────────┘  │
│                                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │  SPC Files                                       │  │
│  │  embedded/  ~/.openlibing/spc/  ./.openlibing/   │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### 项目结构

```
openlibing-cli/
├── cmd/openlibing/main.go         # 入口，组装依赖
├── embedded/spc/                   # 内置 SPC (go:embed 编译进二进制)
│   ├── pipeline-list.spc.yaml
│   ├── pipeline-detail.spc.yaml
│   └── pipeline-logs.spc.yaml
├── internal/
│   ├── engine/                     # SPC 执行引擎
│   │   ├── engine.go              # 7 步管线编排
│   │   ├── parser.go              # YAML → SPCDefinition
│   │   ├── validator.go           # 参数校验（必填/类型/枚举/范围）
│   │   └── resolver.go           # Go template 渲染
│   ├── registry/                   # Super Power 注册中心
│   │   ├── registry.go           # 注册/发现/索引
│   │   ├── loader.go             # 3 层文件系统加载
│   │   └── resolver.go           # 关键词搜索匹配
│   ├── api/                        # OpenLibing API 客户端
│   │   ├── client.go             # HTTP Client + 认证注入 + 重试
│   │   └── cicd.go               # Pipeline 相关接口
│   ├── cli/                        # 命令行 UI
│   │   ├── root.go               # 根命令 + 全局 flag
│   │   ├── run.go                # run <spc-name>
│   │   ├── list.go               # list [--category]
│   │   ├── inspect.go            # inspect <spc-name>
│   │   ├── chat.go               # chat (stub)
│   │   └── output.go             # table/json/yaml/raw 格式化
│   └── config/                     # 配置管理
│       ├── config.go             # ~/.openlibing/config.yaml
│       └── auth.go               # ~/.openlibing/auth.yaml
└── pkg/spc/types.go               # 共享 SPC 类型定义
```

### 核心依赖

| 库 | 用途 |
|----|------|
| `github.com/spf13/cobra` | CLI 框架 |
| `gopkg.in/yaml.v3` | SPC 文件解析 |
| `github.com/tidwall/gjson` | JSON 路径提取 |

## 开发

```bash
# 克隆
git clone <repo-url> && cd openlibing-cli

# 构建
make build

# 运行全部测试
make test

# 覆盖率
make test-cover

# 静态分析
make lint
```

### 测试统计

| 包 | 测试数 |
|---|--------|
| `pkg/spc` | 2 |
| `internal/config` | 6 |
| `internal/api` | 4 |
| `internal/engine` | 19 |
| `internal/registry` | 5 |
| `internal/cli` | 9 |
| `cmd/openlibing` | 1 |
| **合计** | **46** |

## 路线图

| 版本 | 内容 |
|------|------|
| **v0.1** (当前) | 3 内置 SPC (pipeline), 4 元命令, table/json/yaml 输出, Bearer 认证 |
| v0.2 | CodeCheck + SCA SPC, 自定义 format, 分页支持 |
| v0.3 | Chat 模式正式启用, LLM 语义匹配, AI 结果解读 |
| v1.0 | 层级路由, action/workflow SPC 类型, SPC 共享, OAuth2 |

## 许可证

MIT
