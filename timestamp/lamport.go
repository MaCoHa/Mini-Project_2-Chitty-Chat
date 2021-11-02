package timestamp

import (
	"sync/atomic"
)

type Clock struct {
	timestamp int64
}

func NewClock() *Clock {
	return &Clock{
		timestamp: 1,
	}
}

// Time is used to return the current value of the lamport clock
func (l Clock) GetTimestamp() int64 {
	return l.timestamp
}

// The current timestamp is incremented by 1 and the the new value is returned
func (l *Clock) Increment() int64 {
	atomic.AddInt64(&l.timestamp, 1)
	return l.timestamp
}

// Witness is called to update our local clock if necessary after
// witnessing a clock value received from another process
func (l *Clock) Witness(other int64) {
	// If the other value is old, we do not need to do anything
	cur := atomic.LoadInt64(&l.timestamp)

	var val int64
	if cur < other {
		val = other + 1
	} else {
		val = cur + 1
	}

	atomic.SwapInt64(&l.timestamp, val)
}
