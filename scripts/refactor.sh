#!/bin/bash
# 代码重构脚本 - 模块化 main.go

set -e  # 遇到错误立即退出

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

PROJECT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$PROJECT_ROOT"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}代码重构脚本${NC}"
echo -e "${BLUE}目标: 将 main.go 模块化拆分${NC}"
echo -e "${BLUE}========================================${NC}\n"

# 备份原文件
echo -e "${YELLOW}1. 备份当前 main.go...${NC}"
cp main.go main.go.backup.$(date +%Y%m%d_%H%M%S)
echo -e "${GREEN}✓ 备份完成${NC}\n"

# 提示用户
echo -e "${YELLOW}⚠️  重构说明:${NC}"
echo "此脚本将帮助你进行模块化重构，但需要手动完成以下步骤："
echo ""
echo -e "${BLUE}推荐的重构顺序:${NC}"
echo "1. 创建 notifier 包（通知器）"
echo "2. 创建 utils 包（工具函数）"
echo "3. 提取类型定义到 types.go"
echo "4. 创建 monitor 包（核心逻辑）"
echo "5. 简化 main.go"
echo ""
echo -e "${YELLOW}每一步重构后都要运行测试:${NC}"
echo "  go build -o lixiang-monitor"
echo "  ./scripts/test/test-notification.sh"
echo ""
echo -e "${BLUE}详细重构方案请查看:${NC} REFACTORING_PLAN.md"
echo ""

# 询问用户是否继续
read -p "是否要创建基本的目录结构？(y/n) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo -e "${YELLOW}已取消${NC}"
    exit 0
fi

# 创建目录结构
echo -e "${YELLOW}2. 创建目录结构...${NC}"

mkdir -p notifier
mkdir -p monitor
mkdir -p utils

echo -e "${GREEN}✓ 目录结构创建完成${NC}"
echo ""
echo "  notifier/  - 通知器模块"
echo "  monitor/   - 监控器核心逻辑"
echo "  utils/     - 工具函数"
echo ""

# 创建包文件的模板
echo -e "${YELLOW}3. 创建包文件模板...${NC}"

# notifier/notifier.go
cat > notifier/notifier.go << 'EOF'
package notifier

// Notifier 通知接口
type Notifier interface {
	Send(title, content string) error
}
EOF

echo -e "${GREEN}✓ 已创建 notifier/notifier.go${NC}"

# monitor/monitor.go
cat > monitor/monitor.go << 'EOF'
package monitor

import (
	"github.com/robfig/cron/v3"
)

// TODO: 从 main.go 迁移 Monitor 结构和相关方法到这里
EOF

echo -e "${GREEN}✓ 已创建 monitor/monitor.go${NC}"

# utils/time.go
cat > utils/time.go << 'EOF'
package utils

import (
	"fmt"
	"time"
)

const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateTimeShort  = "2006-01-02 15:04"
	DateFormat     = "2006-01-02"
)

// ParseLockOrderTime 解析锁单时间
func ParseLockOrderTime(timeStr string) (time.Time, error) {
	// 支持多种时间格式
	formats := []string{
		DateTimeFormat,
		"2006/01/02 15:04:05",
		DateTimeShort,
		"2006/01/02 15:04",
		DateFormat,
		"2006/01/02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析时间格式: %s", timeStr)
}
EOF

echo -e "${GREEN}✓ 已创建 utils/time.go${NC}"

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✅ 目录结构创建完成！${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${YELLOW}下一步操作:${NC}"
echo ""
echo -e "${BLUE}方式一: 手动逐步重构（推荐）${NC}"
echo "1. 查看 REFACTORING_PLAN.md 了解详细方案"
echo "2. 按照文档逐步迁移代码"
echo "3. 每次迁移后运行 'go build' 测试"
echo ""
echo -e "${BLUE}方式二: 使用自动重构工具${NC}"
echo "1. 安装 gofmt, goimports"
echo "2. 使用 IDE 的重构功能（如 GoLand, VSCode）"
echo ""
echo -e "${YELLOW}重要提示:${NC}"
echo "• 备份文件已保存，可随时恢复"
echo "• 建议使用 Git 提交每个重构步骤"
echo "• 保持测试通过是重构的前提"
echo ""
echo -e "${GREEN}祝重构顺利！${NC}"
