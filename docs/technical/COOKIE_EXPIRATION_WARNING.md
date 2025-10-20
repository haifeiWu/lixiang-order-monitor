# Cookie 过期管理指南

## 📋 概述

从 v1.x 开始，系统支持 Cookie 过期预警功能，可以提前 48 小时提醒你更新 Cookie，避免因 Cookie 失效导致监控中断。

## 🎯 功能特点

- ⏰ **提前预警**: Cookie 过期前 48 小时自动发送提醒通知
- 📊 **状态显示**: 启动时显示当前 Cookie 状态
- 🔄 **定期检查**: 每天凌晨 1 点自动检查 Cookie 过期状态
- 📝 **详细指导**: 通知中包含完整的 Cookie 更新步骤

## ⚙️ 配置说明

在 `config.yaml` 中添加以下配置：

```yaml
# Cookie 过期管理
cookie_valid_days: 7                     # Cookie 有效期（天）
cookie_updated_at: "2025-10-20 10:00:00" # Cookie 最后更新时间
```

### 配置参数

| 参数 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| `cookie_valid_days` | Cookie 有效期（天） | 7 | `7`, `14`, `30` |
| `cookie_updated_at` | Cookie 最后更新时间 | 首次启动时间 | `"2025-10-20 10:00:00"` |

## 📖 使用流程

### 1. 初始配置

首次使用时，在配置文件中设置 Cookie 和更新时间：

```yaml
lixiang_cookies: "你的Cookie字符串"
cookie_valid_days: 7
cookie_updated_at: "2025-10-20 10:00:00"  # 设置为当前时间
```

### 2. 启动查看状态

启动程序后会显示 Cookie 状态：

```
Cookie 状态: 🟢 正常 (剩余 7 天)
```

状态说明：
- 🟢 **正常**: Cookie 有效期充足
- ⚠️ **即将过期**: 48 小时内将过期
- ❌ **已过期**: Cookie 已经失效

### 3. 接收过期预警

当 Cookie 即将过期时（48小时内），系统会发送通知：

```
⚠️ Cookie 即将过期

您的理想汽车 Cookie 即将在 1 天内过期

为避免影响订单监控，请尽快更新 Cookie。

🔄 更新步骤：
1. 访问理想汽车官网并登录
2. 打开浏览器开发者工具 (F12)
3. 切换到 Network 标签页
4. 刷新页面
5. 找到订单详情接口请求
6. 复制 Request Headers 中的 Cookie
7. 更新 config.yaml 中的两个字段：
   - lixiang_cookies: 新的 Cookie 值
   - cookie_updated_at: 当前时间

⏰ Cookie 有效期: 7 天
📅 更新时间: 2025-10-20 10:00:00
🕐 当前时间: 2025-10-27 09:00:00
⏳ 剩余时间: 1 天 1 小时
```

### 4. 更新 Cookie

收到预警后，按照通知中的步骤更新 Cookie：

#### 步骤 1: 获取新的 Cookie

1. 访问 https://www.lixiang.com/ 并登录
2. 按 F12 打开开发者工具
3. 切换到 **Network** 标签
4. 刷新页面或访问订单详情页
5. 找到 `order-web` 相关的请求
6. 点击请求，在 **Headers** 中找到 `Cookie` 字段
7. 复制完整的 Cookie 值

#### 步骤 2: 更新配置文件

编辑 `config.yaml`，更新两个字段：

```yaml
# 更新 Cookie 值
lixiang_cookies: "新复制的Cookie字符串"

# 更新时间为当前时间
cookie_updated_at: "2025-10-27 10:30:00"
```

**重要**: 必须同时更新 `lixiang_cookies` 和 `cookie_updated_at` 两个字段！

#### 步骤 3: 配置自动生效

由于支持配置热加载，无需重启程序，修改后会自动生效。

查看日志确认更新成功：

```
配置文件已重新加载: config.yaml
Cookie 状态: 🟢 正常 (剩余 7 天)
```

## 🧪 测试功能

可以使用测试脚本验证 Cookie 过期预警功能：

```bash
./scripts/test/test-cookie-expiration.sh
```

测试脚本会模拟不同的 Cookie 状态：
- ✅ Cookie 正常（刚更新）
- ⚠️ Cookie 即将过期（48小时内）
- ❌ Cookie 已过期
- 📅 不同有效期设置

## 💡 最佳实践

### 1. 合理设置有效期

根据实际情况设置 `cookie_valid_days`：

- **保守型**: `cookie_valid_days: 7` - 每周更新一次
- **标准型**: `cookie_valid_days: 14` - 每两周更新一次
- **延长型**: `cookie_valid_days: 30` - 每月更新一次

建议初期使用较短的有效期（7天），观察实际过期时间后再调整。

### 2. 及时响应预警

- 收到预警后尽快更新 Cookie，不要拖到最后一刻
- 建议在预警后 24 小时内完成更新
- 可以设置提醒，在收到预警时及时处理

### 3. 记录更新时间

- 每次更新 Cookie 后务必更新 `cookie_updated_at`
- 可以在手机日历中设置提醒，定期检查
- 建议在笔记中记录每次更新的时间

### 4. 配合失效检测

Cookie 过期预警是**预防性措施**，配合原有的失效检测功能使用效果最佳：

- **过期预警**: 提前 48 小时提醒（预防）
- **失效检测**: 实时检测 Cookie 失效（兜底）

即使忘记更新，失效检测也能在 Cookie 真正失效时及时告警。

## 🔍 故障排查

### 问题 1: 启动时未显示 Cookie 状态

**原因**: 配置文件中未配置 `cookie_updated_at`

**解决**: 添加配置项：

```yaml
cookie_updated_at: "2025-10-20 10:00:00"  # 设置为当前时间
```

### 问题 2: 过期预警未触发

**可能原因**:
1. 距离过期时间还超过 48 小时
2. 已经发送过一次预警（不会重复发送）
3. 通知器未正确配置

**检查方法**:
1. 查看启动日志中的 Cookie 状态
2. 确认 `cookie_updated_at` 配置正确
3. 使用测试脚本验证功能

### 问题 3: 更新后状态未变化

**原因**: 未同时更新 `cookie_updated_at`

**解决**: 确保同时更新两个字段：
```yaml
lixiang_cookies: "新的Cookie"
cookie_updated_at: "2025-10-27 10:30:00"  # 更新为当前时间
```

### 问题 4: 配置更新后未生效

**检查项**:
1. 确认配置文件路径正确
2. 查看程序日志是否有热加载提示
3. 检查配置文件格式是否正确（YAML 格式）
4. 如果仍未生效，尝试重启程序

## 📊 检查周期

系统会在以下时机检查 Cookie 过期状态：

1. **启动时**: 立即检查并显示状态
2. **定时检查**: 每天凌晨 1:00 自动检查
3. **配置更新**: 热加载配置后重新评估状态

## 🎯 相关文档

- [Cookie 管理完整指南](./COOKIE_MANAGEMENT.md)
- [Cookie 快速修复指南](../guides/COOKIE_QUICK_FIX.md)
- [Cookie 自动续期可行性分析](./COOKIE_AUTO_RENEWAL_ANALYSIS.md)
- [配置热加载说明](./CONFIG_HOT_RELOAD.md)

## 📝 更新日志

- **v1.x**: 新增 Cookie 过期预警功能
  - 支持提前 48 小时预警
  - 启动时显示 Cookie 状态
  - 每日自动检查
  - 详细的更新指导

---

> 💡 **提示**: Cookie 过期预警可以有效避免因 Cookie 失效导致的监控中断。建议所有用户都配置此功能！
