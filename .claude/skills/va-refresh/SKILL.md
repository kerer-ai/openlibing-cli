---
name: va-refresh
description: Refresh the version availability daily tracking report. Pulls latest data for the last 7 days, regenerates the markdown report, and commits to git. Triggered by "刷新版本可用度", "刷新跟踪报告", "refresh version availability", "/va-refresh".
---

# va-refresh — 版本可用度日跟踪报告刷新

拉取近 7 天版本可用度数据，滚动刷新 `docs/version-availability-tracking.md`。

## 一次执行

```bash
bash scripts/va-refresh.sh
```

脚本自动完成：构建检查 → auth 验证 → 逐日拉取 → 报告生成 → git 提交推送。

### 参数

| 参数 | 说明 |
|------|------|
| `--no-push` | 跳过 git 提交推送，仅生成本地报告 |
| `VA_DAYS=14` | 环境变量控制天数（默认 7） |

## 分步排查

### Auth 过期

若脚本报 `AUTH FAILED`，需更新 `~/.openlibing/auth.yaml`。从浏览器 DevTools 获取：

1. 打开 https://www.openlibing.com 并登录
2. F12 → Network → 找任意 POST 请求 → Request Headers
3. 提取 `Cookie` 和 `Csrf-Token-Open-Li-Bing` 的值
4. 让用户对着当前对话发送这些 headers，助手会自动更新

### 手动单日拉取

```bash
./bin/openlibing run version-availability \
  --start-date 2026-06-26 --end-date 2026-06-26 \
  --page-size 100 --output json > /tmp/va_20260626.json
```

### 手动生成报告

```bash
python3 scripts/va-gen-report.py --days 7 --output docs/version-availability-tracking.md
```

## 项目过滤

报告只展示 7 个关注项目（定义在 `scripts/va-gen-report.py` 的 `KEEP_PROJECTS`）：
FrameworkPTAdapter, MindSpore, MindIE, MindStudio, MindCluster, Ascend-CANN, MindSpeed。

修改过滤列表编辑 `scripts/va-gen-report.py` 中的 `KEEP_PROJECTS` 变量。
