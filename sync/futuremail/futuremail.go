package futuremail

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type NotifyFunc func(key string) (bool, error)

// FutureMail 延时通知
type FutureMail struct {
	ctx        context.Context
	redis      *redis.Client
	flush      chan string
	retry      chan string
	deadletter chan string // TODO: 暂无实现

	config   *Config
	notifies []NotifyFunc
}

func New(config *Config) *FutureMail {
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}

	futureMail := &FutureMail{
		ctx:    context.Background(),
		flush:  make(chan string, 100),
		retry:  make(chan string, 100),
		config: config,
	}
	if config.OpenDeadletter {
		futureMail.deadletter = make(chan string, 100)
	}
	futureMail.redis = redis.NewClient(&redis.Options{
		Addr: config.Addr,
		// Password:     config.Pass, // auth password set
		DB:           config.DB, // use default DB
		PoolSize:     config.PoolSize,
		DialTimeout:  time.Second * 1,
		IdleTimeout:  time.Second * 10,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 1,
	})
	futureMail.config.init()
	return futureMail
}

// RegisterInbox 注册消费消息的所有方法
func (fm *FutureMail) RegisterInbox(notify ...NotifyFunc) {
	fm.notifies = append(fm.notifies, notify...)
}

func (fm *FutureMail) Close() error {
	if err := fm.redis.Close(); err != nil {
		log.Printf("redis close err (%s) \n", err)
	}
	close(fm.flush)
	close(fm.retry)
	if fm.config.OpenDeadletter {
		close(fm.deadletter)
	}
	return nil
}

// ------ 全局监听器 ------

// SendMail 向 [StartInbox] 中发送
func (fm *FutureMail) SendMail(key string, delayDuration time.Duration) error {
	err := fm.redis.Set(fm.ctx, key, true, delayDuration).Err()
	if err != nil {
		return err
	}
	fm.backup(key)
	return nil
}

// StartInbox 收件箱总部
func (fm *FutureMail) StartInbox() {
	fm.redelivery() // redelivery historical data
	pubsub := fm.redis.Subscribe(fm.ctx, fm.config.channel)
	go func() {
		for {
			select {
			case data, ok := <-pubsub.Channel(redis.WithChannelSize(fm.config.ChannelSize)):
				if !ok {
					break
				}
				fm.consume(data.Payload)
			case key, ok := <-fm.retry: // retry channel
				if !ok {
					break
				}
				fm.consume(key)
			case key, ok := <-fm.flush: // flush historical data channel
				if !ok {
					break
				}
				fm.consume(key)
			case <-fm.ctx.Done():
				log.Printf("databus done err (%s) \n", fm.ctx.Err())
				return
			}
		}
		// TODO: case dead := <-fm.deadletter: // dead letter channel
	}()
}

func (fm *FutureMail) consume(data string) {
	for _, notify := range fm.notifies {
		hit, err := notify(data)
		if err != nil {
			fm.incrDeliveryFailed(data) // 递增失败次数
			break
		}
		if !hit {
			break
		}
	}
}

// ------ 独立监听器 ------
// 一个消息绑定一个 channel

// Inbox
type Inbox struct {
	fm            *FutureMail
	cancel        context.CancelFunc
	delayDuration time.Duration
}

// SendMail 向 [Inbox] 中发送
func (i *Inbox) SendMail(key string) error {
	err := i.fm.redis.Set(i.fm.ctx, key, true, i.delayDuration).Err()
	if err != nil {
		i.cancel()
		return err
	}
	i.fm.backup(key, true)
	return nil
}

// Inbox 绑定独立的消息接收器
// 协程清除：消费完成自动清除该协程，若延时消息发送失败也会清除该协程，无泄漏风险
// 协程超时时间：过期时间 + 5min
func (fm *FutureMail) Inbox(delayDuration time.Duration, notify NotifyFunc, errhandler ...func(error)) *Inbox {
	ctx, cancel := context.WithTimeout(context.Background(), delayDuration+5*time.Minute)
	pubsub := fm.redis.Subscribe(ctx, fm.config.channel)
	go func() {
		for {
			select {
			case data := <-pubsub.Channel():
				hit, err := notify(data.Payload)
				if err != nil {
					for _, handler := range errhandler {
						handler(err)
					}
				}
				if hit {
					continue
				}
				cancel()
				return
			case <-ctx.Done():
				log.Printf("databus done; err (%s) \n", ctx.Err())
				return
			}
		}
	}()
	return &Inbox{
		fm:            fm,
		cancel:        cancel,
		delayDuration: delayDuration,
	}
}
