package main

import (
	"fmt"
	"strings"
	"time"
)

// 演示优化后的交付日期计算功能
func DemoDeliveryCalculation() {
	// 创建一个测试用的 Monitor 实例
	lockTime, _ := time.Parse(DateTimeFormat, "2024-09-27 13:08:00")

	monitor := &Monitor{
		OrderID:          "177971759268550919",
		LockOrderTime:    lockTime,
		EstimateWeeksMin: 7,
		EstimateWeeksMax: 9,
	}

	fmt.Println("🚗 理想汽车订单监控 - 交付日期计算优化演示")
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println()

	// 演示基本的交付时间计算
	minDate, maxDate := monitor.calculateEstimatedDelivery()
	fmt.Printf("📅 基于锁单时间的预计交付范围:\n")
	fmt.Printf("   最早: %s\n", minDate.Format(DateFormat))
	fmt.Printf("   最晚: %s\n", maxDate.Format(DateFormat))
	fmt.Println()

	// 演示剩余时间计算
	daysToMin, daysToMax, status := monitor.calculateRemainingDeliveryTime()
	fmt.Printf("⏱️  剩余时间分析:\n")
	fmt.Printf("   距离最早交付: %d 天\n", daysToMin)
	fmt.Printf("   距离最晚交付: %d 天\n", daysToMax)
	fmt.Printf("   当前状态: %s\n", status)
	fmt.Println()

	// 演示进度计算
	progress := monitor.calculateDeliveryProgress()
	fmt.Printf("📊 交付进度: %.1f%%\n", progress)
	fmt.Println()

	// 演示格式化输出
	fmt.Println("📝 格式化交付预测:")
	fmt.Println(monitor.formatDeliveryEstimate())
	fmt.Println()

	// 演示详细信息
	fmt.Println("📋 详细交付信息:")
	fmt.Println(monitor.getDetailedDeliveryInfo())
	fmt.Println()

	// 演示智能分析报告
	fmt.Println(monitor.getDeliveryAnalysisReport())

	// 演示临近交付检查
	isApproaching, msg := monitor.isApproachingDelivery()
	if isApproaching {
		fmt.Printf("⚠️  交付提醒: %s\n", msg)
	} else {
		fmt.Println("✅ 当前无紧急提醒")
	}
}
