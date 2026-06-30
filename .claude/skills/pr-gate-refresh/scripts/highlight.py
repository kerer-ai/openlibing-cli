#!/usr/bin/env python3
"""
Detect cells that exceed product-specific thresholds and output lark-cli
commands to set yellow background on those cells.

Reads the fetch JSON and column map (from stdin or args), checks each row
against the product threshold rules, and outputs the lark-cli +cells-set
commands to apply yellow highlighting.

Usage (called from SKILL.md workflow):
  python3 scripts/highlight.py --data /tmp/pr-gate-data.json \
    --e2e-avg-col D --build-avg-col H --output /tmp/highlight.sh
"""

import argparse
import json
import sys


# Product threshold rules: product_name -> {field: max_value_minutes}
THRESHOLDS = {
    "MindIE":             {"e2e_avg": 30, "build_avg": 10},
    "FrameworkPTAdapter": {"e2e_avg": 60, "build_avg": 20},
    "Ascend-CANN":        {"e2e_avg": 30, "build_avg": 10},
    "MindSpeed":          {"e2e_avg": 60, "build_avg": 10},
    # default (MindCluster, MindSpore, MindStudio, etc.)
    "_default":           {"e2e_avg": 30, "build_avg": 10},
}

SPREADSHEET_TOKEN = "YgAhsy6eHh1xDgt0BBgcC7yTnph"
SHEET_ID = "23b407"
YELLOW = "#FFF2CC"  # soft yellow


def check_thresholds(data, col_map):
    """
    col_map: {col_letter: json_field} — exactly as produced in Step 2.
    Returns list of (row_number, col_letter) tuples that exceed thresholds.
    """
    # Find which columns hold e2e_avg and build_avg
    e2e_col = None
    build_col = None
    for col, field in col_map.items():
        if field == "e2e_avg":
            e2e_col = col
        elif field == "build_avg":
            build_col = col
    if not e2e_col or not build_col:
        print("ERROR: could not find e2e_avg/build_avg in column map", file=sys.stderr)
        return []

    violations = []
    for i, record in enumerate(data):
        product = record.get("product", "")
        row = i + 3  # data starts at row 3

        rules = THRESHOLDS.get(product, THRESHOLDS["_default"])

        e2e_val = _parse_val(record.get("e2e_avg"))
        build_val = _parse_val(record.get("build_avg"))

        if e2e_val is not None and e2e_val > rules["e2e_avg"]:
            violations.append((row, e2e_col))
        if build_val is not None and build_val > rules["build_avg"]:
            violations.append((row, build_col))

    return violations


def _parse_val(v):
    """Parse a value to float, returning None for non-numeric."""
    if v is None or v == "" or v == "-":
        return None
    try:
        return float(v)
    except (ValueError, TypeError):
        return None


def main():
    parser = argparse.ArgumentParser(description="Detect threshold violations")
    parser.add_argument("--data", required=True, help="Fetch JSON file")
    parser.add_argument("--col-map", required=True, help="JSON mapping col_letter->field_name")
    parser.add_argument("--output", default="/tmp/highlight.sh", help="Output shell script")
    args = parser.parse_args()

    data = json.load(open(args.data))
    col_map = json.loads(args.col_map)

    violations = check_thresholds(data, col_map)
    if not violations:
        print("No threshold violations found.")
        with open(args.output, "w") as f:
            f.write("#!/bin/bash\necho 'No violations to highlight.'\n")
        return

    # Group violations by column for batch highlighting
    from collections import defaultdict
    by_col = defaultdict(list)
    for row, col in violations:
        by_col[col].append(row)

    print(f"Found {len(violations)} threshold violations across {len(by_col)} columns:")

    # Build lark-cli commands — one per cell for simplicity
    # (cells-set needs consistent range; use individual calls for scattered cells)
    cmds = ["#!/bin/bash", "# Auto-generated threshold highlight commands"]
    for col, rows in sorted(by_col.items()):
        for row in sorted(rows):
            cell = f"{col}{row}"
            cmds.append(
                f"lark-cli sheets +cells-set"
                f" --spreadsheet-token {SPREADSHEET_TOKEN}"
                f" --sheet-id {SHEET_ID}"
                f" --range {cell}"
                f' --cells \'[[{{"cell_styles":{{"background_color":"{YELLOW}"}}}}]]\''
                f" --as user"
            )
            print(f"  {cell}")

    with open(args.output, "w") as f:
        f.write("\n".join(cmds) + "\n")

    print(f"\nHighlight script written to {args.output}")


if __name__ == "__main__":
    main()
