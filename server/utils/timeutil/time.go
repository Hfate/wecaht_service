package timeutil

import (
	"github.com/golang-module/carbon"
	"strings"
	"time"
)

const (
	MonthEnglish      = "Jan"
	MonthYear         = "Jan 2006"
	YearMonthEN       = "2006 Jan"
	TIME_YMD          = "2006-01-02"
	TIMEYMD           = "20060102"
	FileNamePrefix    = "20060102150405"
	DateTimeShort     = "20060102150405"
	DateTimeLong      = "2006-01-02 15:04:05"
	GroupJoinDateTime = "02-01-2006"
	MinusSign         = "-"
)

var Loc *time.Location

func init() {
	Loc, _ = time.LoadLocation("Asia/Shanghai")
	carbon.SetLocation(Loc)
}

// Now gets current time in millisecond uint
func Now() uint64 {
	return uint64(time.Now().UnixNano()) / uint64(time.Millisecond)
}

func TimeToTimestamp(t time.Duration) uint64 {
	return uint64(t) / uint64(time.Millisecond)
}

// Format timestamp according to current server location
// param
//
//	ts: timestamp in millisecond uint
func DefaultFormat(ts int64) string {
	re := time.Unix(ts/1000, 0).In(Loc).Format(TIME_YMD)
	return re
}

// Format timestamp according to current server location
// param
//
//	ts: timestamp in millisecond uint
func Format(ts int64, format string) string {
	re := time.Unix(ts/1000, 0).In(Loc).Format(format)
	return re
}

func FormatDateTimeLong(ts int64) string {
	return Format(ts, DateTimeLong)
}

func Format2ExcelTime(ts int64) string {
	if ts == 0 {
		return MinusSign
	}
	return Format(ts, DateTimeLong)
}

func Format2ExcelYMDTime(ts int64) string {
	if ts == 0 {
		return MinusSign
	}
	return Format(ts, TIME_YMD)
}

func DiffAbsInDays(t1, t2 int64) int64 {
	return carbon.CreateFromTimestampMilli(t1).DiffAbsInDays(carbon.CreateFromTimestampMilli(t2)) + 1
}

// SubMonthLimit1 计算日期相差多少月
func SubMonthLimit1(t1, t2 int64) int64 {
	return carbon.CreateFromTimestampMilli(t1).DiffInMonths(carbon.CreateFromTimestampMilli(t2)) + 1
}

// SubMonth 计算日期相差多少月
func DiffInMonths(t1, t2 int64) int64 {
	return carbon.CreateFromTimestampMilli(t1).DiffInMonths(carbon.CreateFromTimestampMilli(t2))
}

// SubMonth 计算日期相差多少year
func DiffInYears(t1, t2 int64) int64 {
	return carbon.CreateFromTimestampMilli(t1).DiffAbsInYears(carbon.CreateFromTimestampMilli(t2))
}

// SubMonth 计算日期相差多少月
func SubMonth(t1 int64, months int) int64 {
	return carbon.CreateFromTimestampMilli(t1).SubMonths(months).TimestampMilli()
}

// 计算日期相差多少月
func SubDay(t1 int64, days int) int64 {
	return carbon.CreateFromTimestampMilli(t1).SubDays(days).TimestampMilli()
}

func SubYear(t1 int64) int64 {
	return carbon.CreateFromTimestampMilli(t1).SubYear().TimestampMilli()
}

func SubYears(t1 int64, years int) int64 {
	return carbon.CreateFromTimestampMilli(t1).SubYears(years).TimestampMilli()
}

// IsSameDay reports whether is same day.
// 是否是同一天
func IsSameDay(t1, t2 int64) bool {
	c := carbon.CreateFromTimestampMilli(t1)
	t := carbon.CreateFromTimestampMilli(t2)
	if c.IsInvalid() || t.IsInvalid() {
		return false
	}
	return c.Format("Ymd") == t.Format("Ymd")
}

