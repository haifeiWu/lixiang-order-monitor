#!/bin/bash
# 项目结构查看脚本

echo "🏗️  理想汽车订单监控系统 - 项目结构"
echo "======================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}📁 项目目录结构:${NC}"
echo ""

if command -v tree &> /dev/null; then
    tree -L 3 -I 'node_modules|.git|*.log' --dirsfirst -C
else
    echo "建议安装 tree 命令查看目录结构:"
    echo "  macOS: brew install tree"
    echo "  Linux: sudo apt-get install tree"
    echo ""
    echo "当前目录内容:"
    ls -lah
fi

echo ""
echo -e "${CYAN}📊 统计信息:${NC}"
echo ""

# 统计文件数量
echo "📚 文档文件:"
echo "  - 用户指南: $(ls -1 docs/guides/*.md 2>/dev/null | wc -l | tr -d ' ') 个"
echo "  - 技术文档: $(ls -1 docs/technical/*.md 2>/dev/null | wc -l | tr -d ' ') 个"

echo ""
echo "🔧 脚本文件:"
echo "  - 测试脚本: $(ls -1 scripts/test/* 2>/dev/null | wc -l | tr -d ' ') 个"
echo "  - 部署脚本: $(ls -1 scripts/deploy/* 2>/dev/null | wc -l | tr -d ' ') 个"

echo ""
echo "⚙️  配置文件:"
echo "  - 配置模板: $(ls -1 config/*.yaml 2>/dev/null | wc -l | tr -d ' ') 个"

echo ""
echo -e "${CYAN}🔗 快速访问:${NC}"
echo ""
echo "📖 查看文档导航:"
echo "  cat docs/INDEX.md"
echo ""
echo "🏗️  查看架构文档:"
echo "  cat ARCHITECTURE.md"
echo ""
echo "📋 查看重组总结:"
echo "  cat PROJECT_REORGANIZATION.md"
echo ""
echo "🧪 运行测试:"
echo "  cd scripts/test && ./test-notification.sh"
echo ""
echo "🚀 部署服务:"
echo "  cd scripts/deploy && ./build.sh && ./start.sh"
echo ""
