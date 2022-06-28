package time

import (
	"fmt"
	"strconv"
	"strings"
)

type Clock struct {
	Hour int
	Min  int
	Sec  int

	err error
}

// ParseClock 格式化 hh:mm:ss 格式和空字符串，用作时刻比较
func ParseClock(hhmmss string) *Clock {
	if hhmmss == "" {
		return nil
	}
	pre := strings.Split(hhmmss, ":")
	if len(pre) != 3 {
		panic(fmt.Sprintf("unsupported clock format: %s", hhmmss))
	}
	clock := new(Clock)
	clock.Hour, clock.err = strconv.Atoi(pre[0])
	if clock.err != nil {
		panic(fmt.Errorf("hour format error: %s", clock.err))
	}
	clock.Min, clock.err = strconv.Atoi(pre[1])
	if clock.err != nil {
		panic(fmt.Errorf("min format error: %s", clock.err))
	}
	clock.Sec, clock.err = strconv.Atoi(pre[2])
	if clock.err != nil {
		panic(fmt.Errorf("sec format error: %s", clock.err))
	}
	return clock
}

// 将时分秒换算成秒
func (c *Clock) second() int {
	return c.Hour*3600 + c.Min*60 + c.Sec
}

// Less c less c2 return true, otherwise return false
// c 早于 c2 返回 true
func (c *Clock) Less(c2 *Clock) bool {
	return c.second() < c2.second()
}

// Greater c greater c2 return true, otherwise return false
// c 晚于 c2 返回 true
func (c *Clock) Greater(c2 *Clock) bool {
	return c.second() > c2.second()
}

// GreaterOrEqual c 比 c2 更晚或相等
func (c Clock) GreaterOrEqual(c2 *Clock) bool {
	return !c.Less(c2)
}

// LessOrEqual c 比 c2 更早或相等
func (c Clock) LessOrEqual(c2 *Clock) bool {
	return !c.Greater(c2)
}
