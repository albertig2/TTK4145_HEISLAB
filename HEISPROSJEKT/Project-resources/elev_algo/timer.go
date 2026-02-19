package main

import "time"

var (
	timerEndTime time.Time
	timerActive  bool
)

func timer_start(d time.Duration) {
	timerEndTime = time.Now().Add(d)
	timerActive = true
}

func timer_stop() {
	timerActive = false
}

func timer_timedOut() bool {
	return timerActive && time.Now().After(timerEndTime)
}
