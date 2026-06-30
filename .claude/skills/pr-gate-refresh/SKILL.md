---
name: pr-gate-refresh
description: "Refresh PR gate metrics in the Feishu spreadsheet. Fetches data via openlibing-ops API, reads table headers to understand column mapping, writes metric data preserving A/B identity columns and cell formatting. Trigger: refresh PR gate, update PR metrics."
---

# PR Gate Metrics — Feishu Table Refresher

Two-step workflow: (1) fetch data via script, (2) AI reads table headers,
maps fields, and writes data with lark-cli.

## Quick usage

```bash
# Step 1: fetch data
python3 scripts/fetch.py --start-date 2026-06-12 --end-date 2026-06-30 --output /tmp/pr-gate-data.json
```

Then follow the workflow below to map and write.

## Workflow

### Step 0: Check prerequisites

Both CLIs must be authenticated. If openlibing returns 401, ask user for
fresh cookie/CSRF token → update `~/.openlibing/auth.yaml`. If lark-cli is
not user-authenticated, run split-flow auth login for the wiki domain.

### Step 1: Fetch data

```bash
python3 scripts/fetch.py --start-date <start> --end-date <end> --output /tmp/pr-gate-data.json
```

Produces a JSON array. Each record has:
`product`, `repo`, `branch`, `e2e_p90`, `e2e_avg`, `e2e_p50`,
`build_p50`, `build_p90`, `build_avg`, `build_queue_p90`,
`test_p90`, `test_p50`, `test_avg`, `test_queue_p90`, `check_p90`.

### Step 2: Read table headers

Read rows 1-2 of the spreadsheet to understand the column layout:

```bash
lark-cli sheets +csv-get \
  --spreadsheet-token "YgAhsy6eHh1xDgt0BBgcC7yTnph" \
  --sheet-id "23b407" \
  --range "A1:ZZ2" --as user
```

Columns A (产品) and B (仓库) are identity columns — never modified.
For each data column (C onwards), read the combined header text from
both rows (merge for merged cells). Map each column to a JSON field name
by matching keywords in the header text:

| Header keywords | JSON field |
|---|---|
| E2E, P90 | `e2e_p90` |
| E2E, 平均 | `e2e_avg` |
| E2E, P50 | `e2e_p50` |
| 构建, 排队 | `build_queue_p90` |
| 构建, P50 | `build_p50` |
| 构建, P90 | `build_p90` |
| 构建, 平均 | `build_avg` |
| 测试, 排队 | `test_queue_p90` |
| 测试, P50 | `test_p50` |
| 测试, P90 | `test_p90` |
| 测试, 平均 | `test_avg` |
| 代码检查, P90 | `check_p90` |

Match rule: a column matches when ALL keywords appear in its combined
header text. Check "排队" patterns BEFORE "P90"/"P50" to avoid ambiguity
(e.g. "构建任务排队P90" must match queue, not build_p90).

Build a column map: `{col_letter: json_field_name}`.

### Step 3: Read identity columns

Read A/B from data rows to get the current product/repo order:

```bash
lark-cli sheets +csv-get \
  --spreadsheet-token "YgAhsy6eHh1xDgt0BBgcC7yTnph" \
  --sheet-id "23b407" \
  --range "A3:B200" --as user
```

Parse the CSV: each `[row=N] product,repo` line gives one row. Skip
fully empty rows (both product and repo blank).

### Step 4: Build data CSV

For each product/repo pair from Step 3, look up the record in the fetch
JSON (step 1). For each data column (step 2 map), emit the corresponding
field value. Unmatched pairs get a row of "-" placeholders.

The result is a CSV string with one row per repo, columns in the order
of the column map.

### Step 5: Clear data (content only)

Clear only the data area, preserving formatting:

```bash
lark-cli sheets +cells-clear \
  --spreadsheet-token "YgAhsy6eHh1xDgt0BBgcC7yTnph" \
  --sheet-id "23b407" \
  --range "<start_col>3:<end_col>200" \
  --scope content --as user --yes
```

Use `--scope content` (NOT `--scope all`) to keep cell formatting.

### Step 6: Write data

```bash
lark-cli sheets +csv-put \
  --spreadsheet-token "YgAhsy6eHh1xDgt0BBgcC7yTnph" \
  --sheet-id "23b407" \
  --start-cell "<start_col>3" \
  --csv - --as user < /tmp/data.csv
```

Report: matched count, missing entries, and the written range.

### Step 7: Highlight empty cells and threshold violations

After writing, apply color highlighting:
- **Gray** (`#D9D9D9`): cells with no data (`"-"`)
- **Yellow** (`#FFF2CC`): cells exceeding product thresholds (overrides gray)

Run the highlight script with the column map from Step 2:

```bash
python3 scripts/highlight.py \
  --data /tmp/pr-gate-data.json \
  --col-map '<column map from Step 2 as JSON>' \
  --output /tmp/highlight.sh

bash /tmp/highlight.sh
```

## Threshold rules

| Product | 构建任务平均(min) | 门禁E2E执行平均(min)(不含重试) |
|---------|:----------------:|:---------------------------:|
| MindIE | < 10 | < 30 |
| FrameworkPTAdapter | < 20 | < 60 |
| Ascend-CANN | < 10 | < 30 |
| MindSpeed | < 10 | < 60 |
| Others (MindCluster, MindSpore, MindStudio) | < 10 | < 30 |

Values at or above the threshold → **yellow** (`#FFF2CC`).  
Cells with no data (`"-"`) → **gray** (`#D9D9D9`).  
Yellow takes precedence over gray when both conditions apply.

## Fixed reference

| Setting | Value |
|---------|-------|
| Spreadsheet token | `YgAhsy6eHh1xDgt0BBgcC7yTnph` |
| Sheet ID | `23b407` |
| Identity columns | A (产品), B (仓库) — read-only |
| Header rows | 1-2 (may be merged) |
| Data start row | 3 |
| Project mapping | See `scripts/fetch.py` PROJECTS dict |
