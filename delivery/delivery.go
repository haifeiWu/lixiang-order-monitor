package delivery

import (
	"fmt"
	"strings"
	"time"

	"lixiang-monitor/utils"
)

// Info 交付信息结构
type Info struct {
	LockOrderTime    time.Time
	EstimateWeeksMin int
	EstimateWeeksMax int
	// 缓存计算结果以提升性能
	cachedMinDate time.Time
	cachedMaxDate time.Time
}

// NewInfo 创建交付信息
func NewInfo(lockOrderTime time.Time, estimateWeeksMin, estimateWeeksMax int) *Info {
	info := &Info{
		LockOrderTime:    lockOrderTime,
		EstimateWeeksMin: estimateWeeksMin,
		EstimateWeeksMax: estimateWeeksMax,
	}
	// 预计算并缓存交付日期
	info.cachedMinDate = lockOrderTime.AddDate(0, 0, estimateWeeksMin*7)
	info.cachedMaxDate = lockOrderTime.AddDate(0, 0, estimateWeeksMax*7)
	return info
}

// CalculateEstimatedDelivery 计算预计交付日期范围（使用缓存以提升性能）
func (d *Info) CalculateEstimatedDelivery() (time.Time, time.Time) {
	return d.cachedMinDate, d.cachedMaxDate
}

// CalculateRemainingDeliveryTime 基于当前时间计算剩余交付时间
func (d *Info) CalculateRemainingDeliveryTime() (int, int, string) {
	now := time.Now()
	minDate, maxDate := d.CalculateEstimatedDelivery()

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

// CalculateDeliveryProgress 计算交付进度百分比
func (d *Info) CalculateDeliveryProgress() float64 {
	now := time.Now()

	// 计算从锁单到预计交付的总时间（取最大值）
	_, maxDate := d.CalculateEstimatedDelivery()
	totalDuration := maxDate.Sub(d.LockOrderTime)

	// 计算已经过去的时间
	elapsedDuration := now.Sub(d.LockOrderTime)

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

// FormatDeliveryEstimate 格式化交付日期范围
func (d *Info) FormatDeliveryEstimate() string {
	minDate, maxDate := d.CalculateEstimatedDelivery()
	_, _, status := d.CalculateRemainingDeliveryTime()
	progress := d.CalculateDeliveryProgress()

	baseInfo := ""
	if d.EstimateWeeksMin == d.EstimateWeeksMax {
		baseInfo = fmt.Sprintf("预计 %d 周后交付 (%s 左右)",
			d.EstimateWeeksMin,
			minDate.Format(utils.DateFormat))
	} else {
		baseInfo = fmt.Sprintf("预计 %d-%d 周后交付 (%s 至 %s)",
			d.EstimateWeeksMin,
			d.EstimateWeeksMax,
			minDate.Format(utils.DateFormat),
			maxDate.Format(utils.DateFormat))
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

// GetDetailedDeliveryInfo 获取详细的交付时间信息
func (d *Info) GetDetailedDeliveryInfo() string {
	now := time.Now()
	minDate, maxDate := d.CalculateEstimatedDelivery()
	_, _, status := d.CalculateRemainingDeliveryTime()
	progress := d.CalculateDeliveryProgress()

	// 计算锁单至今的天数
	daysSinceLock := int(now.Sub(d.LockOrderTime).Hours() / 24)

	info := fmt.Sprintf("📅 锁单时间: %s (%d天前)\n",
		d.LockOrderTime.Format(utils.DateTimeShort), daysSinceLock)

	info += fmt.Sprintf("🔮 基于锁单时间预测: %s\n", d.FormatDeliveryEstimate())
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

// GetAnalysisReport 获取交付时间智能分析报告
func (d *Info) GetAnalysisReport() string {
	now := time.Now()
	minDate, maxDate := d.CalculateEstimatedDelivery()
	daysToMin, _, status := d.CalculateRemainingDeliveryTime()
	progress := d.CalculateDeliveryProgress()

	report := "📊 交付时间智能分析报告\n"
	report += "=" + strings.Repeat("=", 30) + "\n\n"

	// 基本信息
	daysSinceLock := int(now.Sub(d.LockOrderTime).Hours() / 24)
	report += fmt.Sprintf("🔐 锁单信息: %s (%d天前)\n",
		d.LockOrderTime.Format(utils.DateTimeShort), daysSinceLock)

	report += fmt.Sprintf("📅 预计交付: %s - %s\n",
		minDate.Format(utils.DateFormat), maxDate.Format(utils.DateFormat))

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
}

// IsApproachingDelivery 检查是否临近预计交付时间
func (d *Info) IsApproachingDelivery() (bool, string) {
	now := time.Now()
	minDate, maxDate := d.CalculateEstimatedDelivery()

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
