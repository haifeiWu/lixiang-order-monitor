# 理想汽车订单交付时间监控工具

这是一个用 Go 语言编写的监控工具，用于监控理想汽车订单的预计交付时间变化，并在变化时通过微信机器人发送通知。

## 功能特性

- 🔍 定时监控理想汽车订单的预计交付时间
- 📱 交付时间变化时自动发送通知
- 📈 **基于锁单时间的交付日期预测** (新功能)
- ⏰ **智能交付提醒** - 临近预计交付时间时主动提醒
- 🔥 **配置热加载** - 修改配置文件后自动生效，无需重启服务
- 🔗 支持多种通知方式：
  - 微信群机器人
  - ServerChan（Server酱）微信推送
- ⚙️ 可配置的检查间隔
- 📊 详细的日志记录
- 🛡️ 错误处理和重试机制
- 🎯 支持同时向多个通道发送通知
- 🔒 线程安全的配置管理

## 安装和配置

### 1. 依赖安装

确保你已经安装了 Go 1.21 或更高版本。

```bash
# 安装依赖
go mod download
```

### 2. 配置文件

编辑 `config.yaml` 文件，配置以下参数：

```yaml
# 订单 ID (必填)
order_id: "177971759268550919"

# 锁单时间配置 (用于计算预计交付时间)
lock_order_time: "2025-09-27 13:08:00"
estimate_weeks_min: 7
estimate_weeks_max: 9

# 检查间隔 (可选，默认每30分钟)
check_interval: "@every 30m"

# 微信群机器人 Webhook URL (可选)
wechat_webhook_url: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_WEBHOOK_KEY"

# ServerChan 配置 (可选)
serverchan_sendkey: "SCT123456T"
serverchan_baseurl: "https://sctapi.ftqq.com/"

# 理想汽车请求的 Cookies (必填)
lixiang_cookies: "你的完整Cookie字符串"
```

**注意**: 至少需要配置一种通知方式（微信群机器人或 ServerChan），否则程序只会记录日志不会发送通知。

### 5. 配置锁单时间

程序支持基于锁单时间的交付预测功能：

- `lock_order_time`: 你的订单锁单时间（格式：YYYY-MM-DD HH:MM:SS）
- `estimate_weeks_min`: 最少预计交付周数（通常为7周）
- `estimate_weeks_max`: 最多预计交付周数（通常为9周）

程序会根据锁单时间自动计算预计交付日期范围，并在临近交付时间时发送提醒。

### 3. 配置通知方式

程序支持两种通知方式，可以单独使用或同时配置：

#### 方式一：微信群机器人（推荐用于团队）

1. 在微信群中添加群机器人
2. 获取 Webhook URL
3. 将 URL 配置到 `config.yaml` 中的 `wechat_webhook_url` 字段
4. 详细步骤请参考 `WECHAT_SETUP.md`

#### 方式二：ServerChan（推荐用于个人）

1. 访问 https://sct.ftqq.com/ 注册账号
2. 获取你的 SendKey
3. 将 SendKey 配置到 `config.yaml` 中的 `serverchan_sendkey` 字段
4. 详细步骤请参考 `SERVERCHAN_SETUP.md`

### 4. 获取理想汽车 Cookies

1. 打开浏览器，登录理想汽车官网
2. 访问你的订单详情页面
3. 打开开发者工具 (F12)
4. 在 Network 标签页中找到订单详情的请求
5. 复制完整的 Cookie 字符串
6. 将 Cookie 配置到 `config.yaml` 中的 `lixiang_cookies` 字段

## 使用方法

### 测试通知功能

在启动监控之前，建议先测试通知功能是否正常：

```bash
# 测试通知配置
./test-notification.sh
```

这个脚本会向所有配置的通知渠道发送测试消息，确保通知功能正常工作。

### 运行程序

```bash
# 直接运行
go run main.go

# 或者编译后运行
go build -o lixiang-monitor
./lixiang-monitor
```

### 后台运行

```bash
# 使用 nohup 在后台运行
nohup ./lixiang-monitor > monitor.log 2>&1 &

# 或者使用 screen
screen -S lixiang-monitor
./lixiang-monitor
# 按 Ctrl+A+D 分离会话
```

## 配置说明

### 检查间隔格式

支持以下格式的检查间隔：

- `@every 30m` - 每30分钟
- `@every 1h` - 每小时
- `@every 30s` - 每30秒
- `"0 */30 * * * *"` - 使用标准 cron 表达式 (带秒)

### 日志输出

程序会输出详细的日志信息，包括：

- 监控启动信息
- 每次检查的结果
- 交付时间变化检测
- 微信通知发送状态
- 错误信息

## 示例输出

```
2024/10/15 10:00:00 main.go:167: 启动监控服务，检查间隔: @every 30m
2024/10/15 10:00:00 main.go:123: 开始检查订单交付时间...
2024/10/15 10:00:01 main.go:139: 当前预计交付时间: 2024-11-15
2024/10/15 10:00:01 main.go:143: 初次检查，记录当前交付时间
2024/10/15 10:00:02 main.go:105: 微信通知发送成功
2024/10/15 10:00:02 main.go:176: 监控服务已启动，等待定时检查...
```

## 故障排除

### 常见问题

1. **Cookie 失效**
   - 症状：API 返回 401 或 403 错误
   - 解决：重新获取 Cookie 并更新配置文件

2. **微信通知发送失败**
   - 症状：日志显示微信通知发送失败
   - 解决：检查 Webhook URL 是否正确，确认机器人是否正常

3. **网络连接问题**
   - 症状：请求超时或连接失败
   - 解决：检查网络连接，可能需要配置代理

### 调试模式

如果遇到问题，可以增加更详细的日志输出，或者减少检查间隔进行调试：

```yaml
# 调试时使用更短的间隔
check_interval: "@every 1m"
```

## 配置热加载

程序支持配置文件的热加载功能，大部分配置项可以在运行时修改并自动生效，无需重启服务：

### 支持热加载的配置项
- ✅ 订单 ID (`order_id`)
- ✅ Cookie (`lixiang_cookies`)
- ✅ 锁单时间相关配置
- ✅ 通知器配置（微信、ServerChan）
- ✅ 通知策略配置

### 需要重启的配置项
- ⚠️ 检查间隔 (`check_interval`) - 修改后需要手动重启服务

### 使用方法
1. 直接编辑 `config.yaml` 文件
2. 保存文件后程序自动检测并加载新配置
3. 查看日志确认配置是否成功加载

**详细说明请参考：** [CONFIG_HOT_RELOAD.md](./CONFIG_HOT_RELOAD.md)

## 注意事项

1. **Cookie 安全性**：请妥善保管你的 Cookie 信息，不要分享给他人
2. **请求频率**：不建议设置过于频繁的检查间隔，以免对服务器造成压力
3. **Cookie 有效期**：Cookie 可能会过期，利用配置热加载功能可以快速更新
4. **网络稳定性**：确保运行环境有稳定的网络连接
5. **配置备份**：修改配置前建议先备份当前配置文件

## 相关文档

- [配置热加载详细说明](./CONFIG_HOT_RELOAD.md) - 配置热加载功能详细文档
- [微信通知配置](./WECHAT_SETUP.md) - 微信群机器人配置指南
- [ServerChan 配置](./SERVERCHAN_SETUP.md) - Server酱配置指南
- [定期通知说明](./PERIODIC_NOTIFICATION.md) - 定期通知功能说明

## 许可证

MIT License