package timers

import "sync"

// Timers ...
type Timers struct {
	Items         map[int32]*Timer
	NextTimeoutID int

	sync.RWMutex
}

// New ...
func New() *Timers {
	t := new(Timers)

	return t
}
