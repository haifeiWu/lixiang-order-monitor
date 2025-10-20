# Cookie 失效检测与处理机制

## 📋 概述

本文档详细说明理想汽车订单监控系统的 Cookie 失效检测和自动处理机制。

## 🔍 Cookie 失效检测

### 检测机制

系统会在以下情况下检测 Cookie 失效：

1. **HTTP 状态码检测**
   - `401 Unauthorized`: 未授权访问
   - `403 Forbidden`: 禁止访问

2. **业务错误码检测**
   - `code: 401/403`: 认证失败
   - `code: 10001/10002`: 理想汽车特定的认证错误

3. **连续失败计数**
   - 记录连续失败次数
   - 失败 3 次后自动发送告警通知

### 关键字段说明

理想汽车使用的关键 Cookie 字段：

| 字段名 | 说明 | 重要性 | 过期特性 |
|--------|------|--------|----------|
| `X-LX-Token` | 会话令牌 | ⭐⭐⭐ | 定期过期（最易失效） |
| `authli_device_id` | 设备标识 | ⭐⭐⭐ | 相对稳定 |
| `X-LX-Deviceid` | 设备 ID | ⭐⭐ | 相对稳定 |
| `X-LX-HeaderData` | 元数据 | ⭐ | 可选 |

## 🚨 失效告警

### 告警触发条件

- 连续失败 **3 次**
- 每次 Cookie 失效只发送 **1 次通知**（避免通知轰炸）

### 告警内容

通知消息包含以下信息：

```
🚨 理想汽车 Cookie 已失效

失效详情：
- 状态码：401
- 错误信息：Unauthorized
- 连续失败次数：3 次
- 检测时间：2025-01-17 15:30:00

Cookie 更新步骤：
1. 打开浏览器访问 https://www.lixiang.com/
2. 登录您的理想汽车账号
3. 按 F12 打开开发者工具
4. 切换到 Network 标签
5. 刷新页面，找到任意请求
6. 在请求头中复制完整的 Cookie 字符串
7. 更新 config.yaml 中的 lixiang_cookies 字段

关键 Cookie 字段说明：
- X-LX-Token: 会话令牌（必需，定期过期）
- authli_device_id: 设备标识（必需）
- X-LX-Deviceid: 设备 ID（必需）
```

## 🔧 Cookie 更新步骤

### 方法 1: 浏览器开发者工具（推荐）

1. **打开理想汽车官网**
   ```
   https://www.lixiang.com/
   ```

2. **登录账号**
   - 使用您的理想汽车账号登录

3. **打开开发者工具**
   - Windows/Linux: 按 `F12` 或 `Ctrl+Shift+I`
   - macOS: 按 `Cmd+Option+I`

4. **获取 Cookie**
   - 切换到 `Network` (网络) 标签
   - 刷新页面 (F5)
   - 点击任意请求（推荐选择 API 请求）
   - 在右侧找到 `Request Headers` (请求头)
   - 复制 `Cookie:` 后面的完整字符串

5. **更新配置文件**
   ```yaml
   # config.yaml
   lixiang_cookies: "X-LX-Deviceid=xxx; X-LX-Token=xxx; authli_device_id=xxx; ..."
   ```

6. **验证更新**
   - 配置热加载会自动生效（无需重启）
   - 查看日志确认 Cookie 已更新

### 方法 2: 浏览器控制台

在理想汽车网站的控制台中执行：

```javascript
document.cookie
```

复制输出的完整 Cookie 字符串到配置文件。

### 方法 3: 使用浏览器扩展

推荐使用 Chrome 扩展：
- **EditThisCookie**
- **Cookie-Editor**

## 🔄 热加载支持

Cookie 更新支持配置热加载：

1. **修改配置文件**
   ```yaml
   lixiang_cookies: "新的 Cookie 字符串"
   ```

2. **自动检测并加载**
   - 程序自动检测 `config.yaml` 变化
   - 1 秒内应用新配置
   - 无需重启程序

3. **验证生效**
   ```
   2025/01/17 15:30:00 main.go:488: 📝 配置已更新 (版本: 2)
   2025/01/17 15:30:00 main.go:489: 📢 发送配置更新通知...
   ```

