package timers

import (
	"sync"
	"time"
)

// Timer ...
type Timer struct {
	ID    int32
	Delay int32
	Done  chan bool

	cb func()

	sync.Once
}

// Clear ...
func (t *Timer) Clear() {
	t.Do(func() {
		close(t.Done)
	})
}

// Start ...
func (t *Timer) Start() {
	go func() {
		defer t.Clear()

		ticker := time.NewTicker(time.Duration(t.Delay) * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-t.Done:
				return
			case <-ticker.C:
				if t.cb != nil {
					t.cb()
				}
			}
		}
	}()
}
