package deadlock

import (
	"src/config"
	"fmt"
	"time"
)

const (
	FSMDeadlockIndex       int = 0
	assignerDeadlockIndex  int = 1
	heartbeatDeadlockIndex int = 2
	peersDeadlockIndex     int = 3
)


var deadlocks [4]time.Time

func Detector(
	ch_FSMDeadlock 			<-chan string,
	ch_assignerDeadlock 	<-chan string,
	ch_heartbeatDeadlock 	<-chan string,
	ch_peersDeadlock 		<-chan string,
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
			if deadlockTime.Add(time.Duration(config.DeadlockTimeOutDurationSec) * time.Second).Before(time.Now()) {
				panic(fmt.Sprintf("DEADLOCK DETECTED IN PROCESS %d", locked))
			}
		}
	}
}
