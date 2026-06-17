package v2ray_admin

import "time"

func MonthlyTrafficLimitStartCreatedAt(now time.Time) int {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		now = now.In(location)
	}
	year, month, _ := now.Date()
	start := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
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

func IsMonthlyTrafficLimitResetDay(now time.Time) bool {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		now = now.In(location)
	}
	return now.Day() == 1
}

func clampDay(year int, month time.Month, day int) int {
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
	if day > lastDay {
		return lastDay
	}
	return day
}
