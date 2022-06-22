package processsignal

import (
	"context"
	"os"
	"syscall"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(0)
}

// kill pid:
/*
	1. ps -Cf | grep tools
	2. find pid
	3. kill [pid]
	4. log print terminated
*/
// use:
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
func TestDeathProcessWatch(t *testing.T) {
	<-DeathProcessWatch().Signal()
	t.Log("kill go test pid success")
}

func TestDiyProcessWatch(t *testing.T) {
	sig := DiyProcessWatch(context.Background(), syscall.SIGTERM)
	defer sig.Close()
	for {
		select {
		case s := <-sig.Signal():
			t.Logf("process signal %s", s)
		}
	}
}
