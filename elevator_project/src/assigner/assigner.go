package assigner

import (
	"driver/backup"
	"driver/config"
	"driver/cost"
	"driver/elevator"
	"driver/elevator_io"
	"encoding/json"
	"time"
)

func RequestAssigner(
	id 							string,
	ch_buttonPressed 			<-chan elevator_io.ButtonEvent,
	ch_completedRequests		<-chan elevator_io.ButtonEvent,
	ch_elevatorStateToAssigner 	<-chan map[string]elevator.ElevatorState,
	ch_hallRequestsIn 			<-chan [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType,
	ch_externalElevators 		<-chan []byte,
	ch_hallRequestsOut 			chan<- [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType,
	ch_localRequests 			chan<- [config.N_FLOORS][config.N_BUTTONS]bool,
	ch_assignerDeadlock 		chan<- int,
) {

	var (
		idleTimeOut        *time.Timer
		allRequests        [config.N_FLOORS][config.N_BUTTONS]elevator.RequestType
		prevAllRequests    [config.N_FLOORS][config.N_BUTTONS]elevator.RequestType
		prevLocalRequests  [config.N_FLOORS][config.N_BUTTONS]bool
		emptyElevatorMap   map[string]elevator.ElevatorState
		hallRequestsOut    [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType
		hallRequests       [config.N_FLOORS][config.N_BUTTONS - 1]bool
		localRequests      [config.N_FLOORS][config.N_BUTTONS]bool
		localElevatorState = map[string]elevator.ElevatorState{id: {Behavior: "idle", Floor: 1, Direction: "stop", CabRequests: []bool{false, false, false, false}}}
	)
	// emptyElevatorMap = map[string]elevator.ElevatorState{}
	externalElevators, _ := json.Marshal(emptyElevatorMap)

	for floor := range allRequests {
		for btn := range allRequests[floor] {
			allRequests[floor][btn] = elevator.NoRequest
			prevAllRequests[floor][btn] = elevator.UndefinedRequest
			prevLocalRequests[floor][btn] = false
		}
	}
	idleTimeOut = time.NewTimer(time.Duration(10) * time.Second)

	for {
		ch_assignerDeadlock <- 1
		select {
		case buttonPressed := <-ch_buttonPressed:
			if buttonPressed.BtnType == elevator_io.BT_Cab {
				allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] = elevator.ConfirmedRequest
				backup.SaveBackupToFile("backup.txt", allRequests)
			} else if allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] != elevator.ConfirmedRequest {
				allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] = elevator.NewRequest
			}
		case completedRequest := <-ch_completedRequests:
			if allRequests[completedRequest.BtnFloor][completedRequest.BtnType] == elevator.CompletedRequest {
				allRequests[completedRequest.BtnFloor][completedRequest.BtnType] = elevator.NoRequest
			} else if allRequests[completedRequest.BtnFloor][completedRequest.BtnType] == elevator.ConfirmedRequest {
				allRequests[completedRequest.BtnFloor][completedRequest.BtnType] = elevator.CompletedRequest
			}
			backup.SaveBackupToFile("backup.txt", allRequests)

		case elevatorState := <-ch_elevatorStateToAssigner:
			localElevatorState = elevatorState

		case currentExternalElevators := <-ch_externalElevators:
			externalElevators = currentExternalElevators

		case updateHallRequest := <-ch_hallRequestsIn:
			for floor := range updateHallRequest {
				for btn := 0; btn < config.N_BUTTONS-1; btn++ {
					switch allRequests[floor][btn] {
					case elevator.NoRequest:
						switch updateHallRequest[floor][btn] {
						case elevator.NewRequest:
							allRequests[floor][btn] = elevator.NewRequest
						case elevator.ConfirmedRequest:
							allRequests[floor][btn] = elevator.ConfirmedRequest
						case elevator.NoRequest, elevator.CompletedRequest:
							// NOP
						}
					case elevator.NewRequest:
						switch updateHallRequest[floor][btn] {
						case elevator.NewRequest, elevator.ConfirmedRequest:
							allRequests[floor][btn] = elevator.ConfirmedRequest
						case elevator.NoRequest, elevator.CompletedRequest:
							// NOP
						}
					case elevator.ConfirmedRequest:
						switch updateHallRequest[floor][btn] {
						case elevator.CompletedRequest:
							allRequests[floor][btn] = elevator.CompletedRequest
						case elevator.NoRequest, elevator.NewRequest, elevator.ConfirmedRequest:
							// NOP
						}
					case elevator.CompletedRequest:
						switch updateHallRequest[floor][btn] {
						case elevator.NoRequest, elevator.CompletedRequest:
							allRequests[floor][btn] = elevator.NoRequest
						case elevator.NewRequest:
							allRequests[floor][btn] = elevator.ConfirmedRequest
						case elevator.ConfirmedRequest:
							//NOP
						}
					}
				}
			}

		default:
			//NOP
		}

		for floor := 0; floor < config.N_FLOORS; floor++ {
			for btn := 0; btn < config.N_BUTTONS-1; btn++ {
				hallRequestsOut[floor][btn] = allRequests[floor][btn]
				if allRequests[floor][btn] == elevator.ConfirmedRequest {
					hallRequests[floor][btn] = true
				} else {
					hallRequests[floor][btn] = false
				}
			}
		}
		ch_hallRequestsOut <- hallRequestsOut
		assignedHallRequests := cost.Cost(id, hallRequests, localElevatorState, externalElevators)
		for floor := 0; floor < config.N_FLOORS; floor++ {
			copy(localRequests[floor][:2], assignedHallRequests[floor][:])

			if allRequests[floor][elevator_io.BT_Cab] == elevator.ConfirmedRequest {
				localRequests[floor][elevator_io.BT_Cab] = true
			} else {
				localRequests[floor][elevator_io.BT_Cab] = false
			}
		}

		if localRequests != prevLocalRequests {
			ch_localRequests <- localRequests
			prevLocalRequests = localRequests
		}
		if allRequests != prevAllRequests {
			elevator.SetAllButtonLights(allRequests)
			prevAllRequests = allRequests
		}

		//if an elevator is idle in more than 'IdleTimeOuutDurationSec' while there are orders, it will take them.
		if localElevatorState[id].Behavior != "idle" {
			idleTimeOut.Reset(time.Duration(config.IdleTimeOutDurationSec) * time.Second)
		} else {
			requestFlag := true
			for floor := 0; floor < config.N_FLOORS; floor++ {
				for btn := 0; btn < config.N_BUTTONS; btn++ {
					if allRequests[floor][btn] == elevator.ConfirmedRequest {
						requestFlag = false
					}
				}
			}
			if requestFlag {
				idleTimeOut.Reset(time.Duration(config.IdleTimeOutDurationSec) * time.Second)
			}
		}
		select {
		case <-idleTimeOut.C:
			for floor := 0; floor < config.N_FLOORS; floor++ {
				for btn := 0; btn < config.N_BUTTONS; btn++ {
					if allRequests[floor][btn] == elevator.ConfirmedRequest {
						localRequests[floor][btn] = true
					} else {
						localRequests[floor][btn] = false
					}
				}
			}
			ch_localRequests <- localRequests
		default:
			//NOP
		}
	}
}
