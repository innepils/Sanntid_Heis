package heartbeat

import (
	"driver/config"
	"driver/elevator"
	"time"
)

type HeartBeat struct {
	SenderID      string
	HallRequests  [config.N_FLOORS][config.N_BUTTONS - 1]int
	ElevatorState elevator.ElevatorState
}

// Send heartbeat to network incl. all info

// func createHeartbeat(
// 	id string,
// 	ch_hallRequestsOut chan [config.N_FLOORS][config.N_BUTTONS - 1]int,
// 	ch_elevatorStateToNetwork chan elevator.ElevatorState) HeartBeat {
// 	for {

// 	}
// }

func Send(
	id string,
	ch_hallRequestsOut chan [config.N_FLOORS][config.N_BUTTONS - 1]int,
	ch_elevatorStateToNetwork chan elevator.ElevatorState,
	ch_msgOut chan HeartBeat) {

	var HallRequests [config.N_FLOORS][config.N_BUTTONS - 1]int = <-ch_hallRequestsOut
	var elevatorState elevator.ElevatorState = <-ch_elevatorStateToNetwork
	go func() {
		for {
			select {
			case newHallRequests := <-ch_hallRequestsOut:
				HallRequests = newHallRequests
			case newElevatorState := <-ch_elevatorStateToNetwork:
				elevatorState = newElevatorState
			default:
				// NOP
			}
		}
	}()
	for {
		newHeartbeat := HeartBeat{
			SenderID:      id,
			HallRequests:  HallRequests,
			ElevatorState: elevatorState,
		}
		ch_msgOut <- newHeartbeat
		time.Sleep(100 * time.Millisecond)
		//fmt.Printf("\n Heartbeat sent:\n")
		//fmt.Println(newHeartbeat)
	}
}
