package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"lixiang-monitor/cfg"
	"lixiang-monitor/cookie"
	"lixiang-monitor/db"
	"lixiang-monitor/delivery"
	"lixiang-monitor/notification"
	"lixiang-monitor/notifier"
	"lixiang-monitor/utils"
	"lixiang-monitor/web"

	"github.com/robfig/cron/v3"
)

type Monitor struct {
	OrderID          string
	LastEstimateTime string
	CheckInterval    string
	LixiangCookies   string
	LixiangHeaders   map[string]string
	Notifiers        []notifier.Notifier
	LockOrderTime    time.Time // 锁单时间
	EstimateWeeksMin int       // 预计交付周数范围（最小）
	EstimateWeeksMax int       // 预计交付周数范围（最大）
	cron             *cron.Cron

	// 定期通知相关字段
	NotificationInterval        time.Duration // 通知间隔（当交付时间未更新时）
	EnablePeriodicNotify        bool          // 是否启用定期通知
	AlwaysNotifyWhenApproaching bool          // 临近交付时总是通知

	// Cookie 管理相关
	LastCookieCheckTime      time.Time // 上次 Cookie 检查时间
	CookieExpiredNotified    bool      // 是否已通知 Cookie 失效
	ConsecutiveCookieFailure int       // 连续 Cookie 失效次数
	CookieUpdatedAt          time.Time // Cookie 更新时间
	CookieValidDays          int       // Cookie 有效天数
	CookieExpirationWarned   bool      // 是否已发送过期预警

	// 配置热加载相关
	mu            sync.RWMutex // 读写锁，保护配置的并发访问
	configVersion int          // 配置版本号，用于跟踪配置变化

	// 包管理器
	deliveryInfo        *delivery.Info        // 交付信息管理器
	cookieManager       *cookie.Manager       // Cookie 管理器
	notificationHandler *notification.Handler // 通知处理器
	database            *db.Database          // 数据库管理器
	webServer           *web.Server           // Web 服务器

	// Web 服务器配置
	WebEnabled  bool   // 是否启用 Web 服务器
	WebPort     int    // Web 服务器端口
	WebBasePath string // Web 服务器根路由
}

// 加载或重新加载配置
func (m *Monitor) loadConfig() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 使用 cfg 包加载配置
	config, err := cfg.Load()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	// 检查检查间隔是否变化
	checkIntervalChanged := false
	if config.CheckInterval != m.CheckInterval && m.CheckInterval != "" {
		checkIntervalChanged = true
	}

	// 更新 Monitor 字段
	m.OrderID = config.OrderID
	m.LixiangCookies = config.LixiangCookies
	m.CheckInterval = config.CheckInterval
	m.LockOrderTime = config.LockOrderTime
	m.EstimateWeeksMin = config.EstimateWeeksMin
	m.EstimateWeeksMax = config.EstimateWeeksMax
	m.EnablePeriodicNotify = config.EnablePeriodicNotify
	m.NotificationInterval = time.Duration(config.NotificationIntervalHours) * time.Hour
	m.AlwaysNotifyWhenApproaching = config.AlwaysNotifyWhenApproaching
	m.Notifiers = config.Notifiers
	m.CookieValidDays = config.CookieValidDays
	m.WebEnabled = config.WebEnabled
	m.WebPort = config.WebPort
	m.WebBasePath = config.WebBasePath

	// Cookie 更新时间处理
	if !config.CookieUpdatedAt.IsZero() {
		m.CookieUpdatedAt = config.CookieUpdatedAt
	} else if m.CookieUpdatedAt.IsZero() {
		m.CookieUpdatedAt = time.Now()
	}

	m.configVersion++
	log.Printf("配置已加载，版本: %d", m.configVersion)

	// 同步更新 deliveryInfo
	if m.deliveryInfo != nil {
		m.deliveryInfo = delivery.NewInfo(m.LockOrderTime, m.EstimateWeeksMin, m.EstimateWeeksMax)
	}

	// 同步更新 cookieManager
	if m.cookieManager != nil {
		m.cookieManager.UpdateCookie(m.LixiangCookies, m.LixiangHeaders)
		m.cookieManager.ValidDays = m.CookieValidDays
		m.cookieManager.UpdatedAt = m.CookieUpdatedAt
	}

	// 同步更新 notificationHandler
	if m.notificationHandler != nil {
		m.notificationHandler.UpdateConfig(
			m.Notifiers,
			m.deliveryInfo,
			m.NotificationInterval,
			m.EnablePeriodicNotify,
			m.AlwaysNotifyWhenApproaching,
		)
	}

	// 如果检查间隔变更且 cron 已经启动，返回错误提示需要重启
	if checkIntervalChanged && m.cron != nil {
		return fmt.Errorf("检查间隔已变更，需要重启服务")
	}

	return nil
}

