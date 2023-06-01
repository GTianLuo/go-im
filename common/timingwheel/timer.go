package timingwheel

import (
	"time"
)

var tw *TimingWheel

func InitTimingWheel() {
	tw = NewTimingWheel(time.Millisecond, 20)
	tw.Start()
}

func StopTimingWheel() {
	tw.Stop()
}

func AfterFunc(d time.Duration, f func()) *Timer {
	return tw.AfterFunc(d, f)
}
