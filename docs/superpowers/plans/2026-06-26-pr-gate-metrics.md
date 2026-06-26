# PR Gate Metrics SPC — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add 2 SPC YAML files for querying PR gate (repo-pr-pipeline) metrics from openlibing-ops.

**Architecture:** Pure SPC-file addition — zero Go code changes. Two SPCs share identical source/parameters but differ in output: one emits full JSON (all ~60 fields), the other emits a curated table (14 fields). Follows existing SPC-First pattern.

**Tech Stack:** YAML only. Validated by existing `engine.ParseSPC` and `engine.Validate`.

## Global Constraints

- SPC `type` must be `query`
- SPC files must live in `embedded/spc/`
- Source endpoint is `gateway/openlibing-ops/common/detail`, method `POST`
- Body template must NOT quote integer params (projectId, page, pageSize)
- All SPC files must parse without error via `engine.ParseSPC`
- `go test ./...` and `go vet ./...` must pass before commit

---

### Task 1: Create pr-gate-metrics SPC (full JSON output)

**Files:**
- Create: `embedded/spc/pr-gate-metrics.spc.yaml`

**Interfaces:**
- Consumes: nothing
- Produces: `pr-gate-metrics` SPC — query type, format `json`, no output fields (engine auto-extracts all keys)

- [ ] **Step 1: Write the SPC file**

```yaml
name: pr-gate-metrics
version: "1.0"
description: >
  Query PR gate (repo-pr-pipeline) metrics for a GitCode project.
  Returns full raw JSON with all ~60 metrics fields including
  E2E execution, build, test, and code check durations.
type: query
category: metrics
tags: [pr, gate, metrics, e2e, build, test]

parameters:
  - name: project_id
    type: integer
    required: true
    description: GitCode project ID

  - name: start_date
    type: string
    required: true
    description: Start date in YYYY-MM-DD format

  - name: end_date
    type: string
    required: true
    description: End date in YYYY-MM-DD format

  - name: sort_field
    type: string
    default: "total"
    description: Sort field

  - name: sort_rule
    type: string
    default: "desc"
    enum: [asc, desc]
    description: Sort direction

  - name: page
    type: integer
    default: 1
    description: Page number
    validation:
      min: 1

  - name: page_size
    type: integer
    default: 10
    description: Results per page
    validation:
      min: 1
      max: 100

source:
  method: POST
  endpoint: gateway/openlibing-ops/common/detail
  headers:
    Content-Type: application/json
  body: |
    {"category":"repo-pr-pipeline","projectId":{{.project_id}},"startDate":"{{.start_date}}","endDate":"{{.end_date}}","pipelineStatus":"","sortField":"{{.sort_field}}","sortRule":"{{.sort_rule}}","page":{{.page}},"pageSize":{{.page_size}}}

output:
  format: json

ai:
  prompt_hint: >
    Use this when users ask about PR gate metrics, PR pipeline efficiency,
    E2E execution time, build/test durations, or CI gate performance for a project.
  natural_language:
    - "show PR gate metrics for project 4 this month"
    - "PR pipeline E2E durations"
    - "how long do PR checks take"
    - "gate metrics for Ascend projects"

examples:
  - command: |
      openlibing run pr-gate-metrics \
        --project-id 4 \
        --start-date 2026-06-01 \
        --end-date 2026-06-26
    description: Query all PR gate metrics for project 4 in June
  - command: |
      openlibing run pr-gate-metrics \
        --project-id 4 \
        --start-date 2026-05-27 \
        --end-date 2026-06-26 \
        --output json
    description: Get raw JSON for further processing
```

- [ ] **Step 2: Verify the SPC parses correctly**

Run:
```bash
go test ./internal/engine/ -run TestParseSPC -v -count=1
```
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add embedded/spc/pr-gate-metrics.spc.yaml
git commit -m "feat: add pr-gate-metrics SPC for full PR gate metrics query

Adds pr-gate-metrics SPC (format: json) that queries the repo-pr-pipeline
category from openlibing-ops common/detail endpoint. Returns all ~60
metrics fields for a given project and date range.

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### Task 2: Create pr-gate-metrics-summary SPC (curated table output)

**Files:**
- Create: `embedded/spc/pr-gate-metrics-summary.spc.yaml`

