# 限流器
> 基于标准库 `golang.org/x/time/rate` 令牌桶的多粒度限流器

## 使用
```golang
    limiter := New(
    	rate.NewLimiter(Per(2, time.Second), 2), // 细度 2个/秒
    	rate.NewLimiter(Per(1000000, time.Hour), 1000000), // 细度 1000000个/小时
    )
    
    // 限流等待
    limiter.Wait() 
```


## 兼容中间件:
##### gin:
```golang
    limter := New(rate.NewLimiter(Per(1, time.Second), 1))
    r := gin.New()
    r.Use(lgin.NewMiddleware(limter))
    // ....
    r.Run()
```
