# OpenLibing CLI — AI 原生架构设计

> 设计日期: 2026-06-22
> 状态: Design Approved

## 1. 项目定位

`openlibing-cli` 是面向 [OpenLibing](https://www.openlibing.com/) CI/CD 可观测性平台的命令行工具。
核心定位是**数据查询**（只读），后期扩展修改能力（写入操作）。

OpenLibing 平台网关架构包含三个核心服务：

| 服务 | 职责 |
|------|------|
| `openlibing-cicd` | Pipeline/Stage/Job/Step 运行详情和执行日志 |
| `openlibing-codecheck` | 代码质量检查（CodeCheck）结果和问题报告 |
| `openlibing-sca` | 软件组成分析（SCA）扫描结果 |

核心实体链：**Project → PipelineRun → Stage → JobRun → StepRun → ExecLog**

## 2. 架构哲学

### SPC-First 架构

**Everything is a Super Power。** CLI 本体是一个极薄的 SPC 执行器，所有能力 — 包括内置子命令 — 都来自 SPC（Skill Pipeline Configuration）文件。

添加新查询能力 = 写一个新的 `.spc.yaml` 文件，零代码。

### 核心理念

- **SPC 即 Super Power** — 每个 `.spc.yaml` 文件定义一个可发现、可组合、可共享的 AI 技能单元
- **SPC 可被 AI 读取** — Chat 模式下 SPC 自动转为 LLM Tool Definition
- **SPC 可被人阅读** — 声明式 YAML，自文档化，自带 examples
- **SPC 分布式管理** — 内置 (embed) / 用户 (~/.openlibing/) / 项目 (./.openlibing/) 三层

## 3. 技术选型

| 决策 | 选择 | 理由 |
|------|------|------|
| 语言 | **Go** 1.22+ | 单一二进制分发，启动快，适合 CLI |
| CLI 框架 | **Cobra** | Go 生态标准，flag 自动生成 |
| SPC 解析 | `gopkg.in/yaml.v3` | YAML 解析标准库 |
| JSON 路径 | `github.com/tidwall/gjson` | 高性能 JSONPath 提取 |
| 模板引擎 | Go `text/template` | 内置，无需外部依赖 |
| 内置 SPC | `//go:embed` | 编译嵌入，零外部文件依赖 |
| LLM SDK | Anthropic Go SDK (后期) | Chat 模式 LLM 调用 |
| HTTP 客户端 | `net/http` | 标准库，无额外依赖 |

## 4. 目录结构

```
openlibing-cli/
├── cmd/openlibing/main.go              # 入口：组装依赖，启动 CLI
├── internal/
│   ├── engine/                          # SPC 执行引擎
│   │   ├── engine.go                   # 核心：解析→校验→执行→输出
│   │   ├── parser.go                   # SPC YAML → 结构化类型
│   │   ├── validator.go                # SPC 完整性 + 参数校验
│   │   └── executor.go                 # 构建 HTTP 请求 + 执行
│   │
│   ├── registry/                        # Super Power 注册中心
│   │   ├── registry.go                 # 注册/发现/索引
│   │   ├── loader.go                   # 文件系统扫描 .spc.yaml
│   │   └── resolver.go                # 意图 → SPC 匹配
│   │
│   ├── api/                             # OpenLibing API 客户端
│   │   ├── client.go                   # HTTP Client + auth 注入
│   │   ├── cicd.go                     # /gateway/openlibing-cicd/*
│   │   ├── codecheck.go                # /gateway/openlibing-codecheck/*
│   │   └── sca.go                     # /gateway/openlibing-sca/*
│   │
│   ├── cli/                             # 命令行 UI 层
│   │   ├── root.go                     # 根命令 + 全局 flag
│   │   ├── run.go                      # run <spc-name> 通用执行
│   │   ├── list.go                     # 列出可用 super powers
│   │   ├── inspect.go                  # 查看 SPC 详情
│   │   ├── chat.go                     # AI 对话模式
│   │   └── output.go                   # table/json/yaml 格式化
│   │
│   ├── chat/                            # AI Chat 模式
│   │   ├── repl.go                     # Read-Eval-Print-Loop
│   │   ├── prompt.go                   # System prompt 构建器
│   │   └── llm.go                     # LLM API 客户端
│   │
│   └── config/                          # 配置
│       ├── config.go                   # ~/.openlibing/config.yaml
│       └── auth.go                     # Token 存储 + 注入
│
├── embed/spc/                           # 内置 SPC（go:embed 编译进二进制）
│   ├── pipeline-list.spc.yaml
│   ├── pipeline-detail.spc.yaml
│   └── pipeline-logs.spc.yaml
│
├── pkg/spc/                             # 公共类型（可被外部引用）
│   └── types.go                        # SPC 结构体定义
│
├── docs/superpowers/
│   ├── specs/                           # 设计文档
│   └── plans/                           # 实现计划
├── go.mod / go.sum / Makefile / README.md
```

### 模块职责一句话

| 模块 | 职责 |
|------|------|
| `engine` | 把 SPC 定义变成 HTTP 请求，把响应变成格式化输出 |
| `registry` | 发现所有可用 SPC（内置 + 用户 + 项目），维护索引 |
| `api` | 封装 OpenLibing 网关的所有 HTTP 接口 |
| `cli` | 用户界面 — 子命令解析、参数绑定、输出打印 |
| `chat` | AI 对话 — NL → SPC 匹配 → 执行 → 结果解释 |
| `config` | 全局配置 + 认证凭据管理 |

**关键设计决策：**

- `embed/spc/` 用 Go 1.16+ `//go:embed` 编译进二进制，内置能力零外部依赖
- `~/.openlibing/spc/` 用户自定义 SPC 目录，启动时自动扫描
- `./.openlibing/spc/` 项目级共享 SPC 目录
- `pkg/spc/types.go` 是唯一的跨模块类型定义，避免循环引用
- `api/` 不依赖 `engine`，可以独立被外部引用

## 5. SPC 文件格式规范

SPC 文件分三层语义：

```
┌─────────────────────────────────────┐
│  Meta Layer    名称/版本/分类/描述    │  ← 发现 & 索引
├─────────────────────────────────────┤
│  Pipeline Layer 参数→请求→变换→输出   │  ← 执行 & 编排
├─────────────────────────────────────┤
│  AI Layer      prompt/hints/examples │  ← Chat 模式增强
└─────────────────────────────────────┘
```

### 完整 Schema

```yaml
# ===== Meta Layer =====
name: pipeline-list               # 唯一标识，对应 openlibing run <name>
version: "1.0"
description: >
  Query pipeline runs for a GitCode project.
  Returns status, branch, duration, and trigger info.
type: query                       # query | action(写操作) | workflow(多步骤)
category: pipeline                # pipeline | codecheck | sca
tags: [ci, status]                # 辅助搜索

# ===== Pipeline Layer =====
parameters:                       # 输入参数 → CLI flag 自动生成
  - name: project_id
    type: string
    required: true
    description: GitCode project identifier

  - name: limit
    type: integer
    default: 10
    validation: { min: 1, max: 100 }

  - name: status
    type: string
    enum: [running, success, failed, pending]
    default: ""

source:                            # 数据源定义
  method: GET
  endpoint: /gateway/openlibing-cicd/project/pipeline/pipeline-run/detail
  query_params:
    projectId: "{{.project_id}}"
    pageSize: "{{.limit}}"
  headers:
    Content-Type: application/json

output:                             # 输出映射
  format: table                     # table | json | yaml | raw
  fields:
    - name: id
      header: "ID"
      path: ".pipelineRunId"
      width: 36
    - name: status
      header: "Status"
      path: ".status"
      transform: upper
    - name: branch
      header: "Branch"
      path: ".ref"
    - name: duration_ms
      header: "Duration"
      path: ".durationMillis"
      transform: duration           # ms → "3m 42s"
    - name: created
      header: "Created"
      path: ".createTime"

# ===== AI Layer =====
ai:
  prompt_hint: >
    Use this when users ask about pipeline runs, CI status, or build history
    for a specific project.
  natural_language:
    - "show pipelines for project 123"
    - "list recent builds"
    - "what's the CI status of project X"
  result_hint: >
    Focus on status and duration. If there are failures, highlight them first.

examples:
  - command: |
      openlibing run pipeline-list --project-id 123
    description: List last 10 pipelines
  - command: |
      openlibing run pipeline-list --project-id 123 --status failed --limit 5
    description: Show 5 most recent failures
```

### 关键规则

**参数与 CLI flag 的映射：**
```
parameters[i].name: project_id
        → CLI flag: --project-id  (snake_case → kebab-case)
        → 模板变量: {{.project_id}}
        → 环境变量: OPENLIBING_PROJECT_ID (fallback)
```

**source 模板变量：**
- `{{.parameter_name}}` — 用户传入的参数值
- `{{.config.key}}` — 配置文件中的值
- `{{.env.VAR}}` — 环境变量
- 不支持复杂逻辑（保证 SPC 的声明性）

**output path 语法** — 使用 `gjson` 路径：
- `.pipelineRunId` — 简单字段
- `.stages.0.name` — 数组索引
- `.stages.#.name` — 数组展开

**type 语义：**
- `query` — 只读，无副作用，不确认
- `action` — 写操作，需要用户确认
- `workflow` — 多步骤编排，展示步骤概览后确认

## 6. Super Power Registry 注册机制

### 加载优先级

```
优先级 1: embed/spc/         编译内置 (MVP 3 个 SPC)
优先级 2: ~/.openlibing/spc/ 用户自定义
优先级 3: ./.openlibing/spc/ 项目级共享

同名覆盖规则: 项目 > 用户 > 内置
```

### 核心接口

```go
type Registry interface {
    LoadAll() error
    Get(name string) (*SPCDefinition, error)
    Search(query string) []*SPCDefinition
    ListByCategory(cat string) []*SPCDefinition
    ListAll() []*SPCDefinition
    Reload() error
}
```

### 发现流程

```
启动 → embed.FS 读取内置 SPC → 扫描 ~/.openlibing/spc/ (覆盖内置)
     → 扫描 ./.openlibing/spc/ (覆盖前面) → 建索引 → Ready
```

### Resolver：NL → SPC 匹配

```
"show me failed pipelines for project 123"
           ↓
  关键词提取: "pipeline" → category hit, "project 123" → parameter hint
           ↓
  索引匹配: name 命中 + description 语义 + tags 交集
           ↓
  排序返回: [pipeline-list: 0.85, pipeline-detail: 0.42]
```

MVP 阶段用关键词 + 索引匹配（确定性，零延迟）。
后期可接入 LLM embedding 做语义匹配。

## 7. 执行引擎数据流

```
用户输入                                     终端输出
┌──────────────────┐                     ┌──────────────────┐
│ openlibing run   │                     │ ID     Status ... │
│ pipeline-list    │                     │ abc123 failed ... │
│ --project-id 456 │                     │ def456 success .. │
└────────┬─────────┘                     └──────────────────┘
         │                                        ▲
         ▼                                        │
┌─────────────────────────────────────────────────────────────┐
│                      Engine.Execute()                       │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │ 1.Parse  │─→│2.Validate│─→│3.Resolve │─→│ 4.Call   │   │
│  │ Registry │  │ 参数校验  │  │ 模板渲染  │  │ HTTP请求 │   │
│  │ 查 SPC   │  │ 类型/必填 │  │ URL/Query │  │          │   │
│  └──────────┘  └──────────┘  └──────────┘  └────┬─────┘   │
│                                                   │         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐        │         │
│  │7.Output │←─│6.Transform│←─│5.Extract │←───────┘         │
│  │ table等  │  │ duration等│  │ gjson路径 │                  │
│  └──────────┘  └──────────┘  └──────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

### 各步骤职责

| 步骤 | 输入 | 动作 | 产物 |
|------|------|------|------|
| Parse | SPC name + CLI flags | Registry.Get() + flag 绑定到 parameters | 完整的 SPCDefinition |
| Validate | parameters + 用户值 | required/type/enum/range 检查 | pass / 明确错误信息 |
| Resolve | source 模板 + 参数值 | `text/template` 渲染 URL/Query/Header | 完整 HTTP 请求定义 |
| Call | HTTP 请求定义 | api.Client 执行，自动注入 auth，3 次重试 | HTTP 响应体 (JSON) |
| Extract | JSON + output.fields | gjson 路径逐字段提取 | []Row (结构化行数据) |
| Transform | Row + field.transform | upper/lower/duration/truncate | 变换后的 Row |
| Output | []Row + format | table/json/yaml 渲染 | 终端输出 |

### 错误处理矩阵

| 错误层 | 处理方式 |
|--------|---------|
| SPC 解析失败 | 启动时报错，指出文件:行号 |
| 参数校验失败 | 显示缺失参数 + usage 示例 |
| 网络超时 | 自动重试 3 次，间隔 1s/2s/4s |
| HTTP 4xx | 解析 body 错误信息给用户 |
| HTTP 5xx | 显示"平台暂时不可用" + 状态码 |
| 响应解析失败 | 降级输出 raw JSON |
| SPC 未找到 | 建议 `openlibing list` |

## 8. CLI 子命令设计

SPC-First 原则：没有硬编码的领域子命令。所有查询能力来自 SPC 文件，CLI 只提供元命令。

```
openlibing
├── run      <spc-name> [flags]     # 执行一个 Super Power
├── list     [--category]            # 列出可用的 Super Powers
├── inspect  <spc-name>             # 查看 SPC 的完整定义
└── chat                             # 进入 AI 对话模式
```

### 使用示例

```bash
# 执行 SPC：参数自动映射为 CLI flags
openlibing run pipeline-list --project-id 123
openlibing run pipeline-list --project-id 123 --status failed --limit 10
openlibing run pipeline-detail --run-id abc-def-123

# 覆写输出格式
openlibing run pipeline-list --project-id 123 -o json

# 发现可用能力
openlibing list
openlibing list --category pipeline
openlibing list --format json

# 查看 SPC 详情
openlibing inspect pipeline-list

# AI 对话模式
openlibing chat
```

### 后期演进：层级子命令路由

```
SPC name: pipeline-list       → openlibing pipeline list       (自动)
SPC name: pipeline-detail     → openlibing pipeline detail
SPC name: codecheck-report    → openlibing codecheck report
# 由 SPC category 驱动第一级，name 剩余部分驱动第二级，CLI 零硬编码
```

## 9. Chat 模式设计

### 架构

```
┌─────────────────────────────────────────────────┐
│              openlibing chat                     │
│                                                  │
│  ┌────────────┐     ┌─────────────────────────┐ │
│  │   REPL      │────→│  LLM Gateway            │ │
│  │  循环       │     │  (Anthropic API / 兼容) │ │
│  └────────────┘     └───────────┬─────────────┘ │
│        ↑                        ▼               │
│        │         ┌─────────────────────────┐    │
│        │         │  Tool Registry           │    │
│        │         │  SPC → LLM Tool Def      │    │
│        │         └───────────┬─────────────┘    │
│        │                     ▼                  │
│        │         ┌─────────────────────────┐    │
│        └─────────│  Engine.Execute()        │    │
│   结果回注+解释   │  (复用 run 的同一引擎)   │    │
│                  └─────────────────────────┘    │
└─────────────────────────────────────────────────┘
```

### 关键设计

**SPC → Tool Definition 自动转换：**
SPC 的 `name` → tool name，`description` → tool description，
`parameters` → JSON Schema `input_schema`。Engine 是 Chat 和 Run 的唯一执行路径。

**确认机制：**
- `type: query` → 直接执行，不确认
- `type: action` → "About to trigger pipeline re-run. Continue? [y/N]"
- `type: workflow` → 展示步骤概览后确认

**上下文记忆：**
最近 5 条交互结果保持在对话上下文中（按轮次），不持久化。

## 10. 配置管理

```
~/.openlibing/
├── config.yaml          # 用户级配置
├── auth.yaml            # 认证凭据 (权限 600 强制)
└── spc/                 # 用户自定义 SPC
```

### 配置优先级

```
CLI flag > 环境变量 > config.yaml > 硬编码默认值
```

## 11. MVP 范围

| 维度 | MVP (v0.1) | 后期 |
|------|-----------|------|
| SPC 类型 | `query` only | `action`, `workflow` |
| 内置 SPC | 3 (pipeline-list/detail/logs) | codecheck, sca 系列 |
| CLI 命令 | `run`, `list`, `inspect` | 层级子命令路由 |
| Chat 模式 | 框架就绪，LLM 接入可后配 | 完整 AI 对话 |
| 配置 | config.yaml + auth.yaml | 多 profile 切换 |
| 输出格式 | table, json | yaml, csv, 自定义模板 |
| 认证 | Bearer token | OAuth2 登录流程 |

## 12. 测试策略

| 层 | 测试类型 | 覆盖重点 |
|----|---------|---------|
| `pkg/spc` | 单元测试 | YAML 解析、类型正确性 |
| `api/` | 单元测试 + mock | HTTP 请求构建、重试、错误映射 |
| `engine/` | 单元测试 (mock HTTP) | 全流程 7 步骤 |
| `registry/` | 单元测试 | 多层加载、覆盖规则、搜索排序 |
| `cli/` | 集成测试 | 端到端二进制执行 |
| `chat/` | 集成测试 | Mock LLM → tool call → engine |

## 13. 架构总览

```
┌─────────────────────────────────────────────────────────┐
│                    openlibing CLI                        │
│                                                         │
│  用户界面     run │ list │ inspect │ chat                │
│                     │                                   │
│  ┌──────────────────┴───────────────────────────────┐  │
│  │              SPC Engine                           │  │
│  │  Parse → Validate → Resolve → Call → Extract     │  │
│  │                    → Transform → Output           │  │
│  └──────┬──────────────────────────────┬───────────┘  │
│         │                              │               │
│  ┌──────┴──────┐              ┌────────┴───────────┐  │
│  │  Registry   │              │   API Client        │  │
│  │  发现/索引   │              │   HTTP + Auth       │  │
│  │  SPC 仓库   │              │   cicd/codecheck/   │  │
│  │             │              │   sca               │  │
│  └─────────────┘              └────────────────────┘  │
│                                                         │
│  ┌──────────────────────────────────────────────────┐  │
│  │  SPC Files (Git-style distributed)                │  │
│  │  embed/  ~/.openlibing/spc/  ./.openlibing/spc/  │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

## 14. 后续演进路径

```
v0.1 MVP         v0.2 扩展          v0.3 AI               v1.0 完整
─────            ─────              ─────                 ─────
3 内置 SPC       + CodeCheck SPC    Chat 模式正式启用     层级路由
run/list/inspect + SCA SPC          LLM 语义匹配          写操作 (action)
query only       + 自定义 format    对话上下文增强         workflow 编排
table/json       + 分页支持         AI 结果解读           SPC 市场/共享
Bearer auth                         多 profile            插件体系
```

---

🤖 Generated with [Claude Code](https://claude.com/claude-code)
