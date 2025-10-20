#!/bin/bash
# Cookie 失效检测测试脚本

set -e

echo "🧪 理想汽车 Cookie 失效检测测试"
echo "=================================="
echo ""

CONFIG_FILE="config.yaml"
BACKUP_FILE="config.yaml.backup"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 备份当前配置
backup_config() {
    if [ -f "$CONFIG_FILE" ]; then
        cp "$CONFIG_FILE" "$BACKUP_FILE"
        echo -e "${GREEN}✅ 已备份当前配置${NC}"
    fi
}

# 恢复配置
restore_config() {
    if [ -f "$BACKUP_FILE" ]; then
        cp "$BACKUP_FILE" "$CONFIG_FILE"
        rm "$BACKUP_FILE"
        echo -e "${GREEN}✅ 已恢复原配置${NC}"
    fi
}

# 设置无效 Cookie
set_invalid_cookie() {
    echo -e "${YELLOW}📝 设置无效的 Cookie...${NC}"
    
    # 使用 sed 替换 Cookie（macOS 兼容）
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' 's/^lixiang_cookies:.*/lixiang_cookies: "INVALID_COOKIE_STRING"/' "$CONFIG_FILE"
    else
        sed -i 's/^lixiang_cookies:.*/lixiang_cookies: "INVALID_COOKIE_STRING"/' "$CONFIG_FILE"
    fi
    
    echo -e "${GREEN}✅ 已设置无效 Cookie${NC}"
}

# 监控日志
monitor_logs() {
    echo -e "${BLUE}📊 监控日志输出（按 Ctrl+C 停止）...${NC}"
    echo ""
    
    # 启动程序并监控日志
    ./lixiang-monitor 2>&1 | grep --line-buffered -E "(Cookie|cookie|失效|验证|Unauthorized|Forbidden)" || true
}

# 测试场景 1: 模拟 Cookie 失效
test_cookie_expiry() {
    echo ""
    echo -e "${YELLOW}📋 测试场景 1: 模拟 Cookie 失效${NC}"
    echo "=================================="
    echo ""
    
    echo "步骤："
    echo "1. 备份当前配置"
    echo "2. 设置无效的 Cookie"
    echo "3. 启动监控程序"
    echo "4. 观察失效检测和告警"
    echo ""
    
    read -p "是否继续? (y/n): " confirm
    if [ "$confirm" != "y" ]; then
        echo "测试已取消"
        return
    fi
    
    backup_config
    set_invalid_cookie
    
    echo ""
    echo -e "${BLUE}🚀 启动监控程序...${NC}"
    echo -e "${YELLOW}预期行为:${NC}"
    echo "  - 第1次失败: 记录 Cookie 验证失败"
    echo "  - 第2次失败: 连续失败计数增加"
    echo "  - 第3次失败: 发送 Cookie 失效告警通知"
    echo ""
    
    # 启动程序（后台运行）
    ./lixiang-monitor > test-cookie-log.txt 2>&1 &
    MONITOR_PID=$!
    
    echo "程序 PID: $MONITOR_PID"
    echo "日志文件: test-cookie-log.txt"
    echo ""
    echo "等待 30 秒以观察行为..."
    
    # 等待并显示日志
    for i in {1..30}; do
        sleep 1
        echo -n "."
        
        # 检查关键日志
        if grep -q "Cookie 验证失败" test-cookie-log.txt 2>/dev/null; then
            echo ""
            echo -e "${RED}⚠️  检测到 Cookie 验证失败${NC}"
            grep "Cookie 验证失败" test-cookie-log.txt | tail -n 1
        fi
        
        if grep -q "Cookie 失效通知已发送" test-cookie-log.txt 2>/dev/null; then
            echo ""
            echo -e "${GREEN}✅ Cookie 失效通知已发送！${NC}"
            break
        fi
    done
    
    echo ""
    echo ""
    
    # 停止程序
    kill $MONITOR_PID 2>/dev/null || true
    wait $MONITOR_PID 2>/dev/null || true
    
    echo -e "${BLUE}📄 完整日志内容:${NC}"
    echo "=================================="
    cat test-cookie-log.txt
    echo "=================================="
    echo ""
    
    restore_config
    
    echo -e "${GREEN}✅ 测试完成${NC}"
}

