package assigner

import (
	"driver/backup"
	"driver/config"
	"driver/cost"
	"driver/elevator"
	"driver/elevator_io"
	"encoding/json"
)

func Assigner(
	id string,
	ch_buttonPressed chan elevator_io.ButtonEvent,
	ch_completedRequests chan elevator_io.ButtonEvent,
	ch_localRequests chan [config.N_FLOORS][config.N_BUTTONS]bool,
	ch_hallRequestsIn chan [config.N_FLOORS][config.N_BUTTONS - 1]int,
	ch_hallRequestsOut chan [config.N_FLOORS][config.N_BUTTONS - 1]int,
	ch_elevatorStateToAssigner chan map[string]elevator.ElevatorState,
	ch_externalElevators chan []byte,
) {
	emptyElevatorMap := map[string]elevator.ElevatorState{}
	externalElevators, _ := json.Marshal(emptyElevatorMap)

	var allRequests [config.N_FLOORS][config.N_BUTTONS]int
	for i := range allRequests {
		for j := range allRequests[i] {
			allRequests[i][j] = 0
		}
	}
	var localElevatorState = map[string]elevator.ElevatorState{id: {Behavior: "idle", Floor: 1, Direction: "stop", CabRequests: []bool{true, false, true, false}}}
	var prevLocalRequests [config.N_FLOORS][config.N_BUTTONS]bool
	for i := range prevLocalRequests {
		for j := range prevLocalRequests[i] {
			prevLocalRequests[i][j] = false
		}
	}
	for {
		//fmt.Printf("Entered assigner loop")
		select {
		// Looks at local buttons and updates request matrix
		case buttonPressed := <-ch_buttonPressed:
			if buttonPressed.BtnType == elevator_io.BT_Cab {
				allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] = 2
				backup.SaveBackupToFile("backup.txt", allRequests)
			} else if allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] != 2 {
				//fmt.Println("req set to 1")
				allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] = 1
			}
		// Data from channel arrives from FSM when local order is complete
		case completedRequest := <-ch_completedRequests: //THIS NEEDS TO BE REVISED
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

		/*
			case updateHallRequest := <-ch_hallRequestsIn:
				for i := range updateHallRequest {
					for j := 0; j < 2; j++ {
						if allRequests[i][j] == 0 {
							if updateHallRequest[i][j] == 0 {
								//NOP
							} else if updateHallRequest[i][j] == 1 {
								allRequests[i][j] = 2
							} else if updateHallRequest[i][j] == 2 {
								allRequests[i][j] = 2
							} else if updateHallRequest[i][j] == 3 {
								//NOP
							}
						} else if allRequests[i][j] == 1 {
							if updateHallRequest[i][j] == 0 {
								//NOP
							} else if updateHallRequest[i][j] == 1 {
								allRequests[i][j] = 2
							} else if updateHallRequest[i][j] == 2 {
								allRequests[i][j] = 2
							} else if updateHallRequest[i][j] == 3 {
								//NOP
							}
						} else if allRequests[i][j] == 2 {
							if updateHallRequest[i][j] == 0 {
								//NOP
							} else if updateHallRequest[i][j] == 1 {
								//NOP
							} else if updateHallRequest[i][j] == 2 {
								//NOP
							} else if updateHallRequest[i][j] == 3 {
								allRequests[i][j] = 3
							}
						} else if allRequests[i][j] == 3 {
							if updateHallRequest[i][j] == 0 {
								allRequests[i][j] = 0
							} else if updateHallRequest[i][j] == 1 {
								allRequests[i][j] = 2
							} else if updateHallRequest[i][j] == 2 {
								//NOP
							} else if updateHallRequest[i][j] == 3 {
								allRequests[i][j] = 0
							}
						}
					}
				}
		*/
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

		//elevator.SetAllButtonLights(allRequests)
		var hallRequestsOut [config.N_FLOORS][config.N_BUTTONS - 1]int
		var hallRequests [config.N_FLOORS][config.N_BUTTONS - 1]bool
		for i := range hallRequests {
			for j := 0; j < 2; j++ {
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
		var localRequests [config.N_FLOORS][config.N_BUTTONS]bool
		for i := range assignedHallRequests {
			for j := 0; j < 2; j++ {
				localRequests[i][j] = assignedHallRequests[i][j]
			}
		}
		for i := range localRequests {
			if allRequests[i][2] == 2 {
				localRequests[i][2] = true
			} else {
				localRequests[i][2] = false
			}
		}

		if localRequests != prevLocalRequests {
			ch_localRequests <- localRequests
			prevLocalRequests = localRequests
		}
	}
}

// backup.SaveBackupToFile("backup.txt", []bool(localElevatorState[id].CabRequests))
