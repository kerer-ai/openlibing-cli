#!/usr/bin/env bash
# CI variant of va-refresh.sh — designed for GitHub Actions.
# Reads auth from environment variables instead of ~/.openlibing/auth.yaml.
#
# Required env vars:
#   OPENLIBING_COOKIE      — full browser cookie string
#   OPENLIBING_CSRF_TOKEN  — CSRF token header value

set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BIN="$PROJECT_DIR/bin/openlibing"
DATA_DIR="/tmp"
DAYS="${VA_DAYS:-7}"

echo "Building openlibing binary..."
cd "$PROJECT_DIR" && make build

# Write auth from env vars
mkdir -p ~/.openlibing
cat > ~/.openlibing/auth.yaml <<EOF
openlibing:
  cookie: "${OPENLIBING_COOKIE}"
  csrf_token: "${OPENLIBING_CSRF_TOKEN}"
EOF

# Verify auth
echo "Checking auth..."
YESTERDAY=$(date -d '1 day ago' +%Y-%m-%d)
AUTH_TEST=$("$BIN" run version-availability \
    --start-date "$YESTERDAY" --end-date "$YESTERDAY" \
    --page-size 1 --output json 2>&1) || true

if echo "$AUTH_TEST" | grep -qi "401\|403\|Error: call"; then
    echo "AUTH FAILED — token expired. Update OPENLIBING_COOKIE / OPENLIBING_CSRF_TOKEN secrets."
    exit 1
fi
echo "Auth OK"

# Pull data
echo "Pulling data for last $DAYS days..."
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
        --start-date "$TARGET_DATE" --end-date "$TARGET_DATE" \
        --page-size 100 --output json > "$CACHE_FILE" 2>/dev/null; then
        count=$(python3 -c "import json; print(len(json.load(open('$CACHE_FILE'))))" 2>/dev/null || echo "?")
        echo "$count pipelines"
    else
        echo "FAILED"
        rm -f "$CACHE_FILE"
    fi
done

# Generate report
echo "Generating report..."
OUTPUT="$PROJECT_DIR/docs/version-availability-tracking.md"
python3 "$SCRIPT_DIR/va-gen-report.py" --days "$DAYS" --data-dir "$DATA_DIR" --output "$OUTPUT"

# Configure git for CI
git config user.name "github-actions[bot]"
git config user.email "github-actions[bot]@users.noreply.github.com"

# Commit & push
cd "$PROJECT_DIR"
git add "$OUTPUT"
if git diff --cached --quiet; then
    echo "No changes to commit."
else
    git commit -m "docs: refresh version availability tracking ($(date +%Y-%m-%d))"
    git push
    echo "Pushed to remote."
fi

echo "Done."
