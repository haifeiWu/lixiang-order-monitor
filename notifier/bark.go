package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// BarkNotifier Bark 推送通知器
type BarkNotifier struct {
	ServerURL string
	Sound     string
	Icon      string
	Group     string
}

// Send 实现 Notifier 接口
func (bark *BarkNotifier) Send(title, content string) error {
	if bark.ServerURL == "" {
		return fmt.Errorf("Bark Server URL 未配置")
	}

	// 构建 Bark 推送数据
	barkData := map[string]interface{}{
		"title": title,
		"body":  content,
	}

	// 添加可选参数
	if bark.Sound != "" {
		barkData["sound"] = bark.Sound
	} else {
		barkData["sound"] = "minuet"
	}

	if bark.Icon != "" {
		barkData["icon"] = bark.Icon
	}

	if bark.Group != "" {
		barkData["group"] = bark.Group
	} else {
		barkData["group"] = "lixiang-monitor"
	}

	// 序列化 JSON
	jsonData, err := json.Marshal(barkData)
	if err != nil {
		return fmt.Errorf("Bark 序列化消息失败: %v", err)
	}

	// 发送 POST 请求
	resp, err := http.Post(bark.ServerURL, "application/json; charset=utf-8", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Bark 发送失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Bark 返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	log.Println("Bark 推送通知发送成功")
	return nil
}
