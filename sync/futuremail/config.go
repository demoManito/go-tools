package futuremail

import "fmt"

type Formatter string

func (f Formatter) Format(params ...interface{}) string {
	return fmt.Sprintf(string(f), params...)
}

// expiredChannel 过期监听
// param: <__keyevent@[db]__:expired>
// example: <__keyevent@0__:expired>
var expiredChannel Formatter = "__keyevent@%d__:expired"

// consumeQueue 延时消费持久化队列
var consumePoolQueue Formatter = "consume:pool:queue:%d"

// consumeSingleQueue 独立消费持久化队列
var consumeSingleQueue Formatter = "consume:single:queue:%d"

func formatChannel(db int) string {
	return expiredChannel.Format(db)
}

func formatPoolQueue(db int) string {
	return consumePoolQueue.Format(db)
}

func formatSingleQueue(db int) string {
	return consumeSingleQueue.Format(db)
}

// Config redis config.
type Config struct {
	Addr           string `xml:"addr" yaml:"addr" json:"addr"`
	Pass           string `xml:"auth" yaml:"auth" json:"auth"`
	DB             int    `xml:"db" yaml:"db" json:"db"`
	PoolSize       int    `xml:"pool_size" yaml:"pool_size" json:"pool_size"`
	ChannelSize    int    `xml:"channel_size" yaml:"channel_size" json:"channel_size"`          // 过期监听 channel 容量
	RetryCount     int    `xml:"retry_count" yaml:"retry_count" json:"retry_count"`             // 重试次数
	OpenDeadletter bool   `xml:"open_deadletter" yaml:"open_deadletter" json:"open_deadletter"` // 是否开启死信队列

	channel     string
	poolQueue   string
	singleQueue string
}

func (c *Config) init() {
	c.channel = formatChannel(c.DB)
	c.poolQueue = formatPoolQueue(c.DB)
	c.singleQueue = formatSingleQueue(c.DB)
	if c.ChannelSize == 0 {
		c.ChannelSize = 100
	}
}
