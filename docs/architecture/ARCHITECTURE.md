# 理想汽车订单监控系统 - 项目架构

> 🎉 **2025-01-XX 更新**: 已完成代码模块化重构（Phase 1-5），notifier、utils、cfg、delivery 和 cookie 包已独立

## 📁 项目结构

```
lixiang-order-monitor/
├── 📦 cfg/                           # 配置管理包（Phase 3）
│   └── config.go                    # 配置加载、热更新、通知器创建
│
├── 📦 delivery/                      # 交付时间计算包（Phase 5 新增）
│   └── delivery.go                  # 交付日期计算、进度分析、智能提醒
│
├── 📦 cookie/                        # Cookie 管理包（Phase 5 新增）
│   └── cookie.go                    # Cookie 验证、过期检测、失效处理
│
├── 📦 notifier/                      # 通知器包（Phase 1）
│   ├── notifier.go                  # 通知器接口定义
│   ├── serverchan.go                # ServerChan 通知器实现
│   ├── wechat.go                    # 微信群机器人通知器实现
│   └── bark.go                      # Bark 推送通知器实现
│
├── 🔧 utils/                         # 工具包（Phase 2）
│   └── time.go                      # 时间工具函数（格式化、解析）
│
├── 📚 docs/                          # 文档目录
│   ├── guides/                       # 用户指南
│   │   ├── COOKIE_QUICK_FIX.md      # Cookie 失效快速修复指南
│   │   ├── WECHAT_SETUP.MD          # 微信群机器人配置指南
│   │   ├── SERVERCHAN_SETUP.md      # ServerChan 配置指南
│   │   ├── BARK_SETUP.md            # Bark 推送配置指南
│   │   ├── HOT_RELOAD_DEMO.md       # 配置热加载使用示例
│   │   └── TESTING_GUIDE.md         # 测试指南
│   │
│   └── technical/                    # 技术文档
│       ├── CONFIG_HOT_RELOAD.md     # 配置热加载技术文档
│       ├── COOKIE_MANAGEMENT.md     # Cookie 管理技术文档
│       ├── COOKIE_IMPLEMENTATION_SUMMARY.md  # Cookie 实现总结
│       ├── IMPLEMENTATION_SUMMARY.md # 热加载实现总结
│       ├── PERIODIC_NOTIFICATION.md  # 定期通知功能文档
│       ├── DELIVERY_OPTIMIZATION.md  # 交付时间优化文档
│       └── PROJECT_FILES.md         # 项目文件说明
│
├── 🔧 scripts/                       # 脚本目录
│   ├── test/                         # 测试脚本
│   │   ├── test-cookie-expiry.sh    # Cookie 失效测试
│   │   ├── test-hot-reload.sh       # 配置热加载测试
│   │   ├── test-notification.sh     # 通知功能测试
│   │   ├── test-bark.sh             # Bark 推送测试
│   │   ├── test-periodic-notification.sh  # 定期通知测试
│   │   └── test_delivery_calc.go    # 交付时间计算测试
│   │
│   ├── deploy/                       # 部署脚本
│   │   ├── build.sh                 # 构建脚本
│   │   ├── start.sh                 # 启动脚本
│   │   ├── stop.sh                  # 停止脚本
│   │   └── status.sh                # 状态查询脚本
│   │
│   ├── refactor.sh                  # 代码重构辅助脚本
│   └── reorganize-project.sh        # 项目重组脚本
│
├── ⚙️ config/                        # 配置模板目录
│   ├── config.example.yaml          # 配置文件示例
│   └── config.enhanced.yaml         # 增强配置示例
│
├── 📝 主要文件
│   ├── main.go                      # 主程序源码（906 行，Phase 3 优化）
│   ├── config.yaml                  # 工作配置文件（不提交到 Git）
│   ├── go.mod                       # Go 模块依赖
│   ├── go.sum                       # Go 依赖校验
│   ├── README.md                    # 项目说明文档
│   ├── ARCHITECTURE.md              # 本文件：架构说明
│   ├── REFACTORING_PLAN.md          # 重构计划文档
│   ├── REFACTORING_SUMMARY.md       # 重构总结报告
│   ├── REFACTORING_PHASE3_COMPLETE.md  # Phase 3 完成报告（新增）
│   └── .gitignore                   # Git 忽略规则
│
└── 🚀 构建产物
    ├── lixiang-monitor              # 编译后的可执行文件
    └── monitor.log                  # 运行日志（不提交到 Git）
```

---

## 🏗️ 系统架构

### 核心组件

