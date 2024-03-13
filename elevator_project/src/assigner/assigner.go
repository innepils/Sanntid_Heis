package assigner

import (
	"driver/backup"
	"driver/config"
	"driver/cost"
	"driver/elevator"
	"driver/elevator_io"
	"encoding/json"
	"fmt"
	"time"
)

func RequestAssigner(
	id string,
	ch_buttonPressed <-chan elevator_io.ButtonEvent,
	ch_completedRequests <-chan elevator_io.ButtonEvent,
	ch_elevatorStateToAssigner <-chan map[string]elevator.ElevatorState,
	ch_hallRequestsIn <-chan [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType,
	ch_externalElevators <-chan []byte,
	ch_hallRequestsOut chan<- [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType,
	ch_localRequests chan<- [config.N_FLOORS][config.N_BUTTONS]bool,
){

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

	for i := range allRequests {
		for j := range allRequests[i] {
			allRequests[i][j] = elevator.NoRequest
			prevAllRequests[i][j] = 4
			prevLocalRequests[i][j] = false
		}
	}
	idleTimeOut = time.NewTimer(time.Duration(10) * time.Second)

	for {
		//fmt.Printf("Entered assigner loop")
		select {
		case buttonPressed := <-ch_buttonPressed:
			// Gets button presses and registers the requests
			if buttonPressed.BtnType == elevator_io.BT_Cab {
				allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] = elevator.ConfirmedRequest
				backup.SaveBackupToFile("backup.txt", allRequests)
			} else if allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] != elevator.ConfirmedRequest {
				//fmt.Println("req set to 1")
				allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] = elevator.NewRequest
			}
		case completedRequest := <-ch_completedRequests:
			// Completed requests from FSM
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
			//fmt.Printf("\nRecieved hallrequest in: ")
			//fmt.Println(updateHallRequest)
			for i := range updateHallRequest {
				for j := 0; j < 2; j++ {
					switch allRequests[i][j] {
					case elevator.NoRequest:
						switch updateHallRequest[i][j] {
						case elevator.NewRequest:
							allRequests[i][j] = elevator.NewRequest
						case elevator.ConfirmedRequest:
							allRequests[i][j] = elevator.ConfirmedRequest
						case elevator.NoRequest, elevator.CompletedRequest:
							// NOP
						}
					case elevator.NewRequest:
						switch updateHallRequest[i][j] {
						case elevator.NewRequest, elevator.ConfirmedRequest:
							allRequests[i][j] = elevator.ConfirmedRequest
						case elevator.NoRequest, elevator.CompletedRequest:
							// NOP
						}
					case elevator.ConfirmedRequest:
						switch updateHallRequest[i][j] {
						case elevator.CompletedRequest:
							allRequests[i][j] = elevator.CompletedRequest
						case elevator.NoRequest, elevator.NewRequest, elevator.ConfirmedRequest:
							// NOP
						}
					case elevator.CompletedRequest:
						switch updateHallRequest[i][j] {
						case elevator.NoRequest, elevator.CompletedRequest:
							allRequests[i][j] = elevator.NoRequest
						case elevator.NewRequest:
							allRequests[i][j] = elevator.ConfirmedRequest
						case elevator.ConfirmedRequest:
							//NOP
						}
					}
				}
			}

		default:
			//NOP
		}

		for i := 0; i < config.N_FLOORS; i++ {
			for j := 0; j < config.N_BUTTONS-1; j++ {
				hallRequestsOut[i][j] = allRequests[i][j]
				if allRequests[i][j] == elevator.ConfirmedRequest {
					hallRequests[i][j] = true
				} else {
					hallRequests[i][j] = false
				}
			}
		}
		ch_hallRequestsOut <- hallRequestsOut
		assignedHallRequests := cost.Cost(id, hallRequests, localElevatorState, externalElevators)
		for i := 0; i < config.N_FLOORS; i++ {
			copy(localRequests[i][:2], assignedHallRequests[i][:])

			if allRequests[i][2] == elevator.ConfirmedRequest {
				localRequests[i][2] = true
			} else {
				localRequests[i][2] = false
			}
		}

		// checks if changes were made, and if so,
		if localRequests != prevLocalRequests {
			ch_localRequests <- localRequests
			prevLocalRequests = localRequests
		}
		if allRequests != prevAllRequests {
			elevator.SetAllButtonLights(allRequests)
			prevAllRequests = allRequests
		}

		if localElevatorState[id].Behavior != "idle" {
			idleTimeOut.Reset(time.Duration(config.IdleTimeOutDurationSec) * time.Second)
		} else {
			requestFlag := true
			for i := 0; i < config.N_FLOORS; i++ {
				for j := 0; j < config.N_BUTTONS; j++ {
					if allRequests[i][j] == elevator.ConfirmedRequest {
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
			fmt.Printf("TIMED OUT!!!!\n")
			for i := 0; i < config.N_FLOORS; i++ {
				for j := 0; j < config.N_BUTTONS; j++ {
					if allRequests[i][j] == elevator.ConfirmedRequest {
						localRequests[i][j] = true
					} else {
						localRequests[i][j] = false
					}
				}
			}
			ch_localRequests <- localRequests
		default:
			//NOP
		}
	}
}
