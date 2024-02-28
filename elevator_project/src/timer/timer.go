package timer

//func DoorTimer(ch_doorTimedOut chan<- bool){

//}

//AT: usikker på om denne er aktuell etter oversetting til go.

// Oversatt kode:
/*
var (
	timerEndTime time.Time
	timerActive  bool
)

func getWallTime() time.Time {
	return time.Now()
}

func TimerStart(duration float64) {
	timerEndTime = getWallTime().Add(time.Duration(duration * float64(time.Second)))
	timerActive = true
}

func TimerStop() {
	timerActive = false
}

func TimerTimedOut() bool {
	return timerActive && getWallTime().After(timerEndTime)
}
*/
