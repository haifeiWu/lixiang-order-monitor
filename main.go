package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

// 时间格式常量
const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateTimeShort  = "2006-01-02 15:04"
	DateFormat     = "2006-01-02"
)

// 理想汽车订单响应结构
type OrderResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Delivery struct {
			EstimateDeliveringAt string `json:"estimateDeliveringAt"`
		} `json:"delivery"`
	} `json:"data"`
}

// 微信机器人消息结构
type WeChatMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

// ServerChan 通知结构
type ServerChanNotifier struct {
	SendKey string
	BaseURL string
}

// 通知接口
type Notifier interface {
	Send(title, content string) error
}

// 微信群机器人通知器
type WeChatWebhookNotifier struct {
	WebhookURL string
}

type Monitor struct {
	OrderID          string
	LastEstimateTime string
	CheckInterval    string
	LixiangCookies   string
	LixiangHeaders   map[string]string
	Notifiers        []Notifier
	LockOrderTime    time.Time // 锁单时间
	EstimateWeeksMin int       // 预计交付周数范围（最小）
	EstimateWeeksMax int       // 预计交付周数范围（最大）
	cron             *cron.Cron

	// 定期通知相关字段
	LastNotificationTime        time.Time     // 上次发送通知的时间
	NotificationInterval        time.Duration // 通知间隔（当交付时间未更新时）
	EnablePeriodicNotify        bool          // 是否启用定期通知
	AlwaysNotifyWhenApproaching bool          // 临近交付时总是通知

	// 配置热加载相关
	mu            sync.RWMutex // 读写锁，保护配置的并发访问
	configVersion int          // 配置版本号，用于跟踪配置变化
}

// ServerChan 通知器实现
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

// 微信群机器人通知器实现
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

// 解析锁单时间
func parseLockOrderTime(timeStr string) (time.Time, error) {
	// 支持多种时间格式
	formats := []string{
		DateTimeFormat,
		"2006/01/02 15:04:05",
		DateTimeShort,
		"2006/01/02 15:04",
		DateFormat,
		"2006/01/02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析时间格式: %s", timeStr)
}

// 计算预计交付日期范围
func (m *Monitor) calculateEstimatedDelivery() (time.Time, time.Time) {
	minDate := m.LockOrderTime.AddDate(0, 0, m.EstimateWeeksMin*7)
	maxDate := m.LockOrderTime.AddDate(0, 0, m.EstimateWeeksMax*7)
	return minDate, maxDate
}

// 基于当前时间计算剩余交付时间
func (m *Monitor) calculateRemainingDeliveryTime() (int, int, string) {
	now := time.Now()
	minDate, maxDate := m.calculateEstimatedDelivery()

	// 计算距离交付时间的天数
	daysToMin := int(minDate.Sub(now).Hours() / 24)
	daysToMax := int(maxDate.Sub(now).Hours() / 24)

	var status string
	if now.After(maxDate) {
		// 已超过预计交付时间
		overdueDays := int(now.Sub(maxDate).Hours() / 24)
		status = fmt.Sprintf("已超期 %d 天", overdueDays)
	} else if now.After(minDate) {
		// 在预计交付时间范围内
		status = "在预计交付时间范围内"
	} else if daysToMin <= 0 {
		// 今天或明天就到交付时间
		status = "即将到达交付时间"
	} else {
		// 还有若干天
		status = fmt.Sprintf("还有 %d-%d 天", daysToMin, daysToMax)
	}

	return daysToMin, daysToMax, status
}

// 计算交付进度百分比
func (m *Monitor) calculateDeliveryProgress() float64 {
	now := time.Now()

	// 计算从锁单到预计交付的总时间（取最大值）
	_, maxDate := m.calculateEstimatedDelivery()
	totalDuration := maxDate.Sub(m.LockOrderTime)

	// 计算已经过去的时间
	elapsedDuration := now.Sub(m.LockOrderTime)

	// 计算进度百分比
	progress := float64(elapsedDuration) / float64(totalDuration) * 100

	// 确保进度在 0-100% 之间
	if progress < 0 {
		progress = 0
	} else if progress > 100 {
		progress = 100
	}

	return progress
}

