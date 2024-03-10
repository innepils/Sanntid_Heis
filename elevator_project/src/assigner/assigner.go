package assigner

import (
	"driver/config"
	"driver/cost"
	"driver/elevator"
	"driver/elevator_io"
)

type requestType int

const (
	none      requestType = 0
	new                   = 1
	confirmed             = 2
	completed             = 3
)

func Assigner(
	ch_buttonPressed chan elevator_io.ButtonEvent,
	ch_completedOrders chan elevator_io.ButtonEvent,
	ch_localOrders chan [config.N_FLOORS][config.N_BUTTONS]bool,
	ch_hallRequestsIn chan [config.N_FLOORS][config.N_BUTTONS - 1]int,
	ch_hallRequestsOut chan [config.N_FLOORS][config.N_BUTTONS - 1]int,
	ch_elevatorStateToAssigner chan map[string]elevator.ElevatorState,
	ch_externalElevators chan map[string]elevator.ElevatorState,
) {
	externalElevators := map[string]elevator.ElevatorState{}
	var allRequests [config.N_FLOORS][config.N_BUTTONS]int
	for i := range allRequests {
		for j := range allRequests[i] {
			allRequests[i][j] = 0
		}
	}
	var localElevatorState = map[string]elevator.ElevatorState{"self": {Behavior: "idle", Floor: 1, Direction: "stop", CabRequests: []bool{true, false, true, false}}}
	var prevLocalRequests [config.N_FLOORS][config.N_BUTTONS]bool
	for i := range prevLocalRequests {
		for j := range prevLocalRequests[i] {
			prevLocalRequests[i][j] = false
		}
	}
	for {
		select {
		case buttonPressed := <-ch_buttonPressed:
			if allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] != 2 {
				allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] = 2
			} else if allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] != 2 {
				allRequests[buttonPressed.BtnFloor][buttonPressed.BtnType] = 1
			}
		case completedOrder := <-ch_completedOrders: //THIS NEEDS TO BE REVISED
			if allRequests[completedOrder.BtnFloor][completedOrder.BtnType] == 3 {
				allRequests[completedOrder.BtnFloor][completedOrder.BtnType] = 0
			} else if allRequests[completedOrder.BtnFloor][completedOrder.BtnType] == 2 {
				allRequests[completedOrder.BtnFloor][completedOrder.BtnType] = 3
			}
		case elevatorState := <-ch_elevatorStateToAssigner:
			localElevatorState = elevatorState
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
							allRequests[i][j] = 2
						} else if updateHallRequest[i][j] == 3 {
							allRequests[i][j] = 0
						}
					}
				}
			}
		*/	
		case updateHallRequest := <-ch_hallRequestsIn:
			for i := range updateHallRequest {
				for j := 0; j < 2; j++ {
					switch allRequests[i][j] {
					case 0:
						switch updateHallRequest[i][j] {
						case 1, 2:
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
						case 1, 2:
							allRequests[i][j] = 2
						}
					}
				}
			}
		
		default:
			//NOP
		}

		elevator.SetAllButtonLights(allRequests)
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

		assignedHallRequests := cost.Cost(hallRequests, localElevatorState, externalElevators)
		var localOrders [config.N_FLOORS][config.N_BUTTONS]bool
		for i := range assignedHallRequests {
			for j := 0; j < 2; j++ {
				localOrders[i][j] = assignedHallRequests[i][j]
			}
		}
		for i := range localOrders {
			if allRequests[i][2] == 2 {
				localOrders[i][2] = true
			} else {
				localOrders[i][2] = false
			}
		}

		if localOrders != prevLocalRequests {
			ch_localOrders <- localOrders
			prevLocalRequests = localOrders
		}
	}
}
