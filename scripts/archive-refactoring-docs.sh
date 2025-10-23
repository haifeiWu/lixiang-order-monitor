#!/bin/bash

# 归档旧的重构文档
# 创建归档目录
mkdir -p docs/refactoring/archive

# 移动阶段性完成报告
echo "📦 归档 Phase 完成报告..."
mv REFACTORING_PHASE2_COMPLETE.md docs/refactoring/archive/ 2>/dev/null
mv REFACTORING_PHASE3_COMPLETE.md docs/refactoring/archive/ 2>/dev/null
mv REFACTORING_PHASE4_COMPLETE.md docs/refactoring/archive/ 2>/dev/null
mv REFACTORING_PHASE5_COMPLETE.md docs/refactoring/archive/ 2>/dev/null
mv REFACTORING_PHASE6_COMPLETE.md docs/refactoring/archive/ 2>/dev/null

# 移动早期计划和总结
echo "📦 归档早期计划文档..."
mv REFACTORING_PLAN.md docs/refactoring/archive/ 2>/dev/null
mv REFACTORING_SUMMARY.md docs/refactoring/archive/ 2>/dev/null
mv REFACTORING_COMPLETE.md docs/refactoring/archive/ 2>/dev/null

# 在归档目录创建 README
cat > docs/refactoring/archive/README.md << 'EOF'
# 重构文档归档

本目录包含了重构过程中的阶段性文档,已被最终报告汇总。

## 归档文件

- `REFACTORING_PHASE2_COMPLETE.md` - Phase 2: 提取 utils 包
- `REFACTORING_PHASE3_COMPLETE.md` - Phase 3: 提取 cfg 包
- `REFACTORING_PHASE4_COMPLETE.md` - Phase 4: 复杂度优化
- `REFACTORING_PHASE5_COMPLETE.md` - Phase 5: 提取 delivery 和 cookie 包
- `REFACTORING_PHASE6_COMPLETE.md` - Phase 6: 提取 notification 包 + Cookie 集成优化
- `REFACTORING_PLAN.md` - 初始重构计划
- `REFACTORING_SUMMARY.md` - 阶段性总结
- `REFACTORING_COMPLETE.md` - 早期完成报告

## 当前有效文档

请查看项目根目录的以下文档:

- **REFACTORING_FINAL_REPORT.md** - 最终汇总报告 (✅ 推荐阅读)
- **ARCHITECTURE.md** - 系统架构说明
- **CHECKDELIVERYTIME_OPTIMIZATION.md** - 复杂度优化细节

---

归档日期: 2025年10月23日
EOF

echo "✅ 归档完成!"
echo ""
echo "📁 归档位置: docs/refactoring/archive/"
echo "📄 当前有效文档:"
echo "   - REFACTORING_FINAL_REPORT.md (最终汇总报告)"
echo "   - ARCHITECTURE.md"
echo "   - CHECKDELIVERYTIME_OPTIMIZATION.md"
