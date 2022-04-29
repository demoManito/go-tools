# liunx 进程信号监听

## 死亡信号监听
```golang
func main() {
    // 省略业务代码

    // 此处会阻塞，等待系统向 go 主进程发送结束信号
    // 例如: kill [pid]
    <-DeathProcessWatch().Signal()    
}
```

## 自定义信号监听
```golang
func main() {
    sig := DiyProcessWatch(context.Background(), syscall.SIGTERM)
    defer sig.Close()
    
    // 模拟一直监听自定义信号接收
    for {
        select {
            case s := <-sig.Signal():
            // 监听到信号后的业务
            // ......
        }
    }
}
```