#!/bin/bash
# check-upstream.sh — 检查 CCX 上游更新对 ccx-cli 的影响
#
# 用法：./scripts/check-upstream.sh [选项]
#   -v, --verbose    显示详细 diff
#   -n, --newest N   只看最近 N 条 commit (默认 20)
#   --add-remote     添加 upstream-ccx remote（首次运行需要）
#   --setup          首次设置（添加 remote + fetch）
#
# 原理：
#   ccx-cli 是 CCX Admin API 的一个完整客户端。
#   它依赖 ccx 服务端的 REST API 路径和响应数据结构。
#   当上游 CCX 更新了 API 路由、请求/响应类型时，ccx-cli 可能需要同步适配。
#   此脚本通过对比关键路径的差异来判断影响范围。

set -e

# 自动检测上游 remote：优先用已有的 upstream/upstream-ccx/up
for candidate in upstream upstream-ccx up; do
  if git -C "$GIT_ROOT" remote get-url "$candidate" &>/dev/null; then
    REMOTE_NAME="$candidate"
    break
  fi
done
REMOTE_NAME="${REMOTE_NAME:-upstream-ccx}"
UPSTREAM_REPO="https://github.com/BenedictKing/ccx.git"
VERBOSE=false
NEWEST=20
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Git 仓库根目录（确保路径解析从 repo 根开始，不受子目录影响）
GIT_ROOT=$(git -C "$PROJECT_DIR" rev-parse --show-toplevel 2>/dev/null || echo "$PROJECT_DIR")

# 颜色
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m'

# ccx-cli 关心的上游路径 — 这些路径的变动意味着 ccx-cli 可能需要适配
#
# 原因：ccx-cli 的 http 客户端（client.go）中所有 API 路径和请求/响应结构
# 与 CCX 后端（backend-go/）的 handler / config / types 一一对应。
#
# 路径对应关系：
#   ccx-cli internal/client/client.go  →  backend-go/internal/handlers/  (API 路径)
#   ccx-cli internal/models/types.go   →  backend-go/internal/types/     (数据结构)
#   ccx-cli internal/models/types.go   →  backend-go/internal/config/    (FullConfig)
#   ccx-cli internal/client/client.go  →  backend-go/internal/middleware/ (认证头)
#
WATCH_PATHS=(
  "backend-go/internal/handlers"     # API handler：路由、URL 参数、请求/响应结构变化 → 影响 client.go
  "backend-go/internal/config"       # 配置结构体：FullConfig 变更 → 影响 types.go
  "backend-go/internal/types"        # 数据类型：ChannelView / UpstreamConfig 等 → 影响 types.go
  "backend-go/internal/middleware"   # 中间件：认证方式变化 → 影响 client.go
  "backend-go/internal/scheduler"    # 调度器：ccx-cli 的 scheduler-stats 命令与此对应
)

# ============================================================

setup_remote() {
  if git -C "$GIT_ROOT" remote get-url "$REMOTE_NAME" &>/dev/null; then
    echo "✓ remote '$REMOTE_NAME' 已存在"
  else
    echo "→ 添加 remote '$REMOTE_NAME' → $UPSTREAM_REPO"
    git -C "$GIT_ROOT" remote add "$REMOTE_NAME" "$UPSTREAM_REPO"
  fi
  echo "→ 拉取 $REMOTE_NAME 分支信息..."
  git -C "$GIT_ROOT" fetch "$REMOTE_NAME"
  echo "✓ 设置完成，当前上游 HEAD：$(git -C "$GIT_ROOT" rev-parse --short "$REMOTE_NAME/main" 2>/dev/null || echo '（无法获取）')"
}

# ============================================================

echo -e "${CYAN}═══════════════════════════════════════════${NC}"
echo -e "${CYAN}   CCX 上游变更检查 — ccx-cli 兼容性评估${NC}"
echo -e "${CYAN}═══════════════════════════════════════════${NC}"
echo ""

