package futuremail

import (
	"testing"
	"time"
)

func TestFutureMail(t *testing.T) {
	notify := New(&Config{Addr: "127.0.01:6379", DB: 0})
	defer notify.Close()

	i1 := notify.Inbox(time.Second, func(key string) (bool, error) {
		t.Log(key)
		if key != "one" {
			return false, nil
		}
		t.Logf("handler notify key: %s", key)
		return true, nil
	})
	i2 := notify.Inbox(time.Second, func(key string) (bool, error) {
		t.Log(key)
		if key != "two" {
			return false, nil
		}
		t.Logf("handler notify key: %s", key)
		return true, nil
	})

	err := i1.SendMail("one")
	if err != nil {
		t.Error(err)
	}
	err = i2.SendMail("two")
	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)
}

func TestFutureMail2(t *testing.T) {
	notify := New(&Config{Addr: "127.0.01:6379", DB: 0})
	defer notify.Close()
	notify.InboxRegister(10, []NotifyFunc{
		func(key string) (bool, error) {
			t.Log(key)
			if key != "one" {
				return false, nil
			}
			t.Logf("handler notify key: %s", key)
			return true, nil
		},
		func(key string) (bool, error) {
			t.Log(key)
			if key != "two" {
				return false, nil
			}
			t.Logf("handler notify key: %s", key)
			return true, nil
		},
	}...)
	notify.InboxPool()

	err := notify.SendMail("one", time.Second)
	if err != nil {
		t.Error(err)
	}
	err = notify.SendMail("two", time.Second)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)
}
