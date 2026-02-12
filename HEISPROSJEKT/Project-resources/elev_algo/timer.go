import "time"

var (
	timerEndTime time.Time
	timerActive  bool
)

func timerStart(d time.Duration) {
	timerEndTime = time.Now().Add(d)
	timerActive = true
}

func timerStop() {
	timerActive = false
}

func timerTimedOut() bool {
	return timerActive && time.Now().After(timerEndTime)
}
