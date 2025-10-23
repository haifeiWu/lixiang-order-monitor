# 理想汽车订单监控系统 - 文档索引

> 完整的文档导航和快速访问指南

## 📚 文档结构

```
docs/
├── INDEX.md                    # 本文件 - 文档索引
├── architecture/               # 架构设计文档
│   ├── ARCHITECTURE.md
│   ├── REFACTORING_FINAL_REPORT.md
│   └── PROJECT_REORGANIZATION.md
├── changelogs/                 # 功能变更日志
│   ├── BARK_FEATURE_CHANGELOG.md
│   ├── COOKIE_EXPIRATION_CHANGELOG.md
│   ├── DATABASE_FEATURE_CHANGELOG.md
│   ├── WEB_BASE_PATH_CHANGELOG.md
│   └── WEB_INTERFACE_CHANGELOG.md
├── guides/                     # 用户指南
│   ├── BARK_SETUP.md
│   ├── COOKIE_EXPIRATION_DEMO.md
│   ├── COOKIE_QUICK_FIX.md
│   ├── HOT_RELOAD_DEMO.md
│   ├── SERVERCHAN_SETUP.md
│   ├── TESTING_GUIDE.md
│   ├── WEB_BASE_PATH.md
│   ├── WEB_INTERFACE.md
│   └── WECHAT_SETUP.md
├── optimization/               # 性能优化文档
│   └── CHECKDELIVERYTIME_OPTIMIZATION.md
├── technical/                  # 技术实现文档
│   ├── CONFIG_HOT_RELOAD.md
│   ├── COOKIE_AUTO_RENEWAL_ANALYSIS.md
│   ├── COOKIE_EXPIRATION_IMPLEMENTATION.md
│   ├── COOKIE_EXPIRATION_WARNING.md
│   ├── COOKIE_IMPLEMENTATION_SUMMARY.md
│   ├── COOKIE_MANAGEMENT.md
│   ├── DATABASE_STORAGE.md
│   ├── DELIVERY_OPTIMIZATION.md
│   ├── IMPLEMENTATION_SUMMARY.md
│   ├── PERIODIC_NOTIFICATION.md
│   └── PROJECT_FILES.md
└── refactoring/                # 重构历史归档
    └── archive/
```

---

## 🚀 快速开始

### 新用户必读
1. [README.md](../README.md) - 项目介绍和快速开始
2. [通知设置指南](guides/WECHAT_SETUP.md) - 配置通知渠道
3. [测试指南](guides/TESTING_GUIDE.md) - 功能测试和验证

### 功能配置
- [Bark 推送配置](guides/BARK_SETUP.md)
- [ServerChan 配置](guides/SERVERCHAN_SETUP.md)
- [企业微信配置](guides/WECHAT_SETUP.md)
- [Web 界面使用](guides/WEB_INTERFACE.md)
- [Web 根路由配置](guides/WEB_BASE_PATH.md)

### 高级功能
- [配置热重载演示](guides/HOT_RELOAD_DEMO.md)
- [Cookie 过期处理](guides/COOKIE_QUICK_FIX.md)
- [Cookie 过期演示](guides/COOKIE_EXPIRATION_DEMO.md)

---

## 📖 文档分类

### 1️⃣ 架构设计 (Architecture)

#### [系统架构文档](architecture/ARCHITECTURE.md)
- 系统整体设计
- 模块划分
- 数据流向
- 技术栈选型

#### [重构总结报告](architecture/REFACTORING_FINAL_REPORT.md)
- 重构目标和原因
- 重构过程记录
- 优化效果评估
- 经验教训总结

#### [项目重组文档](architecture/PROJECT_REORGANIZATION.md)
- 目录结构调整
- 文件组织优化
- 代码模块化改进

---

### 2️⃣ 变更日志 (Changelogs)

#### [Bark 推送功能](changelogs/BARK_FEATURE_CHANGELOG.md)
- **版本**: v1.3.0
- **功能**: iOS Bark 推送通知支持
- **特性**: 自定义音效、图标、分组

#### [Cookie 过期管理](changelogs/COOKIE_EXPIRATION_CHANGELOG.md)
- **版本**: v1.4.0
- **功能**: Cookie 自动过期检测和预警
- **特性**: 7天有效期、3天预警、失效通知

#### [数据库存储](changelogs/DATABASE_FEATURE_CHANGELOG.md)
- **版本**: v1.7.0
- **功能**: SQLite 数据库历史记录
- **特性**: 轻量级、无 CGO、查询工具

