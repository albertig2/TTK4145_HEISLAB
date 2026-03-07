package timer

import "time"

var (
	timerEndTime time.Time
	timerActive  bool
)

func Timer_start(d time.Duration) {
	timerEndTime = time.Now().Add(d)
	timerActive = true
}

func Timer_stop() {
	timerActive = false
}

func Timer_timedOut() bool {
	return timerActive && time.Now().After(timerEndTime)
}
