# ServerChan 配置指南

## 什么是 ServerChan

ServerChan（Server酱）是一个免费的微信推送服务，可以通过简单的API调用向微信发送消息通知。相比微信群机器人，ServerChan 更加简单易用，无需创建群聊即可接收通知。

## 配置步骤

### 1. 注册 ServerChan 账号

1. 访问 ServerChan 官网：https://sct.ftqq.com/
2. 使用微信扫码登录
3. 完成账号注册

### 2. 获取 SendKey

1. 登录后点击左侧菜单的"发送通道"
2. 在页面中找到你的 SendKey，格式类似：`SCT123456T`
3. 复制这个 SendKey

### 3. 配置到程序中

将获取到的 SendKey 填入 `config.yaml` 文件：

```yaml
# ServerChan 配置
serverchan_sendkey: "SCT123456T"  # 替换为你的实际 SendKey
serverchan_baseurl: "https://sctapi.ftqq.com/"
```

### 4. 测试通知

可以使用以下 curl 命令测试 ServerChan 是否正常工作：

```bash
curl -X POST "https://sctapi.ftqq.com/SCT123456T.send" \
     -d "title=测试消息" \
     -d "desp=理想汽车监控机器人已配置完成"
```

记得将 `SCT123456T` 替换为你的实际 SendKey。

## 通知消息格式

程序会发送以下类型的通知：

### 初始启动通知
- **标题**: 🚗 理想汽车订单监控已启动
- **内容**: 
  ```
  订单号: 177971759268550919
  当前预计交付时间: 2024-11-15
  ```

### 交付时间变更通知
- **标题**: 🚗 理想汽车交付时间更新通知
- **内容**: 
  ```
  订单号: 177971759268550919
  原预计时间: 2024-11-15
  新预计时间: 2024-11-20
  变更时间: 2024-10-15 14:30:25
  ```

## 多通道支持

程序支持同时配置 ServerChan 和微信群机器人，会向所有配置的通道发送通知：

```yaml
# 可以同时配置多种通知方式
wechat_webhook_url: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_KEY"
serverchan_sendkey: "SCT123456T"
```

## 优势对比

### ServerChan 优势
- ✅ 配置简单，只需一个 SendKey
- ✅ 无需创建微信群
- ✅ 直接发送到个人微信
- ✅ 支持 Markdown 格式
- ✅ 免费使用

### 微信群机器人优势
- ✅ 支持群组通知
- ✅ 可以 @特定成员
- ✅ 企业微信官方支持
- ✅ 更适合团队协作

## 注意事项

1. **免费限制**: ServerChan 免费版每日有发送限制，具体请查看官网说明
2. **网络要求**: 需要能够访问 `sctapi.ftqq.com` 域名
3. **微信绑定**: 确保你的微信已正确绑定到 ServerChan 账号
4. **消息格式**: ServerChan 支持 Markdown 格式，可以发送更丰富的消息内容

## 故障排除

### SendKey 不工作
1. 检查 SendKey 格式是否正确
2. 确认微信已正确绑定
3. 检查是否超过免费发送限制

### 收不到消息
1. 检查微信的"服务通知"是否被屏蔽
2. 确认 ServerChan 服务状态正常
3. 检查网络连接是否正常

### API 调用失败
1. 检查 SendKey 是否有效
2. 确认 API 地址配置正确
3. 查看程序日志获取详细错误信息