**Interfaces:**
- Consumes: same endpoint/params as `pr-gate-metrics`
- Produces: `pr-gate-metrics-summary` SPC — query type, format `table`, 14 curated fields

- [ ] **Step 1: Write the SPC file**

```yaml
name: pr-gate-metrics-summary
version: "1.0"
description: >
  Curated summary of PR gate (repo-pr-pipeline) metrics — 14 key fields
  including E2E execution (P50/P90/avg), build task time (P50/P90/avg/queue),
  test task time (P50/P90/avg/queue), and code check time (P90).
type: query
category: metrics
tags: [pr, gate, metrics, summary, table]

parameters:
  - name: project_id
    type: integer
    required: true
    description: GitCode project ID

  - name: start_date
    type: string
    required: true
    description: Start date in YYYY-MM-DD format

  - name: end_date
    type: string
    required: true
    description: End date in YYYY-MM-DD format

  - name: sort_field
    type: string
    default: "total"
    description: Sort field

  - name: sort_rule
    type: string
    default: "desc"
    enum: [asc, desc]
    description: Sort direction

  - name: page
    type: integer
    default: 1
    description: Page number
    validation:
      min: 1

  - name: page_size
    type: integer
    default: 10
    description: Results per page
    validation:
      min: 1
      max: 100

source:
  method: POST
  endpoint: gateway/openlibing-ops/common/detail
  headers:
    Content-Type: application/json
  body: |
    {"category":"repo-pr-pipeline","projectId":{{.project_id}},"startDate":"{{.start_date}},"endDate":"{{.end_date}},"pipelineStatus":"","sortField":"{{.sort_field}}","sortRule":"{{.sort_rule}}","page":{{.page}},"pageSize":{{.page_size}}}

output:
  format: table
  fields:
    - name: repo
      header: "Repo"
      path: "repoName"
      width: 16
    - name: branch
      header: "Branch"
      path: "branchName"
      width: 12
    - name: e2e_p90
      header: "E2E P90(min)"
      path: "prE2eAvgTimeP90"
      transform: truncate:8
    - name: e2e_avg
      header: "E2E Avg(min)"
      path: "prE2eAvgTime"
      transform: truncate:8
    - name: e2e_p50
      header: "E2E P50(min)"
      path: "prE2eAvgTimeP50"
      transform: truncate:8
    - name: build_p50
      header: "Build P50(min)"
      path: "buildAvgTimeP50"
      transform: truncate:8
    - name: build_p90
      header: "Build P90(min)"
      path: "buildAvgTimeP90"
      transform: truncate:8
    - name: build_avg
      header: "Build Avg(min)"
      path: "buildAvgTime"
      transform: truncate:8
    - name: build_queue_p90
      header: "BldQ P90(min)"
      path: "buildAvgPendingTimeP90"
      transform: truncate:8
    - name: test_p90
      header: "Test P90(min)"
      path: "dtAvgTimeP90"
      transform: truncate:8
    - name: test_p50
      header: "Test P50(min)"
      path: "dtAvgTimeP50"
      transform: truncate:8
    - name: test_avg
      header: "Test Avg(min)"
      path: "dtAvgTime"
      transform: truncate:8
    - name: test_queue_p90
      header: "TstQ P90(min)"
      path: "dtAvgPendingTimeP90"
      transform: truncate:8
    - name: check_p90
      header: "Chk P90(min)"
      path: "checkAvgTimeP90"
      transform: truncate:8

ai:
  prompt_hint: >
    Use this when users want a curated summary of PR gate metrics —
    key percentiles (P50, P90) and averages for E2E execution, build,
    test, and code check tasks.
  natural_language:
    - "PR gate summary for project 4"
    - "show me PR pipeline key metrics"
    - "summary of PR E2E and build times"
    - "PR门禁核心指标"
  result_hint: >
    Highlight repos with high E2E P90 times (> 60 min). Compare build vs
    test vs code-check durations to identify bottlenecks.

examples:
  - command: |
      openlibing run pr-gate-metrics-summary \
        --project-id 4 \
        --start-date 2026-06-01 \
        --end-date 2026-06-26
    description: Show curated PR gate metrics table for project 4
  - command: |
      openlibing run pr-gate-metrics-summary \
        --project-id 4 \
        --start-date 2026-05-27 \
        --end-date 2026-06-26 \
        --sort-field prE2eAvgTimeP90 \
        --sort-rule desc
    description: Sort by E2E P90 to find slowest repos
```