#### [Web 可视化界面](changelogs/WEB_INTERFACE_CHANGELOG.md)
- **版本**: v1.8.0
- **功能**: Web 管理界面
- **特性**: 实时监控、历史记录、响应式设计

#### [Web 根路由配置](changelogs/WEB_BASE_PATH_CHANGELOG.md)
- **版本**: v1.9.0
- **功能**: 自定义 Web 服务器根路由
- **特性**: 反向代理支持、多实例部署

---

### 3️⃣ 用户指南 (Guides)

#### 通知配置

| 指南 | 描述 | 适用场景 |
|------|------|---------|
| [Bark 配置](guides/BARK_SETUP.md) | iOS Bark 推送设置 | iOS 用户 |
| [ServerChan 配置](guides/SERVERCHAN_SETUP.md) | Server酱推送设置 | 微信通知 |
| [企业微信配置](guides/WECHAT_SETUP.md) | 企业微信群机器人 | 团队协作 |

#### Web 界面

| 指南 | 描述 | 适用场景 |
|------|------|---------|
| [Web 界面使用](guides/WEB_INTERFACE.md) | Web 可视化界面完整指南 | 查看监控数据 |
| [Web 根路由配置](guides/WEB_BASE_PATH.md) | 自定义根路由设置 | 反向代理部署 |

#### 高级功能

| 指南 | 描述 | 适用场景 |
|------|------|---------|
| [配置热重载](guides/HOT_RELOAD_DEMO.md) | 配置文件热重载演示 | 动态修改配置 |
| [Cookie 快速修复](guides/COOKIE_QUICK_FIX.md) | Cookie 过期快速处理 | 紧急修复 |
| [Cookie 过期演示](guides/COOKIE_EXPIRATION_DEMO.md) | Cookie 管理功能演示 | 了解机制 |
| [测试指南](guides/TESTING_GUIDE.md) | 功能测试方法 | 验证配置 |

---

### 4️⃣ 技术实现 (Technical)

#### 核心功能实现

| 文档 | 主题 | 关键技术 |
|------|------|---------|
| [数据库存储](technical/DATABASE_STORAGE.md) | SQLite 集成 | modernc.org/sqlite |
| [配置热重载](technical/CONFIG_HOT_RELOAD.md) | 配置文件监听 | fsnotify, viper |
| [定期通知](technical/PERIODIC_NOTIFICATION.md) | 通知调度机制 | time.Ticker |
| [交付优化](technical/DELIVERY_OPTIMIZATION.md) | 交付时间计算 | 时间算法 |

#### Cookie 管理专题

| 文档 | 内容 |
|------|------|
| [Cookie 管理](technical/COOKIE_MANAGEMENT.md) | Cookie 生命周期管理 |
| [过期检测实现](technical/COOKIE_EXPIRATION_IMPLEMENTATION.md) | 过期检测技术细节 |
| [过期预警](technical/COOKIE_EXPIRATION_WARNING.md) | 预警机制设计 |
| [自动续期分析](technical/COOKIE_AUTO_RENEWAL_ANALYSIS.md) | 自动续期可行性 |
| [实现总结](technical/COOKIE_IMPLEMENTATION_SUMMARY.md) | Cookie 功能总结 |

#### 项目总结

| 文档 | 内容 |
|------|------|
| [实现总结](technical/IMPLEMENTATION_SUMMARY.md) | 整体技术实现总结 |
| [项目文件说明](technical/PROJECT_FILES.md) | 文件组织和职责 |

---

### 5️⃣ 性能优化 (Optimization)

#### [交付时间检查优化](optimization/CHECKDELIVERYTIME_OPTIMIZATION.md)
- 算法优化
- 性能提升
- 资源使用改进

---

### 6️⃣ 重构历史 (Refactoring Archive)

历史重构文档已归档至 `refactoring/archive/` 目录：

- [重构归档说明](refactoring/archive/README.md)
- 各阶段重构完成报告
- 历史参考文档

---

## 🎯 按场景查找文档

### 场景 1: 初次部署

1. 阅读 [README.md](../README.md)
2. 配置通知：[企业微信](guides/WECHAT_SETUP.md) 或 [Bark](guides/BARK_SETUP.md)
3. 参考 [测试指南](guides/TESTING_GUIDE.md) 验证功能

### 场景 2: Cookie 过期处理

1. [Cookie 快速修复指南](guides/COOKIE_QUICK_FIX.md) - 紧急处理
2. [Cookie 管理文档](technical/COOKIE_MANAGEMENT.md) - 了解机制
3. [Cookie 过期演示](guides/COOKIE_EXPIRATION_DEMO.md) - 功能演示

