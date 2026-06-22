# 测试流程指南

## 测试分层

```
make check                              ← CI 门禁 (一键全量)
  ├── make build                        ← 编译
  ├── make test-unit                    ← 单元测试 (快速)
  ├── make test-integration             ← 集成测试 (mock server)
  └── make lint                         ← 静态分析
```

| 目标 | 命令 | 耗时 | 场景 |
|------|------|------|------|
| `make test-unit` | 单元测试 | ~3s | 开发中频繁运行 |
| `make test-integration` | 集成测试 | ~5s | 改 CLI 命令后运行 |
| `make test` | 全量测试 | ~10s | 提交前运行 |
| `make check` | 完整门禁 | ~15s | CI / 合并前 |
| `make test-cover` | 覆盖率报告 | ~10s | 检查测试覆盖 |

## 快速开始

```bash
# 开发中 — 只跑单元测试
make test-unit

# 改完代码 — 完整验证
make check

# 查看覆盖率
make test-cover
```

---

## 新增 CLI 命令时的测试流程

当新增一个 SPC 或 CLI 子命令时，按以下步骤补充测试：

### Step 1: 单元测试 (如果新增了包内逻辑)

在对应包的 `*_test.go` 中添加测试。

**示例 — 新增 SPC 类型字段：**

```go
// pkg/spc/types_test.go
func TestSPCDefinition_Unmarshal_NewField(t *testing.T) {
    input := `
name: test
version: "1.0"
description: Test new field
type: query
category: test
new_field: "value"     # <-- 新增字段
source:
  method: GET
  endpoint: /test
output:
  format: json
`
    var spc SPCDefinition
    err := yaml.Unmarshal([]byte(input), &spc)
    if err != nil {
        t.Fatalf("unmarshal failed: %v", err)
    }
    if spc.NewField != "value" {
        t.Errorf("NewField = %q, want 'value'", spc.NewField)
    }
}
```

### Step 2: Engine 层测试 (如果新增了执行逻辑)

**示例 — 新增 extract 路径：**

```go
// internal/engine/engine_test.go
func TestEngine_Extract_NewFormat(t *testing.T) {
    // ... setup mock server with new response format ...
    result, err := engine.Execute("my-spc", params)
    if err != nil {
        t.Fatalf("Execute failed: %v", err)
    }
    if result.Rows[0]["new_field"] != expected {
        t.Errorf("new_field = %v", result.Rows[0]["new_field"])
    }
}
```

### Step 3: 集成测试 (必须 — 覆盖 CLI 端到端)

在 `cmd/openlibing/integration_test.go` 中按模板添加。

**模板：**

```go
func TestIntegration_MyNewFeature(t *testing.T) {
    // 1. 启动 mock server
    server := testutil.NewMockOpenLibingServer(nil, "")
    defer server.Close()

    // 2. 配置 HOME 指向 mock
    _, cleanup := testutil.SetupTestHome(t, server.URL)
    defer cleanup()

    // 3. 构建最新二进制
    bin := buildBinary(t)

    // 4. 测试正常路径
    t.Run("success", func(t *testing.T) {
        out, code := runCLI(t, bin, "run", "my-new-spc", "--param", "value")
        if code != 0 {
            t.Fatalf("exit %d\n%s", code, out)
        }
        // 断言输出包含预期内容
        if !strings.Contains(out, "expected_output") {
            t.Errorf("output missing expected content\n%s", out)
        }
    })

    // 5. 测试错误路径
    t.Run("missing_required", func(t *testing.T) {
        _, code := runCLI(t, bin, "run", "my-new-spc")
        if code == 0 {
            t.Fatal("expected non-zero exit for missing required param")
        }
    })

    // 6. 测试输出格式
    t.Run("json_output", func(t *testing.T) {
        out, code := runCLI(t, bin, "run", "my-new-spc", "--param", "v", "--output", "json")
        if code != 0 {
            t.Fatalf("exit %d\n%s", code, out)
        }
        var rows []map[string]interface{}
        if err := json.Unmarshal([]byte(out), &rows); err != nil {
            t.Fatalf("invalid JSON: %v", err)
        }
        // ... 断言 rows 内容 ...
    })
}
```

### Step 4: 验证全量回归

```bash
# 确保所有历史测试不受影响
make check
```

---

## 测试辅助工具

### `testutil.NewMockOpenLibingServer`

创建模拟 OpenLibing API 的 HTTP 服务器：

```go
import "github.com/openlibing/openlibing-cli/internal/testutil"

// 使用默认数据
server := testutil.NewMockOpenLibingServer(nil, "")
defer server.Close()

// 使用自定义数据
customRuns := []testutil.PipelineRun{
    {PipelineRunID: "custom-1", Status: "SUCCESS"},
}
server := testutil.NewMockOpenLibingServer(customRuns, "custom log output")
```

### `testutil.SetupTestHome`

创建临时 HOME 目录，其中 `~/.openlibing/config.yaml` 指向 mock 服务器：

```go
home, cleanup := testutil.SetupTestHome(t, server.URL)
defer cleanup()
// 此时 HOME 已指向临时目录，CLI 会读取该目录的配置
```

### `buildBinary` / `runCLI`

```go
bin := buildBinary(t)                              // 编译二进制
out, code := runCLI(t, bin, "list", "--format", "json") // 执行并获取输出
```

---

## 测试覆盖矩阵

| 模块 | 单元测试 | 集成测试 | 覆盖重点 |
|------|---------|---------|---------|
| `pkg/spc/` | ✅ 2 | — | YAML 解析、类型正确性 |
| `internal/config/` | ✅ 6 | — | 配置加载、auth 管理 |
| `internal/api/` | ✅ 4 | — | HTTP 请求、重试、认证注入 |
| `internal/engine/` | ✅ 19 | — | 7 步管线、校验、模板渲染、JSON 提取 |
| `internal/registry/` | ✅ 5 | — | 3 层加载、覆盖、搜索 |
| `internal/cli/` | ✅ 9 | — | 输出格式化、duration 变换 |
| `cmd/openlibing/` | ✅ 1 | ✅ 25 | 全命令端到端、flag 标准化、离线可用 |
| `internal/testutil/` | — | (被引用) | Mock server、临时 HOME |

---

## CI 集成

```yaml
# .gitlab-ci.yml / GitHub Actions 示例
test:
  script:
    - make check
```

或者分阶段：

```yaml
unit:
  script: make test-unit
integration:
  script: make test-integration
lint:
  script: make lint
```

---

## 常见问题

**Q: 集成测试失败 "build failed"？**
A: 确保在项目根目录运行 `make test-integration`，它会先 `make build`。

**Q: 如何只跑某个包的测试？**
A: `go test ./internal/engine/... -v -run TestValidate`

**Q: 新增 SPC 后 `list` 测试失败了？**
A: `TestIntegration_FullPipeline/list_table` 和 `list_json` 会检查内置 SPC 数量 (=3)。如果新增了内置 SPC，需要更新断言中的数量。

**Q: 测试跑得太慢？**
A: 开发中用 `make test-unit`（不含集成测试）。提交前用 `make check` 完整验证。