# 参数解析
while [[ $# -gt 0 ]]; do
  case "$1" in
    --add-remote) setup_remote; exit 0 ;;
    --setup) setup_remote; exit 0 ;;
    -v|--verbose) VERBOSE=true; shift ;;
    -n|--newest) NEWEST="$2"; shift 2 ;;
    *) echo "未知参数: $1"; exit 1 ;;
  esac
done

# 检查 remote 是否存在
if ! git -C "$GIT_ROOT" remote get-url "$REMOTE_NAME" &>/dev/null; then
  echo -e "${YELLOW}⚠ 尚未配置 upstream-ccx remote${NC}"
  echo "首次运行请执行："
  echo "  $0 --setup"
  echo ""
  echo "这会添加 remote 并拉取上游代码，不影响你的 ccx-cli 仓库。"
  exit 1
fi

# 拉取最新
echo -e "${CYAN}→ 拉取上游最新分支信息...${NC}"
git -C "$GIT_ROOT" fetch "$REMOTE_NAME" 2>&1
echo ""

# 获取远程 HEAD
UPSTREAM_HEAD=$(git -C "$GIT_ROOT" rev-parse --short "$REMOTE_NAME/main" 2>/dev/null || true)
if [ -z "$UPSTREAM_HEAD" ]; then
  echo -e "${RED}✗ 无法获取 upstream-ccx 的主分支，请检查 remote 配置${NC}"
  exit 1
fi

echo -e "上游最新 commit: ${CYAN}$UPSTREAM_HEAD${NC}"
echo ""

# ==================== 步骤 1：最新提交列表 ====================
echo -e "${CYAN}─── 最近 $NEWEST 条上游提交 ───${NC}"
git -C "$GIT_ROOT" log "$REMOTE_NAME/main" --oneline -"$NEWEST" \
  --pretty=format:"%C(yellow)%h%Creset %s %Cgreen(%ar)%Creset" 2>/dev/null \
  | head -"$NEWEST"
echo ""

# ==================== 步骤 2：检查关注的路径 ====================
echo -e "${CYAN}─── 关键路径变更分析 ───${NC}"
HAS_CRITICAL=false

# 本地已记录的 last-checked 基线
BASELINE_FILE="$GIT_ROOT/.upstream-baseline"
if [ -f "$BASELINE_FILE" ]; then
  BASELINE=$(cat "$BASELINE_FILE")
  echo -e "上次检查的基线: ${CYAN}$BASELINE${NC}"
else
  BASELINE=""
  echo -e "基线文件不存在（${YELLOW}首次检查${NC}）"
fi
echo ""

MERGE_BASE=""
if [ -n "$BASELINE" ]; then
  # 尝试用基线 hash 做 merge-base，如果它还在历史中
  if git -C "$GIT_ROOT" cat-file -e "$BASELINE" 2>/dev/null; then
    MERGE_BASE="$BASELINE"
  fi
fi

for path in "${WATCH_PATHS[@]}"; do
  # 获取该路径在上游的最新变更
  LAST_CHANGE=$(git -C "$GIT_ROOT" log "$REMOTE_NAME/main" --oneline -1 -- "$path" 2>/dev/null || echo "")
  if [ -z "$LAST_CHANGE" ]; then
    echo -e "  ${YELLOW}?${NC} $path（上游无此路径，或暂无变更历史）"
    continue
  fi

  COMMIT_HASH=$(echo "$LAST_CHANGE" | awk '{print $1}')
  COMMIT_MSG=$(echo "$LAST_CHANGE" | cut -d' ' -f2-)

  # 看这个路径是否有新变更（相对于基线）
  HAS_NEW=false
  if [ -n "$MERGE_BASE" ]; then
    NEW_COMMITS=$(git -C "$GIT_ROOT" rev-list "$MERGE_BASE..$REMOTE_NAME/main" -- "$path" 2>/dev/null | wc -l)
    if [ "$NEW_COMMITS" -gt 0 ]; then
      HAS_NEW=true
    fi
  else
    # 基线不存在，取最近 5 条
    NEW_COMMITS=5
    HAS_NEW=true
  fi

  if [ "$HAS_NEW" = true ]; then
    echo -e "  ${RED}⚠${NC} $path"
    echo -e "    最新: ${CYAN}$COMMIT_HASH${NC} $COMMIT_MSG"
    HAS_CRITICAL=true
  else
    echo -e "  ${GREEN}✓${NC} $path（无新变更）"
  fi
