# PR Gate Metrics SPC — Design Spec

## Overview

新增 2 个 SPC，用于查询 PR 门禁（repo-pr-pipeline）指标数据。数据源为 openlibing-ops 的 `common/detail` 端点，category=repo-pr-pipeline。

## API Mapping

| Endpoint | `gateway/openlibing-ops/common/detail` |
|----------|---------------------------------------|
| Method | POST |
| Category | `repo-pr-pipeline` |

Payload 模板：
```json
{"category":"repo-pr-pipeline","projectId":<project_id>,"startDate":"<start_date>","endDate":"<end_date>","pipelineStatus":"","sortField":"<sort_field>","sortRule":"<sort_rule>","page":<page>,"pageSize":<page_size>}
```

## SPC Definitions

### 1. pr-gate-metrics (base / full data)

- **Purpose**: 拉取全量原始 JSON 数据（~60 个指标字段）
- **Output format**: `json`
- **Fields**: 无（Engine 自动提取所有键）

### 2. pr-gate-metrics-summary (curated view)

- **Purpose**: 精选 12 个核心门禁指标，表格输出
- **Output format**: `table`

Field mapping:

| # | Header | API Path | Description |
|---|--------|----------|-------------|
| 1 | Repo | `repoName` | 仓库名 |
| 2 | Branch | `branchName` | 分支 |
| 3 | E2E P90(min) | `prE2eAvgTimeP90` | 门禁E2E执行(不含重试) P90 |
| 4 | E2E Avg(min) | `prE2eAvgTime` | 门禁E2E执行(不含重试) 平均 |
| 5 | E2E P50(min) | `prE2eAvgTimeP50` | 门禁E2E执行(不含重试) P50 |
| 6 | Build P50(min) | `buildAvgTimeP50` | 构建任务 P50 |
| 7 | Build P90(min) | `buildAvgTimeP90` | 构建任务 P90 |
| 8 | Build Avg(min) | `buildAvgTime` | 构建任务 平均 |
| 9 | Build Queue P90(min) | `buildAvgPendingTimeP90` | 构建任务排队 P90 |
| 10 | Test P90(min) | `dtAvgTimeP90` | 测试任务 P90 |
| 11 | Test P50(min) | `dtAvgTimeP50` | 测试任务 P50 |
| 12 | Test Avg(min) | `dtAvgTime` | 测试任务 平均 |
| 13 | Test Queue P90(min) | `dtAvgPendingTimeP90` | 测试任务排队 P90 |
| 14 | Check P90(min) | `checkAvgTimeP90` | 代码检查任务 P90 |

### Shared Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| project_id | integer | true | — | GitCode 项目 ID |
| start_date | string | true | — | 开始日期 YYYY-MM-DD |
| end_date | string | true | — | 结束日期 YYYY-MM-DD |
| sort_field | string | false | total | 排序字段 |
| sort_rule | string | false | desc | 排序方向 (asc/desc) |
| page | integer | false | 1 | 页码 |
| page_size | integer | false | 10 | 每页条数 (1-100) |

### Shared Source

```yaml
source:
  method: POST
  endpoint: gateway/openlibing-ops/common/detail
  headers:
    Content-Type: application/json
  body: |
    {"category":"repo-pr-pipeline","projectId":{{.project_id}},"startDate":"{{.start_date}}","endDate":"{{.end_date}}","pipelineStatus":"","sortField":"{{.sort_field}}","sortRule":"{{.sort_rule}}","page":{{.page}},"pageSize":{{.page_size}}}
```

## Files Changed

| File | Action |
|------|--------|
| `embedded/spc/pr-gate-metrics.spc.yaml` | Create |
| `embedded/spc/pr-gate-metrics-summary.spc.yaml` | Create |

No Go source changes. Zero code impact.

## Design Decisions

1. **Source duplication over inheritance**: SPC 模型当前不支持 `extends`/`$ref`。复制 source 保持现有模式，且 2 个文件的维护成本极低。
2. **No runtime field selection**: 不扩展 `--fields` 参数。通过不同 SPC 提供不同的字段组合，符合 CLI 的子命令语义。
3. **API values are already in minutes**: 原始数据中 `prE2eAvgTimeP90` 等字段数值单位已是分钟级，无需 transform 转换。

## Testing

- SPC 文件语法正确（可被 `ParseSPC` 解析）
- 参数验证（必填项检查）
- Template 渲染正确（body 中整数参数不产生引号）
