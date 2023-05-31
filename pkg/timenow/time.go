package timenow

import (
	"time"
)

type TimeNow interface {
	Now() time.Time
}

type realTime struct{}

func NewRealTime() TimeNow {
	return &realTime{}
}

func (t *realTime) Now() time.Time {
	return time.Now().UTC()
}

type fixedTime struct {
	time time.Time
}

func NewFixedTime(time time.Time) TimeNow {
	return &fixedTime{time: time}
}

func (t *fixedTime) Now() time.Time {
	return t.time
}
