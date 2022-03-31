package sync

import (
	"encoding/json"

	redis "github.com/go-redis/redis/v8"
)

const (
	actionAll       = "all"
	actionPublish   = "publish"
	actionSubscribe = "subscribe"
)

type DataBus struct {
	message chan *Message

	publish   *redis.Client // redis 发布
	subscribe *redis.Client // redis 订阅
}

type Message struct {
	Topic string          `json:"topic"`
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

// NewDataBus
func NewDataBus(config *DataBusConfig) *DataBus {
	delayed := &DataBus{}
	if config.Action == actionPublish || config.Action == actionAll {
		delayed.publish = redis.NewClient(&redis.Options{})
	}
	if config.Action == actionSubscribe || config.Action == actionAll {
		delayed.subscribe = redis.NewClient(&redis.Options{})
	}
	return delayed
}

// Bus receive message channel
func (db *DataBus) Bus() <-chan *Message {
	return db.message
}

// Send send message
func (db *DataBus) Send() error {
	return nil
}

// SendDelay send delay message
func (db *DataBus) SendDelay() error {
	return nil
}

// Close close databus
func (db *DataBus) Close() (err error) {
	close(db.message)
	switch {
	case db.publish != nil:
		err = db.publish.Close()
	case db.subscribe != nil:
		err = db.subscribe.Close()
	}
	return err
}
