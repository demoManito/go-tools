package syncdatatool

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/demoManito/go-tools/slice/iterator"
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
	iterator, err := iterator.New(users)
	require.NoError(err)
	require.NotNil(iterator)

	sync := New(context.TODO(), func() (interface{}, bool, error) {
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
	return rand.Intn(10) == 8
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
			return 1, true, fmt.Errorf("测试 %d", count)
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
