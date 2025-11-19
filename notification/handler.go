package notification

import (
	"fmt"
	"log"
	"strings"
	"time"

	"lixiang-monitor/delivery"
	"lixiang-monitor/notifier"
	"lixiang-monitor/utils"
)

// 通知相关常量
const (
	WarningPrefix = "\n\n⚠️ "

	// 通知标题
	TitleMonitorStarted    = "🚗 理想汽车订单监控已启动"
	TitleTimeChanged       = "🚗 理想汽车交付时间更新通知"
	TitlePeriodicReport    = "📊 理想汽车订单状态定期报告"
	TitleApproachingRemind = "⏰ 理想汽车交付时间提醒"
)

// Handler 通知处理器
type Handler struct {
	notifiers                   []notifier.Notifier
	deliveryInfo                *delivery.Info
	lastNotificationTime        time.Time
	notificationInterval        time.Duration
	enablePeriodicNotify        bool
	alwaysNotifyWhenApproaching bool
}

// NewHandler 创建通知处理器
func NewHandler(
	notifiers []notifier.Notifier,
	deliveryInfo *delivery.Info,
	notificationInterval time.Duration,
	enablePeriodicNotify bool,
	alwaysNotifyWhenApproaching bool,
) *Handler {
	return &Handler{
		notifiers:                   notifiers,
		deliveryInfo:                deliveryInfo,
		notificationInterval:        notificationInterval,
		enablePeriodicNotify:        enablePeriodicNotify,
		alwaysNotifyWhenApproaching: alwaysNotifyWhenApproaching,
		lastNotificationTime:        time.Time{}, // 初始化为零值
	}
}

// UpdateConfig 更新配置
func (h *Handler) UpdateConfig(
	notifiers []notifier.Notifier,
	deliveryInfo *delivery.Info,
	notificationInterval time.Duration,
	enablePeriodicNotify bool,
	alwaysNotifyWhenApproaching bool,
) {
	h.notifiers = notifiers
	h.deliveryInfo = deliveryInfo
	h.notificationInterval = notificationInterval
	h.enablePeriodicNotify = enablePeriodicNotify
	h.alwaysNotifyWhenApproaching = alwaysNotifyWhenApproaching
}

// HandleFirstCheck 处理首次检查的通知
func (h *Handler) HandleFirstCheck(orderID, currentEstimateTime string, isApproaching bool, approachMsg string) error {
	log.Println("初次检查，记录当前交付时间")

	content := h.buildInitialNotificationContent(orderID, currentEstimateTime)
	if isApproaching {
		content += WarningPrefix + approachMsg
	}

	if err := h.sendNotification(TitleMonitorStarted, content); err != nil {
		return fmt.Errorf("发送初始通知失败: %v", err)
	}

	h.updateLastNotificationTime()
	return nil
}

// buildInitialNotificationContent 构建初始通知内容
func (h *Handler) buildInitialNotificationContent(orderID, currentEstimateTime string) string {
	return fmt.Sprintf("订单号: %s\n官方预计时间: %s\n\n%s",
		orderID,
		currentEstimateTime,
		h.deliveryInfo.GetDetailedDeliveryInfo())
}

// HandleTimeChanged 处理交付时间变化的通知
func (h *Handler) HandleTimeChanged(orderID, currentEstimateTime, lastEstimateTime string, isApproaching bool, approachMsg string) error {
	log.Printf("检测到交付时间变化！从 %s 变更为 %s", lastEstimateTime, currentEstimateTime)

	content := h.buildTimeChangedContent(orderID, lastEstimateTime, currentEstimateTime)
	if isApproaching {
		content += WarningPrefix + approachMsg
	}

	if err := h.sendNotification(TitleTimeChanged, content); err != nil {
		return fmt.Errorf("发送变更通知失败: %v", err)
	}

	h.updateLastNotificationTime()
	return nil
}

// buildTimeChangedContent 构建时间变更通知内容
func (h *Handler) buildTimeChangedContent(orderID, lastEstimateTime, currentEstimateTime string) string {
	return fmt.Sprintf("订单号: %s\n原官方预计时间: %s\n新官方预计时间: %s\n变更时间: %s\n\n%s",
		orderID,
		lastEstimateTime,
		currentEstimateTime,
		time.Now().Format(utils.DateTimeFormat),
		h.deliveryInfo.GetDetailedDeliveryInfo())
}

