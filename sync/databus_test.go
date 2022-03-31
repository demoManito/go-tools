package sync

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	msg := make(chan int)
	msg <- 1
	msg <- 1
	msg <- 1
	msg <- 1
	msg <- 1
	msg <- 1
	msg <- 1
	close(msg)
	for true {
		select {
		case m, ok := <-msg:
			fmt.Println(m, ok)
		default:
			fmt.Println("done")
			return
		}
	}
}
