package utils

import (
	"time"
)

const (
	// Day 天单位
	Day = 24 * time.Hour
	// Week 周单位
	Week = 7 * Day

	month  = "month"
	second = "second"
	day    = "day"
	week   = "week"
	minute = "minute"
	hour   = "hour"
)

var (
	// DefaultLocation 默认时区
	DefaultLocation *time.Location
)

// InitTimeLocation 初始化时区 这里默认机器时间 UTC不是默认值
func InitTimeLocation(location string) error {
	if location == "" {
		DefaultLocation = time.Local
		return nil
	}
	var err error
	DefaultLocation, err = time.LoadLocation(location)
	return err
}

// Now 返回指定时区的时间
func Now() time.Time {
	return time.Now().In(DefaultLocation)
}

// GetNextTime 获取下次执行时间
func GetNextTime(start time.Time, duration int, unit string) time.Time {

	switch unit {
	case second:
		return start.Add(time.Duration(duration) * time.Second)
	case minute:
		return start.Add(time.Duration(duration) * time.Minute)
	case hour:
		return start.Add(time.Duration(duration) * time.Hour)
	case day:
		return start.Add(time.Duration(duration) * Day)
	case week:
		return start.Add(time.Duration(duration) * Week)
	}

	if unit == month {
		start = start.AddDate(0, duration, 0)
	}

	return start
}

// GetNextTimeAfterNow 获取当前时间之后的执行时间
func GetNextTimeAfterNow(start time.Time, duration int, unit string) time.Time {
	now := Now()

	if start.Before(now) {
		du := now.Sub(start)
		var factor int64
		var interval time.Duration
		if unit == month {
			if now.Year() == start.Year() {
				factor = (int64(now.Month()) - int64(start.Month())) / int64(duration)
			} else {
				factor = int64(((now.Year()-start.Year())*12 + int((now.Month() - start.Month()))) / duration)
			}

			start = start.AddDate(0, int(factor)*duration, 0)
			if start.Before(now) {
				return start.AddDate(0, duration, 0)
			}

		} else {
			interval = getInterval(duration, unit)
			factor = int64(du) / int64(interval)
			start = start.Add(time.Duration(factor) * interval)
			if start.Before(now) {
				return start.Add(interval)
			}
		}
	}

	return start
}

func getInterval(duration int, unit string) (t time.Duration) {
	switch unit {
	case second:
		t = time.Duration(duration) * time.Second
	case minute:
		t = time.Duration(duration) * time.Minute
	case hour:
		t = time.Duration(duration) * time.Hour
	case day:
		t = time.Duration(duration) * Day
	case week:
		t = time.Duration(duration) * Week
	}
	return
}

// IsBeforeOrEq 返回a时间是否在b时间之前或等于
func IsBeforeOrEq(a, b time.Time) bool {
	return a.Compare(b) != 1
}
