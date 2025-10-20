#!/bin/bash
# 项目文件重组脚本
# 将文档、脚本和配置文件移动到规范的目录结构中

set -e

echo "🔧 开始重组项目结构..."
echo ""

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 移动用户指南文档到 docs/guides/
echo -e "${BLUE}📚 移动用户指南文档...${NC}"
mv -v COOKIE_QUICK_FIX.md docs/guides/ 2>/dev/null || echo "  - COOKIE_QUICK_FIX.md 已移动或不存在"
mv -v WECHAT_SETUP.md docs/guides/ 2>/dev/null || echo "  - WECHAT_SETUP.md 已移动或不存在"
mv -v SERVERCHAN_SETUP.md docs/guides/ 2>/dev/null || echo "  - SERVERCHAN_SETUP.md 已移动或不存在"
mv -v HOT_RELOAD_DEMO.md docs/guides/ 2>/dev/null || echo "  - HOT_RELOAD_DEMO.md 已移动或不存在"
mv -v TESTING_GUIDE.md docs/guides/ 2>/dev/null || echo "  - TESTING_GUIDE.md 已移动或不存在"

echo ""
echo -e "${BLUE}🔬 移动技术文档...${NC}"
mv -v CONFIG_HOT_RELOAD.md docs/technical/ 2>/dev/null || echo "  - CONFIG_HOT_RELOAD.md 已移动或不存在"
mv -v COOKIE_MANAGEMENT.md docs/technical/ 2>/dev/null || echo "  - COOKIE_MANAGEMENT.md 已移动或不存在"
mv -v COOKIE_IMPLEMENTATION_SUMMARY.md docs/technical/ 2>/dev/null || echo "  - COOKIE_IMPLEMENTATION_SUMMARY.md 已移动或不存在"
mv -v IMPLEMENTATION_SUMMARY.md docs/technical/ 2>/dev/null || echo "  - IMPLEMENTATION_SUMMARY.md 已移动或不存在"
mv -v PERIODIC_NOTIFICATION.md docs/technical/ 2>/dev/null || echo "  - PERIODIC_NOTIFICATION.md 已移动或不存在"
mv -v DELIVERY_OPTIMIZATION.md docs/technical/ 2>/dev/null || echo "  - DELIVERY_OPTIMIZATION.md 已移动或不存在"
mv -v PROJECT_FILES.md docs/technical/ 2>/dev/null || echo "  - PROJECT_FILES.md 已移动或不存在"

echo ""
echo -e "${BLUE}🧪 移动测试脚本...${NC}"
mv -v test-cookie-expiry.sh scripts/test/ 2>/dev/null || echo "  - test-cookie-expiry.sh 已移动或不存在"
mv -v test-hot-reload.sh scripts/test/ 2>/dev/null || echo "  - test-hot-reload.sh 已移动或不存在"
mv -v test-notification.sh scripts/test/ 2>/dev/null || echo "  - test-notification.sh 已移动或不存在"
mv -v test-periodic-notification.sh scripts/test/ 2>/dev/null || echo "  - test-periodic-notification.sh 已移动或不存在"
mv -v test_delivery_calc.go scripts/test/ 2>/dev/null || echo "  - test_delivery_calc.go 已移动或不存在"

echo ""
echo -e "${BLUE}🚀 移动部署脚本...${NC}"
mv -v build.sh scripts/deploy/ 2>/dev/null || echo "  - build.sh 已移动或不存在"
mv -v start.sh scripts/deploy/ 2>/dev/null || echo "  - start.sh 已移动或不存在"
mv -v stop.sh scripts/deploy/ 2>/dev/null || echo "  - stop.sh 已移动或不存在"
mv -v status.sh scripts/deploy/ 2>/dev/null || echo "  - status.sh 已移动或不存在"

echo ""
echo -e "${BLUE}⚙️  移动配置文件...${NC}"
mv -v config.example.yaml config/ 2>/dev/null || echo "  - config.example.yaml 已移动或不存在"
mv -v config.enhanced.yaml config/ 2>/dev/null || echo "  - config.enhanced.yaml 已移动或不存在"

# 保留 config.yaml 在根目录（工作配置）
echo "  - config.yaml 保留在根目录（工作配置）"

echo ""
echo -e "${GREEN}✅ 项目结构重组完成！${NC}"
echo ""
echo "新的目录结构："
echo "├── docs/                    # 📚 文档目录"
echo "│   ├── guides/              # 用户指南"
echo "│   └── technical/           # 技术文档"
echo "├── scripts/                 # 🔧 脚本目录"
echo "│   ├── test/                # 测试脚本"
echo "│   └── deploy/              # 部署脚本"
echo "├── config/                  # ⚙️  配置模板"
echo "├── main.go                  # 主程序"
echo "├── config.yaml              # 工作配置"
echo "└── README.md                # 项目说明"
echo ""
