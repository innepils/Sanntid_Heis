package requests

import (
	"driver/elevator"
	"driver/elevator_io_types"
)

// Assuming the Elevator struct and constants like N_FLOORS, N_BUTTONS, Dirn, ElevatorBehaviour,
// Button, CV_All, CV_InDirn are defined elsewhere in Go.

type DirnBehaviourPair struct {
	Dirn      elevator_io_types.Dirn
	Behaviour elevator.ElevatorBehaviour
}

func Requests_above(e elevator.Elevator) bool {
	for f := e.Floor + 1; f < elevator_io_types.N_FLOORS; f++ {
		for btn := 0; btn < elevator_io_types.N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_below(e elevator.Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < elevator_io_types.N_BUTTONS; btn++ {
			if e.Request[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_here(e elevator.Elevator) bool {
	for btn := 0; btn < elevator_io_types.N_BUTTONS; btn++ {
		if e.Requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

func Requests_chooseDirection(e elevator.Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case elevator_io_types.D_Up:
		return DirnBehaviourPair{
			Dirn: elevator_io_types.D_Up,
			Behaviour: func() elevator.ElevatorBehaviour {
				if requests_above(e) {
					return elevator.EB_Moving
				} else if requests_here(e) {
					return elevator.EB_DoorOpen
				} else if requests_below(e) {
					return elevator.EB_Moving
				} else {
					return elevator.EB_Idle
				}
			}(),
		}
	case elevator_io_types.D_Down:
		return DirnBehaviourPair{
			Dirn: elevator_io_types.D_Down,
			Behaviour: func() elevator.ElevatorBehaviour {
				if requests_below(e) {
					return elevator.EB_Moving
				} else if requests_here(e) {
					return elevator.EB_DoorOpen
				} else if requests_above(e) {
					return elevator.EB_Moving
				} else {
					return elevator.EB_Idle
				}
			}(),
		}
	case elevator_io_types.D_Stop:
		return DirnBehaviourPair{
			Dirn: elevator_io_types.D_Stop,
			Behaviour: func() elevator.ElevatorBehaviour {
				if requests_here(e) {
					return elevator.EB_DoorOpen
				} else if requests_above(e) {
					return EB_Moving
				} else if requests_below(e) {
					return elevator.EB_Moving
				} else {
					return elevator.EB_Idle
				}
			}(),
		}
	default:
		return DirnBehaviourPair{Dirn: elevator_io_types.D_Stop, Behaviour: elevator.EB_Idle}
	}
}

func Requests_shouldStop(e elevator.Elevator) bool {
	switch e.Dirn {
	case elevator_io_types.D_Down:
		return e.Requests[e.Floor][elevator_io_types.B_HallDown] || e.Requests[e.Floor][elevator_io_types.B_Cab] || !requests_below(e)
	case elevator_io_types.D_Up:
		return e.Requests[e.Floor][elevator_io_types.B_HallUp] || e.Requests[e.Floor][elevator_io_types.B_Cab] || !requests_above(e)
	case elevator_io_types.D_Stop:
		fallthrough
	default:
		return true
	}
}

func Requests_shouldClearImmediately(e elevator.Elevator, btnFloor int, btnType elevator_io_types.Button) bool {
	switch e.Config.ClearRequestVariant {
	case elevator.CV_all:
		return e.Floor == btnFloor
	case elevator.CV_InDirn:
		return e.Floor == btnFloor &&
			((e.Dirn == elevator_io_types.D_Up && btnType == elevator_io_types.B_HallUp) || (e.Dirn == elevator_io_types.D_Down && btnType == elevator_io_types.B_HallDown) || e.Dirn == D_Stop || btnType == B_Cab)
	default:
		return false
	}
}

func Requests_clearAtCurrentFloor(e elevator.Elevator) elevator.Elevator {
	switch e.Config.ClearRequestVariant {
	case elevator.CV_all:
		for btn := 0; btn < elevator_io_types.N_BUTTONS; btn++ {
			e.Requests[e.Floor][btn] = false
		}
	case elevator.CV_InDirn:
		e.Requests[e.Floor][elevator_io_types.B_Cab] = false
		switch e.Dirn {
		case elevator_io_types.D_Up:
			if !requests_above(e) && !e.Requests[e.Floor][elevator_io_types.B_HallUp] {
				e.Requests[e.Floor][elevator_io_types.B_HallDown] = false
			}
			e.Requests[e.Floor][elevator_io_types.B_HallUp] = false
		case elevator_io_types.D_Down:
			if !requests_below(e) && !e.Requests[e.Floor][elevator_io_types.B_HallDown] {
				e.Requests[e.Floor][elevator_io_types.B_HallUp] = false
			}
			e.Requests[e.Floor][elevator_io_types.B_HallDown] = false
		case elevator_io_types.D_Stop:
			fallthrough
		default:
			e.Requests[e.Floor][elevator_io_types.B_HallUp] = false
			e.Requests[e.Floor][elevator_io_types.B_HallDown] = false
		}
	}
	return e
}
