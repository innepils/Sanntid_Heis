package requests

import (
	"driver/config"
	"driver/elevator"
	"driver/elevator_io"
)

func Above(e *elevator.Elevator) bool {
	for f := e.Floor + 1; f < config.N_FLOORS; f++ {
		for btn := 0; btn < config.N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}
func Below(e *elevator.Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < config.N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Here(e *elevator.Elevator) bool {
	for btn := 0; btn < config.N_BUTTONS; btn++ {
		if e.Requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

func ChooseDirection(e *elevator.Elevator) {
	switch e.Dirn {
	case elevator_io.MD_Up:
		if Above(e) {
			e.Dirn = elevator_io.MD_Up
			e.Behaviour = elevator.EB_Moving
		} else if Here(e) {
			e.Dirn = elevator_io.MD_Down
			e.Behaviour = elevator.EB_DoorOpen
		} else if Below(e) {
			e.Dirn = elevator_io.MD_Down
			e.Behaviour = elevator.EB_Moving
		} else {
			e.Dirn = elevator_io.MD_Stop
			e.Behaviour = elevator.EB_Idle
		}

	case elevator_io.MD_Down:
		if Below(e) {
			e.Dirn = elevator_io.MD_Down
			e.Behaviour = elevator.EB_Moving
		} else if Here(e) {
			e.Dirn = elevator_io.MD_Up
			e.Behaviour = elevator.EB_DoorOpen
		} else if Above(e) {
			e.Dirn = elevator_io.MD_Up
			e.Behaviour = elevator.EB_Moving
		} else {
			e.Dirn = elevator_io.MD_Stop
			e.Behaviour = elevator.EB_Idle
		}

	case elevator_io.MD_Stop:

		if Here(e) {
			e.Dirn = elevator_io.MD_Stop
			e.Behaviour = elevator.EB_DoorOpen
		} else if Above(e) {
			e.Dirn = elevator_io.MD_Up
			e.Behaviour = elevator.EB_Moving
		} else if Below(e) {
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

func ShouldStop(e *elevator.Elevator) bool {
	switch e.Dirn {
	case elevator_io.MD_Down:
		return e.Requests[e.Floor][elevator_io.BT_HallDown] || e.Requests[e.Floor][elevator_io.BT_Cab] || !Below(e)
	case elevator_io.MD_Up:
		return e.Requests[e.Floor][elevator_io.BT_HallUp] || e.Requests[e.Floor][elevator_io.BT_Cab] || !Above(e)
	default:
		return true
	}
}

// This is a revised version of the one that should erase the edge-condition but also be more effective code (hopefully)
func ClearAtCurrentFloor(e *elevator.Elevator, ch_completedRequests chan<- elevator_io.ButtonEvent) {

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

func AnnounceDirectionChange(e *elevator.Elevator) {
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
