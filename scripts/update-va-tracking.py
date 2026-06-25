#!/usr/bin/env python3
"""
Update version availability tracking report with a new day's data.

Usage:
  python3 scripts/update-va-tracking.py [date]

  If date is omitted, defaults to yesterday.
  The script pulls data from openlibing API, appends the new column
  to the tracking report, and prints the updated markdown.

Example:
  python3 scripts/update-va-tracking.py 2026-06-25
"""

import json, sys, subprocess
from datetime import date, timedelta

# Determine target date
if len(sys.argv) > 1:
    target_date = sys.argv[1]
else:
    target_date = (date.today() - timedelta(days=1)).isoformat()

print(f"Pulling version availability for {target_date}...", file=sys.stderr)

# Run openlibing CLI
result = subprocess.run(
    ["./bin/openlibing", "run", "version-availability",
     "--start-date", target_date,
     "--end-date", target_date,
     "--page-size", "100",
     "--output", "json"],
    capture_output=True, text=True
)

if result.returncode != 0:
    print(f"Error pulling data: {result.stderr}", file=sys.stderr)
    sys.exit(1)

try:
    data = json.loads(result.stdout)
except json.JSONDecodeError as e:
    print(f"Error parsing JSON: {e}", file=sys.stderr)
    print(f"stdout was: {result.stdout[:500]}", file=sys.stderr)
    sys.exit(1)

# Write raw data
out_file = f"/tmp/va_{target_date.replace('-','')}.json"
with open(out_file, 'w') as f:
    json.dump(data, f, indent=2)

print(f"Saved {len(data)} pipelines to {out_file}", file=sys.stderr)
print(f"Done. Re-run gen_tracking.py to regenerate the full report.", file=sys.stderr)
