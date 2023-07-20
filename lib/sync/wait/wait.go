package wait

import (
	"sync"
	"time"
)

// Wait  is similar with sync.WaitGroup
// enhance timeout feature
type Wait struct {
	wg sync.WaitGroup
}

// Add adds delta, which may be negative,to theWaitGroup counter
func (w *Wait) Add(delta int) {
	w.wg.Add(delta)
}

// Done decrements the WaitGroup counter by one
func (w *Wait) Done() {
	w.wg.Done()
}

// Wait blocks until the WaitGroup counter is zero
func (w *Wait) Wait() {
	w.wg.Wait()
}

// WaitWithTimeout blocks until the WaitGroup counter is zero
// timeout return true if timeout
func (w *Wait) WaitWithTimeout(timeout time.Duration) bool {
	c := make(chan bool, 1)
	go func() {
		defer close(c)
		w.wg.Wait()
		c <- true
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}
