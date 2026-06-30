#!/usr/bin/env python3
"""
Fetch version availability data from openlibing-ops API for 7 target products.

Outputs a flat JSON array to stdout, one object per pipeline:
  {"product": "MindIE", "pipeline": "Nightly-CI_MindIE-LLM",
   "versionAvailabilityRate": 85.5, "actualDurationP50Minutes": 45.2, ...}

Usage:
  python3 scripts/fetch.py --start-date 2026-06-12 --end-date 2026-06-30
  python3 scripts/fetch.py --start-date 2026-06-12 --end-date 2026-06-30 --output /tmp/va-data.json
"""

import argparse
import json
import os
import subprocess
import sys


TARGET_PIDS = {
    4:      "FrameworkPTAdapter",
    300030: "MindSpore",
    300036: "MindIE",
    300037: "MindStudio",
    300038: "MindCluster",
    300057: "Ascend-CANN",
    300091: "MindSpeed",
}


def find_openlibing():
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


def fetch_all(start_date, end_date):
    """Paginate through all version-availability-full records."""
    all_records = []
    page = 1
    while True:
        result = subprocess.run([
            find_openlibing(), "run", "version-availability-full",
            "--start-date", start_date,
            "--end-date", end_date,
            "--page-size", "100",
            "--page", str(page),
            "--output", "json",
        ], capture_output=True, text=True)
        if result.returncode != 0:
            print(f"ERROR page {page}: {result.stderr}", file=sys.stderr)
            break
        records = json.loads(result.stdout)
        if not records:
            break
        all_records.extend(records)
        if len(records) < 100:
            break
        page += 1
    return all_records


def fmt_pct(v):
    if v is None: return "-"
    try: return round(float(v), 1)
    except: return v

def fmt_mins(v):
    if v is None: return "-"
    try: return round(float(v), 1)
    except: return v


def main():
    parser = argparse.ArgumentParser(description="Fetch version availability data")
    parser.add_argument("--start-date", required=True)
    parser.add_argument("--end-date", required=True)
    parser.add_argument("--output", "-o", help="Output JSON file (default: stdout)")
    args = parser.parse_args()

    raw = fetch_all(args.start_date, args.end_date)
    print(f"Fetched {len(raw)} raw records", file=sys.stderr)

    output = []
    for r in raw:
        pid = r.get("projectId")
        if pid not in TARGET_PIDS:
            continue
        output.append({
            "product": TARGET_PIDS[pid],
            "pipeline": r.get("pipelineName", "-"),
            "versionAvailabilityRate": fmt_pct(r.get("versionAvailabilityRate")),
            "actualDurationP50Minutes": fmt_mins(r.get("actualDurationP50Minutes")),
            "actualDurationP90Minutes": fmt_mins(r.get("actualDurationP90Minutes")),
            "actualDurationAvgMinutes": fmt_mins(r.get("actualDurationAvgMinutes")),
            "buildTimeP50Minutes": fmt_mins(r.get("buildTimeP50Minutes")),
            "buildTimeP90Minutes": fmt_mins(r.get("buildTimeP90Minutes")),
            "buildTimeAvgMinutes": fmt_mins(r.get("buildTimeAvgMinutes")),
            "testTimeP50Minutes": fmt_mins(r.get("testTimeP50Minutes")),
            "testTimeP90Minutes": fmt_mins(r.get("testTimeP90Minutes")),
            "testTimeAvgMinutes": fmt_mins(r.get("testTimeAvgMinutes")),
            "caseReleaseRateP0": fmt_pct(r.get("caseReleaseRateP0")),
            "casePassRate": fmt_pct(r.get("casePassRate")),
        })

    result = json.dumps(output, ensure_ascii=False, indent=2)
    if args.output:
        with open(args.output, "w") as f:
            f.write(result)
        print(f"Wrote {len(output)} records to {args.output}", file=sys.stderr)
    else:
        print(result)
    print(f"Total: {len(output)} pipelines across {len(TARGET_PIDS)} products", file=sys.stderr)


if __name__ == "__main__":
    main()
