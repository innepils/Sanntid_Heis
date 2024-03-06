package assigner

import (
	"driver/config"
	"driver/elevator_io"
)

type requestType int

const (
	none      requestType = 0
	new                   = 1
	confirmed             = 2
	completed             = 3
)

func Assigner(ch_buttonPressed chan elevator_io.ButtonEvent, ch_completedOrders chan elevator_io.ButtonEvent, ch_localOrders chan elevator_io.ButtonEvent) {
	var allOrders [config.N_FLOORS][config.N_BUTTONS]int
	for i := range allOrders {
		for j := range allOrders[i] {
			allOrders[i][j] = 0
		}
	}

	for {
		select {
		case buttonPressed := <-ch_buttonPressed:
			if allOrders[buttonPressed.BtnFloor][buttonPressed.BtnType] != 2 {
				allOrders[buttonPressed.BtnFloor][buttonPressed.BtnType] = 1
			}
		case completedOrder := <-ch_completedOrders: //THIS NEEDS TO BE REVISED
			if allOrders[completedOrder.BtnFloor][completedOrder.BtnType] == 3 {
				allOrders[completedOrder.BtnFloor][completedOrder.BtnType] = 0
			} else if allOrders[completedOrder.BtnFloor][completedOrder.BtnType] == 2 {
				allOrders[completedOrder.BtnFloor][completedOrder.BtnType] = 3
			}
		}

		//COST
		assignedHallRequests := cost()
		var localOrders [config.N_FLOORS][config.N_BUTTONS]bool
		for i := range assignedHallRequests {
			for j := 0; j < 2; j++ {
				localOrders[i][j] = assignedHallRequests[i][j]
			}
		}
		for i := range localOrders {
			if allOrders[i][2] != 0 {
				localOrders[i][2] = true
			} else {
				localOrders[i][2] = false
			}
		}
		
		for i := range localOrders {
			for j := range localOrders[i] {
				if localOrders[i][j] {
					ch_localOrders <- elevator_io.ButtonEvent{BtnFloor: i, BtnType: elevator_io.ButtonType(j)}
				}
			}
		}
	}
}
