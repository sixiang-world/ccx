#!/bin/bash
# ccx-cli 安装脚本
# 编译并安装 ccx-cli 到系统路径

set -e

BINDIR="${1:-/usr/local/bin}"
BINARY="ccx-cli"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
VERSION_TAG="$(git -C "$PROJECT_DIR" describe --tags --always 2>/dev/null || echo "dev")"
BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
GIT_COMMIT="$(git -C "$PROJECT_DIR" rev-parse HEAD 2>/dev/null || echo "unknown")"

echo "==> 构建 ccx-cli (${VERSION_TAG})..."
cd "$PROJECT_DIR"

go build -ldflags="\
  -X ccx-cli/internal/version.Version=${VERSION_TAG} \
  -X ccx-cli/internal/version.BuildTime=${BUILD_TIME} \
  -X ccx-cli/internal/version.GitCommit=${GIT_COMMIT}" \
  -o "$BINARY" .

echo "==> 安装到 ${BINDIR}/${BINARY}..."
install -d "$BINDIR"
install -m 755 "$BINARY" "${BINDIR}/${BINARY}"

echo "==> 验证安装..."
"${BINDIR}/${BINARY}" version

echo ""
echo "✓ ccx-cli ${VERSION_TAG} 已安装到 ${BINDIR}/${BINARY}"
echo "  运行 'ccx --help' 查看帮助"