// 监听配置文件变化
func (m *Monitor) watchConfig() {
	cfg.Watch(func() {
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
			time.Now().Format(utils.DateTimeFormat),
			m.OrderID,
			m.CheckInterval,
			len(m.Notifiers),
			m.EnablePeriodicNotify,
			m.NotificationInterval.Hours())

		if err := m.notificationHandler.SendCustomNotification(title, content); err != nil {
			log.Printf("发送配置更新通知失败: %v", err)
		}
	})
}

func NewMonitor() *Monitor {
	// 使用 cfg 包初始化配置
	if err := cfg.Init(); err != nil {
		log.Printf("初始化配置失败: %v", err)
	}

	monitor := &Monitor{
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

	// 初始化 delivery 信息管理器
	monitor.deliveryInfo = delivery.NewInfo(monitor.LockOrderTime, monitor.EstimateWeeksMin, monitor.EstimateWeeksMax)

	// 初始化 cookie 管理器
	monitor.cookieManager = cookie.NewManager(
		monitor.LixiangCookies,
		monitor.LixiangHeaders,
		monitor.CookieValidDays,
		monitor.CookieUpdatedAt,
	)

	// 设置 cookie 管理器的回调函数
	monitor.cookieManager.OnCookieExpired = func(statusCode int, message string) {
		title := "❌ 理想汽车 Cookie 已失效"
		content := fmt.Sprintf("检测到 Cookie 已失效,需要立即更新！\n\n"+
			"状态码: %d\n"+
			"错误信息: %s\n"+
			"失败次数: %d\n"+
			"检测时间: %s\n\n"+
			"⚠️  请立即更新 config.yaml 中的 lixiang_cookies 字段！",
			statusCode, message, monitor.cookieManager.ConsecutiveFailure, time.Now().Format(utils.DateTimeFormat))

		if err := monitor.notificationHandler.SendCustomNotification(title, content); err != nil {
			log.Printf("Cookie 失效通知发送失败: %v", err)
		}
	}

	monitor.cookieManager.OnCookieExpirationWarning = func(timeDesc, expireTime, updatedAt string, ageInDays float64) {
		title := "⚠️  理想汽车 Cookie 即将过期"
		content := fmt.Sprintf("您的 Cookie 即将过期,建议提前更新！\n\n"+
			"剩余时间: %s\n"+
			"过期时间: %s\n"+
			"更新时间: %s\n"+
			"已使用: %.1f 天\n\n"+
			"请及时更新 config.yaml 中的 lixiang_cookies 字段，避免监控中断。",
			timeDesc, expireTime, updatedAt, ageInDays)

		if err := monitor.notificationHandler.SendCustomNotification(title, content); err != nil {
			log.Printf("Cookie 过期预警通知发送失败: %v", err)
		}
	}

	// 初始化 notification 处理器
	monitor.notificationHandler = notification.NewHandler(
		monitor.Notifiers,
		monitor.deliveryInfo,
		monitor.NotificationInterval,
		monitor.EnablePeriodicNotify,
		monitor.AlwaysNotifyWhenApproaching,
	)

	// 初始化数据库
	database, err := db.New("./lixiang-monitor.db")
	if err != nil {
		log.Printf("⚠️  数据库初始化失败: %v (历史记录功能将不可用)", err)
	} else {
		monitor.database = database
		log.Println("✅ 数据库初始化成功")
	}

	// 初始化 Web 服务器
	if monitor.WebEnabled && monitor.database != nil {
		webServer, err := web.NewServer(monitor.database, monitor.OrderID, monitor.WebPort, monitor.WebBasePath)
		if err != nil {
			log.Printf("⚠️  Web 服务器初始化失败: %v", err)
		} else {
			monitor.webServer = webServer
			log.Println("✅ Web 服务器初始化成功")
		}
	}

	// 启动配置文件监听
	monitor.watchConfig()

	if len(monitor.Notifiers) == 0 {
		log.Println("⚠️  未配置任何通知器,将不会发送通知")
	} else {
		log.Printf("✅ 已配置 %d 个通知器", len(monitor.Notifiers))
	}

	return monitor
}

// parseOrderResponse 解析订单响应数据
func (m *Monitor) parseOrderResponse(rawData interface{}) (estimateTime string, err error) {
	// 将 interface{} 转换为 map[string]interface{}
	orderDataMap, ok := rawData.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("订单数据格式错误")
	}

	// 解析 code 字段
	code := 0
	if codeVal, ok := orderDataMap["code"].(float64); ok {
		code = int(codeVal)
	}

	if code != 0 {
		message := ""
		if msgVal, ok := orderDataMap["message"].(string); ok {
			message = msgVal
		}
		return "", fmt.Errorf("API 返回错误: %s", message)
	}

	// 解析 EstimateDeliveringAt
	if data, ok := orderDataMap["data"].(map[string]interface{}); ok {
		if delivery, ok := data["delivery"].(map[string]interface{}); ok {
			if estimateTime, ok := delivery["estimateDeliveringAt"].(string); ok {
				return estimateTime, nil
			}
		}
	}

	return "", nil
}

