package deadlockdetector

import (
	"fmt"
	"time"
)

const (
	FSMDeadlockIndex       int = 0
	assignerDeadlockIndex  int = 1
	heartbeatDeadlockIndex int = 2
	peersDeadlockIndex     int = 3
)

var (
	deadlocks [4]time.Time
)

func DeadlockDetector(
	ch_FSMDeadlock 			<-chan int,
	ch_assignerDeadlock 	<-chan int,
	ch_heartbeatDeadlock 	<-chan int,
	ch_peersDeadlock 		<-chan int,
) {
	for i := range deadlocks {
		deadlocks[i] = time.Now()
	}
	for {
		select {
		case <-ch_FSMDeadlock:
			deadlocks[FSMDeadlockIndex] = time.Now()
		case <-ch_assignerDeadlock:
			deadlocks[assignerDeadlockIndex] = time.Now()
		case <-ch_heartbeatDeadlock:
			deadlocks[heartbeatDeadlockIndex] = time.Now()
		case <-ch_peersDeadlock:
			deadlocks[peersDeadlockIndex] = time.Now()
		default:
			// NOP
		}
		for locked, deadlockTime := range deadlocks {
			if deadlockTime.Add(time.Duration(10) * time.Second).Before(time.Now()) {
				panic(fmt.Sprintf("DEADLOCK DETECTED IN PROCESS %d", locked))
			}
		}
	}
}
