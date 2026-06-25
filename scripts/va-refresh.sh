#!/usr/bin/env bash
# va-refresh.sh — Pull version availability data for the last 7 days and regenerate the tracking report.
#
# Usage: bash scripts/va-refresh.sh [--no-push]
#   --no-push   Skip git commit & push (dry-run mode)

set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BIN="$PROJECT_DIR/bin/openlibing"
DATA_DIR="/tmp"
DAYS="${VA_DAYS:-7}"
NO_PUSH="${1:-}"

# ── 1. Ensure binary exists ──────────────────────────────────────
if [ ! -x "$BIN" ]; then
    echo "Building openlibing binary..."
    cd "$PROJECT_DIR" && make build
fi

# ── 2. Verify auth ───────────────────────────────────────────────
echo "Checking auth..."
AUTH_TEST=$("$BIN" run version-availability \
    --start-date "$(date -d '1 day ago' +%Y-%m-%d)" \
    --end-date   "$(date -d '1 day ago' +%Y-%m-%d)" \
    --page-size 1 --output json 2>&1) || true

if echo "$AUTH_TEST" | grep -qi "401\|403\|Error: call"; then
    echo ""
    echo "=============================================="
    echo "  AUTH FAILED — Token may have expired."
    echo ""
    echo "  Please update ~/.openlibing/auth.yaml with"
    echo "  fresh browser headers, then re-run:"
    echo ""
    echo "  bash scripts/va-refresh.sh"
    echo "=============================================="
    exit 1
fi
echo "Auth OK"

# ── 3. Pull data for each of the last N days ─────────────────────
echo "Pulling data for last $DAYS days..."
FAILED=""
for i in $(seq $((DAYS - 1)) -1 0); do
    TARGET_DATE=$(date -d "$i days ago" +%Y-%m-%d)
    SHORT_DATE=$(echo "$TARGET_DATE" | tr -d '-')
    CACHE_FILE="$DATA_DIR/va_$SHORT_DATE.json"

    if [ -s "$CACHE_FILE" ]; then
        count=$(python3 -c "import json; print(len(json.load(open('$CACHE_FILE'))))" 2>/dev/null || echo "0")
        echo "  $TARGET_DATE — cached ($count pipelines)"
        continue
    fi

    echo -n "  $TARGET_DATE — fetching... "
    if "$BIN" run version-availability \
        --start-date "$TARGET_DATE" \
        --end-date   "$TARGET_DATE" \
        --page-size 100 --output json > "$CACHE_FILE" 2>/dev/null; then
        count=$(python3 -c "import json; print(len(json.load(open('$CACHE_FILE'))))" 2>/dev/null || echo "?")
        echo "$count pipelines"
    else
        echo "FAILED"
        FAILED="$FAILED $TARGET_DATE"
        rm -f "$CACHE_FILE"
    fi
done

if [ -n "$FAILED" ]; then
    echo "Warning: failed to fetch:$FAILED"
fi

# ── 4. Generate report ───────────────────────────────────────────
echo "Generating report..."
OUTPUT="$PROJECT_DIR/docs/version-availability-tracking.md"
python3 "$SCRIPT_DIR/va-gen-report.py" --days "$DAYS" --data-dir "$DATA_DIR" --output "$OUTPUT"

# ── 5. Commit & push ─────────────────────────────────────────────
if [ "$NO_PUSH" = "--no-push" ]; then
    echo "Skipping git push (--no-push)"
    echo "Report: $OUTPUT"
else
    cd "$PROJECT_DIR"
    git add "$OUTPUT"
    if git diff --cached --quiet; then
        echo "No changes to commit."
    else
        git commit -m "docs: refresh version availability tracking ($(date +%Y-%m-%d))

Co-Authored-By: Claude <noreply@anthropic.com>"
        git push origin master
        echo "Pushed to remote."
    fi
fi

echo "Done."
