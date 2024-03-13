package deadlockdetector

import (
	"fmt"
	"time"
)

const (
	FSMLifelineIndex       int = 0
	assignerLifeLineIndex  int = 1
	heartbeatLifeLineIndex int = 2
	peersLifeLineIndex     int = 3
)

var (
	lifeLines [4]time.Time
)

func DeadlockDetector(
	ch_FSMLifeline <-chan int,
	ch_assignerLifeLine <-chan int,
	ch_heartbeatLifeLine <-chan int,
	ch_peersLifeLine <-chan int,
) {
	for i := range lifeLines {
		lifeLines[i] = time.Now()
	}
	for {
		select {
		case <-ch_FSMLifeline:
			lifeLines[FSMLifelineIndex] = time.Now()
		case <-ch_assignerLifeLine:
			lifeLines[assignerLifeLineIndex] = time.Now()
		case <-ch_heartbeatLifeLine:
			lifeLines[heartbeatLifeLineIndex] = time.Now()
		case <-ch_peersLifeLine:
			lifeLines[peersLifeLineIndex] = time.Now()
		default:
			// NOP
		}
		for _, lifeLine := range lifeLines {
			if lifeLine.Add(time.Duration(10) * time.Second).Before(time.Now()) {
				fmt.Println("Deadlock!!!")
				panic("DEADLOCK DETECTED")
			}
		}
	}
}
