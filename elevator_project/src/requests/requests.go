package requests

import (
	"driver/config"
	"driver/elevator"
	"driver/elevator_io"
)

func Requests_above(e *elevator.Elevator) bool {
	for f := e.Floor + 1; f < config.N_FLOORS; f++ {
		for btn := 0; btn < config.N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}
func Requests_below(e *elevator.Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < config.N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_here(e *elevator.Elevator) bool {
	for btn := 0; btn < config.N_BUTTONS; btn++ {
		if e.Requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

func Requests_chooseDirection(e *elevator.Elevator) {
	switch e.Dirn {
	case elevator_io.MD_Up:
		if Requests_above(e) {
			e.Dirn = elevator_io.MD_Up
			e.Behaviour = elevator.EB_Moving
		} else if Requests_here(e) {
			e.Dirn = elevator_io.MD_Down
			e.Behaviour = elevator.EB_DoorOpen
		} else if Requests_below(e) {
			e.Dirn = elevator_io.MD_Down
			e.Behaviour = elevator.EB_Moving
		} else {
			e.Dirn = elevator_io.MD_Stop
			e.Behaviour = elevator.EB_Idle
		}

	case elevator_io.MD_Down:
		if Requests_below(e) {
			e.Dirn = elevator_io.MD_Down
			e.Behaviour = elevator.EB_Moving
		} else if Requests_here(e) {
			e.Dirn = elevator_io.MD_Up
			e.Behaviour = elevator.EB_DoorOpen
		} else if Requests_above(e) {
			e.Dirn = elevator_io.MD_Up
			e.Behaviour = elevator.EB_Moving
		} else {
			e.Dirn = elevator_io.MD_Stop
			e.Behaviour = elevator.EB_Idle
		}

	case elevator_io.MD_Stop:

		if Requests_here(e) {
			e.Dirn = elevator_io.MD_Stop
			e.Behaviour = elevator.EB_DoorOpen
		} else if Requests_above(e) {
			e.Dirn = elevator_io.MD_Up
			e.Behaviour = elevator.EB_Moving
		} else if Requests_below(e) {
			e.Dirn = elevator_io.MD_Down
			e.Behaviour = elevator.EB_Moving
		} else {
			e.Dirn = elevator_io.MD_Stop
			e.Behaviour = elevator.EB_Idle
		}

	default:
		e.Dirn = elevator_io.MD_Stop
		e.Behaviour = elevator.EB_Idle
	}
}

func Requests_shouldStop(e *elevator.Elevator) bool {
	switch e.Dirn {
	case elevator_io.MD_Down:
		return e.Requests[e.Floor][elevator_io.BT_HallDown] || e.Requests[e.Floor][elevator_io.BT_Cab] || !Requests_below(e)
	case elevator_io.MD_Up:
		return e.Requests[e.Floor][elevator_io.BT_HallUp] || e.Requests[e.Floor][elevator_io.BT_Cab] || !Requests_above(e)
	default:
		return true
	}
}

/*
// This is the one that is OLD but works except the one edge case.

func Requests_clearAtCurrentFloor(e *elevator.Elevator, ch_completedRequests chan<- elevator_io.ButtonEvent) {

	e.Requests[e.Floor][elevator_io.BT_Cab] = false
	ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_Cab}

	switch e.Dirn {

	case elevator_io.MD_Up:
		if !Requests_above(e) && !e.Requests[e.Floor][elevator_io.BT_HallUp] {
			e.Requests[e.Floor][elevator_io.BT_HallDown] = false
			ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallDown}
		}
		e.Requests[e.Floor][elevator_io.BT_HallUp] = false
		ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallUp}

	case elevator_io.MD_Down:
		if !Requests_below(e) && !e.Requests[e.Floor][elevator_io.BT_HallDown] {
			e.Requests[e.Floor][elevator_io.BT_HallUp] = false
			ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallUp}

		}
		e.Requests[e.Floor][elevator_io.BT_HallDown] = false
		ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallDown}

	case elevator_io.MD_Stop:
		e.Requests[e.Floor][elevator_io.BT_HallUp] = false
		ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallUp}
		e.Requests[e.Floor][elevator_io.BT_HallDown] = false
		ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallDown}
	}
} */
/*
// This is a changed version of the function that should erase the edge-condition.
func Requests_clearAtCurrentFloor(e *elevator.Elevator, ch_completedRequests chan<- elevator_io.ButtonEvent) {

	e.Requests[e.Floor][elevator_io.BT_Cab] = false
	ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_Cab}

	switch e.Dirn {

	case elevator_io.MD_Up:
		if e.Requests[e.Floor][elevator_io.BT_HallUp] {
			e.Requests[e.Floor][elevator_io.BT_HallUp] = false
			ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallUp}
		} else if !Requests_above(e) {
			// Only proceed to clear the hall down button if there are no requests above
			// and this block will only be reached if the hall up button was not active (cleared in a previous iteration or not pressed).
			// This ensures that the function needs to be called again to clear the hall down button.
			if e.Requests[e.Floor][elevator_io.BT_HallDown] {
				e.Requests[e.Floor][elevator_io.BT_HallDown] = false
				ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallDown}
			}
		}

	case elevator_io.MD_Down:
		if e.Requests[e.Floor][elevator_io.BT_HallDown] {
			e.Requests[e.Floor][elevator_io.BT_HallDown] = false
			ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallDown}
		} else if !Requests_below(e) {
			// Only proceed to clear the hall up button if there are no requests below
			// and this block will only be reached if the hall down button was not active (cleared in a previous iteration or not pressed).
			// This ensures that the function needs to be called again to clear the hall up button.
			if e.Requests[e.Floor][elevator_io.BT_HallUp] {
				e.Requests[e.Floor][elevator_io.BT_HallUp] = false
				ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallUp}
			}
		}

	case elevator_io.MD_Stop:
		e.Requests[e.Floor][elevator_io.BT_HallUp] = false
		ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallUp}
		e.Requests[e.Floor][elevator_io.BT_HallDown] = false
		ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallDown}
	}
}

*/

// This is a revised version of the one that should erase the edge-condition but also be more effective code (hopefully)
func Requests_clearAtCurrentFloor(e *elevator.Elevator, ch_completedRequests chan<- elevator_io.ButtonEvent) {

	e.Requests[e.Floor][elevator_io.BT_Cab] = false
	ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_Cab}

	switch e.Dirn {

	case elevator_io.MD_Up:
		if e.Requests[e.Floor][elevator_io.BT_HallUp] {
			e.Requests[e.Floor][elevator_io.BT_HallUp] = false
			ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallUp}
		} else if e.Requests[e.Floor][elevator_io.BT_HallDown] {
			e.Requests[e.Floor][elevator_io.BT_HallDown] = false
			ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallDown}
		}

	case elevator_io.MD_Down:
		if e.Requests[e.Floor][elevator_io.BT_HallDown] {
			e.Requests[e.Floor][elevator_io.BT_HallDown] = false
			ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallDown}
		} else if e.Requests[e.Floor][elevator_io.BT_HallUp] {
			e.Requests[e.Floor][elevator_io.BT_HallUp] = false
			ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallUp}
		}

	case elevator_io.MD_Stop:
		e.Requests[e.Floor][elevator_io.BT_HallUp] = false
		ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallUp}
		e.Requests[e.Floor][elevator_io.BT_HallDown] = false
		ch_completedRequests <- elevator_io.ButtonEvent{BtnFloor: e.Floor, BtnType: elevator_io.BT_HallDown}
	}
}

func Requests_announceDirectionChange(e *elevator.Elevator) {
	println("***** CHANGING DIRCETION *****")

	if e.Dirn == elevator_io.MD_Up {
		println("***** GOING UP *****")
	} else if e.Dirn == elevator_io.MD_Down {
		println("***** GOING DOWN *****")
	} else if e.Dirn == elevator_io.MD_Stop {
		println("***** STAYING *****")
	} else {
		println("Direction is undefined.")
	}
}