// 格式化交付日期范围
func (m *Monitor) formatDeliveryEstimate() string {
	minDate, maxDate := m.calculateEstimatedDelivery()
	_, _, status := m.calculateRemainingDeliveryTime()
	progress := m.calculateDeliveryProgress()

	baseInfo := ""
	if m.EstimateWeeksMin == m.EstimateWeeksMax {
		baseInfo = fmt.Sprintf("预计 %d 周后交付 (%s 左右)",
			m.EstimateWeeksMin,
			minDate.Format(DateFormat))
	} else {
		baseInfo = fmt.Sprintf("预计 %d-%d 周后交付 (%s 至 %s)",
			m.EstimateWeeksMin,
			m.EstimateWeeksMax,
			minDate.Format(DateFormat),
			maxDate.Format(DateFormat))
	}

	// 添加当前时间状态和进度信息
	now := time.Now()
	if now.Before(minDate) {
		// 还未到交付时间
		return fmt.Sprintf("%s\n📅 当前状态: %s\n📊 等待进度: %.1f%%",
			baseInfo, status, progress)
	} else if now.After(maxDate) {
		// 已超过交付时间
		return fmt.Sprintf("%s\n⚠️  当前状态: %s\n📊 进度: %.1f%% (已超期)",
			baseInfo, status, progress)
	} else {
		// 在交付时间范围内
		return fmt.Sprintf("%s\n✅ 当前状态: %s\n📊 进度: %.1f%%",
			baseInfo, status, progress)
	}
}

// 获取详细的交付时间信息
func (m *Monitor) getDetailedDeliveryInfo() string {
	now := time.Now()
	minDate, maxDate := m.calculateEstimatedDelivery()
	_, _, status := m.calculateRemainingDeliveryTime()
	progress := m.calculateDeliveryProgress()

	// 计算锁单至今的天数
	daysSinceLock := int(now.Sub(m.LockOrderTime).Hours() / 24)

	info := fmt.Sprintf("📅 锁单时间: %s (%d天前)\n",
		m.LockOrderTime.Format(DateTimeShort), daysSinceLock)

	info += fmt.Sprintf("🔮 基于锁单时间预测: %s\n", m.formatDeliveryEstimate())
	info += fmt.Sprintf("📊 当前状态: %s (进度: %.1f%%)\n", status, progress)

	// 添加具体的倒计时信息
	if now.Before(minDate) {
		daysToMin := int(minDate.Sub(now).Hours() / 24)
		daysToMax := int(maxDate.Sub(now).Hours() / 24)
		if daysToMin <= 7 {
			info += fmt.Sprintf("⏰ 距离最早交付时间: %d天\n", daysToMin)
		}
		if daysToMax <= 14 {
			info += fmt.Sprintf("⏰ 距离最晚交付时间: %d天\n", daysToMax)
		}
	}

	return info
}

