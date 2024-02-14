package main

import (
    "time"
)

var (
    timerEndTime time.Time
    timerActive  bool
)

// SK: getWallTime returns the current time
func getWallTime() time.Time {
    return time.Now()
}

// SK: timerStart starts the timer with the specified duration in seconds
func timerStart(duration float64) {
    timerEndTime = getWallTime().Add(time.Duration(duration * float64(time.Second)))
    timerActive = true
}

// SK: timerStop stops the timer
func timerStop() {
    timerActive = false
}

// SK: timerTimedOut checks if the timer has timed out
func timerTimedOut() bool {
    return timerActive && getWallTime().After(timerEndTime)
}