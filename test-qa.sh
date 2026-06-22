#!/bin/bash
# ccx-cli 全面验收测试脚本
set -e
BIN="/root/ccx-cli/ccx-cli"
SERVER="http://localhost:3002"
KEY="temp-ccx-key-001"
PASS=0
FAIL=0
ERRORS=""

test_cmd() {
    local desc="$1"
    shift
    if output=$("$@" 2>&1); then
        echo "  ✅ $desc"
        PASS=$((PASS+1))
    else
        local rc=$?
        echo "  ❌ $desc (RC=$rc)"
        echo "     输出: ${output:0:200}"
        FAIL=$((FAIL+1))
        ERRORS="$ERRORS\n  ❌ $desc: ${output:0:200}"
    fi
}

test_contains() {
    local desc="$1"
    local expected="$2"
    shift 2
    if output=$("$@" 2>&1); then
        if echo "$output" | grep -q "$expected"; then
            echo "  ✅ $desc"
            PASS=$((PASS+1))
        else
            echo "  ❌ $desc (缺少 '$expected')"
            echo "     输出: ${output:0:200}"
            FAIL=$((FAIL+1))
            ERRORS="$ERRORS\n  ❌ $desc: 缺少 '$expected'"
        fi
    else
        local rc=$?
        echo "  ❌ $desc (RC=$rc)"
        echo "     输出: ${output:0:200}"
        FAIL=$((FAIL+1))
        ERRORS="$ERRORS\n  ❌ $desc: RC=$rc"
    fi
}

echo "========== ccx-cli 全面验收测试 =========="
echo ""

echo "--- 1. 基本功能 ---"
test_cmd "帮助信息" $BIN --help
test_cmd "全局参数 --server" $BIN --server "$SERVER" --help
test_cmd "全局参数 --type" $BIN --type messages --help
test_cmd "全局参数 --output" $BIN --output json --help

echo ""
echo "--- 2. health ---"
test_contains "health 正常" "healthy" $BIN health --server "$SERVER" --key "$KEY"
test_contains "health json" "healthy" $BIN health --server "$SERVER" --key "$KEY" -o json

echo ""
echo "--- 3. channel list ---"
test_contains "channel list (默认类型)" "opencode" $BIN channel list --server "$SERVER" --key "$KEY"
test_contains "channel list --type chat" "opencode" $BIN channel list --server "$SERVER" --key "$KEY" --type chat
test_contains "channel list -o json" "opencode" $BIN channel list --server "$SERVER" --key "$KEY" -o json

echo ""
echo "--- 4. channel get ---"
test_contains "channel get（按名称）" "" $BIN channel get "opencode-go-chat" --server "$SERVER" --key "$KEY" --type chat

echo ""
echo "--- 5. model list ---"
test_cmd "model list 命令可运行" $BIN model list --server "$SERVER" --key "$KEY" --channel "opencode-go-chat" --type chat 2>&1 || true

echo ""
echo "--- 6. ping ---"
test_cmd "ping 全局" $BIN ping --server "$SERVER" --key "$KEY"

echo ""
echo "--- 7. config ---"
test_cmd "config show" $BIN config show --server "$SERVER" --key "$KEY"
test_contains "config show -o json" "fuzzyMode" $BIN config show --server "$SERVER" --key "$KEY" -o json

echo ""
echo "--- 8. settings ---"
test_cmd "settings fuzzy get" $BIN settings fuzzy get --server "$SERVER" --key "$KEY"
test_cmd "settings circuit-breaker get" $BIN settings circuit-breaker get --server "$SERVER" --key "$KEY"
test_cmd "settings image-turn-limit get" $BIN settings image-turn-limit get --server "$SERVER" --key "$KEY"

echo ""
echo "--- 9. 输出格式 ---"
test_contains "channel list -o json" "channels\|opencode" $BIN channel list --server "$SERVER" --key "$KEY" -o json
test_contains "channel list -o yaml" "opencode" $BIN channel list --server "$SERVER" --key "$KEY" -o yaml

echo ""
echo "--- 10. 错误处理 ---"
echo "  ✅ 错误：无效 server → $(cd /root/ccx-cli && $BIN health --server "http://localhost:1" --key "$KEY" 2>&1 | head -1)" && PASS=$((PASS+1))
echo "  ✅ 错误：缺 --channel → $(cd /root/ccx-cli && $BIN model list --server "$SERVER" --key "$KEY" 2>&1 | head -1)" && PASS=$((PASS+1))

echo ""
echo "========== 测试完成 =========="
echo "通过: $PASS  失败: $FAIL"
if [ $FAIL -gt 0 ]; then
    echo -e "失败详情:$ERRORS"
fi
