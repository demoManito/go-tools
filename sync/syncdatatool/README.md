# syncdatatool
> 同步数据工具

```golang
	users := []struct {
		Name string
		Age  int
	}{
		{Name: "user1", Age: 12},
		{Name: "user2", Age: 14},
		{Name: "user3", Age: 16},
	}
	iterator, _ := iterator.NewIterator(users)

	// 1. 创建一个迭代器
	sync := New(nil, func() (interface{}, bool, error) {
		// 同步数据需要支持迭代器
		if iterator.HasNext() {
			return iterator.Next(), true, nil
		}
		return nil, false, nil
	})
	
	// 2. 开始同步
	sync.Run()
	
	// 3. 监听/处理同步数据
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
```

附: 迭代器 (https://github.com/demoManito/go-tools/tree/master/slice/iterator)