# Bark 推送配置指南

## 📱 什么是 Bark？

[Bark](https://github.com/Finb/Bark) 是一款开源的 iOS/macOS 推送通知应用，支持：
- 📲 即时推送到 iPhone/iPad/Mac
- 🔔 自定义推送声音和图标
- 🏠 支持自建服务器，数据私有
- 🆓 完全免费使用

## 🎯 配置步骤

### 方式一：使用官方服务器（推荐新手）

#### 1. 安装 Bark App

在 App Store 搜索并下载 **Bark** 应用。

#### 2. 获取推送 URL

打开 Bark App，会自动生成一个推送 URL，格式类似：

```
https://api.day.app/your_device_key
```

或者点击"添加服务器"可以获取完整的推送地址。

#### 3. 配置到项目

在 `config.yaml` 中添加：

```yaml
# Bark 推送配置
bark_server_url: "https://api.day.app/your_device_key"
bark_sound: "minuet"              # 推送声音
bark_group: "lixiang-monitor"      # 通知分组
```

### 方式二：使用自建服务器

#### 1. 部署 Bark 服务器

可以使用 Docker 快速部署：

```bash
docker run -d --name bark-server \
  -p 8080:8080 \
  finab/bark-server
```

或者使用 Docker Compose：

```yaml
version: '3'
services:
  bark:
    image: finab/bark-server
    ports:
      - "8080:8080"
    restart: unless-stopped
```

#### 2. 在 Bark App 中添加服务器

打开 Bark App → 设置 → 添加服务器：
- 服务器地址：`http://your_server_ip:8080`

添加后会生成一个 key，例如：`mruLduP9zEnDabDA4uBVLj`

#### 3. 配置到项目

```yaml
# Bark 推送配置（自建服务器）
bark_server_url: "http://your_server_ip:8080/mruLduP9zEnDabDA4uBVLj"
bark_sound: "minuet"
bark_group: "lixiang-monitor"
```

## ⚙️ 配置参数详解

### 必填参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `bark_server_url` | Bark 推送完整 URL | `http://43.128.109.177:8080/mruLduP9zEnDabDA4uBVLj` |

### 可选参数

| 参数 | 说明 | 默认值 | 可选值 |
|------|------|--------|--------|
| `bark_sound` | 推送提示音 | `minuet` | `alarm`, `bell`, `birdsong`, `bloom`, `calypso`, `chime`, `choo`, `descent`, `electronic`, `fanfare`, `glass`, `gotosleep`, `healthnotification`, `horn`, `ladder`, `mailsent`, `minuet`, `multiwayinvitation`, `newmail`, `newsflash`, `noir`, `paymentsuccess`, `shake`, `sherwoodforest`, `silence`, `spell`, `suspense`, `telegraph`, `tiptoes`, `typewriters`, `update` |
| `bark_icon` | 推送图标 URL | 空 | 任意图片 URL |
| `bark_group` | 通知分组 | `lixiang-monitor` | 任意字符串 |

### 完整配置示例

```yaml
# Bark 推送配置
bark_server_url: "http://43.128.109.177:8080/mruLduP9zEnDabDA4uBVLj"
bark_sound: "minuet"                        # 推送声音
bark_icon: "https://www.lixiang.com/favicon.ico"  # 自定义图标
bark_group: "lixiang-monitor"                # 通知分组
```

## 🧪 测试推送

### 方法 1: 使用测试脚本

创建测试脚本 `test-bark.sh`：

```bash
#!/bin/bash

BARK_URL="http://43.128.109.177:8080/mruLduP9zEnDabDA4uBVLj"

curl --location "$BARK_URL" \
--header 'Content-Type: application/json; charset=utf-8' \
--data '{
  "title": "理想汽车订单监控",
  "body": "Bark 推送测试成功！",
  "sound": "minuet",
  "group": "lixiang-monitor"
}'
```

运行测试：

```bash
chmod +x test-bark.sh
./test-bark.sh
```

### 方法 2: 使用通知测试脚本

如果已经配置好，可以使用项目的测试脚本：

```bash
./scripts/test/test-notification.sh
```

## 🎨 推送声音列表

Bark 支持多种内置声音，以下是常用的：

| 声音名称 | 描述 | 适用场景 |
|---------|------|---------|
| `alarm` | 闹钟声 | 紧急通知 |
| `bell` | 铃声 | 重要提醒 |
| `minuet` | 小步舞曲（默认） | 日常通知 |
| `chime` | 钟声 | 温和提醒 |
| `glass` | 玻璃声 | 简洁通知 |
| `horn` | 号角声 | 重要消息 |
| `silence` | 静音 | 仅显示不发声 |

完整声音列表请参考 [Bark 文档](https://github.com/Finb/Bark/blob/master/Sounds.md)。

## 🔍 故障排查

### 问题 1: 推送未收到

**检查项**:
1. ✅ 确认 `bark_server_url` 配置正确
2. ✅ 确认 Bark App 已安装并添加了对应服务器
3. ✅ 检查网络连接（服务器是否可访问）
4. ✅ 查看程序日志是否有错误信息

**测试方法**:
```bash
# 直接测试推送 URL
curl -X POST "http://your_server:8080/your_key" \
  -H "Content-Type: application/json" \
  -d '{"title":"测试","body":"测试消息"}'
```

### 问题 2: 推送有延迟

**可能原因**:
- 网络延迟
- 使用国外服务器
- APNs 推送延迟

**解决方案**:
- 使用自建服务器（国内）
- 检查网络连接质量

### 问题 3: 推送声音不对

**检查**:
```yaml
bark_sound: "minuet"  # 确认声音名称正确（小写）
```

常见拼写错误：
- ❌ `Minuet` → ✅ `minuet`
- ❌ `Alarm` → ✅ `alarm`

### 问题 4: 自建服务器连接失败

**检查项**:
1. 服务器是否正常运行：`docker ps`
2. 端口是否开放：`telnet your_server 8080`
3. 防火墙规则是否正确
4. URL 格式是否正确（包含 key）

## 💡 使用技巧

### 技巧 1: 分组管理通知

使用不同的 `bark_group` 对通知进行分类：

```yaml
bark_group: "lixiang-urgent"     # 紧急通知
bark_group: "lixiang-daily"      # 日常通知
```

在 Bark App 中可以按分组查看。

### 技巧 2: 自定义图标

使用理想汽车的 Logo 作为推送图标：

```yaml
bark_icon: "https://www.lixiang.com/favicon.ico"
```

### 技巧 3: 设置推送优先级

不同类型通知使用不同声音：

```yaml
# Cookie 过期预警 - 使用紧急声音
bark_sound: "alarm"

# 日常订单更新 - 使用温和声音
bark_sound: "minuet"
```

### 技巧 4: 多设备推送

可以在 `config.yaml` 中配置多个 Bark URL：

```yaml
# 虽然当前只支持单个 Bark 配置，
# 但可以结合其他通知渠道实现多端推送
bark_server_url: "http://server/key1"    # iPhone
wechat_webhook_url: "https://..."        # 微信
serverchan_sendkey: "SCT..."             # Server酱
```

## 🔗 相关链接

- [Bark GitHub](https://github.com/Finb/Bark)
- [Bark Server GitHub](https://github.com/Finb/bark-server)
- [Bark 声音列表](https://github.com/Finb/Bark/blob/master/Sounds.md)
- [项目通知测试指南](./TESTING_GUIDE.md)

## 📊 对比其他通知方式

| 特性 | Bark | 微信机器人 | ServerChan |
|------|------|-----------|-----------|
| **平台** | iOS/macOS | 微信 | 微信 |
| **实时性** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| **部署难度** | ⭐⭐ | ⭐ | ⭐ |
| **自定义性** | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐ |
| **费用** | 免费 | 免费 | 免费版限额 |
| **推送数量限制** | 无限制 | 无限制 | 免费版5次/天 |

**推荐方案**: 
- iOS/Mac 用户：优先使用 **Bark**
- 微信用户：配置 **微信机器人** 作为备用
- 多渠道：同时配置多种通知方式，确保不漏消息

## 🎯 总结

Bark 推送的优势：
- ✅ 推送速度快，几乎实时
- ✅ 无推送数量限制
- ✅ 支持自建服务器，数据私有
- ✅ 自定义程度高（声音、图标、分组）
- ✅ 完全免费开源

立即配置 Bark，享受更好的推送体验！

---

> 💡 **提示**: 建议同时配置 Bark 和其他通知渠道（如微信机器人），实现消息多渠道备份，确保不错过重要通知。
