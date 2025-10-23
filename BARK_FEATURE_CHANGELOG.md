# Bark 推送功能更新日志

## 📅 2025-10-22 - Bark 推送功能上线

### 🎉 新功能

新增 **Bark 推送** 通知渠道，为 iOS/macOS 用户提供更好的推送体验。

### ✨ 功能特性

1. **实时推送**: 基于 APNs，推送速度快
2. **自定义声音**: 支持 30+ 种内置提示音
3. **自定义图标**: 支持自定义推送图标
4. **通知分组**: 支持按分组管理通知
5. **自建服务器**: 支持自建 Bark 服务器，数据私有
6. **无限制推送**: 无推送数量限制

### 🔧 技术实现

#### 新增结构体

```go
// Bark 推送通知结构
type BarkNotifier struct {
    ServerURL string // Bark 服务器 URL
    Sound     string // 推送声音
    Icon      string // 推送图标 URL
    Group     string // 通知分组
}
```

#### 实现方法

```go
func (bark *BarkNotifier) Send(title, content string) error
```

发送 HTTP POST 请求到 Bark 服务器，支持以下参数：
- `title`: 推送标题
- `body`: 推送内容
- `sound`: 推送声音（默认 "minuet"）
- `icon`: 推送图标 URL（可选）
- `group`: 通知分组（默认 "lixiang-monitor"）

### 📝 配置说明

#### 新增配置项

```yaml
# Bark 推送配置
bark_server_url: "http://your_server:8080/your_key"  # Bark 服务器 URL（必填）
bark_sound: "minuet"                                 # 推送声音（可选）
bark_icon: ""                                        # 推送图标（可选）
bark_group: "lixiang-monitor"                        # 通知分组（可选）
```

#### 配置示例

官方服务器：
```yaml
bark_server_url: "https://api.day.app/your_device_key"
bark_sound: "minuet"
bark_group: "lixiang-monitor"
```

自建服务器：
```yaml
bark_server_url: "http://43.128.109.177:8080/mruLduP9zEnDabDA4uBVLj"
bark_sound: "alarm"
bark_icon: "https://www.lixiang.com/favicon.ico"
bark_group: "lixiang-urgent"
```

### 📚 文档更新

#### 新增文档

1. **用户指南** (`docs/guides/`)
   - `BARK_SETUP.md` - Bark 推送配置完整指南
     - 官方服务器配置
     - 自建服务器部署
     - 配置参数详解
     - 推送声音列表
     - 故障排查
     - 使用技巧

2. **测试脚本** (`scripts/test/`)
   - `test-bark.sh` - Bark 推送测试脚本

#### 更新文档

1. **README.md**
   - 功能特性列表新增 Bark 推送
   - 配置示例新增 Bark 配置
   - 通知方式新增 Bark 说明

2. **config/config.example.yaml**
   - 新增 Bark 配置示例和注释

3. **config/config.enhanced.yaml**
   - 新增 Bark 配置项

4. **docs/INDEX.md**
   - 配置指南新增 BARK_SETUP.md
   - 测试脚本新增 test-bark.sh

### 🎯 使用指南

#### 快速开始

1. **安装 Bark App**（iOS/macOS）

2. **获取推送 URL**
   - 官方服务器：App 自动生成
   - 自建服务器：部署后添加到 App

3. **配置到项目**
```yaml
bark_server_url: "你的推送URL"
```

4. **测试推送**
```bash
./scripts/test/test-bark.sh
```

#### 推送声音

支持 30+ 种声音，常用的有：

| 声音 | 场景 |
|------|------|
| `alarm` | 紧急通知 |
| `bell` | 重要提醒 |
| `minuet` | 日常通知（默认） |
| `glass` | 简洁通知 |
| `silence` | 静音 |

### 📊 对比其他通知渠道

| 特性 | Bark | 微信机器人 | ServerChan |
|------|------|-----------|-----------|
| **平台** | iOS/macOS | 微信 | 微信 |
| **实时性** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| **自定义性** | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐ |
| **推送限制** | 无限制 | 无限制 | 5次/天（免费） |
| **部署难度** | ⭐⭐ | ⭐ | ⭐ |
| **费用** | 免费 | 免费 | 免费版限额 |

### 💡 推荐配置方案

#### 方案 1: iOS/Mac 用户
```yaml
bark_server_url: "your_bark_url"        # 主要通知
wechat_webhook_url: "your_wechat_url"   # 备用通知
```

#### 方案 2: 多平台用户
```yaml
serverchan_sendkey: "your_sendkey"      # 微信通知
bark_server_url: "your_bark_url"        # iOS 通知
wechat_webhook_url: "your_wechat_url"   # 群通知
```

#### 方案 3: 纯 iOS 用户
```yaml
bark_server_url: "your_bark_url"
bark_sound: "minuet"
bark_group: "lixiang-monitor"
```

### 🧪 测试验证

创建了专门的测试脚本 `scripts/test/test-bark.sh`：

```bash
cd scripts/test
./test-bark.sh
```

测试脚本功能：
- ✅ 自动读取配置
- ✅ 发送测试推送
- ✅ 显示响应结果
- ✅ 错误诊断提示

### 🔄 与现有功能的集成

Bark 推送作为第三种通知渠道，与现有功能完美集成：

```
┌──────────────────────────────┐
│      通知系统架构             │
├──────────────────────────────┤
│                              │
│  通知器 1: 微信机器人         │
│  通知器 2: ServerChan        │
│  通知器 3: Bark (NEW)        │
│                              │
│  • 同时发送到所有已配置渠道  │
│  • 任一渠道失败不影响其他    │
│  • 配置热加载支持            │
│                              │
└──────────────────────────────┘
```

### 🎁 额外优化

1. **默认值处理**: 
   - Sound 默认为 "minuet"
   - Group 默认为 "lixiang-monitor"

2. **错误处理**:
   - 详细的错误信息
   - HTTP 状态码检查
   - 响应体解析

3. **日志记录**:
   - 发送成功日志
   - 发送失败详细错误

### 🚀 未来计划

可能的增强功能：

1. **多设备支持**: 支持配置多个 Bark URL
2. **优先级设置**: 支持设置推送优先级
3. **时间敏感性**: 支持时间敏感通知
4. **图片推送**: 支持发送图片附件
5. **URL 跳转**: 支持点击推送跳转 URL

### 📈 预期效果

- ✅ 为 iOS/macOS 用户提供原生推送体验
- ✅ 降低对微信的依赖
- ✅ 提高推送实时性
- ✅ 增强用户自定义能力
- ✅ 支持数据隐私保护（自建服务器）

### 🔗 相关链接

- [Bark GitHub](https://github.com/Finb/Bark)
- [Bark 配置指南](../docs/guides/BARK_SETUP.md)
- [测试脚本](../scripts/test/test-bark.sh)

---

## 🎯 总结

Bark 推送功能为 iOS/macOS 用户带来了更好的推送体验：

**核心优势**:
- 📱 **原生体验**: 基于 APNs，推送稳定快速
- 🎨 **高度自定义**: 声音、图标、分组随心配置
- 🏠 **数据私有**: 支持自建服务器，掌控数据
- 🆓 **完全免费**: 开源免费，无任何限制

这是一个**为苹果生态用户量身打造**的推送方案，与现有通知渠道形成完美互补！

---

> 💡 **推荐**: iOS/Mac 用户强烈建议使用 Bark 作为主要推送渠道，配合微信机器人作为备用，实现推送双保险！