// logDeliveryInfo 记录交付信息日志
func (m *Monitor) logDeliveryInfo(lockOrderTime time.Time, isApproaching bool, approachMsg string) {
	predictedDelivery := m.deliveryInfo.FormatDeliveryEstimate()
	log.Printf("锁单时间: %s", lockOrderTime.Format(utils.DateTimeFormat))
	log.Printf("基于锁单时间预测: %s", predictedDelivery)
	if isApproaching {
		log.Printf("交付提醒: %s", approachMsg)
	}
}

// handleDeliveryNotification 处理交付通知逻辑
func (m *Monitor) handleDeliveryNotification(orderID, currentEstimateTime, lastEstimateTime string, isApproaching bool, approachMsg string) {
	timeChanged := false
	notificationSent := false

	if lastEstimateTime == "" {
		// 首次检查
		if err := m.notificationHandler.HandleFirstCheck(orderID, currentEstimateTime, isApproaching, approachMsg); err != nil {
			log.Printf("处理首次检查通知失败: %v", err)
		} else {
			notificationSent = true
		}
		m.updateLastEstimateTime(currentEstimateTime)
	} else if currentEstimateTime != lastEstimateTime {
		// 时间发生变化
		timeChanged = true
		if err := m.notificationHandler.HandleTimeChanged(orderID, currentEstimateTime, lastEstimateTime, isApproaching, approachMsg); err != nil {
			log.Printf("处理时间变更通知失败: %v", err)
		} else {
			notificationSent = true
		}
		m.updateLastEstimateTime(currentEstimateTime)
	} else {
		// 时间未变化，检查是否需要定期通知
		log.Println("交付时间未发生变化")
		if err := m.notificationHandler.HandlePeriodicNotification(orderID, currentEstimateTime, isApproaching, approachMsg); err != nil {
			log.Printf("处理定期通知失败: %v", err)
		}
		// 定期通知也算作已发送通知
		notificationSent = m.EnablePeriodicNotify
	}

	// 保存记录到数据库
	m.saveDeliveryRecord(orderID, currentEstimateTime, lastEstimateTime, isApproaching, approachMsg, timeChanged, notificationSent)
}

// updateLastEstimateTime 更新最后的预估时间
func (m *Monitor) updateLastEstimateTime(estimateTime string) {
	m.mu.Lock()
	m.LastEstimateTime = estimateTime
	m.mu.Unlock()
}

