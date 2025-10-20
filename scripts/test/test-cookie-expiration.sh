#!/bin/bash
# Cookie 过期预警功能测试脚本

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Cookie 过期预警功能测试${NC}"
echo -e "${BLUE}========================================${NC}\n"

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CONFIG_FILE="$PROJECT_ROOT/config.yaml"
BACKUP_FILE="$PROJECT_ROOT/config.yaml.backup-expiry-test"

# 检查配置文件
if [ ! -f "$CONFIG_FILE" ]; then
    echo -e "${RED}❌ 配置文件不存在: $CONFIG_FILE${NC}"
    echo -e "${YELLOW}请先创建配置文件${NC}"
    exit 1
fi

# 备份配置文件
echo -e "${YELLOW}📋 备份配置文件...${NC}"
cp "$CONFIG_FILE" "$BACKUP_FILE"
echo -e "${GREEN}✓ 已备份到: $BACKUP_FILE${NC}\n"

# 测试函数
run_test() {
    local test_name=$1
    local cookie_updated_at=$2
    local cookie_valid_days=$3
    local expected_result=$4
    
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}测试: $test_name${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo "Cookie 更新时间: $cookie_updated_at"
    echo "有效期: $cookie_valid_days 天"
    echo "预期结果: $expected_result"
    echo ""
    
    # 更新配置文件
    if grep -q "cookie_updated_at:" "$CONFIG_FILE"; then
        sed -i.tmp "s/cookie_updated_at:.*/cookie_updated_at: \"$cookie_updated_at\"/" "$CONFIG_FILE"
    else
        echo "cookie_updated_at: \"$cookie_updated_at\"" >> "$CONFIG_FILE"
    fi
    
    if grep -q "cookie_valid_days:" "$CONFIG_FILE"; then
        sed -i.tmp "s/cookie_valid_days:.*/cookie_valid_days: $cookie_valid_days/" "$CONFIG_FILE"
    else
        echo "cookie_valid_days: $cookie_valid_days" >> "$CONFIG_FILE"
    fi
    
    rm -f "$CONFIG_FILE.tmp"
    
    echo -e "${YELLOW}配置已更新，启动监控程序（3秒后自动停止）...${NC}"
    echo ""
    
    # 启动程序并捕获输出（3秒后停止）
    timeout 3s "$PROJECT_ROOT/lixiang-monitor" 2>&1 | grep -E "(Cookie 状态|Cookie 即将过期|Cookie 已过期)" || echo "未检测到 Cookie 状态输出"
    
    echo ""
    echo -e "${GREEN}✓ 测试完成${NC}\n"
}

# 测试场景1: Cookie 正常（刚更新）
NOW=$(date "+%Y-%m-%d %H:%M:%S")
echo -e "${YELLOW}准备测试场景...${NC}\n"
run_test "场景1: Cookie 正常（刚更新）" \
    "$NOW" \
    "7" \
    "🟢 正常"

# 测试场景2: Cookie 即将过期（2天后过期）
TWO_DAYS_AGO=$(date -v-5d "+%Y-%m-%d %H:%M:%S" 2>/dev/null || date -d "5 days ago" "+%Y-%m-%d %H:%M:%S")
run_test "场景2: Cookie 即将过期（提前48小时预警）" \
    "$TWO_DAYS_AGO" \
    "7" \
    "⚠️ 即将过期"

# 测试场景3: Cookie 已过期
EIGHT_DAYS_AGO=$(date -v-8d "+%Y-%m-%d %H:%M:%S" 2>/dev/null || date -d "8 days ago" "+%Y-%m-%d %H:%M:%S")
run_test "场景3: Cookie 已过期" \
    "$EIGHT_DAYS_AGO" \
    "7" \
    "❌ 已过期"

# 测试场景4: Cookie 有效期长（30天，刚更新）
run_test "场景4: Cookie 有效期30天（刚更新）" \
    "$NOW" \
    "30" \
    "🟢 正常"

# 恢复配置文件
echo -e "${BLUE}========================================${NC}"
echo -e "${YELLOW}📋 恢复配置文件...${NC}"
mv "$BACKUP_FILE" "$CONFIG_FILE"
echo -e "${GREEN}✓ 配置文件已恢复${NC}"

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✅ 所有测试完成！${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${YELLOW}提示:${NC}"
echo "1. Cookie 过期检查每天凌晨1点自动执行"
echo "2. 启动时也会立即检查一次并显示状态"
echo "3. 提前 48 小时会发送过期预警通知"
echo "4. 更新 Cookie 后请务必更新配置文件中的 cookie_updated_at"
echo ""
