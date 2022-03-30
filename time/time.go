package time

import "time"

type timeNowFunc func() time.Time

var timeNow timeNowFunc

// SetTimeNow 修改当前时间为指定函数
// 传入 nil 重置 TimeNow() 函数
func SetTimeNow(f timeNowFunc) {
	timeNow = f
}

// TimeNow 获取当前时间，可用于 mock 时间
func TimeNow() time.Time {
	if timeNow != nil {
		return timeNow()
	}
	return time.Now()
}
