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

	for i := range allRequests {
		for j := range allRequests[i] {
			allRequests[i][j] = 0
			prevAllRequests[i][j] = 4
			prevLocalRequests[i][j] = false
		}
	}
	idleTimeOut = time.NewTimer(time.Duration(10) * time.Second)
	//var localElevatorState = map[string]elevator.ElevatorState{id: {Behavior: "idle", Floor: 1, Direction: "stop", CabRequests: []bool{true, false, true, false}}}
	//var prevLocalRequests [config.N_FLOORS][config.N_BUTTONS]bool
	// for i := range prevLocalRequests {
	// 	for j := range prevLocalRequests[i] {
	// 		prevLocalRequests[i][j] = false
	// 	}
	// }

	for {
		//fmt.Printf("Entered assigner loop")
		select {
		case buttonPressed := <-ch_buttonPressed:
			// Gets button presses and registers the requests
			if buttonPressed.BtnType == elevator_io.BT_Cab {
				allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] = 2
				backup.SaveBackupToFile("backup.txt", allRequests)
			} else if allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] != 2 {
				//fmt.Println("req set to 1")
				allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] = 1
			}
		case completedRequest := <-ch_completedRequests:
			// Completed requests from FSM
			if allRequests[completedRequest.BtnFloor][completedRequest.BtnType] == 3 {
				allRequests[completedRequest.BtnFloor][completedRequest.BtnType] = 0
			} else if allRequests[completedRequest.BtnFloor][completedRequest.BtnType] == 2 {
				allRequests[completedRequest.BtnFloor][completedRequest.BtnType] = 3
			}
			backup.SaveBackupToFile("backup.txt", allRequests)

		case elevatorState := <-ch_elevatorStateToAssigner:
			localElevatorState = elevatorState

		case currentExternalElevators := <-ch_externalElevators:
			externalElevators = currentExternalElevators

			// case updateHallRequest := <-ch_hallRequestsIn:
			// for i := range updateHallRequest {
			// 	for j := 0; j < 2; j++ {
			// 		if allRequests[i][j] == 0 {
			// 			if updateHallRequest[i][j] == 0 {
			// 				//NOP
			// 			} else if updateHallRequest[i][j] == 1 {
			// 				allRequests[i][j] = 2
			// 			} else if updateHallRequest[i][j] == 2 {
			// 				allRequests[i][j] = 2
			// 			} else if updateHallRequest[i][j] == 3 {
			// 				//NOP
			// 			}
			// 		} else if allRequests[i][j] == 1 {
			// 			if updateHallRequest[i][j] == 0 {
			// 				//NOP
			// 			} else if updateHallRequest[i][j] == 1 {
			// 				allRequests[i][j] = 2
			// 			} else if updateHallRequest[i][j] == 2 {
			// 				allRequests[i][j] = 2
			// 			} else if updateHallRequest[i][j] == 3 {
			// 				//NOP
			// 			}
			// 		} else if allRequests[i][j] == 2 {
			// 			if updateHallRequest[i][j] == 0 {
			// 				//NOP
			// 			} else if updateHallRequest[i][j] == 1 {
			// 				//NOP
			// 			} else if updateHallRequest[i][j] == 2 {
			// 				//NOP
			// 			} else if updateHallRequest[i][j] == 3 {
			// 				allRequests[i][j] = 3
			// 			}
			// 		} else if allRequests[i][j] == 3 {
			// 			if updateHallRequest[i][j] == 0 {
			// 				allRequests[i][j] = 0
			// 			} else if updateHallRequest[i][j] == 1 {
			// 				allRequests[i][j] = 2
			// 			} else if updateHallRequest[i][j] == 2 {
			// 				//NOP
			// 			} else if updateHallRequest[i][j] == 3 {
			// 				allRequests[i][j] = 0
			// 			}
			// 		}
			// 	}
			// }

		case updateHallRequest := <-ch_hallRequestsIn:
			//fmt.Printf("\nRecieved hallrequest in: ")
			//fmt.Println(updateHallRequest)
			for i := range updateHallRequest {
				for j := 0; j < 2; j++ {
					switch allRequests[i][j] {
					case 0:
						switch updateHallRequest[i][j] {
						case 1:
							allRequests[i][j] = 1
						case 2:
							allRequests[i][j] = 2
						case 0, 3:
							// NOP
						}
					case 1:
						switch updateHallRequest[i][j] {
						case 1, 2:
							allRequests[i][j] = 2
						case 0, 3:
							// NOP
						}
					case 2:
						switch updateHallRequest[i][j] {
						case 3:
							allRequests[i][j] = 3
						case 0, 1, 2:
							// NOP
						}
					case 3:
						switch updateHallRequest[i][j] {
						case 0, 3:
							allRequests[i][j] = 0
						case 1:
							allRequests[i][j] = 2
						case 2:
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
				if allRequests[i][j] == 2 {
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

			if allRequests[i][2] == 2 {
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

		if localElevatorState[id].Behavior != "Idle" {
			idleTimeOut.Reset(10 * time.Second)
		} else {
			requestFlag := true
			for i := 0; i < config.N_FLOORS; i++ {
				for j := 0; j < config.N_BUTTONS; j++ {
					if allRequests[i][j] == 2 {
						requestFlag = false
					}
				}
			}
			if requestFlag {
				idleTimeOut.Reset(10 * time.Second)
			}
		}
		select {
		case <-idleTimeOut.C:
			fmt.Printf("TIMED OUT!!!!\n")
			for i := 0; i < config.N_FLOORS; i++ {
				for j := 0; j < config.N_BUTTONS; j++ {
					if allRequests[i][j] == 2 {
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
