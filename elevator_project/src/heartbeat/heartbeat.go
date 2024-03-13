package heartbeat

import (
	"driver/config"
	"driver/elevator"
	"sync"
	"time"
)

type HeartBeat struct {
	SenderID      string
	HallRequests  [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType
	ElevatorState elevator.ElevatorState
}

func Send(
	id string,
	ch_hallRequestsOut <-chan [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType,
	ch_elevatorStateToNetwork <-chan elevator.ElevatorState,
	ch_msgOut chan<- HeartBeat,
	ch_heartbeatLifeLine chan<- int,
) {

	var mtxLock sync.Mutex
	var hallRequests [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType = <-ch_hallRequestsOut
	var elevatorState elevator.ElevatorState = <-ch_elevatorStateToNetwork

	go func() {
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
		ch_heartbeatLifeLine <- 1
		mtxLock.Lock()
		newHeartbeat := HeartBeat{
			SenderID:      id,
			HallRequests:  hallRequests,
			ElevatorState: elevatorState,
		}
		mtxLock.Unlock()
		ch_msgOut <- newHeartbeat
		time.Sleep(time.Duration(config.HeartbeatSleepMillisec) * time.Millisecond)
		//fmt.Printf("\n Heartbeat sent:\n")
		//fmt.Println(newHeartbeat)
	}
}