// HandlePeriodicNotification 处理定期通知和临近提醒
func (h *Handler) HandlePeriodicNotification(orderID, currentEstimateTime string, isApproaching bool, approachMsg string) error {
	shouldNotifyPeriodic := h.shouldSendPeriodicNotification()
	shouldNotifyApproaching := isApproaching && h.alwaysNotifyWhenApproaching

	if !shouldNotifyPeriodic && !shouldNotifyApproaching {
		log.Println("无需发送通知：未到定期通知时间且非临近交付期")
		return nil
	}

	// 确定通知标题和原因
	title, notifyReasons := h.determineNotificationTitleAndReasons(shouldNotifyPeriodic, shouldNotifyApproaching, approachMsg)

	// 构建通知内容
	content := h.buildPeriodicNotificationContent(orderID, currentEstimateTime, notifyReasons, isApproaching, approachMsg, shouldNotifyPeriodic)

	// 发送通知
	if err := h.sendNotification(title, content); err != nil {
		return fmt.Errorf("发送通知失败: %v", err)
	}

	h.updateLastNotificationTime()
	log.Printf("成功发送通知，原因: %s", strings.Join(notifyReasons, "、"))
	return nil
}

// determineNotificationTitleAndReasons 确定通知标题和原因
func (h *Handler) determineNotificationTitleAndReasons(shouldNotifyPeriodic, shouldNotifyApproaching bool, approachMsg string) (string, []string) {
	var title string
	var notifyReasons []string

	if shouldNotifyPeriodic {
		title = TitlePeriodicReport
		notifyReasons = append(notifyReasons, "定期状态更新")
		log.Printf("发送定期通知，距离上次通知已过 %.1f 小时",
			time.Since(h.lastNotificationTime).Hours())
	}

	if shouldNotifyApproaching {
		if title == "" {
			title = TitleApproachingRemind
		}
		notifyReasons = append(notifyReasons, "临近交付时间")
		log.Printf("发送临近交付提醒: %s", approachMsg)
	}

	return title, notifyReasons
}

// buildPeriodicNotificationContent 构建定期通知内容
func (h *Handler) buildPeriodicNotificationContent(orderID, currentEstimateTime string, notifyReasons []string, isApproaching bool, approachMsg string, shouldNotifyPeriodic bool) string {
	// 使用 strings.Builder 优化字符串拼接性能
	var builder strings.Builder
	builder.Grow(512) // 预分配容量
	
	fmt.Fprintf(&builder, "订单号: %s\n官方预计时间: %s\n通知原因: %s\n\n%s",
		orderID,
		currentEstimateTime,
		strings.Join(notifyReasons, "、"),
		h.deliveryInfo.GetDetailedDeliveryInfo())

	if isApproaching {
		builder.WriteString(WarningPrefix)
		builder.WriteString(approachMsg)
	}

	// 添加定期通知的额外信息
	if shouldNotifyPeriodic {
		fmt.Fprintf(&builder, "\n\n📅 通知间隔: 每%.0f小时\n⏰ 下次通知时间: %s",
			h.notificationInterval.Hours(),
			time.Now().Add(h.notificationInterval).Format(utils.DateTimeShort))
	}

	return builder.String()
}

// shouldSendPeriodicNotification 检查是否应该发送定期通知
func (h *Handler) shouldSendPeriodicNotification() bool {
	if !h.enablePeriodicNotify {
		return false
	}

	if h.lastNotificationTime.IsZero() {
		return false
	}

	timeSinceLastNotification := time.Since(h.lastNotificationTime)
	return timeSinceLastNotification >= h.notificationInterval
}

// updateLastNotificationTime 更新最后通知时间
func (h *Handler) updateLastNotificationTime() {
	h.lastNotificationTime = time.Now()
}

// GetLastNotificationTime 获取最后通知时间
func (h *Handler) GetLastNotificationTime() time.Time {
	return h.lastNotificationTime
}

// SetLastNotificationTime 设置最后通知时间
func (h *Handler) SetLastNotificationTime(t time.Time) {
	h.lastNotificationTime = t
}

// sendNotification 发送通知
func (h *Handler) sendNotification(title, content string) error {
	if len(h.notifiers) == 0 {
		log.Println("未配置任何通知器，跳过通知")
		return nil
	}

	var errors []string
	successCount := 0

	for _, n := range h.notifiers {
		if err := n.Send(title, content); err != nil {
			log.Printf("通知发送失败: %v", err)
			errors = append(errors, err.Error())
		} else {
			successCount++
		}
	}

	if successCount == 0 {
		return fmt.Errorf("所有通知器发送失败: %v", errors)
	} else if len(errors) > 0 {
		log.Printf("部分通知器发送失败 (%d/%d 成功): %v", successCount, len(h.notifiers), errors)
	}

	return nil
}

// SendCustomNotification 发送自定义通知
func (h *Handler) SendCustomNotification(title, content string) error {
	return h.sendNotification(title, content)
}