func AddYears(t1 int64, year int) int64 {
	return carbon.CreateFromTimestampMilli(t1).AddYears(year).TimestampMilli()
}

func AddMonth(t1 int64) int64 {
	return carbon.CreateFromTimestampMilli(t1).AddMonth().TimestampMilli()
}

func AddDays(t1 int64, days int) int64 {
	return carbon.CreateFromTimestampMilli(t1).AddDays(days).TimestampMilli()
}

func AddMonths(t1 int64, months int) int64 {
	return carbon.CreateFromTimestampMilli(t1).AddMonths(months).TimestampMilli()
}

func GeTomorrowStartTime() int64 {
	return carbon.Tomorrow().StartOfDay().TimestampMilli()
}

func GetMonthEndTime(t int64) int64 {
	return carbon.CreateFromTimestampMilli(t).EndOfMonth().TimestampMilli()
}

func GetYearStartTime(t int64) int64 {
	return carbon.CreateFromTimestampMilli(t).StartOfYear().TimestampMilli()
}

func GetYearStartTimeByYear(year int) int64 {
	return carbon.CreateFromTimestampMilli(GetCurTime()).SetYear(year).StartOfYear().TimestampMilli()
}

func GetYearTimeRange(year int) (int64, int64) {
	start := carbon.CreateFromTimestampMilli(GetCurTime()).SetYear(year).StartOfYear().TimestampMilli()
	end := carbon.CreateFromTimestampMilli(GetCurTime()).SetYear(year).EndOfYear().TimestampMilli()
	return start, end
}

func GetDayTimeRange(t int64) (int64, int64) {
	start := carbon.CreateFromTimestampMilli(t).StartOfDay().TimestampMilli()
	end := carbon.CreateFromTimestampMilli(t).EndOfDay().TimestampMilli()
	return start, end
}

func GetYearEndTime(t int64) int64 {
	return carbon.CreateFromTimestampMilli(t).EndOfYear().TimestampMilli()
}

func GetYear(t int64) int {
	return carbon.CreateFromTimestampMilli(t).Year()
}

func GetDateStartTime(t int64) int64 {
	return carbon.CreateFromTimestampMilli(t).StartOfDay().TimestampMilli()
}

func GetDateEndTime(t int64) int64 {
	return carbon.CreateFromTimestampMilli(t).EndOfDay().TimestampMilli()
}

func GetMonthStartTime(t int64) int64 {
	return carbon.CreateFromTimestampMilli(t).StartOfMonth().TimestampMilli()
}

func GetMonthStartAndEnd(t int64) (int64, int64) {
	return GetMonthStartTime(t), GetMonthEndTime(t)
}

func GetCurTime() int64 {
	return GetNano() / 1e6
}

func GetCurDate() string {
	return strings.ReplaceAll(carbon.Now().ToDateString(), MinusSign, "")
}

func GetCurMonth() int64 {
	return GetMonthStartTime(GetCurTime())
}

func GetCurYear() int {
	return carbon.CreateFromTimestampMilli(GetCurTime()).Year()
}

func GetNano() int64 {
	return time.Now().UnixNano()
}

func StrToTimeStamp(timeStr string, format string) (int64, error) {
	// 默认按 "-" 分隔
	sep := MinusSign
	if !strings.Contains(timeStr, sep) {
		// 否则按 "/"
		sep = "/"
	}
	timeArr := strings.Split(timeStr, sep)

	for index, item := range timeArr {
		// 月 或者 日不带0
		if len(item) == 1 {
			timeArr[index] = "0" + item
		}
	}
	newTime := strings.Join(timeArr, sep)
	c := carbon.ParseByFormat(newTime, format)

	return c.TimestampMilli(), c.Error
}

func StrToTime(timeStr string, format string) (time.Time, error) {
	return time.ParseInLocation(format, timeStr, Loc)
}
