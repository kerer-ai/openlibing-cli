#!/usr/bin/env python3
"""
Detect threshold violations and empty cells, output lark-cli commands to
apply yellow/gray background highlighting. Matches data to table rows by
product+repo key (not by JSON array index) so row order is always correct.

Usage:
  python3 scripts/highlight.py --data /tmp/pr-gate-data.json \
    --col-map '{"C":"e2e_p50",...}' --ab-csv /tmp/ab.csv --output /tmp/highlight.sh
"""

import argparse
import json
import sys
from collections import defaultdict


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
    if v is None or v == "" or v == "-" or v == "/":
        return None
    try:
        return float(v)
    except (ValueError, TypeError):
        return None


def load_lookup(data):
    """Build {(product, repo): record} lookup from fetch JSON."""
    return {(r.get("product", ""), r.get("repo", "")): r for r in data}


def parse_ab_csv(csv_text):
    """Parse A/B CSV to [(row_number, product, repo), ...] skipping empties."""
    rows = []
    for line in csv_text.strip().split("\n"):
        content = line.split(" ", 1)[1]  # strip [row=N] prefix
        parts = content.split(",", 1)
        product = parts[0].strip()
        repo = parts[1].strip() if len(parts) > 1 else ""
        if product or repo:
            rows.append((product, repo))
    return rows


def find_violations(ab_rows, lookup, col_map):
    """Find yellow cells. Uses ab_rows for correct row numbers."""
    e2e_col = build_col = None
    for col, field in col_map.items():
        if field == "e2e_avg":   e2e_col = col
        if field == "build_avg": build_col = col

    violations = []
    for i, (product, repo) in enumerate(ab_rows):
        row = DATA_START_ROW + i
        record = lookup.get((product, repo), {})
        rules = THRESHOLDS.get(product, THRESHOLDS["_default"])
        if e2e_col is not None:
            val = _parse_val(record.get("e2e_avg"))
            if val is not None and val >= rules["e2e_avg"]:
                violations.append((row, e2e_col))
        if build_col is not None:
            val = _parse_val(record.get("build_avg"))
            if val is not None and val >= rules["build_avg"]:
                violations.append((row, build_col))
    return violations


def find_empty(ab_rows, lookup, col_map):
    """Find gray cells (no data). Uses ab_rows for correct row numbers."""
    empty = []
    for i, (product, repo) in enumerate(ab_rows):
        row = DATA_START_ROW + i
        record = lookup.get((product, repo), {})
        for col, field in col_map.items():
            if col in ("A", "B"):
                continue
            v = record.get(field)
            if v is None or v == "" or v == "-" or v == "/":
                empty.append((row, col))
    return empty


def build_script(cells, color, label):
    if not cells:
        return []
    cmds = [f"# {label} ({len(cells)} cells)"]
    for row, col in sorted(cells, key=lambda x: (x[1], x[0])):
        cmds.append(
            f"lark-cli sheets +cells-set"
            f" --spreadsheet-token {SPREADSHEET_TOKEN}"
            f" --sheet-id {SHEET_ID}"
            f" --range {col}{row}"
            f' --cells \'[[{{"cell_styles":{{"background_color":"{color}"}}}}]]\''
            f" --as user"
        )
    return cmds


def main():
    parser = argparse.ArgumentParser(description="Generate highlight commands")
    parser.add_argument("--data", required=True, help="Fetch JSON file")
    parser.add_argument("--ab-csv", required=True, help="A/B CSV from table (row-ordered)")
    parser.add_argument("--col-map", required=True, help='JSON {col: field}')
    parser.add_argument("--output", default="/tmp/highlight.sh")
    args = parser.parse_args()

    data = json.load(open(args.data))
    col_map = json.loads(args.col_map)

    # Load A/B order from table
    ab_rows = parse_ab_csv(open(args.ab_csv).read() if args.ab_csv != "-" else sys.stdin.read())
    # Build lookup by (product, repo)
    lookup = load_lookup(data)

    # Detect
    empty = find_empty(ab_rows, lookup, col_map)
    violations = find_violations(ab_rows, lookup, col_map)

    # Yellow overrides gray
    yellow_set = set(violations)
    empty = [e for e in empty if e not in yellow_set]

    # Data range
    data_cols = [c for c in col_map.keys() if c not in ("A", "B")]
    start_col = min(data_cols)
    end_col = max(data_cols)
    data_range = f"{start_col}{DATA_START_ROW}:{end_col}200"

    cmds = ["#!/bin/bash", "# Auto-generated highlight commands"]
    cmds.append(
        f"lark-cli sheets +cells-clear"
        f" --spreadsheet-token {SPREADSHEET_TOKEN}"
        f" --sheet-id {SHEET_ID}"
        f" --range {data_range}"
        f" --scope formats --as user --yes"
    )
    cmds += build_script(empty, GRAY, "empty cells")
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
