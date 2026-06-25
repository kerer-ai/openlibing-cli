#!/usr/bin/env python3
"""
Version availability tracking report generator.

Reads daily JSON files from /tmp/va_YYYYMMDD.json, filters target projects,
and generates a rolling-7-day markdown tracking report.

Usage:
  python3 scripts/va-gen-report.py [--days 7] [--output docs/version-availability-tracking.md]
"""

import json, datetime, sys, argparse
from collections import defaultdict, OrderedDict

# ── Project metadata ──────────────────────────────────────────────
PROJECT_META = {
    2:  ("openUBMC", "openUBMC"),
    3:  ("openLiBing", "openLiBing"),
    4:  ("FrameworkPTAdapter", "Ascend"),
    300030: ("MindSpore", "MindSpore"),
    300036: ("MindIE", "Ascend"),
    300037: ("MindStudio", "Ascend"),
    300038: ("MindCluster", "Ascend"),
    300057: ("Ascend-CANN", "CANN"),
    300059: ("boostkit", "Kunpeng"),
    300073: ("HPCKit", "Kunpeng"),
    300077: ("Triton", "Ascend"),
    300088: ("LQ_FrameworkPTAdapter", "LingQu"),
    300091: ("MindSpeed", "Ascend"),
}

KEEP_PROJECTS = {4, 300030, 300036, 300037, 300038, 300057, 300091}
BASE_URL = "https://www.openlibing.com/apps/nightlyPipelineDashboard?projectId="

# ── Helpers ───────────────────────────────────────────────────────
def proj_name(pid):
    m = PROJECT_META.get(pid)
    return m[0] if m else f"p{pid}"

def product_name(pid):
    m = PROJECT_META.get(pid)
    return m[1] if m else "?"

def date_label(d):
    dt = datetime.date.fromisoformat(d)
    wd = ["一","二","三","四","五","六","日"][dt.weekday()]
    return f"{d[5:]}({wd})"

def fmt_va(v):
    if v is None: return "-"
    if v == 100: return f"**{v}%**"
    return f"<span style=\"color:red\">{v}%</span>"

def trend_icon(vals):
    va_list = [v for v in vals if v is not None]
    if not va_list: return "⚪ 无数据"
    if all(v == 100 for v in va_list): return "🟢 持续100%"
    if all(v == 0 for v in va_list): return "🔴 持续0%"
    if len(va_list) >= 2:
        first100 = next((i for i, v in enumerate(vals) if v == 100), None)
        if first100 is not None:
            if any(v == 0 for v in vals[:first100] if v is not None):
                return "📈 新达标"
        last100 = None
        for i in range(len(vals)-1, -1, -1):
            if vals[i] == 100: last100 = i; break
        if last100 is not None and last100 < len(vals)-1:
            if any(v == 0 for v in vals[last100+1:] if v is not None):
                return "📉 已丢失"
    has_100 = any(v == 100 for v in va_list)
    return "🟡 波动(有100%)" if has_100 else "🟠 波动"

