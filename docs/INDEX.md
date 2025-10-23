# 📚 文档导航 (Documentation Navigator)

快速找到您需要的文档。

---

## 🚀 快速开始

| 文档 | 描述 | 适合人群 |
|------|------|----------|
| [README.md](../README.md) | 项目概述和快速开始指南 | 所有人 |
| [ARCHITECTURE.md](../ARCHITECTURE.md) | 完整的系统架构文档 | 开发者 |

---

## 📖 用户指南 (guides/)

### 配置指南

| 文档 | 内容 | 难度 |
|------|------|------|
| [SERVERCHAN_SETUP.md](./guides/SERVERCHAN_SETUP.md) | ServerChan（Server酱）配置指南 | ⭐️ 简单 |
| [WECHAT_SETUP.md](./guides/WECHAT_SETUP.md) | 微信群机器人配置指南 | ⭐️ 简单 |
| [BARK_SETUP.md](./guides/BARK_SETUP.md) | 🆕 Bark 推送（iOS/macOS）配置指南 | ⭐️ 简单 |

### 使用指南

| 文档 | 内容 | 难度 |
|------|------|------|
| [COOKIE_QUICK_FIX.md](./guides/COOKIE_QUICK_FIX.md) | 🔥 Cookie 失效快速修复（5分钟） | ⭐️ 简单 |
| [COOKIE_EXPIRATION_DEMO.md](./guides/COOKIE_EXPIRATION_DEMO.md) | 🆕 Cookie 过期预警功能演示 | ⭐️ 简单 |
| [HOT_RELOAD_DEMO.md](./guides/HOT_RELOAD_DEMO.md) | 配置热加载使用示例 | ⭐️⭐️ 中等 |
| [TESTING_GUIDE.md](./guides/TESTING_GUIDE.md) | 完整的功能测试指南 | ⭐️⭐️ 中等 |

---

## 🔬 技术文档 (technical/)

### 核心功能

| 文档 | 内容 | 适合人群 |
|------|------|----------|
| [CONFIG_HOT_RELOAD.md](./technical/CONFIG_HOT_RELOAD.md) | 配置热加载技术实现 | 开发者 |
| [COOKIE_MANAGEMENT.md](./technical/COOKIE_MANAGEMENT.md) | Cookie 管理和失效检测机制 | 开发者 |
| [COOKIE_EXPIRATION_WARNING.md](./technical/COOKIE_EXPIRATION_WARNING.md) | 🆕 Cookie 过期预警功能 | 所有用户 |
| [COOKIE_AUTO_RENEWAL_ANALYSIS.md](./technical/COOKIE_AUTO_RENEWAL_ANALYSIS.md) | Cookie 自动续期可行性分析 | 开发者 |
| [PERIODIC_NOTIFICATION.md](./technical/PERIODIC_NOTIFICATION.md) | 定期通知功能说明 | 开发者 |
| [DELIVERY_OPTIMIZATION.md](./technical/DELIVERY_OPTIMIZATION.md) | 交付时间预测和优化 | 开发者 |

### 实现总结

| 文档 | 内容 | 适合人群 |
|------|------|----------|
| [IMPLEMENTATION_SUMMARY.md](./technical/IMPLEMENTATION_SUMMARY.md) | 热加载功能实现总结 | 开发者 |
| [COOKIE_IMPLEMENTATION_SUMMARY.md](./technical/COOKIE_IMPLEMENTATION_SUMMARY.md) | Cookie 功能实现总结 | 开发者 |
| [COOKIE_EXPIRATION_IMPLEMENTATION.md](./technical/COOKIE_EXPIRATION_IMPLEMENTATION.md) | 🆕 Cookie 过期预警实现总结 | 开发者 |
| [PROJECT_FILES.md](./technical/PROJECT_FILES.md) | 项目文件说明 | 开发者 |

---

## 🧪 测试脚本 (scripts/test/)

| 脚本 | 功能 | 用途 |
|------|------|------|
| `test-notification.sh` | 测试通知功能 | 验证 ServerChan/微信通知 |
| `test-bark.sh` | 🆕 测试 Bark 推送 | 验证 Bark 推送功能 |
| `test-cookie-expiry.sh` | 测试 Cookie 失效检测 | 验证 Cookie 管理功能 |
| `test-cookie-expiration.sh` | 🆕 测试 Cookie 过期预警 | 验证过期预警功能 |
| `test-hot-reload.sh` | 测试配置热加载 | 验证配置自动重载 |
| `test-periodic-notification.sh` | 测试定期通知 | 验证定期通知功能 |
| `test_delivery_calc.go` | 测试交付时间计算 | 单元测试 |

**使用方法**:
```bash
# 进入测试脚本目录
cd scripts/test/

# 运行测试
./test-notification.sh
./test-bark.sh              # 新增
./test-cookie-expiry.sh
./test-cookie-expiration.sh  # 新增
./test-hot-reload.sh
```

---

## 🚀 部署脚本 (scripts/deploy/)