## 📊 监控状态

### 成功状态

```
2025/01/17 15:00:00 main.go:670: 当前预计交付时间: 2025-12-01
2025/01/17 15:00:00 main.go:672: ✅ Cookie 验证成功
```

### 失效状态

```
2025/01/17 15:00:00 main.go:636: ⚠️  Cookie 验证失败 (状态码: 401, 连续失败: 1 次): Unauthorized
2025/01/17 15:15:00 main.go:636: ⚠️  Cookie 验证失败 (状态码: 401, 连续失败: 2 次): Unauthorized
2025/01/17 15:30:00 main.go:636: ⚠️  Cookie 验证失败 (状态码: 401, 连续失败: 3 次): Unauthorized
2025/01/17 15:30:00 main.go:660: ✅ Cookie 失效通知已发送
```

## 🛡️ 安全建议

### Cookie 保护

1. **不要分享 Cookie**
   - Cookie 包含您的登录凭证
   - 不要在公开场合泄露

2. **定期更新**
   - 建议每周检查 Cookie 有效性
   - 发现失效立即更新

3. **配置文件安全**
   - 不要将 `config.yaml` 提交到版本控制系统
   - 使用 `.gitignore` 排除配置文件

### 最佳实践

```bash
# .gitignore
config.yaml
lixiang-monitor
*.log
```

## 🔧 故障排查

### 问题 1: Cookie 频繁失效

**症状**：
- 每天都需要更新 Cookie
- 连续失败通知频繁

**原因**：
- 理想汽车会话过期策略较严格
- 使用了临时会话（未勾选"记住我"）

**解决方案**：
1. 登录时勾选"记住我"
2. 使用常用浏览器登录
3. 避免频繁更换设备

### 问题 2: 更新后仍然失效

**症状**：
- 更新 Cookie 后仍提示失效

**排查步骤**：
1. 确认复制了完整的 Cookie 字符串
2. 检查是否包含 `X-LX-Token` 字段
3. 确认配置文件格式正确（YAML 格式）
4. 查看程序日志确认配置已重新加载

### 问题 3: 没有收到失效通知

**症状**：
- Cookie 失效但没有收到通知

**排查步骤**：
1. 检查通知器配置（ServerChan 或 WeChat）
2. 查看日志确认是否达到 3 次失败阈值
3. 验证通知器 SendKey 或 Webhook URL 有效

## 📈 技术实现

### 数据结构

```go
type Monitor struct {
    // Cookie 管理相关
    LastCookieCheckTime      time.Time // 上次检查时间
    CookieExpiredNotified    bool      // 是否已通知
    ConsecutiveCookieFailure int       // 连续失败次数
    
    // 配置热加载
    mu sync.RWMutex // 读写锁
}

type CookieExpiredError struct {
    StatusCode int
    Message    string
}
```

### 核心逻辑

```go
// 1. Cookie 失效检测
if resp.StatusCode == 401 || resp.StatusCode == 403 {
    m.handleCookieExpired(resp.StatusCode, string(body))
    return nil, &CookieExpiredError{...}
}

// 2. 失败计数与告警
func (m *Monitor) handleCookieExpired(statusCode int, message string) {
    m.ConsecutiveCookieFailure++
    
    if m.ConsecutiveCookieFailure >= 3 && !m.CookieExpiredNotified {
        m.sendNotification("Cookie 已失效", "...")
        m.CookieExpiredNotified = true
    }
}

// 3. 成功后重置
m.ConsecutiveCookieFailure = 0
m.CookieExpiredNotified = false
```

## 🎯 总结

本系统提供了完善的 Cookie 失效检测和处理机制：

- ✅ **自动检测**：多维度检测 Cookie 失效
- ✅ **智能告警**：连续失败 3 次后发送通知
- ✅ **热加载支持**：更新 Cookie 无需重启
- ✅ **详细指引**：告警消息包含详细更新步骤
- ✅ **状态追踪**：记录失败次数和检查时间

确保您的订单监控服务持续稳定运行！
