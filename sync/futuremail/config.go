package futuremail

import (
	"fmt"
)

type Formatter string

// expiredChannel 过期监听
// param: <__keyevent@[db]__:expired>
// example: <__keyevent@0__:expired>
var expiredChannel Formatter = "__keyevent@%d__:expired"

func (f Formatter) Format(params ...interface{}) string {
	return fmt.Sprintf(string(f), params...)
}

func formatChannel(db int) string {
	return expiredChannel.Format(db)
}

// Config redis config.
type Config struct {
	Addr     string `xml:"addr" yaml:"addr" json:"addr"`
	Pass     string `xml:"auth" yaml:"auth" json:"auth"`
	DB       int    `xml:"db" yaml:"db" json:"db"`
	PoolSize int    `xml:"pool_size" yaml:"pool_size" json:"pool_size"`
	channel  string
}