done
echo ""

# ==================== 步骤 3：如有关键变更，展示 diff ====================
if [ "$HAS_CRITICAL" = true ]; then
  echo -e "${YELLOW}─── 建议 REVIEW 的变更 ───${NC}"

  for path in "${WATCH_PATHS[@]}"; do
    if [ -n "$MERGE_BASE" ]; then
      COMMITS=$(git -C "$GIT_ROOT" rev-list "$MERGE_BASE..$REMOTE_NAME/main" -- "$path" 2>/dev/null)
    else
      COMMITS=$(git -C "$GIT_ROOT" rev-list "$REMOTE_NAME/main" --oneline -5 -- "$path" | awk '{print $1}' 2>/dev/null)
    fi

    if [ -n "$COMMITS" ]; then
      echo ""
      echo -e "${CYAN}$path${NC}"
      echo "  ~~~~~ 相关 commit ~~~~~"
      for c in $COMMITS; do
        echo -e "    ${YELLOW}$(git log --oneline -1 "$c" 2>/dev/null)${NC}"
      done

      if [ "$VERBOSE" = true ]; then
        # 显示这个路径的 diff 摘要
        RANGE="${MERGE_BASE:-$(git rev-list --max-parents=0 "$REMOTE_NAME/main" 2>/dev/null)}..$REMOTE_NAME/main"
        DIFF_STAT=$(git -C "$GIT_ROOT" diff "$RANGE" -- "$path" --stat 2>/dev/null | tail -1)
        if [ -n "$DIFF_STAT" ]; then
          echo "    差异统计: $DIFF_STAT"
        fi
      fi
    fi
  done
  echo ""

  # 给出判断建议
  echo -e "${RED}═══ 结论：建议人工审核 ═══${NC}"
  echo ""
  echo "上游 CCX 在关键路径上有变更，需检查："
  echo "  1. API 路由是否变化 → 更新 client.go"
  echo "  2. 请求/响应结构是否变化 → 更新 types.go"
  echo "  3. 是否有新增功能需要 ccx-cli 支持"
  echo ""
  echo -e "查看详细日志："
  echo -e "  ${CYAN}git fetch upstream-ccx${NC}"
  echo -e "  ${CYAN}git log upstream-ccx/main --oneline -40${NC}"
  echo -e "  ${CYAN}git diff ${BASELINE:-<old_hash>}..upstream-ccx/main -- internal/router/ internal/api/${NC}"

  # 更新基线
  echo ""
  echo -e "看完后运行："
  echo -e "  ${CYAN}echo $(git -C "$GIT_ROOT" rev-parse "$REMOTE_NAME/main") > .upstream-baseline${NC}"
else
  echo -e "${GREEN}═══ 结论：ccx-cli 无需适配 ═══${NC}"
  echo "上游最近 $NEWEST 条 commit 没有涉及 ccx-cli 关注的关键路径。"

  # 自动更新基线
  git -C "$GIT_ROOT" rev-parse "$REMOTE_NAME/main" > "$BASELINE_FILE"
  echo "✓ 基线已更新"
fi

echo ""
echo "──────────────────────────────────"
echo -e "运行 ${CYAN}./scripts/check-upstream.sh -v${NC} 查看详细 diff"
