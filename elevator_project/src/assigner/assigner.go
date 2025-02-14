package assigner

import (
	"encoding/json"
	"src/backup"
	"src/config"
	"src/cost"
	"src/elevator"
	"src/elevator_io"
	"time"
)

func RequestAssigner(
	id 							string,
	ch_buttonPressed 			<-chan elevator_io.ButtonEvent,
	ch_completedRequests		<-chan elevator_io.ButtonEvent,
	ch_elevatorStateToAssigner 	<-chan map[string]elevator.HRAElevatorState,
	ch_hallRequestsIn 			<-chan [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType,
	ch_externalElevators 		<-chan []byte,
	ch_hallRequestsOut 			chan<- [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType,
	ch_localRequests 			chan<- [config.N_FLOORS][config.N_BUTTONS]bool,
	ch_assignerDeadlock 		chan<- string,
) {

	var (
		idleTimeOut        *time.Timer
		allRequests        [config.N_FLOORS][config.N_BUTTONS]elevator.RequestType
		prevAllRequests    [config.N_FLOORS][config.N_BUTTONS]elevator.RequestType
		prevLocalRequests  [config.N_FLOORS][config.N_BUTTONS]bool
		emptyElevatorMap   map[string]elevator.HRAElevatorState
		hallRequestsOut    [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType
		hallRequests       [config.N_FLOORS][config.N_BUTTONS - 1]bool
		localRequests      [config.N_FLOORS][config.N_BUTTONS]bool
		localElevatorState = map[string]elevator.HRAElevatorState{id: {Behavior: "idle", Floor: 1, Direction: "stop", CabRequests: []bool{false, false, false, false}}}
	)
	externalElevators, _ := json.Marshal(emptyElevatorMap)

	for floor := range allRequests {
		for btn := range allRequests[floor] {
			allRequests[floor][btn] = elevator.NoRequest
			prevAllRequests[floor][btn] = elevator.UndefinedRequest
			prevLocalRequests[floor][btn] = false
		}
	}
	idleTimeOut = time.NewTimer(time.Duration(config.IdleTimeOutDurationSec) * time.Second)

	for {
		ch_assignerDeadlock <- "requestAssigner Alive"
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
			// Hall requests recieved from the network
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
					} // switch
				} // for btn
			} // for floor
		default:
			//NOP
		} // select

		// Preparing hall requests for network and cost function
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

		// Assigning requests for local elevator
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

		// If an elevator is idle for more than 'IdleTimeOuutDurationSec' while there are orders, it will take them.
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
			// Assigns all requests to local elevator
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
	} // for
} // RequestAssigner
