# 架构设计文档目录

本目录包含系统架构设计、重构历史和项目组织相关的文档。

## 📋 文档列表

### [系统架构文档](ARCHITECTURE.md)
**内容概要**:
- 系统整体架构设计
- 模块划分和职责
- 数据流向和交互
- 技术栈选型和说明
- 设计原则和模式

**适合阅读对象**:
- 新加入的开发者
- 想要了解系统全貌的用户
- 需要进行二次开发的开发者

---

### [重构总结报告](REFACTORING_FINAL_REPORT.md)
**内容概要**:
- 重构的背景和动机
- 重构目标和计划
- 重构过程详细记录
- 重构前后对比
- 优化效果评估
- 经验教训总结

**重构亮点**:
- ✅ 模块化设计
- ✅ 代码组织优化
- ✅ 可维护性提升
- ✅ 扩展性增强
- ✅ 性能优化

---

### [项目重组文档](PROJECT_REORGANIZATION.md)
**内容概要**:
- 目录结构调整
- 文件重新组织
- 代码模块化
- 文档体系建立
- 最佳实践应用

**重组成果**:
- ✅ 清晰的包结构
- ✅ 合理的职责划分
- ✅ 完整的文档体系
- ✅ 易于维护和扩展

---

## 🏗️ 系统架构概览

### 核心包结构

```
lixiang-monitor/
├── main.go              # 程序入口，监控协调
├── cfg/                 # 配置管理包
│   └── config.go
├── cookie/              # Cookie 管理包
│   └── cookie.go
├── db/                  # 数据库管理包
│   └── database.go
├── delivery/            # 交付信息管理包
│   └── delivery.go
├── notification/        # 通知处理包
│   └── handler.go
├── notifier/            # 通知器实现包
│   ├── notifier.go
│   ├── bark.go
│   ├── serverchan.go
│   └── wechat.go
├── utils/               # 工具函数包
│   └── time.go
└── web/                 # Web 服务器包
    ├── server.go
    └── templates/
```

### 模块职责

| 包名 | 职责 | 依赖 |
|------|------|------|
| `main` | 程序入口、监控协调 | 所有包 |
| `cfg` | 配置文件管理、热重载 | viper, fsnotify |
| `cookie` | Cookie 生命周期管理 | utils |
| `db` | 数据持久化、历史记录 | sqlite |
| `delivery` | 交付时间计算、状态管理 | utils |
| `notification` | 通知逻辑协调 | notifier, delivery |
| `notifier` | 具体通知渠道实现 | http |
| `utils` | 时间处理等工具函数 | time |
| `web` | Web 界面、API 服务 | db, template |

### 数据流向

```
配置文件 (config.yaml)
    ↓
cfg 包 (配置管理)
    ↓
main.go (监控器初始化)
    ↓
├─→ delivery 包 (交付信息)
├─→ cookie 包 (Cookie 管理)
├─→ notification 包 (通知协调)
│       ↓
│   notifier 包 (通知发送)
│       ↓
│   [微信/Bark/ServerChan]
│
├─→ db 包 (数据存储)
│       ↓
│   SQLite 数据库
│
└─→ web 包 (Web 服务)
        ↓
    HTTP API / Web 界面
```

## 🎯 设计原则

### 1. 单一职责原则 (SRP)
每个包只负责一个特定功能：
- `cfg` 只管理配置
- `db` 只管理数据
- `notifier` 只管理通知

### 2. 依赖倒置原则 (DIP)
使用接口定义抽象：
```go
type Notifier interface {
    SendNotification(message string) error
}
```

### 3. 开闭原则 (OCP)
对扩展开放，对修改关闭：
- 添加新通知器：实现 `Notifier` 接口
- 添加新存储：实现数据库接口

### 4. 最小知识原则
模块间通过明确的接口通信，减少耦合

## 📊 技术栈

### 核心技术

| 技术 | 用途 | 版本 |
|------|------|------|
| Go | 编程语言 | 1.21+ |
| Viper | 配置管理 | v1.x |
| fsnotify | 文件监听 | v1.x |
| robfig/cron | 定时任务 | v3.x |
| modernc.org/sqlite | 数据库 | v1.39+ |

### 标准库

- `net/http` - HTTP 服务器
- `html/template` - 模板渲染
- `embed` - 资源嵌入
- `time` - 时间处理
- `sync` - 并发控制

## 🔄 架构演进

### v1.0 - 单文件架构
```
main.go (所有功能)
```

### v1.5 - 初步模块化
```
main.go
notifier/ (通知器)
utils/ (工具函数)
```

### v2.0 - 完整模块化（当前）
```
main.go
cfg/ (配置)
cookie/ (Cookie)
db/ (数据库)
delivery/ (交付)
notification/ (通知协调)
notifier/ (通知器)
utils/ (工具)
web/ (Web服务)
```

## 📖 阅读建议

### 新开发者入门顺序

1. **[系统架构文档](ARCHITECTURE.md)** - 了解整体设计
2. **[项目重组文档](PROJECT_REORGANIZATION.md)** - 了解文件组织
3. **[重构总结报告](REFACTORING_FINAL_REPORT.md)** - 了解演进历程
4. **[技术文档](../technical/)** - 深入具体实现

### 深入研究路径

**理解核心流程**:
1. 阅读 `main.go` 了解程序入口
2. 查看 `cfg` 包了解配置管理
3. 研究 `delivery` 包了解核心业务逻辑
4. 探索 `notification` 包了解通知机制

**理解数据流**:
1. 配置加载 → `cfg` 包
2. 数据采集 → `main.go`
3. 数据存储 → `db` 包
4. 数据展示 → `web` 包

**理解扩展点**:
1. 新增通知器 → 实现 `Notifier` 接口
2. 新增存储 → 扩展 `db` 包
3. 新增 API → 扩展 `web` 包

## 🛠️ 开发指南

### 添加新功能

1. **确定模块归属**
   - 通知相关 → `notifier` 包
   - 数据相关 → `db` 包
   - Web 相关 → `web` 包

2. **遵循现有模式**
   - 参考同类功能实现
   - 保持代码风格一致
   - 添加必要的注释

3. **更新文档**
   - 更新相关技术文档
   - 添加使用示例
   - 编写变更日志

### 代码组织原则

```go
// 包级别文档注释
package mypackage

// 导入标准库
import (
    "fmt"
    "time"
)

// 导入第三方库
import (
    "github.com/spf13/viper"
)

// 导入本项目包
import (
    "lixiang-monitor/utils"
)

// 常量定义
const (
    DefaultTimeout = 10 * time.Second
)

// 类型定义
type MyStruct struct {
    field1 string
    field2 int
}

// 函数实现
func NewMyStruct() *MyStruct {
    return &MyStruct{}
}
```

## 🔗 相关链接

- [文档索引](../INDEX.md)
- [用户指南](../guides/)
- [技术文档](../technical/)
- [变更日志](../changelogs/)

---

**目录维护**: ✅ 活跃维护  
**最后更新**: 2025-10-23  
**文档数量**: 3 个架构文档