// 获取交付时间智能分析报告
func (m *Monitor) getDeliveryAnalysisReport() string {
	now := time.Now()
	minDate, maxDate := m.calculateEstimatedDelivery()
	daysToMin, _, status := m.calculateRemainingDeliveryTime()
	progress := m.calculateDeliveryProgress()

	report := "📊 交付时间智能分析报告\n"
	report += "=" + strings.Repeat("=", 30) + "\n\n"

	// 基本信息
	daysSinceLock := int(now.Sub(m.LockOrderTime).Hours() / 24)
	report += fmt.Sprintf("🔐 锁单信息: %s (%d天前)\n",
		m.LockOrderTime.Format(DateTimeShort), daysSinceLock)

	report += fmt.Sprintf("📅 预计交付: %s - %s\n",
		minDate.Format(DateFormat), maxDate.Format(DateFormat))

	report += fmt.Sprintf("📈 当前进度: %.1f%%\n", progress)
	report += fmt.Sprintf("⏱️  剩余时间: %s\n\n", status)

	// 时间状态分析
	if now.Before(minDate) {
		if daysToMin <= 3 {
			report += "🚨 紧急提醒: 即将进入交付时间窗口！\n"
		} else if daysToMin <= 7 {
			report += "⚡ 重要提醒: 距离交付时间不到一周\n"
		} else if daysToMin <= 14 {
			report += "📢 提前提醒: 距离交付时间不到两周\n"
		} else {
			report += "😌 状态良好: 还有充足的等待时间\n"
		}
	} else if now.After(minDate) && now.Before(maxDate) {
		report += "🎯 关键时期: 正处于预计交付时间范围内\n"
		report += "👀 建议: 密切关注官方通知\n"
	} else if now.After(maxDate) {
		overdueDays := int(now.Sub(maxDate).Hours() / 24)
		report += "⚠️  延期状态: 已超过预计交付时间\n"
		if overdueDays <= 7 {
			report += "💡 建议: 可联系客服了解具体情况\n"
		} else {
			report += "📞 建议: 强烈建议联系客服获取最新进展\n"
		}
	}

	return report
} // 检查是否临近预计交付时间
func (m *Monitor) isApproachingDelivery() (bool, string) {
	now := time.Now()
	minDate, maxDate := m.calculateEstimatedDelivery()

	// 计算距离最早预计交付时间的天数
	daysToMin := int(minDate.Sub(now).Hours() / 24)
	daysToMax := int(maxDate.Sub(now).Hours() / 24)

	// 如果在预计交付时间范围内
	if now.After(minDate) && now.Before(maxDate) {
		return true, "当前处于预计交付时间范围内"
	}

	// 如果距离最早交付时间不到7天
	if daysToMin <= 7 && daysToMin > 0 {
		return true, fmt.Sprintf("距离最早预计交付时间还有 %d 天", daysToMin)
	}

	// 如果距离最晚交付时间不到7天
	if daysToMax <= 7 && daysToMax > 0 {
		return true, fmt.Sprintf("距离最晚预计交付时间还有 %d 天", daysToMax)
	}

	// 如果已经超过预计交付时间
	if now.After(maxDate) {
		overdueDays := int(now.Sub(maxDate).Hours() / 24)
		return true, fmt.Sprintf("已超过预计交付时间 %d 天", overdueDays)
	}

	return false, ""
}

// 检查是否应该发送定期通知
func (m *Monitor) shouldSendPeriodicNotification() bool {
	if !m.EnablePeriodicNotify {
		return false
	}

	now := time.Now()

	// 如果从未发送过通知，返回false（初始通知会单独处理）
	if m.LastNotificationTime.IsZero() {
		return false
	}

	// 检查是否已经超过通知间隔
	return now.Sub(m.LastNotificationTime) >= m.NotificationInterval
}

// 更新最后通知时间
func (m *Monitor) updateLastNotificationTime() {
	m.LastNotificationTime = time.Now()
}

// 加载或重新加载配置
func (m *Monitor) loadConfig() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 解析锁单时间
	lockOrderTimeStr := viper.GetString("lock_order_time")
	lockOrderTime, err := parseLockOrderTime(lockOrderTimeStr)
	if err != nil {
		log.Printf("锁单时间解析失败: %v, 保持当前时间", err)
		if m.LockOrderTime.IsZero() {
			lockOrderTime, _ = time.Parse(DateTimeFormat, "2025-09-27 13:08:00")
		} else {
			lockOrderTime = m.LockOrderTime
		}
	}

	// 更新配置
	m.OrderID = viper.GetString("order_id")
	m.LixiangCookies = viper.GetString("lixiang_cookies")
	m.LockOrderTime = lockOrderTime
	m.EstimateWeeksMin = viper.GetInt("estimate_weeks_min")
	m.EstimateWeeksMax = viper.GetInt("estimate_weeks_max")
	m.EnablePeriodicNotify = viper.GetBool("enable_periodic_notify")
	m.NotificationInterval = time.Duration(viper.GetInt("notification_interval_hours")) * time.Hour
	m.AlwaysNotifyWhenApproaching = viper.GetBool("always_notify_when_approaching")

	// 检查检查间隔是否变化
	newCheckInterval := viper.GetString("check_interval")
	checkIntervalChanged := false
	if newCheckInterval != m.CheckInterval && m.CheckInterval != "" {
		// 只有在非首次加载且间隔发生变化时才记录
		checkIntervalChanged = true
	}
	m.CheckInterval = newCheckInterval

	// 重新初始化通知器
	var notifiers []Notifier

	// 添加微信群机器人通知器
	wechatWebhookURL := viper.GetString("wechat_webhook_url")
	if wechatWebhookURL != "" {
		notifiers = append(notifiers, &WeChatWebhookNotifier{
			WebhookURL: wechatWebhookURL,
		})
	}

	// 添加 ServerChan 通知器
	serverChanSendKey := viper.GetString("serverchan_sendkey")
	if serverChanSendKey != "" {
		notifiers = append(notifiers, &ServerChanNotifier{
			SendKey: serverChanSendKey,
			BaseURL: viper.GetString("serverchan_baseurl"),
		})
	}

	m.Notifiers = notifiers
	m.configVersion++

	log.Printf("配置已加载，版本: %d", m.configVersion)

	// 如果检查间隔变更且 cron 已经启动，返回错误提示需要重启
	if checkIntervalChanged && m.cron != nil {
		return fmt.Errorf("检查间隔已变更，需要重启服务")
	}

	return nil
}

