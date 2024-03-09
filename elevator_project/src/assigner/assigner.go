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
	var allOrders [config.N_FLOORS][config.N_BUTTONS]int
	for i := range allOrders {
		for j := range allOrders[i] {
			allOrders[i][j] = 0
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
			if allOrders[buttonPressed.BtnFloor][buttonPressed.BtnType] != 2 {
				allOrders[buttonPressed.BtnFloor][buttonPressed.BtnType] = 2
			} else if allOrders[buttonPressed.BtnFloor][buttonPressed.BtnType] != 2 {
				allOrders[buttonPressed.BtnFloor][buttonPressed.BtnType] = 1
			}
		case completedOrder := <-ch_completedOrders: //THIS NEEDS TO BE REVISED
			if allOrders[completedOrder.BtnFloor][completedOrder.BtnType] == 3 {
				allOrders[completedOrder.BtnFloor][completedOrder.BtnType] = 0
			} else if allOrders[completedOrder.BtnFloor][completedOrder.BtnType] == 2 {
				allOrders[completedOrder.BtnFloor][completedOrder.BtnType] = 3
			}
		case elevatorState := <-ch_elevatorStateToAssigner:
			localElevatorState = elevatorState
		case updateHallRequest := <-ch_hallRequestsIn:
			for i := range updateHallRequest {
				for j := 0; j < 2; j++ {
					if allOrders[i][j] == 0 {
						if updateHallRequest[i][j] == 0 {
							//NOP
						} else if updateHallRequest[i][j] == 1 {
							allOrders[i][j] = 2
						} else if updateHallRequest[i][j] == 2 {
							allOrders[i][j] = 2
						} else if updateHallRequest[i][j] == 3 {
							//NOP
						}
					} else if allOrders[i][j] == 1 {
						if updateHallRequest[i][j] == 0 {
							//NOP
						} else if updateHallRequest[i][j] == 1 {
							allOrders[i][j] = 2
						} else if updateHallRequest[i][j] == 2 {
							allOrders[i][j] = 2
						} else if updateHallRequest[i][j] == 3 {
							//NOP
						}
					} else if allOrders[i][j] == 2 {
						if updateHallRequest[i][j] == 0 {
							//NOP
						} else if updateHallRequest[i][j] == 1 {
							//NOP
						} else if updateHallRequest[i][j] == 2 {
							//NOP
						} else if updateHallRequest[i][j] == 3 {
							allOrders[i][j] = 3
						}
					} else if allOrders[i][j] == 3 {
						if updateHallRequest[i][j] == 0 {
							allOrders[i][j] = 0
						} else if updateHallRequest[i][j] == 1 {
							allOrders[i][j] = 2
						} else if updateHallRequest[i][j] == 2 {
							allOrders[i][j] = 2
						} else if updateHallRequest[i][j] == 3 {
							allOrders[i][j] = 0
						}
					}
				}
			}
		default:
			//NOP
		}

		elevator.SetAllButtonLights(allOrders)
		var hall_requestsOut [config.N_FLOORS][config.N_BUTTONS - 1]int
		var hall_requests [config.N_FLOORS][config.N_BUTTONS - 1]bool
		for i := range hall_requests {
			for j := 0; j < 2; j++ {
				hall_requestsOut[i][j] = allOrders[i][j]
				if allOrders[i][j] == 2 {
					hall_requests[i][j] = true
				} else {
					hall_requests[i][j] = false
				}
			}
		}
		ch_hallRequestsOut <- hall_requestsOut

		assignedHallRequests := cost.Cost(hall_requests, localElevatorState, externalElevators)
		var localOrders [config.N_FLOORS][config.N_BUTTONS]bool
		for i := range assignedHallRequests {
			for j := 0; j < 2; j++ {
				localOrders[i][j] = assignedHallRequests[i][j]
			}
		}
		for i := range localOrders {
			if allOrders[i][2] == 2 {
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
