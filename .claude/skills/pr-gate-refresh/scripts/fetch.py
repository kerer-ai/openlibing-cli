#!/usr/bin/env python3
"""
Fetch PR gate metrics from openlibing-ops API for all 7 projects.

Outputs a flat JSON array to stdout, one object per repo:
  {"product": "Ascend-CANN", "repo": "pytorch", "branch": "master",
   "e2e_p90": "80.8", "e2e_avg": "44.1", ...}

Usage:
  python3 scripts/fetch.py --start-date 2026-06-12 --end-date 2026-06-30
  python3 scripts/fetch.py --start-date 2026-06-12 --end-date 2026-06-30 --output /tmp/data.json
"""

import argparse
import json
import os
import subprocess
import sys


PROJECTS = [
    ("Ascend-CANN",        "300057"),
    ("FrameworkPTAdapter", "4"),
    ("MindCluster",        "300038"),
    ("MindIE",             "300036"),
    ("MindSpeed",          "300091"),
    ("MindSpore",          "300030"),
    ("MindStudio",         "300037"),
]


def find_openlibing():
    """Locate the openlibing binary."""
    for d in os.environ.get("PATH", "").split(os.pathsep):
        p = os.path.join(d, "openlibing")
        if os.path.isfile(p):
            return p
    candidates = [
        "./bin/openlibing",
        os.path.expanduser("~/workspace/openlibing/openlibing-cli/bin/openlibing"),
    ]
    for p in candidates:
        if os.path.isfile(p):
            return p
    return "openlibing"


def fmt(v):
    """Format a metric value: None/empty/dash -> '-', otherwise 1 decimal."""
    if v is None or v == "" or v == "-":
        return "-"
    try:
        fv = float(v)
        return "0.0" if fv == 0.0 else f"{fv:.1f}"
    except (ValueError, TypeError):
        return str(v) if v else "-"


def fetch_project(project_id, start_date, end_date):
    """Query one project, return list of enriched records."""
    result = subprocess.run([
        find_openlibing(), "run", "pr-gate-metrics-summary",
        "--project-id", project_id,
        "--start-date", start_date,
        "--end-date", end_date,
        "--page-size", "100",
        "--output", "json",
    ], capture_output=True, text=True)
    if result.returncode != 0:
        print(f"WARNING: project {project_id} failed: {result.stderr.strip()}", file=sys.stderr)
        return []
    return json.loads(result.stdout)


def main():
    parser = argparse.ArgumentParser(
        description="Fetch PR gate metrics from openlibing-ops API"
    )
    parser.add_argument("--start-date", required=True, help="Start date YYYY-MM-DD")
    parser.add_argument("--end-date", required=True, help="End date YYYY-MM-DD")
    parser.add_argument("--output", "-o", help="Output JSON file (default: stdout)")
    args = parser.parse_args()

    all_records = []
    for product, pid in PROJECTS:
        records = fetch_project(pid, args.start_date, args.end_date)
        for r in records:
            # Normalize: flatten SPC output keys with formatted values
            all_records.append({
                "product": product,
                "repo": r.get("repo", "-"),
                "branch": r.get("branch", "-") or "-",
                "e2e_p90": fmt(r.get("e2e_p90")),
                "e2e_avg": fmt(r.get("e2e_avg")),
                "e2e_p50": fmt(r.get("e2e_p50")),
                "build_p50": fmt(r.get("build_p50")),
                "build_p90": fmt(r.get("build_p90")),
                "build_avg": fmt(r.get("build_avg")),
                "build_queue_p90": fmt(r.get("build_queue_p90")),
                "test_p90": fmt(r.get("test_p90")),
                "test_p50": fmt(r.get("test_p50")),
                "test_avg": fmt(r.get("test_avg")),
                "test_queue_p90": fmt(r.get("test_queue_p90")),
                "check_p90": fmt(r.get("check_p90")),
            })

    output = json.dumps(all_records, ensure_ascii=False, indent=2)
    if args.output:
        with open(args.output, "w") as f:
            f.write(output)
        print(f"Wrote {len(all_records)} records to {args.output}", file=sys.stderr)
    else:
        print(output)

    print(f"Total: {len(all_records)} repos across {len(PROJECTS)} projects", file=sys.stderr)


if __name__ == "__main__":
    main()
