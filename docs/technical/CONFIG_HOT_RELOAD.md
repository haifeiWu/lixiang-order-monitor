# 配置热加载功能说明

## 功能概述

项目已实现配置文件（`config.yaml`）的热加载功能。当配置文件被修改并保存后，程序会自动检测变化并重新加载配置，无需手动重启服务。

## 支持热加载的配置项

以下配置项支持实时热加载：

### 1. 订单配置
- `order_id` - 订单ID
- `lixiang_cookies` - 理想汽车请求的 Cookies

### 2. 交付时间配置
- `lock_order_time` - 锁单时间
- `estimate_weeks_min` - 预计交付周数范围（最小值）
- `estimate_weeks_max` - 预计交付周数范围（最大值）

### 3. 通知器配置
- `wechat_webhook_url` - 微信群机器人 Webhook URL
- `serverchan_sendkey` - ServerChan 发送密钥
- `serverchan_baseurl` - ServerChan API 基础 URL

### 4. 通知策略配置
- `enable_periodic_notify` - 是否启用定期通知
- `notification_interval_hours` - 通知间隔（小时）
- `always_notify_when_approaching` - 临近交付时是否总是通知

### 5. 检查间隔配置（需要重启）
- `check_interval` - 检查间隔（cron 表达式）
  - ⚠️ **注意**：此配置项的修改需要手动重启服务才能生效

## 使用方式

### 1. 启动服务
```bash
./lixiang-monitor
```

或使用启动脚本：
```bash
./start.sh
```

### 2. 修改配置
服务运行期间，直接编辑 `config.yaml` 文件：

```bash
vim config.yaml
# 或使用你喜欢的编辑器
```

### 3. 保存配置
保存配置文件后，程序会自动：
1. 检测到配置文件变化
2. 重新读取配置文件
3. 验证并应用新配置
4. 发送配置更新通知（如果配置了通知器）

### 4. 查看日志
在日志中可以看到类似以下的输出：

```
2025/10/17 10:30:15 检测到配置文件变化: config.yaml
2025/10/17 10:30:15 配置已加载，版本: 2
2025/10/17 10:30:15 ✅ 配置已成功热加载
2025/10/17 10:30:15 ServerChan 通知发送成功
```

## 配置更新通知

当配置成功热加载后，如果配置了通知器（微信机器人或 ServerChan），会收到以下内容的通知：

```
⚙️ 监控服务配置已更新

配置版本: 2
更新时间: 2025-10-17 10:30:15

当前配置:
订单ID: 177971759268550919
检查间隔: @every 12h
通知器数量: 1
定期通知: true
通知间隔: 24小时
```

## 线程安全

配置热加载使用了读写锁（`sync.RWMutex`）来保证并发安全：

- **读取配置**：使用读锁（`RLock`），允许多个 goroutine 同时读取
- **更新配置**：使用写锁（`Lock`），确保配置更新时的原子性

这确保了在配置更新过程中，正在进行的订单检查不会读取到不一致的配置数据。

## 配置版本跟踪

每次配置加载成功后，内部的 `configVersion` 计数器会自动递增，用于：
- 跟踪配置变化次数
- 调试和日志记录
- 通知消息中显示版本信息

## 错误处理

### 配置文件格式错误
如果配置文件格式有误（如 YAML 语法错误），会在日志中输出错误信息，配置不会被更新，程序继续使用旧配置运行：

```
2025/10/17 10:30:15 检测到配置文件变化: config.yaml
2025/10/17 10:30:15 重新读取配置文件失败: yaml: line 5: mapping values are not allowed in this context
```

### 配置项验证失败
如果某些配置项无效（如时间格式错误），会使用默认值或保持原值，并在日志中记录警告：

```
2025/10/17 10:30:15 锁单时间解析失败: 无法解析时间格式: 2025-13-40 25:70:00, 保持当前时间
```

### 检查间隔变更
如果修改了 `check_interval`（检查间隔），需要手动重启服务：

```
2025/10/17 10:30:15 重新加载配置失败: 检查间隔已变更，需要重启服务
2025/10/17 10:30:15 ⚠️  检测到检查间隔变更，请手动重启服务以应用新的检查间隔
```

## 最佳实践

### 1. 配置备份
在修改配置前，建议先备份当前配置：
```bash
cp config.yaml config.yaml.backup
```

### 2. 逐项修改
建议一次只修改一个配置项，便于定位问题

### 3. 验证格式
使用 YAML 验证工具检查格式：
```bash
# 使用 yamllint（如果已安装）
yamllint config.yaml
```

### 4. 监控日志
修改配置后，及时查看日志确认是否成功加载：
```bash
tail -f lixiang-monitor.log
```

## 示例场景

### 场景 1：更新 Cookies
当理想汽车的登录 Cookie 过期时：

1. 从浏览器获取新的 Cookie
2. 编辑 `config.yaml`，更新 `lixiang_cookies` 字段
3. 保存文件
4. 程序自动加载新 Cookie，下次检查时使用新的认证信息

### 场景 2：添加通知渠道
需要增加微信通知：

1. 获取微信群机器人 Webhook URL
2. 编辑 `config.yaml`，设置 `wechat_webhook_url`
3. 保存文件
4. 程序自动添加微信通知器，后续通知会同时发送到 ServerChan 和微信

### 场景 3：调整通知频率
觉得通知太频繁：

1. 编辑 `config.yaml`，修改 `notification_interval_hours` 从 24 改为 48
2. 保存文件
3. 程序自动调整通知间隔为 48 小时

## 技术实现

### 核心依赖
- `github.com/spf13/viper` - 配置管理
- `github.com/fsnotify/fsnotify` - 文件系统监听（viper 内置）
- `sync.RWMutex` - 并发安全控制

### 关键代码
```go
// 监听配置文件变化
viper.OnConfigChange(func(e fsnotify.Event) {
    // 重新读取和加载配置
})
viper.WatchConfig()

// 使用读写锁保护配置访问
m.mu.RLock()
config := m.ConfigField
m.mu.RUnlock()
```

## 故障排查

### 问题：配置修改后没有生效
**排查步骤：**
1. 检查日志是否有 "检测到配置文件变化" 消息
2. 确认配置文件保存成功（检查文件修改时间）
3. 验证 YAML 格式是否正确
4. 查看是否有错误日志

### 问题：收不到配置更新通知
**可能原因：**
1. 通知器配置有误
2. 网络问题
3. Webhook URL 或 SendKey 失效

**解决方法：**
- 检查日志中的错误信息
- 验证通知器配置是否正确
- 测试通知渠道是否可用

### 问题：程序频繁重新加载配置
**可能原因：**
- 文件被其他程序频繁修改
- 编辑器的自动保存功能
- 云同步软件的影响

**解决方法：**
- 关闭编辑器的自动保存
- 将配置文件排除在云同步之外
- 使用 `vim` 或 `nano` 等命令行编辑器

## 版本历史

- **v1.0** - 初始版本，支持基本的配置热加载功能
- 配置版本在运行时从 0 开始递增

## 相关文档

- [README.md](./README.md) - 项目主文档
- [WECHAT_SETUP.md](./WECHAT_SETUP.md) - 微信通知配置
- [SERVERCHAN_SETUP.md](./SERVERCHAN_SETUP.md) - ServerChan 配置
