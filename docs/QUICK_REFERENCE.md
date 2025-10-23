# 📚 文档导航快速参考

## 🎯 快速开始

| 你想要... | 查看这个文档 |
|----------|------------|
| 🚀 快速部署项目 | [README.md](../README.md) |
| 📖 浏览所有文档 | [docs/INDEX.md](INDEX.md) |
| 🏗️ 了解系统架构 | [docs/architecture/ARCHITECTURE.md](architecture/ARCHITECTURE.md) |
| 📝 查看功能更新 | [docs/changelogs/](changelogs/) |
| 📘 学习如何使用 | [docs/guides/](guides/) |
| 🔬 深入技术细节 | [docs/technical/](technical/) |

---

## 📂 文档目录结构

```
docs/
├── 📑 INDEX.md                          # 主索引 - 从这里开始
│
├── 🏗️ architecture/                    # 架构设计
│   ├── README.md                        # 架构文档导航
│   ├── ARCHITECTURE.md                  # 系统架构文档 ⭐
│   ├── REFACTORING_FINAL_REPORT.md     # 重构总结报告
│   └── PROJECT_REORGANIZATION.md       # 项目重组文档
│
├── 📝 changelogs/                      # 变更日志
│   ├── README.md                        # 变更日志导航
│   ├── BARK_FEATURE_CHANGELOG.md       # v1.3.0 - Bark推送
│   ├── COOKIE_EXPIRATION_CHANGELOG.md  # v1.4.0 - Cookie管理
│   ├── DATABASE_FEATURE_CHANGELOG.md   # v1.7.0 - 数据库
│   ├── WEB_INTERFACE_CHANGELOG.md      # v1.8.0 - Web界面
│   └── WEB_BASE_PATH_CHANGELOG.md      # v1.9.0 - Web根路由
│
├── 📘 guides/                          # 用户指南
│   ├── WECHAT_SETUP.md                 # 企业微信配置 🔥
│   ├── BARK_SETUP.md                   # Bark推送配置
│   ├── SERVERCHAN_SETUP.md             # ServerChan配置
│   ├── WEB_INTERFACE.md                # Web界面使用 ⭐
│   ├── WEB_BASE_PATH.md                # Web根路由配置
│   ├── COOKIE_QUICK_FIX.md             # Cookie快速修复 🔥
│   ├── COOKIE_EXPIRATION_DEMO.md       # Cookie演示
│   ├── HOT_RELOAD_DEMO.md              # 热重载演示
│   └── TESTING_GUIDE.md                # 测试指南
│
├── ⚡ optimization/                    # 性能优化
│   └── CHECKDELIVERYTIME_OPTIMIZATION.md
│
├── 🔬 technical/                       # 技术文档
│   ├── DATABASE_STORAGE.md             # 数据库实现 ⭐
│   ├── CONFIG_HOT_RELOAD.md            # 配置热重载
│   ├── PERIODIC_NOTIFICATION.md        # 定期通知
│   ├── DELIVERY_OPTIMIZATION.md        # 交付优化
│   ├── COOKIE_MANAGEMENT.md            # Cookie管理
│   ├── COOKIE_EXPIRATION_IMPLEMENTATION.md
│   ├── COOKIE_EXPIRATION_WARNING.md
│   ├── COOKIE_AUTO_RENEWAL_ANALYSIS.md
│   ├── COOKIE_IMPLEMENTATION_SUMMARY.md
│   ├── IMPLEMENTATION_SUMMARY.md       # 实现总结
│   └── PROJECT_FILES.md                # 文件说明
│
└── 🔄 refactoring/                     # 重构归档
    └── archive/                         # 历史文档
```

---

## 🎯 按场景查找

### 🆕 场景 1: 我是新用户，刚接触这个项目
```
1. README.md                          # 了解项目
2. docs/guides/WECHAT_SETUP.md       # 配置通知
3. docs/guides/TESTING_GUIDE.md      # 测试功能
```

### 🔧 场景 2: Cookie 失效了，需要紧急处理
```
1. docs/guides/COOKIE_QUICK_FIX.md        # 5分钟快速修复 🔥
2. docs/technical/COOKIE_MANAGEMENT.md    # 了解机制
```

### 🌐 场景 3: 想要启用 Web 界面
```
1. docs/guides/WEB_INTERFACE.md           # 使用指南 ⭐
2. docs/guides/WEB_BASE_PATH.md          # 高级配置
3. docs/changelogs/WEB_INTERFACE_CHANGELOG.md  # 技术细节
```

### 🔄 场景 4: 需要用 Nginx 反向代理部署
```
1. docs/guides/WEB_BASE_PATH.md           # 根路由配置
2. docs/architecture/ARCHITECTURE.md      # 架构理解
3. 查看 Web 根路由配置的 Nginx 示例
```

