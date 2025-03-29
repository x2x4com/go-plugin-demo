package main

import (
	"time"
)

// AddDays 日期加减
func AddDays(date time.Time, days int) time.Time {
	return date.AddDate(0, 0, days)
}

// Format 日期格式化
func Format(date time.Time, layout string) string {
	return date.Format(layout)
}

// Parse 日期解析
func Parse(dateStr, layout string) (time.Time, error) {
	return time.Parse(layout, dateStr)
}

// Between 计算日期差值
func Between(start, end time.Time) int {
	duration := end.Sub(start)
	return int(duration.Hours() / 24)
}
