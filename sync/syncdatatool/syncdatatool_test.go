package syncdatatool

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"go-tools/slice/iterator"
)

func TestSyncdataToolMockData(t *testing.T) {
	require := require.New(t)

	users := []struct {
		Name string
		Age  int
	}{
		{Name: "user1", Age: 12},
		{Name: "user2", Age: 14},
		{Name: "user3", Age: 16},
	}
	iterator, err := iterator.NewIterator(users)
	require.NoError(err)
	require.NotNil(iterator)

	sync := New(nil, func() (interface{}, bool, error) {
		if iterator.HasNext() {
			return iterator.Next(), true, nil
		}
		return nil, false, nil
	})
	sync.Run()
	for {
		select {
		case data := <-sync.Data():
			if data.Error != nil {
				t.Error(data.Error)
			}
			t.Logf("%d: %v", data.Offset, data.Data)
		case <-sync.Done():
			return
		}
	}
}

func mockHasMore() bool {
	if rand.Intn(10) == 8 {
		return false
	}
	return true
}

func TestSynchronizationData_Run(t *testing.T) {
	sync := New(context.Background(), func() (interface{}, bool, error) {
		if !mockHasMore() {
			return nil, false, nil
		}
		return 1, true, nil
	})
	sync.Run()
loop1:
	for {
		select {
		case data := <-sync.Data():
			t.Logf("data: %d, %d", data.Offset, data.Data)
		case <-sync.Done():
			t.Logf("done; err: %s", sync.Error())
			break loop1
		}
	}

	t.Log("===============")

	count := 0
	sync = New(context.Background(),
		func() (interface{}, bool, error) {
			count++
			if !mockHasMore() {
				return nil, false, nil
			}
			return 1, true, errors.New(fmt.Sprintf("测试 %d", count))
		},
		func(err error) bool {
			t.Logf("err: %s", err)
			return false
		},
	)
	sync.Run()
loop2:
	for {
		select {
		case data := <-sync.Data():
			if data.Error != nil {
				t.Logf("data err: %s", data.Error)
				continue
			}
			t.Logf("data: %d, %d", data.Offset, data.Data)
		case <-sync.Done():
			t.Logf("done; err: %s", sync.Error())
			break loop2
		}
	}
}