### 📊 场景 5: 想查看历史数据
```
1. docs/guides/WEB_INTERFACE.md           # Web界面查看
2. docs/technical/DATABASE_STORAGE.md     # 数据库说明
3. 使用 scripts/query-db.sh             # 命令行查询
```

### 🔔 场景 6: 配置 iOS 推送通知
```
1. docs/guides/BARK_SETUP.md              # Bark配置指南
2. docs/changelogs/BARK_FEATURE_CHANGELOG.md  # 功能详情
```

### 💻 场景 7: 我是开发者，想了解技术实现
```
1. docs/architecture/ARCHITECTURE.md      # 系统架构 ⭐
2. docs/technical/                        # 技术文档目录
3. docs/changelogs/                       # 变更历史
```

### 📈 场景 8: 想了解项目的发展历程
```
1. docs/changelogs/README.md              # 版本时间线
2. docs/architecture/REFACTORING_FINAL_REPORT.md  # 重构历史
3. docs/architecture/PROJECT_REORGANIZATION.md    # 项目演进
```

---

## 🔍 按关键词查找

| 关键词 | 相关文档 |
|--------|---------|
| **Cookie** | [快速修复](guides/COOKIE_QUICK_FIX.md) · [管理文档](technical/COOKIE_MANAGEMENT.md) · [过期演示](guides/COOKIE_EXPIRATION_DEMO.md) |
| **Web** | [界面指南](guides/WEB_INTERFACE.md) · [根路由](guides/WEB_BASE_PATH.md) · [变更日志](changelogs/WEB_INTERFACE_CHANGELOG.md) |
| **通知** | [微信](guides/WECHAT_SETUP.md) · [Bark](guides/BARK_SETUP.md) · [ServerChan](guides/SERVERCHAN_SETUP.md) · [定期通知](technical/PERIODIC_NOTIFICATION.md) |
| **数据库** | [存储实现](technical/DATABASE_STORAGE.md) · [变更日志](changelogs/DATABASE_FEATURE_CHANGELOG.md) |
| **配置** | [热重载](technical/CONFIG_HOT_RELOAD.md) · [热重载演示](guides/HOT_RELOAD_DEMO.md) |
| **架构** | [系统架构](architecture/ARCHITECTURE.md) · [重构报告](architecture/REFACTORING_FINAL_REPORT.md) |
| **测试** | [测试指南](guides/TESTING_GUIDE.md) |

---

## 📊 文档统计

| 分类 | 数量 | 位置 |
|------|------|------|
| 📑 主文档 | 2 | 根目录 + docs/ |
| 🏗️ 架构设计 | 3 + README | docs/architecture/ |
| 📝 变更日志 | 5 + README | docs/changelogs/ |
| 📘 用户指南 | 9 | docs/guides/ |
| 🔬 技术文档 | 11 | docs/technical/ |
| ⚡ 性能优化 | 1 | docs/optimization/ |
| 🔄 重构归档 | 7+ | docs/refactoring/archive/ |
| **总计** | **38+** | - |

---

## 💡 使用技巧

### 🎯 Tip 1: 从 INDEX.md 开始
[docs/INDEX.md](INDEX.md) 提供了完整的文档导航，包括：
- 完整文档树
- 按分类浏览
- 按场景查找
- 关键词搜索

### 🎯 Tip 2: 查看目录 README
每个子目录都有 README.md，提供：
- 目录说明
- 文档列表
- 使用建议

### 🎯 Tip 3: 利用场景导航
不知道看什么？在 [INDEX.md](INDEX.md) 中找到你的使用场景。

### 🎯 Tip 4: 关注标记
- ⭐ = 重要文档，建议优先阅读
- 🔥 = 常用文档，快速参考
- 📝 = 详细文档，深入了解

---

## 🔗 快速链接

### 最常用的 5 个文档
1. [README.md](../README.md) - 项目主页
2. [Web 界面使用](guides/WEB_INTERFACE.md) - Web 可视化
3. [Cookie 快速修复](guides/COOKIE_QUICK_FIX.md) - 紧急处理
4. [系统架构](architecture/ARCHITECTURE.md) - 架构设计
5. [文档索引](INDEX.md) - 完整导航

### 最新的功能
- [Web 根路由配置 (v1.9.0)](changelogs/WEB_BASE_PATH_CHANGELOG.md)
- [Web 可视化界面 (v1.8.0)](changelogs/WEB_INTERFACE_CHANGELOG.md)
- [数据库存储 (v1.7.0)](changelogs/DATABASE_FEATURE_CHANGELOG.md)

---

## 📮 反馈

文档有问题或建议？
- 📧 提交 GitHub Issue
- 💬 发起 Pull Request
- 📝 联系项目维护者

---

**快速参考版本**: v1.0  
**更新时间**: 2025-10-23  
**维护状态**: ✅ 活跃维护

---

> 💡 **提示**: 将此文档添加到浏览器书签，方便快速查找文档！
