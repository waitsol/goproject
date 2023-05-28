package com

import "time"

func isHoliday(year int, month time.Month, day int) bool {
	// 判断是否是元旦节
	if month == time.January && day == 1 {
		return true
	}

	// 判断是否是春节
	if month == time.February && (day == 11 || day == 12 || day == 13 || day == 14 || day == 15 || day == 16 || day == 17) {
		return true
	}

	// 判断是否是清明节
	if month == time.April && (day == 4 || day == 5 || day == 6) {
		return true
	}

	// 判断是否是劳动节
	if month == time.May && (day == 1 || day == 2 || day == 3) {
		return true
	}

	// 判断是否是端午节
	if month == time.June && (day == 9 || day == 10 || day == 11) {
		return true
	}

	// 判断是否是中秋节
	if month == time.September && (day == 19 || day == 20 || day == 21) {
		return true
	}

	// 判断是否是国庆节
	if month == time.October && (day == 1 || day == 2 || day == 3 || day == 4 || day == 5 || day == 6 || day == 7) {
		return true
	}

	return false
}

func IsSend() bool {
	// 获取当前日期
	today := time.Now()

	// 获取今天是周几
	dayOfWeek := int(today.Weekday())

	// 判断今天是否是周六或周日
	if dayOfWeek == 0 || dayOfWeek == 6 {
		return false
	} else {
		// 获取今天的年、月、日
		year, month, day := today.Date()

		// 判断今天是否是法定节假日
		if isHoliday(year, month, day) {
			return false
		} else {
			return true
		}
	}
}
