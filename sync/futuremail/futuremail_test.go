package futuremail

import (
	"errors"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	os.Exit(0)
	// PASS: 需要本地 redis 服务器
}

func TestFutureMail(t *testing.T) {
	notify := New(&Config{Addr: "127.0.0.1:6379", DB: 0})
	defer notify.Close()

	// 注册延时消息处理方法
	i1 := notify.Inbox(time.Second, func(key string) (bool, error) {
		t.Log(key)
		if key != "one" {
			return true, nil
		}
		t.Logf("handler notify key: %s", key)
		notify.Ack(key, true)
		return false, nil
	})
	// 发送消息
	err := i1.SendMail("one")
	if err != nil {
		t.Error(err)
	}

	i2 := notify.Inbox(time.Second, func(key string) (bool, error) {
		t.Log(key)
		if key != "two" {
			return true, nil
		}
		t.Logf("handler notify key: %s", key)
		notify.Ack(key, true)
		return false, nil
	})
	err = i2.SendMail("two")
	if err != nil {
		t.Error(err)
	}

	time.Sleep(2 * time.Second)
}

func TestFutureMail2(t *testing.T) {
	fm := New(&Config{Addr: "127.0.0.1:6379", DB: 0})
	defer fm.Close()

	// 注册需要监听的延时消息
	fm.RegisterInbox([]NotifyFunc{
		func(key string) (bool, error) {
			t.Log(key)
			if key != "one" {
				return true, nil
			}
			t.Logf("handler fm key: %s", key)
			fm.Ack(key)
			return false, nil
		},
		func(key string) (bool, error) {
			t.Log(key)
			if key != "two" {
				return true, nil
			}
			t.Logf("handler fm key: %s", key)
			fm.Ack(key)
			return false, nil
		},
	}...)
	// 启动收件箱
	fm.StartInbox()

	// 发送消息
	err := fm.SendMail("one", time.Second)
	if err != nil {
		t.Error(err)
	}
	err = fm.SendMail("two", time.Second)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(3 * time.Second)
}

func TestFutureMail3(t *testing.T) {
	fm := New(&Config{Addr: "127.0.0.1:6379", DB: 0})
	defer fm.Close()

	// NOTICE:  重启后刷新重试 < 3 次消息，这里只会消费 two 的消息
	fm.backup("one")
	fm.backup("two")

	// 注册需要监听的延时消息
	fm.RegisterInbox([]NotifyFunc{
		func(key string) (bool, error) {
			t.Log(key)
			if key != "one" {
				return true, nil
			}
			t.Logf("handler fm key: %s", key)
			return false, errors.New("test one")
		},
		func(key string) (bool, error) {
			t.Log(key)
			if key != "two" {
				return true, nil
			}
			t.Logf("handler fm key: %s", key)
			fm.Ack(key)
			return false, nil
		},
	}...)
	// 启动收件箱
	fm.StartInbox()

	time.Sleep(3 * time.Second)
}