- [ ] **Step 2: Verify the SPC parses correctly**

Run:
```bash
go test ./internal/engine/ -run TestParseSPC -v -count=1
```
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add embedded/spc/pr-gate-metrics-summary.spc.yaml
git commit -m "feat: add pr-gate-metrics-summary SPC for curated PR gate table

Adds pr-gate-metrics-summary SPC (format: table) with 14 curated fields:
repo/branch context + 12 key metrics across E2E execution, build, test,
and code check categories. Shares the same source/parameters as the
full pr-gate-metrics SPC.

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### Task 3: Write SPC parsing test for the new files

**Files:**
- Create: `internal/engine/pr_gate_metrics_test.go`

**Interfaces:**
- Consumes: `embedded/spc/` (via `embedded.SPCs`)
- Produces: test coverage verifying both SPCs parse, params validate, and templates resolve correctly

- [ ] **Step 1: Write the test file**

```go
package engine

import (
	"os"
	"strings"
	"testing"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

func loadSPCFile(t *testing.T, path string) *spc.SPCDefinition {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	def, err := ParseSPC(strings.NewReader(string(data)))
	if err != nil {
		t.Fatalf("ParseSPC %s: %v", path, err)
	}
	return def
}

func TestPrGateMetrics_ParsesSuccessfully(t *testing.T) {
	def := loadSPCFile(t, "../../embedded/spc/pr-gate-metrics.spc.yaml")

	if def.Name != "pr-gate-metrics" {
		t.Errorf("name = %q, want pr-gate-metrics", def.Name)
	}
	if def.Type != "query" {
		t.Errorf("type = %q, want query", def.Type)
	}
	if def.Output.Format != "json" {
		t.Errorf("output.format = %q, want json", def.Output.Format)
	}
	if len(def.Output.Fields) != 0 {
		t.Errorf("fields = %d, want 0 (full output)", len(def.Output.Fields))
	}
	if len(def.Parameters) != 7 {
		t.Errorf("parameters = %d, want 7", len(def.Parameters))
	}
}

func TestPrGateMetricsSummary_ParsesSuccessfully(t *testing.T) {
	def := loadSPCFile(t, "../../embedded/spc/pr-gate-metrics-summary.spc.yaml")

	if def.Name != "pr-gate-metrics-summary" {
		t.Errorf("name = %q, want pr-gate-metrics-summary", def.Name)
	}
	if def.Type != "query" {
		t.Errorf("type = %q, want query", def.Type)
	}
	if def.Output.Format != "table" {
		t.Errorf("output.format = %q, want table", def.Output.Format)
	}
	if len(def.Output.Fields) != 14 {
		t.Errorf("fields = %d, want 14", len(def.Output.Fields))
	}

	// Verify the 14 field names are as expected
	expectedFields := []string{
		"repo", "branch",
		"e2e_p90", "e2e_avg", "e2e_p50",
		"build_p50", "build_p90", "build_avg", "build_queue_p90",
		"test_p90", "test_p50", "test_avg", "test_queue_p90",
		"check_p90",
	}
	for i, expected := range expectedFields {
		if def.Output.Fields[i].Name != expected {
			t.Errorf("field[%d].name = %q, want %q", i, def.Output.Fields[i].Name, expected)
		}
	}
}

func TestPrGateMetrics_Validate_AllRequired(t *testing.T) {
	def := loadSPCFile(t, "../../embedded/spc/pr-gate-metrics.spc.yaml")

	params := map[string]interface{}{
		"project_id": 4,
		"start_date": "2026-06-01",
		"end_date":   "2026-06-26",
	}
	err := Validate(def, params)
	if err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestPrGateMetrics_Validate_MissingRequired(t *testing.T) {
	def := loadSPCFile(t, "../../embedded/spc/pr-gate-metrics.spc.yaml")

	params := map[string]interface{}{
		"project_id": 4,
		// start_date and end_date missing
	}
	err := Validate(def, params)
	if err == nil {
		t.Fatal("expected validation error for missing required params")
	}
}

func TestPrGateMetrics_Resolve_BodyTemplate(t *testing.T) {
	def := loadSPCFile(t, "../../embedded/spc/pr-gate-metrics.spc.yaml)

	params := map[string]interface{}{
		"project_id": 4,
		"start_date": "2026-06-01",
		"end_date":   "2026-06-26",
		"sort_field": "total",
		"sort_rule":  "desc",
		"page":       1,
		"page_size":  10,
	}

	resolved, err := Resolve(&def.Source, params)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	// Integer params must NOT be quoted in the body
	if strings.Contains(resolved.Body, `"projectId":"4"`) {
		t.Error("projectId should not be quoted (integer)")
	}
	if strings.Contains(resolved.Body, `"page":"1"`) {
		t.Error("page should not be quoted (integer)")
	}
	if strings.Contains(resolved.Body, `"pageSize":"10"`) {
		t.Error("pageSize should not be quoted (integer)")
	}

	// String params must be quoted
	if !strings.Contains(resolved.Body, `"startDate":"2026-06-01"`) {
		t.Error("startDate missing or incorrectly formatted")
	}
	if !strings.Contains(resolved.Body, `"category":"repo-pr-pipeline"`) {
		t.Error("category missing or incorrect")
	}
}

func TestPrGateMetricsSummary_SameParamsAsBase(t *testing.T) {
	base := loadSPCFile(t, "../../embedded/spc/pr-gate-metrics.spc.yaml")
	summary := loadSPCFile(t, "../../embedded/spc/pr-gate-metrics-summary.spc.yaml")

	// Both SPCs should have identical parameters
	if len(base.Parameters) != len(summary.Parameters) {
		t.Fatalf("param count mismatch: base=%d summary=%d",
			len(base.Parameters), len(summary.Parameters))
	}
	for i := range base.Parameters {
		b, s := base.Parameters[i], summary.Parameters[i]
		if b.Name != s.Name || b.Type != s.Type || b.Required != s.Required {
			t.Errorf("param[%d] mismatch: base={%s,%s,%v} summary={%s,%s,%v}",
				i, b.Name, b.Type, b.Required, s.Name, s.Type, s.Required)
		}
	}

	// Both SPCs should have identical source (endpoint, method, headers, body)
	if base.Source.Endpoint != summary.Source.Endpoint {
		t.Error("endpoint mismatch")
	}
	if base.Source.Method != summary.Source.Method {
		t.Error("method mismatch")
	}
	if base.Source.Body != summary.Source.Body {
		t.Error("body mismatch")
	}
}
```