// 监听配置文件变化
func (m *Monitor) watchConfig() {
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("检测到配置文件变化: %s", e.Name)

		// 重新读取配置
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("重新读取配置文件失败: %v", err)
			return
		}

		// 重新加载配置
		if err := m.loadConfig(); err != nil {
			log.Printf("重新加载配置失败: %v", err)
			if err.Error() == "检查间隔已变更，需要重启服务" {
				log.Println("⚠️  检测到检查间隔变更，请手动重启服务以应用新的检查间隔")
			}
			return
		}

		log.Println("✅ 配置已成功热加载")

		// 发送配置更新通知
		title := "⚙️ 监控服务配置已更新"
		content := fmt.Sprintf("配置版本: %d\n更新时间: %s\n\n当前配置:\n订单ID: %s\n检查间隔: %s\n通知器数量: %d\n定期通知: %v\n通知间隔: %.0f小时",
			m.configVersion,
			time.Now().Format(DateTimeFormat),
			m.OrderID,
			m.CheckInterval,
			len(m.Notifiers),
			m.EnablePeriodicNotify,
			m.NotificationInterval.Hours())

		if err := m.sendNotification(title, content); err != nil {
			log.Printf("发送配置更新通知失败: %v", err)
		}
	})

	viper.WatchConfig()
	log.Println("✅ 配置文件监听已启动")
}

func NewMonitor() *Monitor {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("配置文件读取失败: %v", err)
	}

	// 设置默认值
	viper.SetDefault("order_id", "177971759268550919")
	viper.SetDefault("check_interval", "@every 30m") // 每30分钟检查一次
	viper.SetDefault("wechat_webhook_url", "")
	viper.SetDefault("serverchan_sendkey", "")
	viper.SetDefault("serverchan_baseurl", "https://sctapi.ftqq.com/")
	viper.SetDefault("lock_order_time", "2025-09-27 13:08:00")
	viper.SetDefault("estimate_weeks_min", 7)
	viper.SetDefault("estimate_weeks_max", 9)
	// 新增定期通知配置
	viper.SetDefault("enable_periodic_notify", true)         // 启用定期通知
	viper.SetDefault("notification_interval_hours", 24)      // 24小时发送一次状态通知
	viper.SetDefault("always_notify_when_approaching", true) // 临近交付时总是通知

	monitor := &Monitor{
		LastNotificationTime: time.Time{}, // 初始化为零值，表示从未发送过通知
		LixiangHeaders: map[string]string{
			"accept":             "application/json, text/plain, */*",
			"accept-language":    "en-US,en;q=0.9,zh-CN;q=0.8,zh-TW;q=0.7,zh;q=0.6",
			"origin":             "https://www.lixiang.com",
			"priority":           "u=1, i",
			"referer":            "https://www.lixiang.com/",
			"sec-ch-ua":          `"Google Chrome";v="141", "Not?A_Brand";v="8", "Chromium";v="141"`,
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua-platform": `"macOS"`,
			"sec-fetch-dest":     "empty",
			"sec-fetch-mode":     "cors",
			"sec-fetch-site":     "same-site",
			"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36",
			"x-chj-devicetype":   "1",
			"x-chj-metadata":     `{"code":"102002"}`,
			"x-chj-sourceurl":    "https://www.lixiang.com/?chjchannelcode=102002",
			"x-chj-traceid":      "75697683-7eae-0fbe-ae8e-86bfa4aab99d",
		},
		cron:          cron.New(cron.WithSeconds()),
		configVersion: 0,
	}

	// 加载初始配置
	if err := monitor.loadConfig(); err != nil {
		log.Printf("加载初始配置失败: %v", err)
	}

	// 启动配置文件监听
	monitor.watchConfig()

	if len(monitor.Notifiers) == 0 {
		log.Println("⚠️  未配置任何通知器，将不会发送通知")
	} else {
		log.Printf("✅ 已配置 %d 个通知器", len(monitor.Notifiers))
	}

	return monitor
}

