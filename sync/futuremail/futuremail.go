package futuremail

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type NotifyFunc func(key string) (bool, error)

// FutureMail 延时通知
type FutureMail struct {
	ctx   context.Context
	redis *redis.Client

	config *Config

	notifies    []NotifyFunc
	channelSize int
}

func New(config *Config) *FutureMail {
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}

	futureMail := &FutureMail{
		ctx:    context.Background(),
		config: config,
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
	futureMail.config.channel = formatChannel(config.DB)
	// futureMail.subscribe = futureMail.redis.Subscribe(futureMail.ctx, futureMail.config.channel).Channel()
	return futureMail
}

// InboxRegister 注册消费消息的所有方法
func (fm *FutureMail) InboxRegister(channelSize int, notify ...NotifyFunc) {
	fm.channelSize = channelSize
	if fm.channelSize == 0 {
		fm.channelSize = 10
	}
	fm.notifies = append(fm.notifies, notify...)
}

func (fm *FutureMail) Close() error {
	err := fm.redis.Close()
	if err != nil {
		return err
	}
	return nil
}

// ------ 收件箱 - 总部 ------

// SendMail 向 [InboxPool] 中发送
func (fm *FutureMail) SendMail(key string, delayDuration time.Duration) error {
	err := fm.redis.Set(fm.ctx, key, true, delayDuration).Err()
	if err != nil {
		return err
	}
	return nil
}

// InboxPool 收件箱总部
func (fm *FutureMail) InboxPool() {
	pubsub := fm.redis.Subscribe(fm.ctx, fm.config.channel)
	go func() {
		for {
			select {
			case data := <-pubsub.ChannelSize(fm.channelSize):
				for _, notify := range fm.notifies {
					execute, err := notify(data.Payload)
					if err != nil {
						continue
					}
					if execute {
						break
					}
				}
			case <-fm.ctx.Done():
				log.Printf("databus done; err (%s)", fm.ctx.Err())
				return
			}
		}
	}()
}

// ------ 收件箱 - 分部 ------

// Inbox
type Inbox struct {
	fn            *FutureMail
	cancel        context.CancelFunc
	delayDuration time.Duration
}

// SendMail 向 [Inbox] 中发送
func (i *Inbox) SendMail(key string) error {
	err := i.fn.redis.Set(i.fn.ctx, key, true, i.delayDuration).Err()
	if err != nil {
		i.cancel()
		return err
	}
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
				exit, err := notify(data.Payload)
				if err != nil {
					for _, handler := range errhandler {
						handler(err)
					}
				}
				if !exit {
					continue
				}
				cancel()
				return
			case <-ctx.Done():
				log.Printf("databus done; err (%s)", ctx.Err())
				return
			}
		}
	}()
	return &Inbox{
		fn:            fm,
		cancel:        cancel,
		delayDuration: delayDuration,
	}
}
