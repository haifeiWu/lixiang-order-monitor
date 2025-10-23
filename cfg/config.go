package cfg

import (
	"log"
	"time"

	"lixiang-monitor/notifier"
	"lixiang-monitor/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	// 订单信息
	OrderID          string
	LixiangCookies   string
	CheckInterval    string
	LockOrderTime    time.Time
	EstimateWeeksMin int
	EstimateWeeksMax int

	// 通知相关
	Notifiers                   []notifier.Notifier
	EnablePeriodicNotify        bool
	NotificationIntervalHours   int
	AlwaysNotifyWhenApproaching bool

	// Cookie 管理
	CookieUpdatedAt time.Time
	CookieValidDays int

	// Web 服务器
	WebEnabled  bool
	WebPort     int
	WebBasePath string
}

// Init 初始化配置系统
func Init() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("配置文件读取失败: %v", err)
	}

	// 设置默认值
	setDefaults()

	return nil
}

// setDefaults 设置配置默认值
func setDefaults() {
	viper.SetDefault("order_id", "177971759268550919")
	viper.SetDefault("check_interval", "@every 30m")
	viper.SetDefault("wechat_webhook_url", "")
	viper.SetDefault("serverchan_sendkey", "")
	viper.SetDefault("serverchan_baseurl", "https://sctapi.ftqq.com/")
	viper.SetDefault("bark_server_url", "")
	viper.SetDefault("bark_sound", "minuet")
	viper.SetDefault("bark_icon", "")
	viper.SetDefault("bark_group", "lixiang-monitor")
	viper.SetDefault("lock_order_time", "2025-09-27 13:08:00")
	viper.SetDefault("estimate_weeks_min", 7)
	viper.SetDefault("estimate_weeks_max", 9)
	viper.SetDefault("enable_periodic_notify", true)
	viper.SetDefault("notification_interval_hours", 24)
	viper.SetDefault("always_notify_when_approaching", true)
	viper.SetDefault("cookie_valid_days", 7)
	viper.SetDefault("web_enabled", true)
	viper.SetDefault("web_port", 8080)
	viper.SetDefault("web_base_path", "")
}

// Load 加载配置并返回 Config 结构
func Load() (*Config, error) {
	cfg := &Config{}

	// 解析锁单时间
	lockOrderTimeStr := viper.GetString("lock_order_time")
	lockOrderTime, err := utils.ParseLockOrderTime(lockOrderTimeStr)
	if err != nil {
		log.Printf("锁单时间解析失败: %v, 使用默认时间", err)
		lockOrderTime, _ = time.Parse(utils.DateTimeFormat, "2025-09-27 13:08:00")
	}

	// 基本配置
	cfg.OrderID = viper.GetString("order_id")
	cfg.LixiangCookies = viper.GetString("lixiang_cookies")
	cfg.CheckInterval = viper.GetString("check_interval")
	cfg.LockOrderTime = lockOrderTime
	cfg.EstimateWeeksMin = viper.GetInt("estimate_weeks_min")
	cfg.EstimateWeeksMax = viper.GetInt("estimate_weeks_max")

	// 通知配置
	cfg.EnablePeriodicNotify = viper.GetBool("enable_periodic_notify")
	cfg.NotificationIntervalHours = viper.GetInt("notification_interval_hours")
	cfg.AlwaysNotifyWhenApproaching = viper.GetBool("always_notify_when_approaching")
	cfg.Notifiers = loadNotifiers()

	// Cookie 配置
	cfg.CookieValidDays = viper.GetInt("cookie_valid_days")
	if cfg.CookieValidDays == 0 {
		cfg.CookieValidDays = 7
	}

	cookieUpdatedStr := viper.GetString("cookie_updated_at")
	if cookieUpdatedStr != "" {
		if parsedTime, err := time.Parse(utils.DateTimeFormat, cookieUpdatedStr); err == nil {
			cfg.CookieUpdatedAt = parsedTime
		} else {
			cfg.CookieUpdatedAt = time.Now()
		}
	} else {
		cfg.CookieUpdatedAt = time.Now()
	}

	// Web 服务器配置
	cfg.WebEnabled = viper.GetBool("web_enabled")
	cfg.WebPort = viper.GetInt("web_port")
	if cfg.WebPort == 0 {
		cfg.WebPort = 8080
	}
	cfg.WebBasePath = viper.GetString("web_base_path")

	return cfg, nil
}

// loadNotifiers 加载所有配置的通知器
func loadNotifiers() []notifier.Notifier {
	var notifiers []notifier.Notifier

	// 微信群机器人
	wechatWebhookURL := viper.GetString("wechat_webhook_url")
	if wechatWebhookURL != "" {
		notifiers = append(notifiers, &notifier.WeChatWebhookNotifier{
			WebhookURL: wechatWebhookURL,
		})
	}

	// ServerChan
	serverChanSendKey := viper.GetString("serverchan_sendkey")
	if serverChanSendKey != "" {
		notifiers = append(notifiers, &notifier.ServerChanNotifier{
			SendKey: serverChanSendKey,
			BaseURL: viper.GetString("serverchan_baseurl"),
		})
	}

	// Bark
	barkServerURL := viper.GetString("bark_server_url")
	if barkServerURL != "" {
		notifiers = append(notifiers, &notifier.BarkNotifier{
			ServerURL: barkServerURL,
			Sound:     viper.GetString("bark_sound"),
			Icon:      viper.GetString("bark_icon"),
			Group:     viper.GetString("bark_group"),
		})
	}

	return notifiers
}

// Watch 监听配置文件变化
func Watch(callback func()) {
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("配置文件已更新: %s", e.Name)

		if err := viper.ReadInConfig(); err != nil {
			log.Printf("重新读取配置文件失败: %v", err)
			return
		}

		if callback != nil {
			callback()
		}
	})

	viper.WatchConfig()
	log.Println("✅ 配置文件监听已启动")
}

// GetString 获取字符串配置
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt 获取整数配置
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool 获取布尔配置
func GetBool(key string) bool {
	return viper.GetBool(key)
}
