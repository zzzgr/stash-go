package time_util

import (
	"time"
)

func FormatDateTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05.000")
}

func RoundTimeToHour(t time.Time, offset time.Duration) time.Time {
	// 获取整点时间
	roundedTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())

	// 判断分钟是否大于等于30
	if t.Minute() >= 30 {
		// 将小时加一
		roundedTime = roundedTime.Add(time.Hour)
	}

	roundedTime = roundedTime.Add(offset)
	return roundedTime
}

func RoundTimeToMinute(t time.Time, offset time.Duration) time.Time {
	// 获取整点时间
	roundedTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())

	// 判断分钟是否大于等于30
	if t.Second() >= 30 {
		// 将小时加一
		roundedTime = roundedTime.Add(time.Minute)
	}

	roundedTime = roundedTime.Add(offset)
	return roundedTime
}