# 测试场景 2: Cookie 热更新
test_cookie_hot_reload() {
    echo ""
    echo -e "${YELLOW}📋 测试场景 2: Cookie 热更新${NC}"
    echo "=================================="
    echo ""
    
    echo "此测试将演示:"
    echo "1. 启动监控程序（使用无效 Cookie）"
    echo "2. 动态更新配置文件（恢复有效 Cookie）"
    echo "3. 验证配置热加载生效"
    echo ""
    
    read -p "是否继续? (y/n): " confirm
    if [ "$confirm" != "y" ]; then
        echo "测试已取消"
        return
    fi
    
    backup_config
    set_invalid_cookie
    
    echo ""
    echo -e "${BLUE}🚀 启动监控程序（使用无效 Cookie）...${NC}"
    
    # 启动程序
    ./lixiang-monitor > test-reload-log.txt 2>&1 &
    MONITOR_PID=$!
    
    echo "程序 PID: $MONITOR_PID"
    sleep 3
    
    echo ""
    echo -e "${YELLOW}📝 10 秒后将恢复有效 Cookie...${NC}"
    sleep 10
    
    echo -e "${GREEN}恢复有效 Cookie...${NC}"
    restore_config
    
    echo "等待配置重新加载..."
    sleep 5
    
    # 检查日志
    echo ""
    echo -e "${BLUE}📄 日志内容:${NC}"
    echo "=================================="
    tail -n 20 test-reload-log.txt
    echo "=================================="
    echo ""
    
    # 停止程序
    kill $MONITOR_PID 2>/dev/null || true
    wait $MONITOR_PID 2>/dev/null || true
    
    echo -e "${GREEN}✅ 测试完成${NC}"
}

# 显示菜单
show_menu() {
    echo ""
    echo "请选择测试场景:"
    echo ""
    echo "1) 测试 Cookie 失效检测和告警"
    echo "2) 测试 Cookie 热更新"
    echo "3) 查看 Cookie 管理文档"
    echo "4) 手动监控日志（实时）"
    echo "5) 清理测试日志"
    echo "0) 退出"
    echo ""
}

# 主菜单循环
main() {
    while true; do
        show_menu
        read -p "请输入选项 (0-5): " choice
        
        case $choice in
            1)
                test_cookie_expiry
                ;;
            2)
                test_cookie_hot_reload
                ;;
            3)
                if [ -f "COOKIE_MANAGEMENT.md" ]; then
                    echo ""
                    less COOKIE_MANAGEMENT.md
                else
                    echo -e "${RED}❌ 文档文件不存在${NC}"
                fi
                ;;
            4)
                monitor_logs
                ;;
            5)
                rm -f test-cookie-log.txt test-reload-log.txt
                echo -e "${GREEN}✅ 已清理测试日志${NC}"
                ;;
            0)
                echo ""
                echo "退出测试"
                exit 0
                ;;
            *)
                echo -e "${RED}❌ 无效选项${NC}"
                ;;
        esac
    done
}

# 检查程序是否存在
if [ ! -f "./lixiang-monitor" ]; then
    echo -e "${RED}❌ 错误: lixiang-monitor 程序不存在${NC}"
    echo "请先运行: go build -o lixiang-monitor main.go"
    exit 1
fi

# 检查配置文件
if [ ! -f "$CONFIG_FILE" ]; then
    echo -e "${RED}❌ 错误: config.yaml 配置文件不存在${NC}"
    exit 1
fi

# 陷阱处理：确保退出时恢复配置
trap restore_config EXIT INT TERM

# 运行主程序
main
