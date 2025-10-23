package utils

import (
	"fmt"
	"time"
)

const (
	// DateTimeFormat 日期时间格式
	DateTimeFormat = "2006-01-02 15:04:05"
	// DateTimeShort 短日期时间格式
	DateTimeShort = "2006-01-02 15:04"
	// DateFormat 日期格式
	DateFormat = "2006-01-02"
)

// ParseLockOrderTime 解析锁单时间，支持多种时间格式
func ParseLockOrderTime(timeStr string) (time.Time, error) {
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
