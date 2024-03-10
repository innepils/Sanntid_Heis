package requests

import (
	"driver/config"
	"driver/elevator"
	"driver/elevator_io"
)

type DirnBehaviourPair struct {
	Dirn      elevator_io.MotorDirection
	Behaviour elevator.ElevatorBehaviour
}

func Requests_above(e elevator.Elevator) bool {
	for f := e.Floor + 1; f < config.N_FLOORS; f++ {
		for btn := 0; btn < config.N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_below(e elevator.Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < config.N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_here(e elevator.Elevator) bool {
	for btn := 0; btn < config.N_BUTTONS; btn++ {
		if e.Requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

func Requests_chooseDirection(e elevator.Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case elevator_io.MD_Up:
		if Requests_above(e) {
			return DirnBehaviourPair{elevator_io.MD_Up, elevator.EB_Moving}
		} else if Requests_here(e) {
			return DirnBehaviourPair{elevator_io.MD_Down, elevator.EB_DoorOpen}
		} else if Requests_below(e) {
			return DirnBehaviourPair{elevator_io.MD_Down, elevator.EB_Moving}
		} else {
			return DirnBehaviourPair{elevator_io.MD_Stop, elevator.EB_Idle}
		}

	case elevator_io.MD_Down:
		if Requests_below(e) {
			return DirnBehaviourPair{elevator_io.MD_Down, elevator.EB_Moving}
		} else if Requests_here(e) {
			return DirnBehaviourPair{elevator_io.MD_Up, elevator.EB_DoorOpen}
		} else if Requests_above(e) {
			return DirnBehaviourPair{elevator_io.MD_Up, elevator.EB_Moving}
		} else {
			return DirnBehaviourPair{elevator_io.MD_Stop, elevator.EB_Idle}
		}

	case elevator_io.MD_Stop:

		if Requests_here(e) {
			return DirnBehaviourPair{elevator_io.MD_Stop, elevator.EB_DoorOpen}
		} else if Requests_above(e) {
			return DirnBehaviourPair{elevator_io.MD_Up, elevator.EB_Moving}
		} else if Requests_below(e) {
			return DirnBehaviourPair{elevator_io.MD_Down, elevator.EB_Moving}
		} else {
			return DirnBehaviourPair{elevator_io.MD_Stop, elevator.EB_Idle}
		}

	default:
		return DirnBehaviourPair{Dirn: elevator_io.MD_Stop, Behaviour: elevator.EB_Idle}
	}
}

func Requests_shouldStop(e elevator.Elevator) bool {
	switch e.Dirn {
	case elevator_io.MD_Down:
		return e.Requests[e.Floor][elevator_io.BT_HallDown] || e.Requests[e.Floor][elevator_io.BT_Cab] || !Requests_below(e)
	case elevator_io.MD_Up:
		return e.Requests[e.Floor][elevator_io.BT_HallUp] || e.Requests[e.Floor][elevator_io.BT_Cab] || !Requests_above(e)
	case elevator_io.MD_Stop:
		fallthrough
	default:
		return true
	}
}

// NOT IN USE?
func Requests_shouldClearImmediately(e elevator.Elevator, btnFloor int, btnType elevator_io.ButtonType) bool {
		return e.Floor == btnFloor &&
			((e.Dirn == elevator_io.MD_Up && btnType == elevator_io.BT_HallUp) || (e.Dirn == elevator_io.MD_Down && btnType == elevator_io.BT_HallDown) || e.Dirn == elevator_io.MD_Stop || btnType == elevator_io.BT_Cab)

}

func Requests_clearAtCurrentFloor(e elevator.Elevator, ch_completedRequests chan<- elevator_io.ButtonEvent) elevator.Elevator {

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
	return e
}
