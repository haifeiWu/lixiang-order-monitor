# Cookie 过期预警功能实现总结

## 📋 概述

本文档总结了 Cookie 过期预警功能的完整实现，包括技术细节、配置说明和使用指导。

**实现日期**: 2025-10-20  
**版本**: v1.x  
**状态**: ✅ 已完成并测试

## 🎯 功能目标

虽然 Cookie 自动续期在技术上不可行（需要手动短信验证码），但我们可以通过**提前预警**的方式，让用户在 Cookie 过期前及时更新，从而避免监控中断。

### 核心目标

1. **预防性监控**: 提前 48 小时发现即将过期的 Cookie
2. **用户友好**: 提供详细的更新步骤和指导
3. **自动化检查**: 无需人工干预，系统自动定期检查
4. **状态可见**: 启动时显示 Cookie 健康状态

## 🏗️ 架构设计

### 数据结构

在 `Monitor` 结构体中新增三个字段：

```go
type Monitor struct {
    // ... 现有字段 ...
    
    // Cookie 过期管理
    CookieUpdatedAt        time.Time  // Cookie 最后更新时间
    CookieValidDays        int        // Cookie 有效期（天）
    CookieExpirationWarned bool       // 是否已发送过期警告
}
```

### 核心方法

#### 1. `checkCookieExpiration()` - Cookie 过期检查

**功能**: 检查 Cookie 是否即将过期（48小时内），如果是则发送预警通知

**触发条件**:
- 距离过期时间 <= 48 小时
- 未发送过预警（避免重复通知）

**通知内容**:
- 过期倒计时
- 详细的更新步骤
- 配置示例
- 关键时间点

**代码逻辑**:
```go
func (m *Monitor) checkCookieExpiration() {
    // 1. 计算过期时间
    expirationTime := m.CookieUpdatedAt.Add(time.Duration(m.CookieValidDays) * 24 * time.Hour)
    
    // 2. 计算剩余时间
    remainingDuration := time.Until(expirationTime)
    
    // 3. 判断是否需要预警（48小时内且未预警过）
    warningThreshold := 48 * time.Hour
    if remainingDuration <= warningThreshold && remainingDuration > 0 && !m.CookieExpirationWarned {
        // 4. 发送预警通知
        m.sendNotification(title, content)
        m.CookieExpirationWarned = true
    }
    
    // 5. Cookie 已过期的情况
    if remainingDuration <= 0 {
        // 发送已过期通知
    }
}
```

#### 2. `getCookieStatus()` - Cookie 状态查询

**功能**: 返回格式化的 Cookie 状态字符串

**返回值**:
- `🟢 正常 (剩余 X 天)` - Cookie 有效期充足
- `⚠️ 即将过期 (剩余 X 小时)` - 48 小时内将过期
- `❌ 已过期 (已过期 X 天)` - Cookie 已失效

**代码逻辑**:
```go
func (m *Monitor) getCookieStatus() string {
    expirationTime := m.CookieUpdatedAt.Add(time.Duration(m.CookieValidDays) * 24 * time.Hour)
    remainingDuration := time.Until(expirationTime)
    
    if remainingDuration <= 0 {
        return fmt.Sprintf("❌ 已过期 (已过期 %.0f 天)", math.Abs(remainingDuration.Hours()/24))
    } else if remainingDuration <= 48*time.Hour {
        return fmt.Sprintf("⚠️ 即将过期 (剩余 %.0f 小时)", remainingDuration.Hours())
    } else {
        return fmt.Sprintf("🟢 正常 (剩余 %.0f 天)", remainingDuration.Hours()/24)
    }
}
```

#### 3. `loadConfig()` - 配置加载增强

