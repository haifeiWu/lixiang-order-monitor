# 理想汽车订单交付时间监控工具

这是一个用 Go 语言编写的监控工具，用于监控理想汽车订单的预计交付时间变化，并在变化时通过微信机器人发送通知。

> 📖 **完整架构文档**: 查看 [ARCHITECTURE.md](./ARCHITECTURE.md) 了解详细的系统架构和技术实现

## 功能特性

- 🔍 定时监控理想汽车订单的预计交付时间
- 📱 交付时间变化时自动发送通知
- 📈 **基于锁单时间的交付日期预测**
- ⏰ **智能交付提醒** - 临近预计交付时间时主动提醒
- 🔥 **配置热加载** - 修改配置文件后自动生效，无需重启服务
- 🍪 **Cookie 失效自动检测** - 智能检测并告警 Cookie 失效
- ⏳ **Cookie 过期预警** - 提前 48 小时提醒 Cookie 即将过期
- 💾 **历史数据存储** - 使用 SQLite 数据库持久化保存所有监控记录
- 🌐 **Web 可视化界面** - 提供美观的 Web 管理界面，实时查看监控状态
- � 支持多种通知方式：
  - 微信群机器人
  - ServerChan（Server酱）微信推送
  - Bark 推送（iOS/macOS）
- ⚙️ 可配置的检查间隔
- 📊 详细的日志记录
- 🛡️ 错误处理和重试机制
- 🎯 支持同时向多个通道发送通知
- 🔒 线程安全的配置管理

## 📁 项目结构

```
lixiang-order-monitor/
├── docs/                    # 📚 文档目录
│   ├── guides/             # 用户指南
│   └── technical/          # 技术文档
├── scripts/                # 🔧 脚本目录
│   ├── test/              # 测试脚本
│   └── deploy/            # 部署脚本
├── config/                 # ⚙️ 配置模板
├── main.go                 # 主程序
├── config.yaml            # 工作配置
├── README.md              # 项目说明
└── ARCHITECTURE.md        # 架构文档
```

详细目录结构请查看 [ARCHITECTURE.md](./ARCHITECTURE.md)

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

# Bark 推送配置 (可选)
bark_server_url: "http://your_server:8080/your_key"
bark_sound: "minuet"
bark_group: "lixiang-monitor"

# 理想汽车请求的 Cookies (必填)
lixiang_cookies: "你的完整Cookie字符串"

# Cookie 过期管理 (可选，但强烈建议配置)
cookie_valid_days: 7                     # Cookie 有效期，默认 7 天
cookie_updated_at: "2025-10-20 10:00:00" # Cookie 最后更新时间
```

**注意**: 至少需要配置一种通知方式（微信群机器人、ServerChan 或 Bark），否则程序只会记录日志不会发送通知。

### 5. 配置 Cookie 过期管理（推荐）

为了及时发现 Cookie 过期问题，建议配置 Cookie 过期管理：

- `cookie_valid_days`: Cookie 有效期（天），默认 7 天。根据实际情况调整
- `cookie_updated_at`: Cookie 最后更新时间，每次更新 Cookie 后请务必更新此字段

**系统会在 Cookie 过期前 48 小时自动发送提醒通知**，通知内容包括详细的更新步骤。

### 6. 配置锁单时间

程序支持基于锁单时间的交付预测功能：

- `lock_order_time`: 你的订单锁单时间（格式：YYYY-MM-DD HH:MM:SS）
- `estimate_weeks_min`: 最少预计交付周数（通常为7周）
- `estimate_weeks_max`: 最多预计交付周数（通常为9周）

程序会根据锁单时间自动计算预计交付日期范围，并在临近交付时间时发送提醒。

### 3. 配置通知方式

程序支持三种通知方式，可以单独使用或同时配置：

#### 方式一：微信群机器人（推荐用于团队）

1. 在微信群中添加群机器人
2. 获取 Webhook URL
3. 将 URL 配置到 `config.yaml` 中的 `wechat_webhook_url` 字段
4. 详细步骤请参考 [WECHAT_SETUP.md](./docs/guides/WECHAT_SETUP.md)

#### 方式二：ServerChan（推荐用于个人）

1. 访问 https://sct.ftqq.com/ 注册账号
2. 获取你的 SendKey
3. 将 SendKey 配置到 `config.yaml` 中的 `serverchan_sendkey` 字段
4. 详细步骤请参考 [SERVERCHAN_SETUP.md](./docs/guides/SERVERCHAN_SETUP.md)

#### 方式三：Bark 推送（推荐 iOS/macOS 用户）

1. 在 App Store 下载 Bark App
2. 获取推送 URL（自动生成或自建服务器）
3. 将 URL 配置到 `config.yaml` 中的 `bark_server_url` 字段
4. 详细步骤请参考 [BARK_SETUP.md](./docs/guides/BARK_SETUP.md)

**推荐组合**：
- iOS/Mac 用户：Bark + 微信机器人（双保险）
- 其他用户：ServerChan + 微信机器人

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
./scripts/test/test-notification.sh
```

