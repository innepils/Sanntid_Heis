package heartbeat

import (
	"src/config"
	"src/elevator"
	"sync"
	"time"
)

type HeartBeat struct {
	SenderID     	string
	HallRequests  	[config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType
	ElevatorState 	elevator.HRAElevatorState
}

func Send(
	nodeID 						string,
	ch_hallRequestsOut 		  	<-chan [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType,
	ch_elevatorStateToNetwork	<-chan elevator.HRAElevatorState,
	ch_msgOut 					chan<- HeartBeat,
	ch_heartbeatDeadlock 	  	chan<- string,
) {
	var (
		mtxLock 		sync.Mutex
		hallRequests 	[config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType = <-ch_hallRequestsOut
		elevatorState 	elevator.HRAElevatorState = <-ch_elevatorStateToNetwork

	)

	go func() { 
		// This updates in a "go func" due to heartbeat sleep
		for {
			select {
			case newHallRequests := <-ch_hallRequestsOut:
				mtxLock.Lock()
				hallRequests = newHallRequests
				mtxLock.Unlock()
			case newElevatorState := <-ch_elevatorStateToNetwork:
				mtxLock.Lock()
				elevatorState = newElevatorState
				mtxLock.Unlock()
			default:
				// NOP
			}
		}
	}()

	for {
		ch_heartbeatDeadlock <- "Heartbeat Alive"
		mtxLock.Lock()
		newHeartbeat := HeartBeat{
			SenderID:      nodeID,
			HallRequests:  hallRequests,
			ElevatorState: elevatorState,
		}
		mtxLock.Unlock()
		ch_msgOut <- newHeartbeat
		time.Sleep(time.Duration(config.HeartbeatSleepMillisec) * time.Millisecond)
	}
}
