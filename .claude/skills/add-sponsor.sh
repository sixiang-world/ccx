#!/usr/bin/env bash
# CCX Add Sponsor Skill
# 自动化添加新赞助商的完整集成流程

set -euo pipefail

SKILL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SKILL_DIR/../.." && pwd)"

echo "🎯 CCX Add Sponsor Integration"
echo "================================"
echo ""

# 交互式收集赞助商信息
echo "📝 请提供赞助商信息："
echo ""

read -p "赞助商 ID (英文，小写，使用连字符，如 unity2): " SPONSOR_ID
read -p "赞助商名称 (英文，如 Unity2.ai): " SPONSOR_NAME_EN
read -p "赞助商名称 (中文，如 Unity2.ai): " SPONSOR_NAME_ZH
read -p "Base URL (如 https://unity2.ai/v1): " BASE_URL
read -p "控制台 URL (如 https://unity2.ai/dashboard): " CONSOLE_URL
read -p "推广链接 (可选，如 https://unity2.ai/register?source=ccx): " PROMO_URL
read -p "Order 值 (如 45，在哪个位置): " ORDER
read -p "前一个赞助商 ID (在哪个赞助商之后，如 runapi): " PREV_SPONSOR
read -p "图标文件路径 (如 ~/Downloads/sponsor.jpg): " ICON_PATH

echo ""
read -p "简短描述 (英文，一句话): " DESC_SHORT_EN
echo ""
read -p "简短描述 (中文，一句话): " DESC_SHORT_ZH
echo ""
read -p "详细描述 (英文，多行，按 Ctrl+D 结束):"
DESC_LONG_EN=$(cat)
echo ""
read -p "详细描述 (中文，多行，按 Ctrl+D 结束):"
DESC_LONG_ZH=$(cat)

# 确认信息
echo ""
echo "📋 请确认以下信息："
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "ID: $SPONSOR_ID"
echo "名称: $SPONSOR_NAME_EN / $SPONSOR_NAME_ZH"
echo "Base URL: $BASE_URL"
echo "控制台: $CONSOLE_URL"
echo "推广链接: ${PROMO_URL:-无}"
echo "Order: $ORDER"
echo "插入位置: 在 $PREV_SPONSOR 之后"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

read -p "确认无误？(y/N): " CONFIRM
if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
    echo "❌ 已取消"
    exit 1
fi

echo ""
echo "🚀 开始集成..."
echo ""

# 保存到临时文件供 Claude 读取
TEMP_INFO="/tmp/ccx-sponsor-info-$$.json"
cat > "$TEMP_INFO" <<EOF
{
  "sponsorId": "$SPONSOR_ID",
  "nameEn": "$SPONSOR_NAME_EN",
  "nameZh": "$SPONSOR_NAME_ZH",
  "baseUrl": "$BASE_URL",
  "consoleUrl": "$CONSOLE_URL",
  "promoUrl": "$PROMO_URL",
  "order": $ORDER,
  "prevSponsor": "$PREV_SPONSOR",
  "iconPath": "$ICON_PATH",
  "descShortEn": "$DESC_SHORT_EN",
  "descShortZh": "$DESC_SHORT_ZH",
  "descLongEn": $(echo "$DESC_LONG_EN" | jq -Rs .),
  "descLongZh": $(echo "$DESC_LONG_ZH" | jq -Rs .)
}
EOF

echo "✅ 赞助商信息已保存到: $TEMP_INFO"
echo ""
echo "📌 下一步："
echo "   请 Claude 读取此文件并执行集成操作。"
echo ""
echo "   提示 Claude："
echo "   请读取 $TEMP_INFO 文件，并按照 .claude/skills/add-sponsor.md 中的流程"
echo "   完成新赞助商 $SPONSOR_NAME_EN 的完整集成。"
echo ""
