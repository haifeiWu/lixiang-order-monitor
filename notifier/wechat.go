package notifier

import (
"bytes"
"encoding/json"
"fmt"
"io"
"log"
"net/http"
)

// WeChatWebhookNotifier 微信群机器人通知器
type WeChatWebhookNotifier struct {
	WebhookURL string
}

// WeChatMessage 微信消息结构
type WeChatMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

// Send 实现 Notifier 接口
func (wc *WeChatWebhookNotifier) Send(title, content string) error {
	if wc.WebhookURL == "" {
		return fmt.Errorf("微信 Webhook URL 未配置")
	}

	// 组合标题和内容
	message := title
	if content != "" {
		message += "\n\n" + content
	}

	wechatMsg := WeChatMessage{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: message,
		},
	}

	jsonData, err := json.Marshal(wechatMsg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	resp, err := http.Post(wc.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("发送微信通知失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("微信通知返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	log.Println("微信群机器人通知发送成功")
	return nil
}