### 场景 3: 启用 Web 界面

1. [Web 界面使用指南](guides/WEB_INTERFACE.md) - 基础配置
2. [Web 根路由配置](guides/WEB_BASE_PATH.md) - 高级配置（可选）
3. [Web 界面变更日志](changelogs/WEB_INTERFACE_CHANGELOG.md) - 技术细节

### 场景 4: 反向代理部署

1. [Web 根路由配置指南](guides/WEB_BASE_PATH.md)
2. [系统架构文档](architecture/ARCHITECTURE.md)
3. 查看 Nginx 配置示例

### 场景 5: 多实例部署

1. [Web 根路由配置](guides/WEB_BASE_PATH.md) - 路径隔离
2. [系统架构文档](architecture/ARCHITECTURE.md) - 架构理解
3. [配置热重载](guides/HOT_RELOAD_DEMO.md) - 动态管理

### 场景 6: 故障排查

1. [测试指南](guides/TESTING_GUIDE.md) - 功能验证
2. [Cookie 快速修复](guides/COOKIE_QUICK_FIX.md) - Cookie 问题
3. [Web 界面使用指南](guides/WEB_INTERFACE.md) - Web 问题

### 场景 7: 了解技术实现

1. [系统架构文档](architecture/ARCHITECTURE.md) - 整体设计
2. [数据库存储](technical/DATABASE_STORAGE.md) - 存储实现
3. [实现总结](technical/IMPLEMENTATION_SUMMARY.md) - 技术汇总

### 场景 8: 查看功能演进

1. [变更日志目录](changelogs/) - 所有功能更新
2. [重构总结报告](architecture/REFACTORING_FINAL_REPORT.md) - 重构历史
3. [项目重组文档](architecture/PROJECT_REORGANIZATION.md) - 结构变化

---

## 📊 文档统计

| 分类 | 文档数量 | 说明 |
|------|---------|------|
| 架构设计 | 3 | 系统设计和重构 |
| 变更日志 | 5 | 功能更新记录 |
| 用户指南 | 9 | 配置和使用说明 |
| 技术实现 | 11 | 技术细节文档 |
| 性能优化 | 1 | 优化相关 |
| 重构归档 | 7+ | 历史参考 |
| **总计** | **36+** | 完整文档体系 |

---

## 🔍 文档搜索技巧

### 按关键词查找

- **Cookie**: `Cookie 管理`、`Cookie 过期`、`Cookie 预警`
- **Web**: `Web 界面`、`Web 根路由`、`反向代理`
- **通知**: `Bark`、`ServerChan`、`企业微信`、`定期通知`
- **数据库**: `SQLite`、`历史记录`、`数据存储`
- **配置**: `配置文件`、`热重载`、`config.yaml`

### 按技术栈查找

- **Go 语言**: 所有 technical/ 目录文档
- **SQLite**: `DATABASE_STORAGE.md`
- **Viper**: `CONFIG_HOT_RELOAD.md`
- **fsnotify**: `CONFIG_HOT_RELOAD.md`
- **HTTP**: `WEB_INTERFACE.md`、`WEB_BASE_PATH.md`

---

## 📝 文档维护

### 文档更新原则

1. **用户指南** (guides/): 面向用户，注重实用性
2. **技术实现** (technical/): 面向开发者，注重技术细节
3. **变更日志** (changelogs/): 记录功能演进，注重完整性
4. **架构设计** (architecture/): 系统层面，注重全局性

### 新增文档规范

- **用户指南**: 放在 `guides/` 目录
- **技术文档**: 放在 `technical/` 目录
- **功能变更**: 放在 `changelogs/` 目录
- **架构设计**: 放在 `architecture/` 目录

### 文档命名规范

- 使用 `UPPERCASE_WITH_UNDERSCORE.md` 格式
- 名称简洁明了，体现文档主题
- 避免过长的文件名

---

## 🔗 相关链接

- [项目主页](../README.md)
- [GitHub 仓库](https://github.com/haifeiWu/lixiang-order-monitor)
- [问题反馈](https://github.com/haifeiWu/lixiang-order-monitor/issues)

---

## 📮 反馈建议

如果您对文档有任何建议或发现错误，欢迎：

1. 提交 GitHub Issue
2. 发起 Pull Request
3. 联系项目维护者

---

**文档索引版本**: v2.0  
**最后更新**: 2025-10-23  
**维护状态**: ✅ 活跃维护
