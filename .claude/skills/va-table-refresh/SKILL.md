---
name: va-table-refresh
description: "Refresh version availability data in the Feishu spreadsheet. Fetches from openlibing-ops nightly-dashboard API, reads table headers to map columns, writes metric data preserving identity columns and formatting. Trigger: refresh version availability table, update VA data, /va-table-refresh."
---

# Version Availability — Feishu Table Refresher

Two-step workflow: (1) fetch data via script, (2) AI reads table headers,
maps fields, and writes data with lark-cli.

## Quick usage

```bash
python3 scripts/fetch.py --start-date 2026-06-12 --end-date 2026-06-30 --output /tmp/va-data.json
```

Then follow the workflow below to map and write.

## Workflow

### Step 0: Check prerequisites

Both CLIs must be authenticated. If openlibing returns 401, ask user for
fresh cookie/CSRF token and update `~/.openlibing/auth.yaml`.
If lark-cli is not user-authenticated, run split-flow auth login.

### Step 1: Fetch data

```bash
python3 scripts/fetch.py --start-date <start> --end-date <end> --output /tmp/va-data.json
```

Produces a JSON array. Each record has:
`product`, `pipeline`, `versionAvailabilityRate`,
`actualDurationP50Minutes`, `actualDurationP90Minutes`, `actualDurationAvgMinutes`,
`buildTimeP50Minutes`, `buildTimeP90Minutes`, `buildTimeAvgMinutes`,
`testTimeP50Minutes`, `testTimeP90Minutes`, `testTimeAvgMinutes`,
`caseReleaseRateP0`, `casePassRate`.

### Step 2: Read table headers

Read rows 1-2 of the spreadsheet to understand column layout:

```bash
lark-cli sheets +csv-get \
  --spreadsheet-token "YgAhsy6eHh1xDgt0BBgcC7yTnph" \
  --sheet-id "Ou0Fnk" \
  --range "A1:ZZ2" --as user
```

Columns A (产品) and B (流水线) are identity columns — never modified.
For each data column (C onwards), merge header text from both rows, then
map to a JSON field by matching keywords:

| Header keywords | JSON field |
|---|---|
| 版本可用度 | `versionAvailabilityRate` |
| E2E, P50 | `actualDurationP50Minutes` |
| E2E, P90 | `actualDurationP90Minutes` |
| E2E, 平均 | `actualDurationAvgMinutes` |
| 编译, P50 | `buildTimeP50Minutes` |
| 编译, P90 | `buildTimeP90Minutes` |
| 编译, 平均 | `buildTimeAvgMinutes` |
| 测试, P50 | `testTimeP50Minutes` |
| 测试, P90 | `testTimeP90Minutes` |
| 测试, 平均 | `testTimeAvgMinutes` |
| P0, 执行率 | `caseReleaseRateP0` |
| 通过率 | `casePassRate` |

Match rule: a column matches when ALL keywords appear in its combined
header text (case-insensitive). Check more specific patterns first.

Build a column map: `{col_letter: json_field_name}`.

### Step 3: Read identity columns

Read A/B from data rows to get the current product/pipeline order:

```bash
lark-cli sheets +csv-get \
  --spreadsheet-token "YgAhsy6eHh1xDgt0BBgcC7yTnph" \
  --sheet-id "Ou0Fnk" \
  --range "A3:B200" --as user
```

Parse each `[row=N] product,pipeline` line. Skip fully empty rows.

### Step 4: Build data CSV

For each product/pipeline pair from Step 3, look up the record in the
fetch JSON. Values formatted as: percentages → `"XX.X%"`, minutes → `"XX.X"`,
null/missing → `"-"`. Build CSV in column map order.

### Step 5: Clear data (content only)

```bash
lark-cli sheets +cells-clear \
  --spreadsheet-token "YgAhsy6eHh1xDgt0BBgcC7yTnph" \
  --sheet-id "Ou0Fnk" \
  --range "<start_col>3:<end_col>200" \
  --scope content --as user --yes
```

Use `--scope content` to preserve cell formatting.

### Step 6: Write data

```bash
lark-cli sheets +csv-put \
  --spreadsheet-token "YgAhsy6eHh1xDgt0BBgcC7yTnph" \
  --sheet-id "Ou0Fnk" \
  --start-cell "<start_col>3" \
  --csv - --as user < /tmp/va-csv.txt
```

Report: matched count, missing entries, written range.

## Fixed reference

| Setting | Value |
|---------|-------|
| Spreadsheet token | `YgAhsy6eHh1xDgt0BBgcC7yTnph` |
| Sheet ID | `Ou0Fnk` |
| Identity columns | A (产品), B (流水线) — read-only |
| Header rows | 1-2 (may be merged) |
| Data start row | 3 |
| Products | FrameworkPTAdapter, MindSpeed, MindIE, Ascend-CANN, MindStudio, MindCluster, MindSpore (7 total) |