```
┌─────────────────────────────────────────────────────────────┐
│                      理想汽车订单监控系统                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
        ┌─────────────────────────────────────────┐
        │         main.go (主程序)                 │
        │  - Monitor 核心结构                      │
        │  - 业务逻辑协调                          │
        │  - 定时任务 (Cron)                       │
        └─────────────────────────────────────────┘
                              │
        ┌─────────┬───────────┼───────────┬─────────────┬─────────────┐
        │         │           │           │             │             │
        ▼         ▼           ▼           ▼             ▼             ▼
┌───────────┐ ┌─────────┐ ┌─────────┐ ┌─────────────┐ ┌──────────┐ ┌──────────┐
│  cfg 包    │ │delivery │ │ cookie  │ │ notifier 包 │ │ utils 包  │ │HTTP请求  │
├───────────┤ ├─────────┤ ├─────────┤ ├─────────────┤ ├──────────┤ ├──────────┤
│• 配置加载  │ │• 交付计算│ │• Cookie │ │ • 微信机器人 │ │ • 时间格式│ │• API调用 │
│• 热加载    │ │• 进度分析│ │  验证   │ │ • ServerChan │ │ • 时间解析│ │• 错误处理│
│• 通知器   │ │• 日期预测│ │• 过期   │ │ • Bark推送   │ │          │ │          │
│  创建     │ │• 智能提醒│ │  检测   │ │ • 接口定义   │ │          │ │          │
│• 并发安全 │ │          │ │• 失效   │ │              │ │          │ │          │
│          │ │          │ │  处理   │ │              │ │          │ │          │
└───────────┘ └─────────┘ └─────────┘ └─────────────┘ └──────────┘ └──────────┘
        │         │           │           │             │             │
        └─────────┴───────────┼───────────┴─────────────┴─────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │   理想汽车 API   │
                    │  订单数据接口    │
                    └─────────────────┘
```

### 数据流图

```
1. 定时触发
   └─> Cron (@every 12h)
       └─> checkDeliveryTime()

2. 获取数据
   └─> fetchOrderData()
       ├─> HTTP GET 请求
       │   ├─> Headers (Cookie, User-Agent, etc.)
       │   └─> https://api-web.lixiang.com/.../orders/...
       │
       ├─> 状态码检测 (cookie 包)
       │   ├─> 401/403 → Cookie 失效
       │   └─> 200 → 继续
       │
       └─> JSON 解析
           └─> OrderResponse

3. 数据分析
   └─> 比对交付时间
       ├─> 交付计算 (delivery 包)
       │   ├─> 计算预计日期
       │   ├─> 计算剩余时间
       │   ├─> 计算进度百分比
       │   └─> 检查是否临近
       │
       ├─> 时间变化？
       │   ├─> 是 → 发送变更通知
       │   └─> 否 → 检查定期通知条件
       │
       └─> 临近交付？
           └─> 是 → 发送提醒通知

4. 发送通知
   └─> sendNotification(title, content)
       ├─> ServerChan
       │   └─> POST https://sctapi.ftqq.com/{sendkey}.send
       │
       └─> WeChat Webhook
           └─> POST https://qyapi.weixin.qq.com/...

5. Cookie 失效处理
   └─> handleCookieExpired()
       ├─> 失败计数 +1
       ├─> 达到 3 次？
       │   └─> 发送 Cookie 失效告警
       │
       └─> 记录日志
```

---

## 🔧 核心功能模块

### 1. 配置管理 (Configuration Management) - cfg 包

**文件**: `cfg/config.go`

**功能**:
- ✅ 使用 Viper 管理配置
- ✅ 支持配置热加载（fsnotify）
- ✅ 线程安全
- ✅ 通知器创建和管理

**关键方法**:
```go
func Init() error                    // 初始化配置
func Load() (*Config, error)         // 加载配置
func Watch(onConfigChange func())    // 监听配置变化
```

**相关文档**:
- `docs/technical/CONFIG_HOT_RELOAD.md`
- `docs/guides/HOT_RELOAD_DEMO.md`

---

### 2. 交付时间计算 (Delivery Calculation) - delivery 包

**文件**: `delivery/delivery.go`

**功能**:
- ✅ 预计交付日期计算
- ✅ 剩余时间分析
- ✅ 进度百分比计算
- ✅ 智能提醒判断

**核心结构**:
```go
type Info struct {
    LockOrderTime    time.Time
    EstimateWeeksMin int
    EstimateWeeksMax int
}
```

**关键方法**:
```go
func (d *Info) CalculateEstimatedDelivery() (time.Time, time.Time)
func (d *Info) CalculateRemainingDeliveryTime() (int, int, string)
func (d *Info) CalculateDeliveryProgress() float64
func (d *Info) FormatDeliveryEstimate() string
func (d *Info) GetDetailedDeliveryInfo() string
func (d *Info) GetAnalysisReport() string
func (d *Info) IsApproachingDelivery() (bool, string)
```