func (m *Monitor) fetchOrderData() (*OrderResponse, error) {
	m.mu.RLock()
	orderID := m.OrderID
	cookies := m.LixiangCookies
	headers := m.LixiangHeaders
	m.mu.RUnlock()

	url := fmt.Sprintf("https://api-web.lixiang.com/vehicle-api/v1-0/orders/pointer/vehicleOrderDetail_PC/%s", orderID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 设置 Cookie
	if cookies != "" {
		req.Header.Set("Cookie", cookies)
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

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API 返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var orderResp OrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v", err)
	}

	return &orderResp, nil
}

func (m *Monitor) sendNotification(title, content string) error {
	m.mu.RLock()
	notifiers := m.Notifiers
	m.mu.RUnlock()

	if len(notifiers) == 0 {
		log.Println("未配置任何通知器，跳过通知")
		return nil
	}

	var errors []string
	successCount := 0

	for _, notifier := range notifiers {
		if err := notifier.Send(title, content); err != nil {
			log.Printf("通知发送失败: %v", err)
			errors = append(errors, err.Error())
		} else {
			successCount++
		}
	}

	if successCount == 0 {
		return fmt.Errorf("所有通知器发送失败: %v", errors)
	} else if len(errors) > 0 {
		log.Printf("部分通知器发送失败: %v", errors)
	}

	log.Printf("成功发送 %d/%d 个通知", successCount, len(notifiers))
	return nil
}

func (m *Monitor) checkDeliveryTime() {
	log.Println("开始检查订单交付时间...")

	orderData, err := m.fetchOrderData()
	if err != nil {
		log.Printf("获取订单数据失败: %v", err)
		return
	}

	if orderData.Code != 0 {
		log.Printf("API 返回错误: %s", orderData.Message)
		return
	}

	currentEstimateTime := orderData.Data.Delivery.EstimateDeliveringAt
	log.Printf("当前预计交付时间: %s", currentEstimateTime)

	// 读取配置（加锁保护）
	m.mu.RLock()
	orderID := m.OrderID
	lockOrderTime := m.LockOrderTime
	lastEstimateTime := m.LastEstimateTime
	m.mu.RUnlock()

	// 计算基于锁单时间的预测
	predictedDelivery := m.formatDeliveryEstimate()
	isApproaching, approachMsg := m.isApproachingDelivery()

	log.Printf("锁单时间: %s", lockOrderTime.Format(DateTimeFormat))
	log.Printf("基于锁单时间预测: %s", predictedDelivery)
	if isApproaching {
		log.Printf("交付提醒: %s", approachMsg)
	}

	// 如果是第一次检查，记录当前时间
	if lastEstimateTime == "" {
		m.mu.Lock()
		m.LastEstimateTime = currentEstimateTime
		m.mu.Unlock()

		log.Println("初次检查，记录当前交付时间")

		// 发送初始通知
		title := "🚗 理想汽车订单监控已启动"
		content := fmt.Sprintf("订单号: %s\n官方预计时间: %s\n\n%s",
			orderID,
			currentEstimateTime,
			m.getDetailedDeliveryInfo())

		if isApproaching {
			content += "\n\n⚠️ " + approachMsg
		}

		if err := m.sendNotification(title, content); err != nil {
			log.Printf("发送初始通知失败: %v", err)
		} else {
			// 初始通知成功后也要更新通知时间
			m.updateLastNotificationTime()
		}
		return
	}

	// 检查时间是否发生变化
	if currentEstimateTime != lastEstimateTime {
		log.Printf("检测到交付时间变化！从 %s 变更为 %s", lastEstimateTime, currentEstimateTime)

		title := "🚗 理想汽车交付时间更新通知"
		content := fmt.Sprintf("订单号: %s\n原官方预计时间: %s\n新官方预计时间: %s\n变更时间: %s\n\n%s",
			orderID,
			lastEstimateTime,
			currentEstimateTime,
			time.Now().Format(DateTimeFormat),
			m.getDetailedDeliveryInfo())

		if isApproaching {
			content += "\n\n⚠️ " + approachMsg
		}

		if err := m.sendNotification(title, content); err != nil {
			log.Printf("发送变更通知失败: %v", err)
		}

		// 更新记录的时间和通知时间
		m.mu.Lock()
		m.LastEstimateTime = currentEstimateTime
		m.mu.Unlock()
		m.updateLastNotificationTime()
	} else {
		log.Println("交付时间未发生变化")

		// 检查是否需要发送定期通知或临近交付提醒
		shouldNotifyPeriodic := m.shouldSendPeriodicNotification()

		m.mu.RLock()
		alwaysNotifyWhenApproaching := m.AlwaysNotifyWhenApproaching
		notificationInterval := m.NotificationInterval
		m.mu.RUnlock()

		shouldNotifyApproaching := isApproaching && alwaysNotifyWhenApproaching

		if shouldNotifyPeriodic || shouldNotifyApproaching {
			var title string
			var notifyReasons []string

			if shouldNotifyPeriodic {
				title = "📊 理想汽车订单状态定期报告"
				notifyReasons = append(notifyReasons, "定期状态更新")
				log.Printf("发送定期通知，距离上次通知已过 %.1f 小时",
					time.Since(m.LastNotificationTime).Hours())
			}

			if shouldNotifyApproaching {
				if title == "" {
					title = "⏰ 理想汽车交付时间提醒"
				}
				notifyReasons = append(notifyReasons, "临近交付时间")
				log.Printf("发送临近交付提醒: %s", approachMsg)
			}

			content := fmt.Sprintf("订单号: %s\n官方预计时间: %s\n通知原因: %s\n\n%s",
				orderID,
				currentEstimateTime,
				strings.Join(notifyReasons, "、"),
				m.getDetailedDeliveryInfo())

			if isApproaching {
				content += "\n\n⚠️ " + approachMsg
			}

			// 添加定期通知的额外信息
			if shouldNotifyPeriodic {
				content += fmt.Sprintf("\n\n📅 通知间隔: 每%.0f小时\n⏰ 下次通知时间: %s",
					notificationInterval.Hours(),
					time.Now().Add(notificationInterval).Format(DateTimeShort))
			}

			if err := m.sendNotification(title, content); err != nil {
				log.Printf("发送通知失败: %v", err)
			} else {
				// 只有成功发送通知后才更新时间
				m.updateLastNotificationTime()
				log.Printf("成功发送通知，原因: %s", strings.Join(notifyReasons, "、"))
			}
		} else {
			log.Println("无需发送通知：未到定期通知时间且非临近交付期")
		}
	}
}

func (m *Monitor) Start() error {
	log.Printf("启动监控服务，检查间隔: %s", m.CheckInterval)

	// 立即执行一次检查
	m.checkDeliveryTime()

	// 添加定时任务
	_, err := m.cron.AddFunc(m.CheckInterval, m.checkDeliveryTime)
	if err != nil {
		return fmt.Errorf("添加定时任务失败: %v", err)
	}

	m.cron.Start()
	log.Println("监控服务已启动，等待定时检查...")

	// 保持程序运行
	select {}
}

func (m *Monitor) Stop() {
	log.Println("停止监控服务...")
	m.cron.Stop()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	monitor := NewMonitor()

	// 检查配置
	if len(monitor.Notifiers) == 0 {
		log.Println("警告: 未配置任何通知器，将不会发送通知")
	}

	if monitor.LixiangCookies == "" {
		log.Println("警告: 未配置理想汽车 Cookies，可能导致请求失败")
	}

	// 启动监控
	if err := monitor.Start(); err != nil {
		log.Fatalf("启动监控服务失败: %v", err)
		os.Exit(1)
	}
}