// saveDeliveryRecord 保存交付记录到数据库
func (m *Monitor) saveDeliveryRecord(orderID, currentEstimateTime, previousEstimate string, isApproaching bool, approachMsg string, timeChanged, notificationSent bool) {
	// 如果数据库未初始化，跳过保存
	if m.database == nil {
		return
	}

	record := &db.DeliveryRecord{
		OrderID:          orderID,
		EstimateTime:     currentEstimateTime,
		LockOrderTime:    m.LockOrderTime,
		CheckTime:        time.Now(),
		IsApproaching:    isApproaching,
		ApproachMessage:  approachMsg,
		TimeChanged:      timeChanged,
		PreviousEstimate: previousEstimate,
		NotificationSent: notificationSent,
		CreatedAt:        time.Now(),
	}

	if err := m.database.SaveDeliveryRecord(record); err != nil {
		log.Printf("保存交付记录失败: %v", err)
	}
}

func (m *Monitor) checkDeliveryTime() {
	log.Println("开始检查订单交付时间...")

	// 获取订单数据
	m.mu.RLock()
	orderID := m.OrderID
	m.mu.RUnlock()

	rawData, err := m.cookieManager.FetchOrderData(orderID)
	if err != nil {
		if _, isCookieError := err.(*cookie.CookieExpiredError); isCookieError {
			log.Printf("⚠️  Cookie 已失效，跳过本次检查: %v", err)
			return
		}
		log.Printf("获取订单数据失败: %v", err)
		return
	}

	// 解析订单响应
	currentEstimateTime, err := m.parseOrderResponse(rawData)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	log.Printf("当前预计交付时间: %s", currentEstimateTime)

	// 读取配置信息
	m.mu.RLock()
	lockOrderTime := m.LockOrderTime
	lastEstimateTime := m.LastEstimateTime
	m.mu.RUnlock()

	// 计算交付预测和临近状态
	isApproaching, approachMsg := m.deliveryInfo.IsApproachingDelivery()

	// 记录交付信息
	m.logDeliveryInfo(lockOrderTime, isApproaching, approachMsg)

	// 处理通知逻辑
	m.handleDeliveryNotification(orderID, currentEstimateTime, lastEstimateTime, isApproaching, approachMsg)
}

func (m *Monitor) Start() error {
	log.Printf("启动监控服务，检查间隔: %s", m.CheckInterval)

	// 立即执行一次检查
	m.checkDeliveryTime()

	// 立即检查 Cookie 过期状态并显示状态
	log.Printf("Cookie 状态: %s", m.cookieManager.GetStatus())
	m.cookieManager.CheckExpiration()

	// 添加定时任务 - 订单检查
	_, err := m.cron.AddFunc(m.CheckInterval, m.checkDeliveryTime)
	if err != nil {
		return fmt.Errorf("添加定时任务失败: %v", err)
	}

	// 添加定时任务 - 每日检查 Cookie 过期（凌晨 1 点）
	_, err = m.cron.AddFunc("0 0 1 * * *", func() {
		log.Printf("执行定期 Cookie 过期检查")
		m.cookieManager.CheckExpiration()
	})
	if err != nil {
		log.Printf("警告: 添加 Cookie 过期检查任务失败: %v", err)
	}

	m.cron.Start()

	// 启动 Web 服务器
	if m.webServer != nil {
		if err := m.webServer.Start(); err != nil {
			log.Printf("⚠️  Web 服务器启动失败: %v", err)
		}
	}

	log.Println("监控服务已启动，等待定时检查...")

	// 保持程序运行
	select {}
}

func (m *Monitor) Stop() {
	log.Println("停止监控服务...")
	m.cron.Stop()

	// 停止 Web 服务器
	if m.webServer != nil {
		if err := m.webServer.Stop(); err != nil {
			log.Printf("关闭 Web 服务器失败: %v", err)
		}
	}

	// 关闭数据库连接
	if m.database != nil {
		if err := m.database.Close(); err != nil {
			log.Printf("关闭数据库连接失败: %v", err)
		}
	}
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