**新增配置加载**:
```go
func (m *Monitor) loadConfig() {
    // ... 现有配置加载 ...
    
    // 加载 Cookie 有效期配置
    m.CookieValidDays = viper.GetInt("cookie_valid_days")
    if m.CookieValidDays == 0 {
        m.CookieValidDays = 7 // 默认 7 天
    }
    
    // 加载 Cookie 更新时间
    cookieUpdatedAtStr := viper.GetString("cookie_updated_at")
    if cookieUpdatedAtStr != "" {
        if t, err := time.Parse("2006-01-02 15:04:05", cookieUpdatedAtStr); err == nil {
            m.CookieUpdatedAt = t
        }
    }
    if m.CookieUpdatedAt.IsZero() {
        m.CookieUpdatedAt = time.Now() // 默认使用当前时间
    }
    
    // 重置警告标志（配置更新后重新评估）
    m.CookieExpirationWarned = false
}
```

### 集成点

#### 1. 启动时检查

在 `Start()` 方法中添加：

```go
func (m *Monitor) Start() error {
    // ... 现有启动逻辑 ...
    
    // 显示 Cookie 状态
    log.Printf("Cookie 状态: %s", m.getCookieStatus())
    
    // 立即执行一次过期检查
    m.checkCookieExpiration()
    
    // ... 添加定时任务 ...
}
```

#### 2. 定时检查任务

添加每日凌晨 1 点的检查任务：

```go
// 添加定时任务 - 每日检查 Cookie 过期（凌晨 1 点）
_, err = m.cron.AddFunc("0 1 * * *", func() {
    log.Printf("执行定期 Cookie 过期检查")
    m.checkCookieExpiration()
})
```

## 📝 配置说明

### 配置文件格式

```yaml
# Cookie 过期管理
cookie_valid_days: 7                     # Cookie 有效期（天）
cookie_updated_at: "2025-10-20 10:00:00" # Cookie 最后更新时间
```

### 配置参数详解

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `cookie_valid_days` | int | 否 | 7 | Cookie 有效期（天），建议根据实际情况调整 |
| `cookie_updated_at` | string | 否 | 当前时间 | Cookie 最后更新时间，格式: "2006-01-02 15:04:05" |

### 热加载支持

配置文件修改后会自动重新加载，无需重启服务：

1. 自动检测配置文件变化
2. 重新加载 Cookie 配置
3. 重置警告状态
4. 重新评估过期状态

## 🔔 通知机制

### 预警通知内容

```
标题: ⚠️ Cookie 即将过期

内容:
⚠️ Cookie 即将过期

您的理想汽车 Cookie 即将在 X 天内过期

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

### 通知触发逻辑

```
┌─────────────────┐
│  配置加载完成    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  计算过期时间    │
│  和剩余时间      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐      No
│ 剩余时间 <= 48h? │──────────► 无需操作
└────────┬────────┘
         │ Yes
         ▼
┌─────────────────┐      Yes
│  已发送过预警?   │──────────► 跳过（避免重复）
└────────┬────────┘
         │ No
         ▼
┌─────────────────┐
│  发送预警通知    │
│  标记已预警      │
└─────────────────┘
```

## 🧪 测试方案

### 测试脚本

创建了 `scripts/test/test-cookie-expiration.sh` 测试脚本，包含以下测试场景：

1. **场景1**: Cookie 正常（刚更新）
   - 设置: 当前时间，有效期 7 天
   - 预期: 🟢 正常

2. **场景2**: Cookie 即将过期
   - 设置: 5天前更新，有效期 7 天（剩余 2 天）
   - 预期: ⚠️ 即将过期，发送预警

3. **场景3**: Cookie 已过期
   - 设置: 8天前更新，有效期 7 天
   - 预期: ❌ 已过期

4. **场景4**: 长有效期
   - 设置: 当前时间，有效期 30 天
   - 预期: 🟢 正常

### 测试方法

```bash
# 运行测试脚本
./scripts/test/test-cookie-expiration.sh
```

测试脚本会：
1. 备份当前配置
2. 依次测试各个场景
3. 恢复原始配置

## 📊 技术实现细节

### 时间计算

```go
// 过期时间 = 更新时间 + 有效期
expirationTime := m.CookieUpdatedAt.Add(
    time.Duration(m.CookieValidDays) * 24 * time.Hour
)

