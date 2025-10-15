# 微信机器人配置指南

## 如何在微信群中添加机器人

### 1. 创建群聊机器人

1. 在 PC 端微信中，右键点击群聊
2. 选择"添加群机器人"
3. 创建新的机器人，设置机器人名称（如"理想汽车交付监控"）
4. 复制生成的 Webhook URL

### 2. 获取 Webhook URL

Webhook URL 的格式通常为：
```
https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_WEBHOOK_KEY
```

### 3. 配置到程序中

将获取到的 Webhook URL 填入 `config.yaml` 文件：

```yaml
wechat_webhook_url: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_WEBHOOK_KEY"
```

### 4. 测试通知

可以使用以下 curl 命令测试微信机器人是否正常工作：

```bash
curl -X POST "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_WEBHOOK_KEY" \
     -H "Content-Type: application/json" \
     -d '{
         "msgtype": "text",
         "text": {
             "content": "测试消息：理想汽车监控机器人已配置完成"
         }
     }'
```

## 通知消息格式

程序会发送以下类型的通知：

### 初始启动通知
```
🚗 理想汽车订单监控已启动
订单号: 177971759268550919
当前预计交付时间: 2024-11-15
```

### 交付时间变更通知
```
🚗 理想汽车交付时间更新通知
订单号: 177971759268550919
原预计时间: 2024-11-15
新预计时间: 2024-11-20
变更时间: 2024-10-15 14:30:25
```

## 注意事项

1. **安全性**：请不要将 Webhook URL 分享给不相关的人员
2. **频率限制**：微信机器人有发送频率限制，建议监控间隔不要设置过短
3. **群成员**：只有群成员才能看到机器人发送的消息
4. **权限**：确保机器人在群中有发送消息的权限