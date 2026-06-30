#!/usr/bin/env python3
"""
Detect threshold violations and empty cells, output lark-cli commands to
apply yellow/gray background highlighting.

Usage:
  python3 scripts/highlight.py --data /tmp/pr-gate-data.json \
    --col-map '{"C":"e2e_p50",...}' --output /tmp/highlight.sh
"""

import argparse
import json
import sys
from collections import defaultdict


# Product threshold rules: product_name -> {field: max_value_minutes}
THRESHOLDS = {
    "MindIE":             {"e2e_avg": 30, "build_avg": 10},
    "FrameworkPTAdapter": {"e2e_avg": 60, "build_avg": 20},
    "Ascend-CANN":        {"e2e_avg": 30, "build_avg": 10},
    "MindSpeed":          {"e2e_avg": 60, "build_avg": 10},
    "_default":           {"e2e_avg": 30, "build_avg": 10},
}

SPREADSHEET_TOKEN = "YgAhsy6eHh1xDgt0BBgcC7yTnph"
SHEET_ID = "23b407"
YELLOW = "#FFF2CC"
GRAY = "#D9D9D9"

DATA_START_ROW = 3


def _parse_val(v):
    if v is None or v == "" or v == "-":
        return None
    try:
        return float(v)
    except (ValueError, TypeError):
        return None


def find_violations(data, col_map):
    """Find yellow cells (exceed thresholds). Returns [(row, col), ...]."""
    # Locate e2e_avg and build_avg columns
    e2e_col = build_col = None
    for col, field in col_map.items():
        if field == "e2e_avg":   e2e_col = col
        if field == "build_avg": build_col = col

    violations = []
    for i, record in enumerate(data):
        product = record.get("product", "")
        row = DATA_START_ROW + i
        rules = THRESHOLDS.get(product, THRESHOLDS["_default"])
        for field, col in [(e2e_col, "e2e_avg"), (build_col, "build_avg")]:
            if col is None:
                continue
            key = {"e2e_avg": "e2e_avg", "build_avg": "build_avg"}[field] if field in ("e2e_avg","build_avg") else field
            val = _parse_val(record.get(key))
            if val is not None and val > rules[key]:
                violations.append((row, col))
    return violations


def find_empty(data, col_map):
    """Find gray cells (value is '-'). Returns [(row, col), ...]."""
    empty = []
    for i, record in enumerate(data):
        row = DATA_START_ROW + i
        for col, field in col_map.items():
            if col in ("A", "B"):
                continue  # skip identity columns
            v = record.get(field)
            if v is None or v == "" or v == "-":
                empty.append((row, col))
    return empty


def build_script(cells, color, label):
    """Generate lark-cli commands to set background color on given cells."""
    if not cells:
        return []
    cmds = [f"# {label} ({len(cells)} cells)"]
    for row, col in sorted(cells, key=lambda x: (x[1], x[0])):
        cell_ref = f"{col}{row}"
        cmds.append(
            f"lark-cli sheets +cells-set"
            f" --spreadsheet-token {SPREADSHEET_TOKEN}"
            f" --sheet-id {SHEET_ID}"
            f" --range {cell_ref}"
            f' --cells \'[[{{"cell_styles":{{"background_color":"{color}"}}}}]]\''
            f" --as user"
        )
    return cmds


def main():
    parser = argparse.ArgumentParser(description="Generate highlight commands")
    parser.add_argument("--data", required=True)
    parser.add_argument("--col-map", required=True, help='JSON {col: field}')
    parser.add_argument("--output", default="/tmp/highlight.sh")
    args = parser.parse_args()

    data = json.load(open(args.data))
    col_map = json.loads(args.col_map)

    # 1. Gray: empty cells
    empty = find_empty(data, col_map)
    # 2. Yellow: threshold violations (overrides gray)
    violations = find_violations(data, col_map)

    # Remove yellow cells from gray list (yellow wins)
    yellow_set = set(violations)
    empty = [e for e in empty if e not in yellow_set]

    cmds = ["#!/bin/bash", "# Auto-generated highlight commands"]
    cmds += build_script(empty, GRAY, "empty cells (no data)")
    cmds += build_script(violations, YELLOW, "threshold violations")

    if not empty and not violations:
        cmds.append("echo 'Nothing to highlight.'")

    with open(args.output, "w") as f:
        f.write("\n".join(cmds) + "\n")

    print(f"Gray (empty):  {len(empty)} cells")
    print(f"Yellow (over): {len(violations)} cells")
    print(f"Script: {args.output}")


if __name__ == "__main__":
    main()
