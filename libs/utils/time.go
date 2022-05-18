package utils

import (
	"time"
)

// NowTimeStringCN get now time format cn 2006-01-02 15:04:05
func NowTimeStringCN() string {
	var cstZone = time.FixedZone("CST", 8*3600) // UTC/GMT +08:00
	t := time.Now()
	return t.In(cstZone).Format("2006-01-02 15:04:05")
}

// NowTimeStringCN2 get now time format cn 20060102150405
func NowTimeStringCN2() string {
	var cstZone = time.FixedZone("CST", 8*3600) // UTC/GMT +08:00
	t := time.Now()
	return t.In(cstZone).Format("20060102150405")
}

// NowDateStringCN get now time format cn 20060102
func NowDateStringCN() string {
	var cstZone = time.FixedZone("CST", 8*3600) // UTC/GMT +08:00
	t := time.Now()
	return t.In(cstZone).Format("20060102")
}

// UnixTimstampSecond 返回秒级时间戳
func UnixTimstampSecond() int64 {
	return time.Now().Unix()
}

// UnixTimstampMillisecond 返回毫秒级时间戳
func UnixTimstampMillisecond() int64 {
	return time.Now().UnixNano() / 1000000
}

// UnixTimstampNanosecond 返回Nano级时间戳
func UnixTimstampNanosecond() int64 {
	return time.Now().UnixNano()
}

// ParseTimeFromString 把字符串日期解析成 time.Time 类型
func ParseTimeFromString(t string) (time.Time, error) {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	return time.ParseInLocation("2006-01-02 15:04:05", t, cstZone)
}
