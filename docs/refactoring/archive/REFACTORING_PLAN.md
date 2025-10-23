# 代码重构方案 - main.go 模块化拆分

## 📊 当前状态

- **main.go**: 1172 行，包含所有功能
- **问题**: 文件过大，难以维护，职责不清晰

## 🎯 重构目标

将 main.go 拆分成多个职责清晰的模块文件，提高代码可维护性和可读性。

## 📋 模块化方案

### 文件结构

```
lixiang-order-monitor/
├── main.go              # 程序入口 (~50 行)
├── types.go             # 数据结构定义 (~80 行)
├── notifier/            # 通知器模块
│   ├── notifier.go      # 通知器接口定义
│   ├── wechat.go        # 微信通知器
│   ├── serverchan.go    # ServerChan 通知器
│   └── bark.go          # Bark 通知器
├── monitor/             # 监控器模块
│   ├── monitor.go       # Monitor 结构和核心方法
│   ├── delivery.go      # 交付时间计算
│   ├── cookie.go        # Cookie 管理
│   ├── config.go        # 配置管理
│   └── api.go           # API 请求
└── utils/               # 工具函数
    └── time.go          # 时间处理工具
```

### 模块职责划分

#### 1. main.go (程序入口)
```go
package main

import (
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	monitor := NewMonitor()

	// 检查配置
	if len(monitor.Notifiers) == 0 {
		log.Println("警告: 未配置任何通知器，将不会发送通知")
	}

	if monitor.LixiangCookies == "" {
		log.Println("警告: 未配置理想汽车 Cookies，可能导致请求失败")
	}

	// 启动监控
	if err := monitor.Start(); err != nil {
		log.Fatalf("启动监控服务失败: %v", err)
		os.Exit(1)
	}
}
```

#### 2. types.go (数据结构定义)
**内容**:
- 时间格式常量
- OrderResponse 结构
- CookieExpiredError 结构
- Monitor 结构

**行数**: ~80 行

#### 3. notifier/ 目录 (通知器模块)

**notifier/notifier.go** (~20 行):
```go
package notifier

// Notifier 通知接口
type Notifier interface {
	Send(title, content string) error
}
```

**notifier/wechat.go** (~70 行):
- WeChatMessage 结构
- WeChatWebhookNotifier 结构
- Send 方法实现

**notifier/serverchan.go** (~50 行):
- ServerChanNotifier 结构
- Send 方法实现

**notifier/bark.go** (~80 行):
- BarkNotifier 结构
- Send 方法实现

#### 4. monitor/ 目录 (监控器模块)

**monitor/monitor.go** (~100 行):
- NewMonitor 函数
- Start 方法
- Stop 方法
- sendNotification 方法

**monitor/delivery.go** (~250 行):
- calculateEstimatedDelivery 方法
- calculateRemainingDeliveryTime 方法
- calculateDeliveryProgress 方法
- formatDeliveryEstimate 方法
- getDetailedDeliveryInfo 方法
- getDeliveryAnalysisReport 方法
- isApproachingDelivery 方法
- shouldSendPeriodicNotification 方法
- updateLastNotificationTime 方法
- checkDeliveryTime 方法

**monitor/cookie.go** (~200 行):
- checkCookieExpiration 方法
- getCookieStatus 方法
- handleCookieExpired 方法

**monitor/config.go** (~200 行):
- loadConfig 方法
- watchConfig 方法
- 配置初始化逻辑

**monitor/api.go** (~100 行):
- fetchOrderData 方法
- HTTP 请求处理

#### 5. utils/ 目录 (工具函数)

**utils/time.go** (~50 行):
- parseLockOrderTime 函数
- 其他时间处理工具函数

## 🔧 重构步骤

### 方案 A: 手动重构（推荐）

按照以下顺序逐步重构，每次重构后都要编译测试：

