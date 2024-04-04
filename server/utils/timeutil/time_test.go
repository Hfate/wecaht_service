package timeutil

import (
	"fmt"
	"strconv"
	"testing"
)

func TestNow(t *testing.T) {
	ts := Now()
	if len(strconv.FormatUint(ts, 10)) != 13 {
		t.Errorf("Now() returns incorrect value: %v", ts)
	}
}

func TestStrToTimeStamp(t *testing.T) {
	//got := StrToTimeStamp("2006-01-02", "2006-01-02 15:04:05")
	fmt.Println(Format(GetCurTime(), FileNamePrefix))
	fmt.Println(Format(GetCurTime(), DateTimeLong))
	fmt.Println(Format(GetCurTime(), DateTimeShort))
	fmt.Println(Format(1577807999999, DateTimeLong))
	fmt.Println(Format(1577807999999, DateTimeLong))

	fmt.Println(StrToTimeStamp("2022-12-31 00:00:00", DateTimeLong))
	fmt.Println(StrToTime("2022-12-31 00:00:00", DateTimeLong))

}

func TestParseTime(t *testing.T) {
	fmt.Println(ParseTime("2020-01-01"))
	fmt.Println(ParseTime("2020/01/01"))
	fmt.Println(ParseTime("20200101"))
	fmt.Println(ParseTime("2020.01.01"))
	fmt.Println(ParseTime("43831"))
}

func TestAddMonths(t *testing.T) {

	fmt.Println(AddMonths(GetYearStartTime(GetCurTime()), 6))

}
