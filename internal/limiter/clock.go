package limiter

import "time"

// Clock abstracts time for testability.
type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

// realClock uses the real system clock.
type realClock struct{}

func (realClock) Now() time.Time                        { return time.Now() }
func (realClock) Sleep(d time.Duration)                 { time.Sleep(d) }
