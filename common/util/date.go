package util

import (
	"strconv"
	"time"
)

// 获取日期的年月日
func GetDateItem(t *time.Time) (year int, month int, day int) {
	year = t.Year()
	month = int(t.Month())
	day = t.Day()
	return
}

//获取时间的时分秒
func GetTimeItem(t *time.Time) (hour int, min int, sec int) {
	hour = t.Hour()
	min = t.Minute()
	sec = t.Second()
	return
}

func GetNowDateFormat() string {
	t := time.Now()
	return t.Format("2006-01-02")
}

func GetNowTimeStap() int64 {
	return time.Now().Unix()
}

func GetNowTimestapByString(dataStr string, format string) int64 {
	t, _ := time.Parse(format, dataStr)
	return t.Unix()
}

func GetNowDateTimeFormat() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

func GetNowDateDayFormat() string {
	t := time.Now()
	return t.Format("2006-01-02")
}

func GetTimeFromDefaultString(str string) time.Time {
	layout := "2006-01-02"

	t, _ := time.Parse(layout, str)
	return t
}

func GetNowDateTimeFormatCustom(format string) string {
	t := time.Now()
	return t.Format(format)
}

func GetTimeFormat(ts int64) string {
	t := time.Unix(ts, 0)
	return t.Format("2006-01-02 15:04:05")
}

func GetTimeFormatCustom(ts int64, format string) string {
	t := time.Unix(ts, 0)
	return t.Format(format)
}

func GetDateFormat(t time.Time) string {
	return t.Format("2006-01-02")
}

//判断是否为闰年
func IsLeapYear(year int) bool {
	//判断是否为闰年
	if year%4 == 0 && year%100 != 0 || year%400 == 0 {
		return true
	}

	return false
}

//获取某月有多少天
func GetMonthDays(year int, month int) int {
	days := 0
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30
		} else {
			days = 31
		}
	} else {
		if IsLeapYear(year) {
			days = 29
		} else {
			days = 28
		}
	}
	return days
}

//时间字符串
func GetNowTime() string {
	t := time.Now()
	year := strconv.Itoa(t.Year())
	month := strconv.Itoa(int(t.Month()))
	day := strconv.Itoa(t.Day())
	hour := strconv.Itoa(t.Hour())
	min := strconv.Itoa(t.Minute())
	sec := strconv.Itoa(t.Second())

	if len(month) == 1 {
		month = "0" + month
	}

	if len(day) == 1 {
		day = "0" + day
	}
	return year + month + day + hour + min + sec
}

// 报表时间: 年月日
func GetTimeString() (y, m, d string) {
	t := time.Now()
	y = t.Format("2006")
	m = t.Format("2006-01")
	d = t.Format("2006-01-02")
	return
}

//时间字符串
func GetNowTime4Day() string {
	t := time.Now()
	day := t.Format("20060102")
	return day
}

// 获取
func GetTimeArray(startTime, endTime string) []string {
	layout := "2006-01-02 15:04:05"
	start, err := time.Parse(layout, startTime)
	resArr := make([]string, 0)
	if err != nil {
		return resArr
	}
	str1 := start.Format("20060102")
	resArr = append(resArr, str1)
	end, err := time.Parse(layout, endTime)
	if err != nil {
		return resArr
	}
	timeBuild(start, end, &resArr)
	str2 := end.Format("20060102")

	if str2 != str1 && str2 != resArr[len(resArr)-1] {
		resArr = append(resArr, str2)
	}
	return resArr
}

func timeBuild(start time.Time, end time.Time, between *[]string) {
	d, err := time.ParseDuration("24h")
	if err != nil {
		return
	}
	nextDay := start.Add(d)
	if end.After(nextDay) {
		dayStr := nextDay.Format("20060102")
		*between = append(*between, dayStr)
		timeBuild(nextDay, end, between)
	}
}

// 获取某时刻多少天前的时间
func GetTimeBeforNowDay(start string, num int) (string, string) {
	t := time.Now()
	if start != "" {
		layout := "2006-01-02 15:04:05"
		startTime, err := time.Parse(layout, start)
		if err != nil {
			t = startTime
		}
	} else {
		start = t.Format("2006-01-02 15:04:05")
	}

	allHour := num * 24
	allHourStr := "-" + strconv.Itoa(allHour) + "h"

	d, err := time.ParseDuration(allHourStr)
	if err != nil {
		return "", start
	}
	nextDay := t.Add(d)
	day := nextDay.Format("2006-01-02 15:04:05")
	return day, start
}

//年+月 时间字符串
func GetNowYearMoth() []string {
	t := time.Now()
	day := t.Format("2006")
	resArr := make([]string, 0)
	fMonth, err := strconv.Atoi(day + "01")
	if err != nil {
		return resArr
	}
	for i := 0; i < 12; i++ {
		mStr := strconv.Itoa(fMonth + i)
		resArr = append(resArr, mStr)
	}
	return resArr
}
