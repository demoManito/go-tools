package futuremail

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// Ack ack
// isSingle 是否 ack 独立收件箱消息，不传递默认 ack 总收件箱消息
func (fm *FutureMail) Ack(key string, isEngross ...bool) {
	var cmd *redis.IntCmd
	if len(isEngross) != 0 && isEngross[0] == true {
		cmd = fm.redis.SRem(fm.ctx, fm.config.singleQueue, key) // SET
	} else {
		cmd = fm.redis.ZRem(fm.ctx, fm.config.poolQueue, key) // ZSET
	}
	if cmd.Err() != nil {
		log.Printf("ack [%s], err (%s) \n", cmd.String(), cmd.Err())
	}
}

// backup 添加消息记录
func (fm *FutureMail) backup(key string, isEngross ...bool) {
	var cmd *redis.IntCmd
	if len(isEngross) != 0 && isEngross[0] == true {
		cmd = fm.redis.SAdd(fm.ctx, fm.config.singleQueue, key) // SET
	} else {
		cmd = fm.redis.ZAdd(fm.ctx, fm.config.poolQueue, &redis.Z{Member: key, Score: 0}) // ZSET
	}
	if cmd.Err() != nil {
		log.Printf("backup record [%s], err (%s) \n", cmd.String(), cmd.Err())
	}
}

// incrDeliveryFailed 递增消息投递失败次数
func (fm *FutureMail) incrDeliveryFailed(key string) {
	pipe := fm.redis.Pipeline()
	pipe.ZIncrBy(fm.ctx, fm.config.poolQueue, 1, key)
	cmd := pipe.ZScore(fm.ctx, fm.config.poolQueue, key)
	if _, err := pipe.Exec(fm.ctx); err != nil {
		log.Printf("incr delivery failed err (%s) \n", err.Error())
		return
	}
	go func() {
		switch score, _ := cmd.Result(); {
		case score <= float64(fm.config.RetryCount):
			// 重试队列
			fm.retry <- key
		case score > float64(fm.config.RetryCount):
			// 死信队列
			if fm.config.OpenDeadletter {
				fm.deadletter <- key
			}
		}
	}()
}

// redelivery 重新投递到期数据
// 投递重试次数小于 retry_count 消息
func (fm *FutureMail) redelivery() {
	keys, err := fm.redis.ZRevRangeByScore(fm.ctx, fm.config.poolQueue,
		&redis.ZRangeBy{Min: "0", Max: fmt.Sprint(fm.config.RetryCount)}).Result()
	if err != nil {
		log.Printf("redelivery err (%s) \n", err.Error())
		return
	}
	if len(keys) == 0 {
		return
	}

	flushCmds := make(map[string]*redis.DurationCmd)
	pipe := fm.redis.Pipeline()
	for _, key := range keys {
		flushCmds[key] = pipe.TTL(fm.ctx, key)
	}
	if _, err = pipe.Exec(fm.ctx); err != nil {
		log.Printf("redelivery ttl err (%s) \n", err.Error())
		return
	}
	if len(flushCmds) == 0 {
		return
	}

	go func() {
		for key, cmd := range flushCmds {
			duration, _ := cmd.Result()
			if duration == -2*time.Nanosecond {
				fm.flush <- key
			}
		}
	}()
}