| 脚本 | 功能 | 用途 |
|------|------|------|
| `build.sh` | 构建程序 | 编译 Go 程序 |
| `start.sh` | 启动服务 | 后台启动监控服务 |
| `stop.sh` | 停止服务 | 停止运行中的服务 |
| `status.sh` | 查看状态 | 查看服务运行状态 |

**使用方法**:
```bash
# 进入部署脚本目录
cd scripts/deploy/

# 构建并启动
./build.sh
./start.sh

# 查看状态
./status.sh

# 停止服务
./stop.sh
```

---

## ⚙️ 配置文件 (config/)

| 文件 | 说明 | 用途 |
|------|------|------|
| `config.example.yaml` | 配置文件示例 | 参考模板 |
| `config.enhanced.yaml` | 增强配置示例 | 高级配置参考 |

**注意**: 工作配置文件 `config.yaml` 位于项目根目录。

---

## 🔍 按场景查找文档

### 场景 1: 首次使用
1. 阅读 [README.md](../README.md) - 了解项目
2. 查看 [SERVERCHAN_SETUP.md](./guides/SERVERCHAN_SETUP.md) 或 [WECHAT_SETUP.md](./guides/WECHAT_SETUP.md) - 配置通知
3. 运行 `scripts/test/test-notification.sh` - 测试通知

### 场景 2: Cookie 失效了
1. 查看 [COOKIE_QUICK_FIX.md](./guides/COOKIE_QUICK_FIX.md) - 5分钟快速修复
2. 如需详细了解，查看 [COOKIE_MANAGEMENT.md](./technical/COOKIE_MANAGEMENT.md)

### 场景 3: 修改配置
1. 直接编辑 `config.yaml`
2. 保存文件（自动生效）
3. 查看 [HOT_RELOAD_DEMO.md](./guides/HOT_RELOAD_DEMO.md) 了解更多

### 场景 4: 部署到服务器
1. 查看 [ARCHITECTURE.md](../ARCHITECTURE.md) - 了解部署架构
2. 使用 `scripts/deploy/` 中的部署脚本
3. 参考 [README.md](../README.md) 的"使用方法"章节

### 场景 5: 开发和调试
1. 查看 [ARCHITECTURE.md](../ARCHITECTURE.md) - 了解系统架构
2. 查看 [TESTING_GUIDE.md](./guides/TESTING_GUIDE.md) - 测试指南
3. 查看 `docs/technical/` 中的技术文档

### 场景 6: 理解实现原理
1. [ARCHITECTURE.md](../ARCHITECTURE.md) - 整体架构
2. [CONFIG_HOT_RELOAD.md](./technical/CONFIG_HOT_RELOAD.md) - 热加载原理
3. [COOKIE_MANAGEMENT.md](./technical/COOKIE_MANAGEMENT.md) - Cookie 管理原理
4. [IMPLEMENTATION_SUMMARY.md](./technical/IMPLEMENTATION_SUMMARY.md) - 实现总结

---

## 📝 文档编写规范

### 用户指南 (guides/)
- **目标受众**: 最终用户
- **内容风格**: 简明易懂，步骤清晰
- **包含内容**: 
  - 问题描述
  - 解决步骤
  - 常见问题
  - 示例和截图

### 技术文档 (technical/)
- **目标受众**: 开发者
- **内容风格**: 详细技术说明
- **包含内容**:
  - 技术原理
  - 架构设计
  - 代码实现
  - API 文档
  - 性能指标

---

## 🆘 获取帮助

### 问题分类

| 问题类型 | 查看文档 | 运行脚本 |
|----------|----------|----------|
| Cookie 失效 | [COOKIE_QUICK_FIX.md](./guides/COOKIE_QUICK_FIX.md) | `test-cookie-expiry.sh` |
| 通知不工作 | [SERVERCHAN_SETUP.md](./guides/SERVERCHAN_SETUP.md) / [WECHAT_SETUP.md](./guides/WECHAT_SETUP.md) | `test-notification.sh` |
| 配置不生效 | [HOT_RELOAD_DEMO.md](./guides/HOT_RELOAD_DEMO.md) | `test-hot-reload.sh` |
| 理解架构 | [ARCHITECTURE.md](../ARCHITECTURE.md) | - |

### 支持渠道
1. 📖 查看相关文档
2. 🧪 运行测试脚本诊断
3. 📝 查看程序日志 (`monitor.log`)
4. 💬 提交 GitHub Issue

---

## 📊 文档更新记录

| 日期 | 更新内容 | 版本 |
|------|----------|------|
| 2025-10-20 | 重组项目结构，创建文档导航 | v1.1.0 |
| 2025-10-20 | 添加 Cookie 失效处理功能 | v1.1.0 |
| 2025-09-27 | 添加配置热加载功能 | v1.0.0 |

---

**最后更新**: 2025-10-20  
**维护者**: haifeiWu

---

## 🔗 快速链接

- 🏠 [返回项目首页](../README.md)
- 🏗️ [查看系统架构](../ARCHITECTURE.md)
- 🔥 [Cookie 快速修复](./guides/COOKIE_QUICK_FIX.md)
- 🧪 [测试指南](./guides/TESTING_GUIDE.md)