**相关文档**:
- `docs/technical/DELIVERY_OPTIMIZATION.md`

---

### 3. Cookie 管理 (Cookie Management) - cookie 包

**文件**: `cookie/cookie.go`

**功能**:
- ✅ 自动检测 Cookie 失效
- ✅ 过期预警（提前 2 天）
- ✅ 连续失败计数
- ✅ 智能告警（3 次失败）
- ✅ 回调机制

**核心结构**:
```go
type Manager struct {
    Cookies                   string
    Headers                   map[string]string
    ValidDays                 int
    UpdatedAt                 time.Time
    ExpirationWarned          bool
    ConsecutiveFailure        int
    ExpiredNotified           bool
    LastCheckTime             time.Time
    OnCookieExpired           func(statusCode int, message string)
    OnCookieExpirationWarning func(timeDesc, expireTime, updatedAt string, ageInDays float64)
}
```

**关键方法**:
```go
func (cm *Manager) FetchOrderData(orderID string) (interface{}, error)
func (cm *Manager) CheckExpiration()
func (cm *Manager) GetStatus() string
func (cm *Manager) UpdateCookie(cookies string, headers map[string]string)
```

**错误类型**:
```go
type CookieExpiredError struct {
    StatusCode int
    Message    string
}
```

**相关文档**:
- `docs/technical/COOKIE_MANAGEMENT.md`
- `docs/guides/COOKIE_QUICK_FIX.md`
- `docs/technical/COOKIE_IMPLEMENTATION_SUMMARY.md`

---

### 4. 通知系统 (Notification System) - notifier 包

**文件**: `notifier/*.go`

**功能**:
- ✅ 多通道支持
- ✅ 通知接口抽象
- ✅ 错误处理和重试

**接口定义**:
```go
type Notifier interface {
    Send(title, content string) error
}
```

**实现**:
- `ServerChanNotifier` - Server酱微信推送
- `WeChatWebhookNotifier` - 微信群机器人

**相关文档**:
- `docs/guides/SERVERCHAN_SETUP.md`
- `docs/guides/WECHAT_SETUP.md`

---

### 4. 定时任务 (Scheduled Tasks)

**文件**: `main.go` (cron setup)

**功能**:
- ✅ 使用 robfig/cron 实现
- ✅ 可配置检查间隔
- ✅ 优雅停止

**Cron 表达式**:
- `@every 30m` - 每 30 分钟
- `@every 1h` - 每小时
- `@every 12h` - 每 12 小时

**相关文档**:
- `docs/technical/PERIODIC_NOTIFICATION.md`

---

### 5. 交付时间预测 (Delivery Prediction)

**文件**: `main.go` (checkDeliveryTime, calculateDeliveryRange)

**功能**:
- ✅ 基于锁单时间预测
- ✅ 交付日期范围计算
- ✅ 临近提醒

**配置参数**:
- `lock_order_time` - 锁单时间
- `estimate_weeks_min` - 最少周数
- `estimate_weeks_max` - 最多周数

**相关文档**:
- `docs/technical/DELIVERY_OPTIMIZATION.md`

---

## 🔐 安全机制

### 线程安全

```go
type Monitor struct {
    mu sync.RWMutex  // 读写锁
    // ...
}

// 读取配置
m.mu.RLock()
value := m.SomeField
m.mu.RUnlock()

// 修改配置
m.mu.Lock()
m.SomeField = newValue
m.mu.Unlock()
```

### Cookie 保护

- ❌ 不在日志中记录完整 Cookie
- ❌ 不在告警中包含敏感信息
- ✅ 建议添加到 `.gitignore`

### 错误隔离

- Cookie 失效不影响其他功能
- 程序继续运行等待 Cookie 更新
- 优雅处理网络错误

---

## 📊 数据结构

### Monitor 核心结构

```go
type Monitor struct {
    // 基础配置
    OrderID          string
    CheckInterval    string
    LixiangCookies   string
    LixiangHeaders   map[string]string
    
    // 交付预测
    LockOrderTime    time.Time
    EstimateWeeksMin int
    EstimateWeeksMax int
    LastEstimateTime string
    
    // 通知系统
    Notifiers        []Notifier
    
    // 定期通知
    LastNotificationTime        time.Time
    NotificationInterval        time.Duration
    EnablePeriodicNotify        bool
    AlwaysNotifyWhenApproaching bool
    
    // Cookie 管理
    LastCookieCheckTime      time.Time
    CookieExpiredNotified    bool
    ConsecutiveCookieFailure int
    
    // 配置热加载
    mu            sync.RWMutex
    configVersion int
    
    // 定时任务
    cron *cron.Cron
}
```

