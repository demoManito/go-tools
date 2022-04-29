package processsignal

import (
	"context"
	"syscall"
	"testing"
)

func TestDeathProcessWatch(t *testing.T) {
	/*
		example：
		```go
			func main() {
				// business code
				// ...
				<-DeathProcessWatch().Signal()
			}
		```

	*/
	<-DeathProcessWatch().Signal()
	t.Log("kill go test pid success")
}

func TestDiyProcessWatch(t *testing.T) {
	/*
		1. ps -Cf | grep tools
		2. find pid
		3. kill [pid]
		4. log print terminated
	*/
	sig := DiyProcessWatch(context.Background(), syscall.SIGTERM)
	defer sig.Close()
	for {
		select {
		case s := <-sig.Signal():
			t.Logf("process signal %s", s)
		}
	}
}
