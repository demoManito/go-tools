package processsignal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// ProcessSignal ps
type ProcessSignal struct {
	signal chan os.Signal
}

// signal
func (ps *ProcessSignal) Signal() <-chan os.Signal {
	return ps.signal
}

// Close close
func (ps *ProcessSignal) Close() {
	close(ps.signal)
}

func newPS() *ProcessSignal {
	return &ProcessSignal{signal: make(chan os.Signal)}
}

// DeathProcessWatch used to block the mainline program
func DeathProcessWatch() *ProcessSignal {
	ps := newPS()
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM)
		ps.signal <- <-sig
	}()
	return ps
}

// DiyProcessWatch
func DiyProcessWatch(ctx context.Context, sigs ...os.Signal) *ProcessSignal {
	ps := newPS()
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, sigs...)
		for {
			select {
			case s, ok := <-sig:
				if !ok {
					return
				}
				ps.signal <- s
			case <-ctx.Done():
				return
			}
		}
	}()
	return ps
}
