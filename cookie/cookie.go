package cookie

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"lixiang-monitor/utils"
)

// CookieExpiredError Cookie 失效错误
type CookieExpiredError struct {
	StatusCode int
	Message    string
}

func (e *CookieExpiredError) Error() string {
	return fmt.Sprintf("Cookie 已失效 (状态码: %d): %s", e.StatusCode, e.Message)
}

// Manager Cookie 管理器
type Manager struct {
	Cookies                   string
	Headers                   map[string]string
	ValidDays                 int
	UpdatedAt                 time.Time
	ExpirationWarned          bool
	ConsecutiveFailure        int
	ExpiredNotified           bool
	LastCheckTime             time.Time
	OnCookieExpired           func(statusCode int, message string)
	OnCookieExpirationWarning func(timeDesc, expireTime, updatedAt string, ageInDays float64)
}

// NewManager 创建 Cookie 管理器
func NewManager(cookies string, headers map[string]string, validDays int, updatedAt time.Time) *Manager {
	return &Manager{
		Cookies:   cookies,
		Headers:   headers,
		ValidDays: validDays,
		UpdatedAt: updatedAt,
	}
}

// FetchOrderData 获取订单数据
func (cm *Manager) FetchOrderData(orderID string) (interface{}, error) {
	url := fmt.Sprintf("https://api-web.lixiang.com/vehicle-api/v1-0/orders/pointer/vehicleOrderDetail_PC/%s", orderID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	for key, value := range cm.Headers {
		req.Header.Set(key, value)
	}

	// 设置 Cookie
	if cm.Cookies != "" {
		req.Header.Set("Cookie", cm.Cookies)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检测 Cookie 失效的常见状态码
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		cm.handleExpired(resp.StatusCode, string(body))
		return nil, &CookieExpiredError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API 返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var orderResp map[string]interface{}
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v", err)
	}

	// 检查业务层错误码（理想汽车可能返回 200 但 code != 0）
	if code, ok := orderResp["code"].(float64); ok && code != 0 {
		message := ""
		if msg, ok := orderResp["message"].(string); ok {
			message = msg
		}

		// 常见的认证失败错误码
		if int(code) == 401 || int(code) == 403 ||
			int(code) == 10001 || int(code) == 10002 {
			cm.handleExpired(int(code), message)
			return nil, &CookieExpiredError{
				StatusCode: int(code),
				Message:    message,
			}
		}
		return nil, fmt.Errorf("API 返回业务错误: code=%d, message=%s", int(code), message)
	}

	// 请求成功，重置失败计数器
	cm.ConsecutiveFailure = 0
	cm.ExpiredNotified = false
	cm.LastCheckTime = time.Now()

	return orderResp, nil
}

// CheckExpiration 检查 Cookie 是否即将过期
func (cm *Manager) CheckExpiration() {
	if cm.ValidDays == 0 {
		return // 未配置有效期，跳过检查
	}

	// 计算 Cookie 年龄和剩余时间
	cookieAge := time.Since(cm.UpdatedAt)
	expireTime := cm.UpdatedAt.Add(time.Duration(cm.ValidDays) * 24 * time.Hour)
	remaining := time.Until(expireTime)

	// 提前 2 天开始预警（48 小时）
	warningThreshold := 48 * time.Hour

	if remaining > 0 && remaining < warningThreshold && !cm.ExpirationWarned {
		// 计算剩余天数和小时数
		remainingDays := int(remaining.Hours() / 24)
		remainingHours := int(remaining.Hours()) % 24

		var timeDesc string
		if remainingDays > 0 {
			timeDesc = fmt.Sprintf("%d 天 %d 小时", remainingDays, remainingHours)
		} else {
			timeDesc = fmt.Sprintf("%d 小时", remainingHours)
		}

		if cm.OnCookieExpirationWarning != nil {
			cm.OnCookieExpirationWarning(
				timeDesc,
				expireTime.Format(utils.DateTimeFormat),
				cm.UpdatedAt.Format(utils.DateTimeFormat),
				cookieAge.Hours()/24,
			)
			cm.ExpirationWarned = true
			log.Printf("✅ Cookie 过期预警通知已发送（剩余: %s）", timeDesc)
		}
	} else if remaining < 0 {
		// Cookie 已过期
		if !cm.ExpirationWarned {
			log.Printf("⚠️  Cookie 已过期 %s", time.Since(expireTime))
		}
	} else if remaining > warningThreshold && cm.ExpirationWarned {
		// Cookie 已更新，重置预警状态
		cm.ExpirationWarned = false
	}
}

// GetStatus 获取 Cookie 状态信息
func (cm *Manager) GetStatus() string {
	if cm.ValidDays == 0 {
		return "未配置过期检测"
	}

	expireTime := cm.UpdatedAt.Add(time.Duration(cm.ValidDays) * 24 * time.Hour)
	remaining := time.Until(expireTime)

	if remaining < 0 {
		return fmt.Sprintf("❌ 已过期 %s", time.Since(expireTime).Round(time.Hour))
	} else if remaining < 24*time.Hour {
		return fmt.Sprintf("⚠️  即将过期（剩余 %d 小时）", int(remaining.Hours()))
	} else if remaining < 48*time.Hour {
		return fmt.Sprintf("⚠️  即将过期（剩余 %.1f 天）", remaining.Hours()/24)
	} else {
		return fmt.Sprintf("🟢 正常（剩余 %.1f 天）", remaining.Hours()/24)
	}
}

// handleExpired 处理 Cookie 失效的情况
func (cm *Manager) handleExpired(statusCode int, message string) {
	cm.ConsecutiveFailure++

	log.Printf("⚠️  Cookie 验证失败 (状态码: %d, 连续失败: %d 次): %s",
		statusCode, cm.ConsecutiveFailure, message)

	// 连续失败 3 次且未通知过，则发送告警
	if cm.ConsecutiveFailure >= 3 && !cm.ExpiredNotified {
		if cm.OnCookieExpired != nil {
			cm.OnCookieExpired(statusCode, message)
			cm.ExpiredNotified = true
			log.Println("✅ Cookie 失效通知已发送")
		}
	}
}

// ResetFailureCount 重置失败计数器
func (cm *Manager) ResetFailureCount() {
	cm.ConsecutiveFailure = 0
	cm.ExpiredNotified = false
}

// UpdateCookie 更新 Cookie
func (cm *Manager) UpdateCookie(cookies string, headers map[string]string) {
	cm.Cookies = cookies
	cm.Headers = headers
	cm.UpdatedAt = time.Now()
	cm.ExpirationWarned = false
	cm.ConsecutiveFailure = 0
	cm.ExpiredNotified = false
}