1. **第一步**: 创建 notifier 包
   ```bash
   mkdir notifier
   # 创建 notifier/notifier.go
   # 创建 notifier/wechat.go
   # 创建 notifier/serverchan.go
   # 创建 notifier/bark.go
   # 从 main.go 移除对应代码
   go build && go test
   ```

2. **第二步**: 创建 utils 包
   ```bash
   mkdir utils
   # 创建 utils/time.go
   # 从 main.go 移除对应代码
   go build && go test
   ```

3. **第三步**: 创建 monitor 包
   ```bash
   mkdir monitor
   # 创建 monitor/monitor.go（核心结构）
   # 创建 monitor/delivery.go
   # 创建 monitor/cookie.go
   # 创建 monitor/config.go
   # 创建 monitor/api.go
   # 从 main.go 移除对应代码
   go build && go test
   ```

4. **第四步**: 创建 types.go
   ```bash
   # 创建 types.go
   # 从 main.go 移除结构定义
   go build && go test
   ```

5. **第五步**: 简化 main.go
   ```bash
   # 只保留 main 函数
   go build && go test
   ```

### 方案 B: 使用重构脚本（快速）

创建一个自动化重构脚本：

```bash
./scripts/refactor.sh
```

## 📝 注意事项

### 1. 包的可见性
- 导出的类型/函数: 首字母大写（如 `Monitor`, `NewMonitor`）
- 内部使用: 首字母小写（如 `checkCookieExpiration`）

### 2. 循环依赖
避免包之间的循环依赖：
```
main -> monitor -> notifier ✅
main -> notifier -> monitor ❌
```

### 3. 接口设计
```go
// 好的设计：在使用的地方定义接口
package monitor

type Notifier interface {
    Send(title, content string) error
}
```

### 4. 向后兼容
- 保持公开 API 不变
- 测试脚本无需修改
- 配置文件无需修改

## 🧪 测试计划

每一步重构后都要执行：

```bash
# 1. 编译测试
go build -o lixiang-monitor main.go

# 2. 运行测试脚本
./scripts/test/test-notification.sh
./scripts/test/test-bark.sh
./scripts/test/test-cookie-expiration.sh

# 3. 短暂运行验证
timeout 5s ./lixiang-monitor
```

## 📊 重构前后对比

### 重构前
```
main.go (1172 lines)
├── 所有功能混在一起
├── 难以定位代码
├── 修改影响面大
└── 测试困难
```

### 重构后
```
main.go (50 lines)
notifier/ (200 lines)
monitor/ (650 lines)
utils/ (50 lines)
types.go (80 lines)
├── 职责清晰
├── 易于维护
├── 修改隔离
└── 便于测试
```

## 🎯 优先级建议

### 高优先级（必须）
1. ✅ 创建 notifier 包 - 通知器独立性高
2. ✅ 提取 types.go - 减少 main.go 体积

### 中优先级（建议）
3. ⭐ 创建 monitor 包 - 核心逻辑模块化
4. ⭐ 创建 utils 包 - 工具函数复用

### 低优先级（可选）
5. 💡 进一步细分 monitor 子模块
6. 💡 添加单元测试
7. 💡 添加接口文档

## 📋 实施检查清单

- [ ] 备份当前 main.go
- [ ] 创建 notifier 包并测试
- [ ] 创建 utils 包并测试
- [ ] 创建 monitor 包并测试
- [ ] 创建 types.go 并测试
- [ ] 简化 main.go 并测试
- [ ] 运行所有测试脚本
- [ ] 更新文档（如需要）
- [ ] Git 提交

## 🔗 相关文档

- [Go 项目布局标准](https://github.com/golang-standards/project-layout)
- [Effective Go - 包设计](https://golang.org/doc/effective_go#package-names)
- [代码审查指南](./docs/technical/CODE_REVIEW_GUIDELINES.md)

---

> 💡 **建议**: 采用渐进式重构，每次重构一个模块，确保测试通过后再继续下一步。