### API 响应结构

```go
type OrderResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    struct {
        Delivery struct {
            EstimateDeliveringAt string `json:"estimateDeliveringAt"`
        } `json:"delivery"`
    } `json:"data"`
}
```

---

## 🚀 部署架构

### 开发环境

```bash
# 安装依赖
go mod download

# 编译
go build -o lixiang-monitor main.go

# 运行
./lixiang-monitor
```

### 生产环境

```bash
# 使用部署脚本
./scripts/deploy/build.sh    # 构建
./scripts/deploy/start.sh    # 启动
./scripts/deploy/status.sh   # 查看状态
./scripts/deploy/stop.sh     # 停止
```

### 进程管理

**推荐使用 systemd**:

```ini
[Unit]
Description=理想汽车订单监控服务
After=network.target

[Service]
Type=simple
User=yourusername
WorkingDirectory=/path/to/lixiang-order-monitor
ExecStart=/path/to/lixiang-order-monitor/lixiang-monitor
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
```

---

## 📈 性能指标

### 资源占用

| 指标 | 值 |
|------|-----|
| 内存占用 | ~20-30 MB |
| CPU 占用 | < 1% (空闲时) |
| 磁盘占用 | < 10 MB (程序 + 日志) |
| 网络流量 | ~1 KB/次请求 |

### 响应时间

| 操作 | 时间 |
|------|------|
| 配置热加载 | < 1 秒 |
| API 请求 | 1-3 秒 |
| 通知发送 | 1-2 秒 |
| Cookie 检测 | < 1 毫秒 |

---

## 🧪 测试策略

### 单元测试

```bash
# 交付时间计算测试
go run scripts/test/test_delivery_calc.go
```

### 集成测试

```bash
# 通知功能测试
./scripts/test/test-notification.sh

# 配置热加载测试
./scripts/test/test-hot-reload.sh

# Cookie 失效测试
./scripts/test/test-cookie-expiry.sh

# 定期通知测试
./scripts/test/test-periodic-notification.sh
```

**相关文档**:
- `docs/guides/TESTING_GUIDE.md`

---

## 📚 文档导航

### 快速开始
1. 阅读 `README.md` - 项目概述和快速开始
2. 配置通知方式:
   - `docs/guides/SERVERCHAN_SETUP.md` - ServerChan 配置
   - `docs/guides/WECHAT_SETUP.md` - 微信机器人配置

### 使用指南
- `docs/guides/COOKIE_QUICK_FIX.md` - Cookie 失效快速修复
- `docs/guides/HOT_RELOAD_DEMO.md` - 配置热加载使用示例
- `docs/guides/TESTING_GUIDE.md` - 测试指南

### 技术文档
- `docs/technical/CONFIG_HOT_RELOAD.md` - 配置热加载技术实现
- `docs/technical/COOKIE_MANAGEMENT.md` - Cookie 管理机制
- `docs/technical/PERIODIC_NOTIFICATION.md` - 定期通知功能
- `docs/technical/DELIVERY_OPTIMIZATION.md` - 交付时间优化

### 实现总结
- `docs/technical/IMPLEMENTATION_SUMMARY.md` - 热加载实现总结
- `docs/technical/COOKIE_IMPLEMENTATION_SUMMARY.md` - Cookie 功能实现总结

---

## 🔄 版本历史

### v1.1.0 (2025-10-20)
- ✅ 新增 Cookie 失效自动检测
- ✅ 新增智能告警通知
- ✅ 优化项目目录结构
- ✅ 完善文档体系

### v1.0.0 (2025-09-27)
- ✅ 配置热加载功能
- ✅ 定期通知功能
- ✅ 交付时间预测
- ✅ 多通道通知支持

---

## 🤝 贡献指南

### 目录规范

- `docs/guides/` - 用户指南（面向最终用户）
- `docs/technical/` - 技术文档（面向开发者）
- `scripts/test/` - 测试脚本
- `scripts/deploy/` - 部署脚本
- `config/` - 配置模板

### 文档规范

- 使用 Markdown 格式
- 中文文档优先
- 包含清晰的示例代码
- 添加目录索引

### 代码规范

- 遵循 Go 官方编码规范
- 添加必要的注释
- 使用有意义的变量名
- 错误处理要完善

---

## 📞 支持

### 文档
- 项目架构: `ARCHITECTURE.md` (本文件)
- 使用说明: `README.md`
- 用户指南: `docs/guides/`
- 技术文档: `docs/technical/`

### 问题反馈
- GitHub Issues
- 邮件支持

---

**最后更新**: 2025-10-20  
**维护者**: haifeiWu
