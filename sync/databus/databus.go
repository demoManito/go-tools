package databus

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis"
)

const (
	redisAll       = "all"
	redisPublish   = "publish"
	redisSubscribe = "subscribe"

	databusOpen   = int32(0)
	databusClosed = int32(1)
)

// DataBus
type DataBus struct {
	message chan *Message

	publish   *redis.Client // redis 发布
	subscribe *redis.Client // redis 订阅

	lock   sync.RWMutex
	marked map[int32]int64 // 消费标记
	closed int32           // databus 是否被关闭

	config *DataBusConfig
}

// Message 消息体
type Message struct {
	Topic     string          `json:"topic"`
	Key       string          `json:"key"`
	Value     json.RawMessage `json:"value"`
	Partition int32           `json:"partition"`
	Offset    int64           `json:"offset"`

	databus *DataBus // ACK 需要记录偏移量
}

// Ack 记录偏移量
func (m *Message) Ack() {
	m.databus.lock.Lock()
	defer m.databus.lock.Unlock()
	if m.Offset > m.databus.marked[m.Partition] {
		m.databus.marked[m.Partition] = m.Offset
	}
}

// NewDataBus
func NewDataBus(config *DataBusConfig) *DataBus {
	if config.Buffer == 0 {
		config.Buffer = 1000
	}
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}
	delayed := &DataBus{
		message: make(chan *Message, config.Buffer),
		marked:  make(map[int32]int64),
		closed:  databusOpen,
		config:  config,
	}
	if config.Action == redisPublish || config.Action == redisAll {
		delayed.publish = redis.NewClient(&redis.Options{
			Addr:         config.Addr,
			Password:     config.Pass,
			DB:           0, // use default DB
			DialTimeout:  time.Second * 1,
			IdleTimeout:  time.Second * 10,
			ReadTimeout:  time.Second * 2,
			WriteTimeout: time.Second * 1,
			PoolSize:     1,
		})
	}
	if config.Action == redisSubscribe || config.Action == redisAll {
		delayed.subscribe = redis.NewClient(&redis.Options{
			Addr:         config.Addr,
			Password:     config.Pass,
			DB:           0, // use default DB
			DialTimeout:  time.Second * 1,
			IdleTimeout:  time.Second * 10,
			ReadTimeout:  time.Second * 2,
			WriteTimeout: time.Second * 1,
			PoolSize:     config.PoolSize,
		})
		go delayed.waitingTake()
	}
	return delayed
}

// waitingTake waiting take redis subscribe in channel data
func (db *DataBus) waitingTake() {
	var (
		err       error
		retry     int
		redisConn *redis.Conn
		acked     = make(map[int32]int64)
		ack       = make(map[int32]int64)
	)
	for {
		if atomic.LoadInt32(&db.closed) == databusClosed {
			if redisConn != nil {
				_ = redisConn.Close()
				redisConn = nil
			}
			close(db.message)
			return
		}
		if err != nil {
			// TODO: 退避算法
			// time.Sleep(bk.Backoff(retry))
			retry++
		} else {
			retry = 0
		}

		if redisConn == nil {
			db.lock.Lock()
			db.marked = make(map[int32]int64)
			db.lock.Unlock()
			// db.subscribe.Conn(nil)
			err = nil
		}
		db.lock.RLock()
		for k, v := range db.marked {
			// 是否在已经消费完的 map 里；如果不在就需要存到 ack 里
			if o, ok := acked[k]; o != v || !ok {
				ack[k] = v
			}
		}
		db.lock.RUnlock()
		// TODO: redis: get pb? where did this <pb> key come from?
	}
}

// Bus receive message channel
func (db *DataBus) Bus() <-chan *Message {
	return db.message
}

// Send send message
func (db *DataBus) Send(ctx context.Context, key string, val interface{}) error {
	return db.send(ctx, key, val, "")
}

// SendDelay send delay message
func (db *DataBus) SendDelay(ctx context.Context, key string, val interface{}, deliverAt time.Time) error {
	return db.send(ctx, key, val, strconv.FormatInt(deliverAt.UnixNano()/1e6, 10))
}

// Close close databus
func (db *DataBus) Close() (err error) {
	// CompareAndSwapInt32(addr *int32, old, new int32)
	ok := atomic.CompareAndSwapInt32(&db.closed, databusOpen, databusClosed)
	if !ok {
		return
	}
	if db.publish != nil {
		err = db.publish.Close()
	}
	if db.subscribe != nil {
		err = db.subscribe.Close()
	}
	return err
}

func (db *DataBus) send(ctx context.Context, key string, val interface{}, deliverAt string) error {
	value, err := json.Marshal(val)
	if err != nil {
		return err
	}
	switch {
	case len(deliverAt) == 0: // 非延时
		err = db.publish.Do(ctx, "SET", key, value).Err()
		if err != nil {
			return err
		}
		return nil
	case len(deliverAt) != 0: // 延时
		// TODO: proto data
		field, err := json.Marshal(deliverAt)
		if err != nil {
			return err
		}
		err = db.publish.Do(ctx, "HSET", key, field, value).Err()
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
