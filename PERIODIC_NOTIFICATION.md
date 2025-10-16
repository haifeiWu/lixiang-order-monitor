# 定期通知功能优化

## 🎯 优化目标
解决交付时间长期不更新时缺乏通知的问题，增加定期状态报告功能。

## ✨ 新增功能

### 1. 定期通知机制
```yaml
# 配置文件新增项
enable_periodic_notify: true            # 启用定期通知
notification_interval_hours: 24         # 通知间隔（小时）
always_notify_when_approaching: true    # 临近交付强制通知
```

### 2. 智能通知策略
- **交付时间变化**: 立即通知（原有功能）
- **定期状态报告**: 按配置间隔发送，即使无变化
- **临近交付提醒**: 可配置是否强制通知
- **多重触发条件**: 支持同时满足多个通知条件

### 3. 增强的通知内容
```
📊 理想汽车订单状态定期报告
订单号: 177971759268550919
官方预计时间: 预计7-9周内交付
通知原因: 定期状态更新

📅 锁单时间: 2024-09-27 13:08 (50天前)
🔮 基于锁单时间预测: 预计 7-9 周后交付 (2024-11-15 至 2024-11-29)
📊 当前状态: 还有 5-19 天 (进度: 78.5%)

📅 通知间隔: 每24小时
⏰ 下次通知时间: 2024-11-17 10:30
```

## 🔧 技术实现

### 新增字段
```go
type Monitor struct {
    // ... 原有字段
    LastNotificationTime   time.Time     // 上次通知时间
    NotificationInterval   time.Duration // 通知间隔
    EnablePeriodicNotify   bool         // 启用定期通知
    AlwaysNotifyWhenApproaching bool    // 临近交付强制通知
}
```

### 核心方法
- `shouldSendPeriodicNotification()`: 检查是否需要定期通知
- `updateLastNotificationTime()`: 更新通知时间记录
- 优化 `checkDeliveryTime()`: 增强通知逻辑

## 📊 通知频率控制

| 场景 | 通知条件 | 频率 |
|------|----------|------|
| 交付时间变化 | 立即触发 | 实时 |
| 定期状态报告 | 配置间隔 | 12-48小时 |
| 临近交付提醒 | 距离交付≤7天 | 每次检查 |
| 初始通知 | 启动时 | 一次 |

## 🎮 使用示例

```bash
# 1. 更新配置文件
cp config.enhanced.yaml config.yaml

# 2. 编辑配置
vim config.yaml

# 3. 启动监控
./lixiang-monitor
```

## 🔍 日志示例
```
2024/11/16 10:30:00 交付时间未发生变化
2024/11/16 10:30:00 发送定期通知，距离上次通知已过 24.0 小时
2024/11/16 10:30:01 成功发送通知，原因: 定期状态更新
```

## ⚙️ 配置建议
- **notification_interval_hours**: 建议24-48小时
- **enable_periodic_notify**: 建议开启
- **always_notify_when_approaching**: 建议开启，确保重要时间不遗漏