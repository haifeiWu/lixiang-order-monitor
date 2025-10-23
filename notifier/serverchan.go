package notifier

import (
"fmt"
"io"
"log"
"net/http"
"net/url"
)

// ServerChanNotifier ServerChan 通知器
type ServerChanNotifier struct {
	SendKey string
	BaseURL string
}

// Send 实现 Notifier 接口
func (sc *ServerChanNotifier) Send(title, content string) error {
	if sc.SendKey == "" {
		return fmt.Errorf("ServerChan SendKey 未配置")
	}

	// 构建请求数据
	data := url.Values{}
	data.Set("title", title)
	data.Set("desp", content)

	// 构建正确的 ServerChan API URL
	apiURL := sc.BaseURL + sc.SendKey + ".send"

	// 发送请求
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("ServerChan 发送失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ServerChan 返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	log.Println("ServerChan 通知发送成功")
	return nil
}
