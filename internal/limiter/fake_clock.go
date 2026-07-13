package limiter

import (
	"sync"
	"time"
)

// fakeClock is a monotonic clock for deterministic tests.
type fakeClock struct {
	mu    sync.Mutex
	now   time.Time
	cond  *sync.Cond
	sleep []sleepReq
}

type sleepReq struct {
	until time.Time
	ch    chan struct{}
}

func newFakeClock(t time.Time) *fakeClock {
	fc := &fakeClock{now: t}
	fc.cond = sync.NewCond(&fc.mu)
	return fc
}

func (fc *fakeClock) Now() time.Time {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	return fc.now
}

// advance moves the clock forward and wakes waiting goroutines.
func (fc *fakeClock) advance(d time.Duration) {
	fc.mu.Lock()
	fc.now = fc.now.Add(d)
	for _, sr := range fc.sleep {
		if fc.now.Compare(sr.until) >= 0 {
			select {
			case <-sr.ch:
			default:
				close(sr.ch)
			}
		}
	}
	fc.cond.Broadcast()
	fc.mu.Unlock()
}

func (fc *fakeClock) Sleep(d time.Duration) {
	fc.mu.Lock()
	for {
		until := fc.now.Add(d)
		if fc.now.Compare(until) >= 0 {
			fc.mu.Unlock()
			return
		}
		ch := make(chan struct{})
		fc.sleep = append(fc.sleep, sleepReq{until: until, ch: ch})
		fc.cond.Wait()
		found := false
		for i, sr := range fc.sleep {
			if sr.ch == ch {
				fc.sleep = append(fc.sleep[:i], fc.sleep[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			continue
		}
		break
	}
	fc.mu.Unlock()
}