这个脚本会向所有配置的通知渠道发送测试消息，确保通知功能正常工作。

### 运行程序

**开发环境**:
```bash
# 直接运行
go run main.go

# 或者编译后运行
go build -o lixiang-monitor
./lixiang-monitor
```

**生产环境** (推荐使用部署脚本):
```bash
# 构建程序
./scripts/deploy/build.sh

# 启动服务
./scripts/deploy/start.sh

# 查看状态
./scripts/deploy/status.sh

# 停止服务
./scripts/deploy/stop.sh
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

## 历史数据查询

### 查看监控历史记录

程序会自动将每次检查的结果保存到 SQLite 数据库中，可以使用以下方法查询历史记录：

```bash
# 使用提供的查询脚本
./scripts/query-db.sh
```

输出示例：
```
📊 理想汽车订单监控 - 历史记录
================================

📈 记录统计
--------------------------------
总记录数: 25
订单数: 1
时间变更次数: 3
通知发送次数: 12

📋 最近 10 条记录
--------------------------------
id  order_id  estimate_time      check_time           approaching  changed  notified
--  --------  ----------------   ------------------   -----------  -------  --------
25  550919    预计6-8周内交付    2025-10-23 16:34     否           否       是
```

### 手动查询数据库

如果你熟悉 SQL，可以直接查询数据库：

```bash
# 进入数据库
sqlite3 lixiang-monitor.db

# 查询所有记录
SELECT * FROM delivery_records ORDER BY check_time DESC LIMIT 10;

# 查询时间变更记录
SELECT check_time, previous_estimate, estimate_time 
FROM delivery_records 
WHERE time_changed = 1;