# ── Main ──────────────────────────────────────────────────────────
def main():
    parser = argparse.ArgumentParser(description="Generate version availability tracking report")
    parser.add_argument("--days", type=int, default=7, help="Number of recent days to include (default: 7)")
    parser.add_argument("--output", type=str, default=None, help="Output file path (default: stdout)")
    parser.add_argument("--data-dir", type=str, default="/tmp", help="Directory for cached JSON data files")
    args = parser.parse_args()

    # Compute date range: last N days
    today = datetime.date.today()
    dates = [(today - datetime.timedelta(days=i)).isoformat() for i in range(args.days - 1, -1, -1)]

    # Load data
    all_data = {}
    missing = []
    for d in dates:
        path = f"{args.data_dir}/va_{d.replace('-','')}.json"
        try:
            records = json.load(open(path))
            records = [r for r in records if r["project_id"] in KEEP_PROJECTS]
            all_data[d] = records
        except (FileNotFoundError, json.JSONDecodeError):
            all_data[d] = []
            missing.append(d)

    if missing:
        print(f"Warning: missing data for: {', '.join(missing)}", file=sys.stderr)

    # Build pipeline index: (project_id, pipeline_name) → {date: version_avail}
    pipeline_map = OrderedDict()
    for d in dates:
        for r in all_data[d]:
            key = (r["project_id"], r["pipeline_name"])
            if key not in pipeline_map:
                pipeline_map[key] = {"project_id": r["project_id"], "name": r["pipeline_name"], "dates": {}}
            pipeline_map[key]["dates"][d] = r["version_avail"]

    def sort_key(item):
        (pid, name), info = item
        return (proj_name(pid).lower(), name)

    sorted_pipelines = sorted(pipeline_map.items(), key=sort_key)

    # ── Build markdown ──
    lines = []
    lines.append("# 版本可用度日跟踪报告")
    lines.append("")
    date_cols = " | ".join(date_label(d) for d in dates)
    lines.append(f"> 统计周期: {date_label(dates[0])} ~ {date_label(dates[-1])} | 共 {len(pipeline_map)} 条流水线")
    lines.append("")

    # ── Project summary ──
    proj_stats = defaultdict(lambda: {"total": 0, "daily100": [0]*len(dates)})
    for (pid, name), info in sorted_pipelines:
        proj_stats[pid]["total"] += 1
        for di, d in enumerate(dates):
            if info["dates"].get(d) == 100:
                proj_stats[pid]["daily100"][di] += 1

    lines.append("## 按项目汇总")
    lines.append("")
    lines.append(f"| 项目 | 所属产品 | 流水线数 | {date_cols} | 持续达标 | 趋势 |")
    lines.append(f"|------|----------|----------|{'---|' * len(dates)}----------|------|")

    for pid in sorted(proj_stats.keys()):
        s = proj_stats[pid]
        d100 = s["daily100"]
        always = sum(1 for (p, n), info in sorted_pipelines
                     if p == pid and all(info["dates"].get(d) == 100 for d in dates))
        if max(d100) == 0: pt = "🔴 无达标"
        elif d100[-1] > d100[0]: pt = "📈 上升"
        elif d100[-1] < d100[0]: pt = "📉 下降"
        else: pt = "➡️ 持平"
        pn = proj_name(pid)
        prod = product_name(pid)
        link = f"[{pn}]({BASE_URL}{pid})"
        cells = " | ".join(str(x) for x in d100)
        lines.append(f"| {link} | {prod} | {s['total']} | {cells} | {always} | {pt} |")

    lines.append("")

    # ── Trend table ──
    lines.append("## 版本可用度变化趋势")
    lines.append("")
    lines.append(f"| 项目 | 流水线 | 所属产品 | {date_cols} | 趋势 |")
    lines.append(f"|------|--------|----------|{'---|' * len(dates)}------|")

    for (pid, name), info in sorted_pipelines:
        vals = [info["dates"].get(d) for d in dates]
        t = trend_icon(vals)
        pn = proj_name(pid)
        prod = product_name(pid)
        display_name = name
        same_name = sum(1 for k in pipeline_map if k[1] == name)
        if same_name > 1:
            display_name = f"{name} (p{pid})"
        cells = " | ".join(fmt_va(v) for v in vals)
        lines.append(f"| {pn} | {display_name} | {prod} | {cells} | {t} |")

    lines.append("")
    lines.append("---")
    lines.append("")
    lines.append(f"> 生成时间: {datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')} | 数据源: openlibing-ops API")

    output = "\n".join(lines)

    if args.output:
        with open(args.output, 'w') as f:
            f.write(output)
        print(f"Report written to {args.output} ({len(pipeline_map)} pipelines, {len(dates)} days)", file=sys.stderr)
    else:
        print(output)


if __name__ == "__main__":
    main()
