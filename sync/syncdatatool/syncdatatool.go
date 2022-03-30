package syncdatatool

import (
	"context"
	"sync/atomic"
)

type (
	pullDataFunc  func() (interface{}, bool, error)
	errHandleFunc func(error) bool
)

// SyncdataTool sync data tool
type SyncdataTool struct {
	ctx       context.Context
	data      chan *DataCarrier
	done      chan struct{}
	error     error
	pull      pullDataFunc  // pull data
	errHandle errHandleFunc // err handle when return <true> quit sync data now, it is not recommended to output logs here!

	exit bool // when <true> quit sync data
}

// DataCarrier 数据载体
type DataCarrier struct {
	Data   interface{}
	Error  error
	Offset int64
}

// New new sync data tool
// ctx: input nil will create a new context
// pull: pull data func
// errHandle (optional): errHandleFunc [func(error) bool] return <true> represent terminate data synchronization
func New(ctx context.Context, pull pullDataFunc, errHandle ...errHandleFunc) *SyncdataTool {
	if ctx == nil {
		ctx = context.Background()
	}
	synchronization := &SyncdataTool{
		ctx:  ctx,
		data: make(chan *DataCarrier),
		done: make(chan struct{}),
		pull: pull,
	}
	if len(errHandle) != 0 {
		synchronization.errHandle = errHandle[0]
	}
	return synchronization
}

// Data data chan
func (s *SyncdataTool) Data() <-chan *DataCarrier {
	return s.data
}

// Done done signal
func (s *SyncdataTool) Done() <-chan struct{} {
	return s.done
}

// Error sync error
func (s *SyncdataTool) Error() error {
	return s.error
}

func (s *SyncdataTool) errHandler(err error) {
	s.error = err
	if s.errHandle != nil {
		s.exit = s.errHandle(err)
		if !s.exit {
			s.error = nil
		}
	}
}

// Run start run sync data
func (s *SyncdataTool) Run() {
	var count int64
	go func() {
		defer close(s.data)
		defer close(s.done)
		for {
			data, hasMore, err := s.pull()
			if err != nil {
				s.errHandler(err)
			}
			if s.exit {
				break
			}
			if !hasMore {
				break
			}
			s.data <- &DataCarrier{
				Data:   data,
				Error:  err,
				Offset: atomic.AddInt64(&count, 1),
			}
		}
		s.done <- struct{}{}
	}()
}