# 退出
.quit
```

详细的数据库说明请参考：[DATABASE_STORAGE.md](./docs/technical/DATABASE_STORAGE.md)

## Web 可视化界面

### 访问界面

启动程序后，打开浏览器访问：

```
http://localhost:8080
```

### 界面功能

- **实时统计**: 查看总检查次数、时间变更次数、通知发送次数
- **最新状态**: 显示当前预计交付时间、锁单时间、临近状态
- **时间变更历史**: 追踪交付时间的历史变化
- **检查记录**: 查看最近的所有检查记录
- **自动刷新**: 每 30 秒自动更新数据

### 配置

在 `config.yaml` 中配置：

```yaml
# Web 管理界面配置
web_enabled: true       # 是否启用 Web 界面
web_port: 8080          # Web 服务器端口
```

### 特点

- 🎨 美观的现代化界面设计
- 📱 响应式布局，支持移动端
- ⚡ 实时数据展示
- 📊 直观的数据可视化
- 🔄 自动刷新机制

详细说明请参考：[WEB_INTERFACE.md](./docs/guides/WEB_INTERFACE.md)

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

**详细说明请参考：** [docs/technical/CONFIG_HOT_RELOAD.md](./docs/technical/CONFIG_HOT_RELOAD.md)

## 注意事项

1. **Cookie 安全性**：请妥善保管你的 Cookie 信息，不要分享给他人
2. **请求频率**：不建议设置过于频繁的检查间隔，以免对服务器造成压力
3. **Cookie 有效期**：Cookie 可能会过期，利用配置热加载功能可以快速更新
4. **网络稳定性**：确保运行环境有稳定的网络连接
5. **配置备份**：修改配置前建议先备份当前配置文件

## Cookie 管理 🍪

本系统内置了智能 Cookie 失效检测和处理机制：

### 自动检测 Cookie 失效

- ✅ 自动检测 HTTP 401/403 状态码
- ✅ 检测理想汽车 API 业务错误码
- ✅ 连续失败 3 次后自动发送告警通知
- ✅ 告警消息包含详细的 Cookie 更新步骤

### Cookie 失效处理

**收到 Cookie 失效告警时，请按以下步骤操作：**

1. 访问 https://www.lixiang.com/ 并登录
2. 按 F12 打开开发者工具
3. 在 Network 标签中找到任意请求
4. 复制 Request Headers 中的完整 Cookie 字符串
5. 更新 `config.yaml` 中的 `lixiang_cookies` 字段
6. 保存文件（配置会自动热加载，无需重启）

**快速参考：**
- 📖 [Cookie 失效快速修复指南](./docs/guides/COOKIE_QUICK_FIX.md) - 5 分钟快速解决
- 📚 [Cookie 管理完整文档](./docs/technical/COOKIE_MANAGEMENT.md) - 详细技术说明

**测试 Cookie 功能：**
```bash
# 运行 Cookie 失效检测测试
./scripts/test/test-cookie-expiry.sh
```

### 关键 Cookie 字段

| 字段 | 说明 | 重要性 |
|------|------|--------|
| `X-LX-Token` | 会话令牌 | ⭐⭐⭐ (最易失效) |
| `authli_device_id` | 设备标识 | ⭐⭐⭐ |
| `X-LX-Deviceid` | 设备 ID | ⭐⭐ |

## 📚 相关文档

### 📖 项目架构
- [ARCHITECTURE.md](./ARCHITECTURE.md) - **完整的系统架构文档**

### 📘 用户指南 (docs/guides/)
- [Cookie 失效快速修复](./docs/guides/COOKIE_QUICK_FIX.md) - 🔥 5 分钟快速解决 Cookie 失效问题
- [微信通知配置](./docs/guides/WECHAT_SETUP.md) - 微信群机器人配置指南
- [ServerChan 配置](./docs/guides/SERVERCHAN_SETUP.md) - Server酱配置指南
- [配置热加载示例](./docs/guides/HOT_RELOAD_DEMO.md) - 配置热加载使用示例
- [测试指南](./docs/guides/TESTING_GUIDE.md) - 完整的测试指南

### 🔬 技术文档 (docs/technical/)
- [配置热加载技术文档](./docs/technical/CONFIG_HOT_RELOAD.md) - 配置热加载功能详细实现
- [Cookie 管理技术文档](./docs/technical/COOKIE_MANAGEMENT.md) - Cookie 失效检测和处理机制详解
- [Cookie 实现总结](./docs/technical/COOKIE_IMPLEMENTATION_SUMMARY.md) - Cookie 功能实现总结
- [热加载实现总结](./docs/technical/IMPLEMENTATION_SUMMARY.md) - 热加载实现总结
- [定期通知功能](./docs/technical/PERIODIC_NOTIFICATION.md) - 定期通知功能说明
- [交付时间优化](./docs/technical/DELIVERY_OPTIMIZATION.md) - 交付时间优化文档

### 🔧 脚本工具
- **测试脚本** (scripts/test/):
  - `test-notification.sh` - 通知功能测试
  - `test-cookie-expiry.sh` - Cookie 失效测试
  - `test-hot-reload.sh` - 配置热加载测试
  - `test-periodic-notification.sh` - 定期通知测试

- **部署脚本** (scripts/deploy/):
  - `build.sh` - 构建脚本
  - `start.sh` - 启动脚本
  - `stop.sh` - 停止脚本
  - `status.sh` - 状态查询脚本

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

请遵循项目的目录结构规范：
- 用户指南放在 `docs/guides/`
- 技术文档放在 `docs/technical/`
- 测试脚本放在 `scripts/test/`
- 部署脚本放在 `scripts/deploy/`

## 📄 许可证

MIT License