// 剩余时间 = 过期时间 - 当前时间
remainingDuration := time.Until(expirationTime)

// 预警阈值: 48 小时
warningThreshold := 48 * time.Hour

// 判断逻辑
needsWarning := remainingDuration <= warningThreshold && 
                remainingDuration > 0 && 
                !m.CookieExpirationWarned
```

### 状态格式化

```go
// 时间转换
days := remainingDuration.Hours() / 24
hours := remainingDuration.Hours()

// 状态图标
expired := "❌ 已过期"
warning := "⚠️ 即将过期"
normal := "🟢 正常"
```

### 并发安全

配置加载使用 `sync.RWMutex` 保证线程安全：

```go
type Monitor struct {
    configMu sync.RWMutex
    // ...
}

func (m *Monitor) loadConfig() {
    m.configMu.Lock()
    defer m.configMu.Unlock()
    // 配置加载逻辑
}
```

## 🔄 与现有功能的关系

### Cookie 管理体系

```
┌────────────────────────────────────────────────┐
│           Cookie 管理完整体系                   │
├────────────────────────────────────────────────┤
│                                                │
│  1. 过期预警 (NEW)                             │
│     - 提前 48 小时预警                          │
│     - 定期检查 (每日凌晨1点)                    │
│     - 启动时状态显示                            │
│                                                │
│  2. 失效检测 (现有)                             │
│     - 实时检测 HTTP 401/403                    │
│     - 业务错误码检测                            │
│     - 3次失败告警机制                           │
│                                                │
│  3. 热加载支持                                  │
│     - 配置文件变化自动加载                      │
│     - 无需重启服务                              │
│     - 立即生效                                  │
│                                                │
└────────────────────────────────────────────────┘
```

### 防御层次

1. **第一层 - 过期预警** (预防)
   - 提前 48 小时提醒
   - 给用户充足的更新时间

2. **第二层 - 失效检测** (兜底)
   - 实时检测失效
   - 即使忘记更新也能及时发现

3. **第三层 - 热加载** (便捷)
   - 更新后立即生效
   - 无需重启服务

## 📈 优化建议

### 已实现的优化

✅ 避免重复通知 - 使用 `CookieExpirationWarned` 标志  
✅ 详细的更新指导 - 通知中包含完整步骤  
✅ 状态可见性 - 启动时显示状态  
✅ 灵活配置 - 支持自定义有效期  
✅ 热加载支持 - 配置更新自动生效  

### 未来可能的增强

💡 支持邮件通知  
💡 添加 Cookie 更新历史记录  
💡 提供 Web UI 查看状态  
💡 支持多用户 Cookie 管理  

## 📚 相关文档

- [Cookie 过期管理指南](./COOKIE_EXPIRATION_WARNING.md) - 用户使用指南
- [Cookie 自动续期分析](./COOKIE_AUTO_RENEWAL_ANALYSIS.md) - 可行性分析
- [Cookie 管理机制](./COOKIE_MANAGEMENT.md) - 失效检测机制
- [配置热加载](./CONFIG_HOT_RELOAD.md) - 热加载技术实现

## 🎯 总结

Cookie 过期预警功能通过**预防性监控**的方式，有效避免了因 Cookie 失效导致的监控中断：

1. **提前预警**: 48 小时提前量，给用户充足时间
2. **自动化**: 无需人工干预，系统自动检查
3. **用户友好**: 详细的更新步骤和指导
4. **可靠性**: 多层防御，配合失效检测形成完整保障

虽然无法实现完全自动化的 Cookie 续期（需要短信验证码），但通过智能预警机制，我们将用户的操作负担降到了最低，同时保持了良好的用户体验和系统可靠性。

---

**实现完成**: ✅  
**测试通过**: ✅  
**文档完善**: ✅  
