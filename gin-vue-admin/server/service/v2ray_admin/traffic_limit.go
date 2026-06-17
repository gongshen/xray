package v2ray_admin

import "time"

func TrafficLimitStartCreatedAt(now time.Time, resetDay int) int {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		now = now.In(location)
	}
	if resetDay <= 0 {
		resetDay = 1
	}

	year, month, _ := now.Date()
	currentResetDay := clampDay(year, month, resetDay)
	var start time.Time
	if now.Day() >= currentResetDay {
		start = time.Date(year, month, currentResetDay, 0, 0, 0, 0, now.Location())
	} else {
		prev := now.AddDate(0, -1, 0)
		prevYear, prevMonth, _ := prev.Date()
		prevResetDay := clampDay(prevYear, prevMonth, resetDay)
		start = time.Date(prevYear, prevMonth, prevResetDay, 0, 0, 0, 0, now.Location())
	}
	return start.Year()*10000 + int(start.Month())*100 + start.Day()
}

func ShouldLimitTraffic(traffic uint64, trafficLimitGB float64) bool {
	if trafficLimitGB < 0 {
		return false
	}
	if trafficLimitGB == 0 {
		return true
	}
	return traffic >= uint64(trafficLimitGB*1024*1024*1024)
}

func IsTrafficResetDay(now time.Time, resetDay int) bool {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		now = now.In(location)
	}
	if resetDay <= 0 {
		resetDay = 1
	}
	return now.Day() == clampDay(now.Year(), now.Month(), resetDay)
}

func clampDay(year int, month time.Month, day int) int {
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
	if day > lastDay {
		return lastDay
	}
	return day
}
