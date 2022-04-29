# futuremail
> 基于 redis 延时消息系统

> `redis.conf` 打开：`notify-keyspace-events Ex` 配置

## 独立的延时消息
> `notify.Inbox()` 会启动一个协程监听程序，处理完成消息或超时未处理则会自动关闭改协程，不会产生内存泄漏的风险

```golang
	fm := New(&Config{Addr: "127.0.01:6379", DB: 0})
	defer fm.Close()

	// 1. 注册延时消息处理方法
	i1 := fm.Inbox(time.Second, func(key string) (bool, error) {
		t.Log(key)
		if key != "one" {
			return false, nil
		}
		t.Logf("handler notify key: %s", key)
		// 3. ack
		fm.Ack(key, true)
		return true, nil
	})
	// 2. 发送消息
	err := i1.SendMail("one")
	if err != nil {
		t.Error(err)
	}
```

## 全局延时消息
> 需要在项目初始化时开启监听程序 `fm.StartInbox()`，发送消息会放到全局监听器中
1. 启动会重试所有的 KEY
2. 重试队列默认重试 3 次
3. 重试超过 3 次，将放入死信队列（TODO）
```golang
	fm := New(&Config{Addr: "127.0.01:6379", DB: 0})
	defer fm.Close()

	// 1. 注册需要监听的延时消息 
	fm.RegisterInbox([]NotifyFunc{
		func(key string) (bool, error) {
			t.Log(key)
			if key != "one" {
				return false, nil
			}
			t.Logf("handler notify key: %s", key)
			// 4. ack
			fm.Ack(key)
			return true, nil
		},
		func(key string) (bool, error) {
			t.Log(key)
			if key != "two" {
				return false, nil
			}
			t.Logf("handler fm key: %s", key)
			// 4. ack 
			fm.Ack(key)
			return true, nil
		},
	}...)
	// 2. 启动收件箱 
	fm.StartInbox()

	// 3. 发送消息
	err := fm.SendMail("one", time.Second)
	if err != nil {
		t.Error(err)
	}
	err = fm.SendMail("two", time.Second)
	if err != nil {
		t.Error(err)
	}
```
