package repository

import "time"

// Clock interface
// Need to mock time in tests
type Clock interface {
	Now() time.Time
	AddTime(seconds int)
}

type RealClock struct{}

func (c RealClock) Now() time.Time {
	return time.Now().UTC()
}

func (c RealClock) AddTime(seconds int) {}
