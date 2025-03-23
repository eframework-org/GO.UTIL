// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XTime

import (
	"time"
)

const (
	FormatFull = "2006-01-02 15:04:05 +0800 CST" // 时间格式化模板（标准）
	FormatLite = "2006-01-02 15:04:05"           // 时间格式化模板（简易）
	FormatFile = "2006-01-02_15_04_05"           // 时间格式化模板（文件）
)

const (
	Second1  int = 1       // 1秒
	Second5  int = 5       // 5秒
	Second10 int = 10      // 10秒
	Second15 int = 15      // 15秒
	Second20 int = 20      // 20秒
	Second25 int = 25      // 25秒
	Second30 int = 30      // 30秒
	Second35 int = 35      // 35秒
	Second40 int = 40      // 40秒
	Second45 int = 45      // 45秒
	Second50 int = 50      // 50秒
	Second55 int = 55      // 55秒
	Minute1  int = 60      // 1分钟
	Minute2  int = 120     // 2分钟
	Minute3  int = 180     // 3分钟
	Minute4  int = 240     // 4分钟
	Minute5  int = 300     // 5分钟
	Minute6  int = 360     // 6分钟
	Minute7  int = 420     // 7分钟
	Minute8  int = 480     // 8分钟
	Minute9  int = 540     // 9分钟
	Minute10 int = 600     // 10分钟
	Minute12 int = 720     // 12分钟
	Minute15 int = 900     // 15分钟
	Minute20 int = 1200    // 20分钟
	Minute25 int = 1500    // 25分钟
	Minute30 int = 1800    // 30分钟
	Minute35 int = 2100    // 35分钟
	Minute40 int = 2400    // 40分钟
	Minute45 int = 2700    // 45分钟
	Minute50 int = 3000    // 50分钟
	Minute55 int = 3300    // 55分钟
	Hour1    int = 3600    // 1小时
	Hour2    int = 7200    // 2小时
	Hour3    int = 10800   // 3小时
	Hour4    int = 14400   // 4小时
	Hour5    int = 18000   // 5小时
	Hour6    int = 21600   // 6小时
	Hour7    int = 25200   // 7小时
	Hour8    int = 28800   // 8小时
	Hour9    int = 32400   // 9小时
	Hour10   int = 36000   // 10小时
	Hour11   int = 39600   // 11小时
	Hour12   int = 43200   // 12小时
	Hour13   int = 46800   // 13小时
	Hour14   int = 50400   // 14小时
	Hour15   int = 54000   // 15小时
	Hour16   int = 57600   // 16小时
	Hour17   int = 61200   // 17小时
	Hour18   int = 64800   // 18小时
	Hour19   int = 68400   // 19小时
	Hour20   int = 72000   // 20小时
	Hour21   int = 75600   // 21小时
	Hour22   int = 79200   // 22小时
	Hour23   int = 82800   // 23小时
	Day1     int = 86400   // 1天
	Day2     int = 172800  // 2天
	Day3     int = 259200  // 3天
	Day4     int = 345600  // 4天
	Day5     int = 432000  // 5天
	Day6     int = 518400  // 6天
	Day7     int = 604800  // 7天
	Day8     int = 691200  // 8天
	Day9     int = 777600  // 9天
	Day10    int = 864000  // 10天
	Day15    int = 1296000 // 15天
	Day20    int = 1728000 // 20天
	Day30    int = 2592000 // 30天
)

// GetMicrosecond 获取当前时间的微秒级时间戳。
// 返回以微秒为单位的时间戳。
func GetMicrosecond() int {
	ltime := time.Now().UnixNano() / 1e3
	return int(ltime)
}

// GetMillisecond 获取当前时间的毫秒级时间戳。
// 返回以毫秒为单位的时间戳。
func GetMillisecond() int {
	ltime := time.Now().UnixNano() / 1e6
	return int(ltime)
}

// GetTimestamp 获取当前时间的秒级时间戳。
// 返回以秒为单位的时间戳。
func GetTimestamp() int {
	ltime := time.Now().Unix()
	time.Unix(0, 0).Format("")
	return int(ltime)
}

// NowTime 获取当前时间的time.Time对象。
// 返回表示当前时间的time.Time对象。
func NowTime() time.Time {
	return time.Now()
}

// ToTime 将秒级时间戳转换为time.Time对象。
// timestamp 是要转换的秒级时间戳。
// 返回对应时间戳的time.Time对象。
func ToTime(timestamp int) time.Time {
	return time.Unix(int64(timestamp), 0)
}

// TimeToZero 计算指定时间到下一个零点的秒数。
// timestamp 是可选的指定时间戳，如果不提供则使用当前时间。
// 返回到下一个零点的秒数。
func TimeToZero(timestamp ...int) int {
	t := 0
	if len(timestamp) == 1 {
		t = timestamp[0]
	} else {
		t = GetTimestamp()
	}
	return Day1 - (t+Hour8)%Day1
}

// ZeroTime 获取指定时间的零点时间戳。
// timestamp 是可选的指定时间戳，如果不提供则使用当前时间。
// 返回零点时间的秒级时间戳。
func ZeroTime(timestamp ...int) int {
	t := 0
	if len(timestamp) == 1 {
		t = timestamp[0]
	} else {
		t = GetTimestamp()
	}
	return t - (t+Hour8)%Day1
}

// Format 将时间戳格式化为指定格式的字符串。
// timestamp 是要格式化的秒级时间戳。
// format 是格式化模板，可以使用预定义的 FormatFull、FormatLite 或 FormatFile。
// 返回格式化后的时间字符串。
func Format(timestamp int, format string) string {
	return time.Unix(int64(timestamp), 0).Format(format)
}