- [ ] **Step 2: Run the new tests (expect FAIL — files don't exist yet)**

Wait — skip this step. The test file reads from disk and the SPC files will be created in Tasks 1 & 2. This test validates the already-created files.

- [ ] **Step 3: Run the new tests (after Tasks 1 & 2)**

Run:
```bash
go test ./internal/engine/ -run TestPrGate -v -count=1
```
Expected: PASS (all 6 tests)

- [ ] **Step 4: Run full test suite**

Run:
```bash
go test ./... -v -count=1
```
Expected: PASS (all packages)

- [ ] **Step 5: Run go vet**

Run:
```bash
go vet ./...
```
Expected: no output (no issues)

- [ ] **Step 6: Commit**

```bash
git add internal/engine/pr_gate_metrics_test.go
git commit -m "test: add SPC parsing and validation tests for pr-gate-metrics

Verifies both SPC files parse correctly, parameters validate as expected,
body templates resolve without quoting integer params, and both SPCs
share identical parameters and source definitions.

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### Task 4: Final verification and integration

**Files:**
- None (verification only)

- [ ] **Step 1: Verify SPCs appear in registry**

Run:
```bash
go build -o bin/openlibing ./cmd/openlibing/ && ./bin/openlibing list | grep pr-gate
```
Expected: both `pr-gate-metrics` and `pr-gate-metrics-summary` appear

- [ ] **Step 2: Verify SPC inspection works**

Run:
```bash
./bin/openlibing inspect pr-gate-metrics-summary
```
Expected: shows parameters, fields, and examples

- [ ] **Step 3: Final check**

Run:
```bash
make check
```
Expected: "All checks passed — ready to commit"
