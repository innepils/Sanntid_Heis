package main

// Assuming the Elevator struct and constants like N_FLOORS, N_BUTTONS, Dirn, ElevatorBehaviour,
// Button, CV_All, CV_InDirn are defined elsewhere in Go.

type DirnBehaviourPair struct {
	Dirn      Dirn
	Behaviour ElevatorBehaviour
}

func Requests_above(e Elevator) bool {
	for f := e.Floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_below(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func Requests_here(e Elevator) bool {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.Requests[e.Floor][btn] {
			return true
		}
	}
	return false
}

func Requests_chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case D_Up:
		return DirnBehaviourPair{
			Dirn: D_Up,
			Behaviour: func() ElevatorBehaviour {
				if requests_above(e) {
					return EB_Moving
				} else if requests_here(e) {
					return EB_DoorOpen
				} else if requests_below(e) {
					return EB_Moving
				} else {
					return EB_Idle
				}
			}(),
		}
	case D_Down:
		return DirnBehaviourPair{
			Dirn: D_Down,
			Behaviour: func() ElevatorBehaviour {
				if requests_below(e) {
					return EB_Moving
				} else if requests_here(e) {
					return EB_DoorOpen
				} else if requests_above(e) {
					return EB_Moving
				} else {
					return EB_Idle
				}
			}(),
		}
	case D_Stop:
		return DirnBehaviourPair{
			Dirn: D_Stop,
			Behaviour: func() ElevatorBehaviour {
				if requests_here(e) {
					return EB_DoorOpen
				} else if requests_above(e) {
					return EB_Moving
				} else if requests_below(e) {
					return EB_Moving
				} else {
					return EB_Idle
				}
			}(),
		}
	default:
		return DirnBehaviourPair{Dirn: D_Stop, Behaviour: EB_Idle}
	}
}

func Requests_shouldStop(e Elevator) bool {
	switch e.Dirn {
	case D_Down:
		return e.Requests[e.Floor][B_HallDown] || e.Requests[e.Floor][B_Cab] || !requests_below(e)
	case D_Up:
		return e.Requests[e.Floor][B_HallUp] || e.Requests[e.Floor][B_Cab] || !requests_above(e)
	case D_Stop:
		fallthrough
	default:
		return true
	}
}

func Requests_shouldClearImmediately(e Elevator, btnFloor int, btnType Button) bool {
	switch e.Config.ClearRequestVariant {
	case CV_All:
		return e.Floor == btnFloor
	case CV_InDirn:
		return e.Floor == btnFloor &&
			((e.Dirn == D_Up && btnType == B_HallUp) || (e.Dirn == D_Down && btnType == B_HallDown) || e.Dirn == D_Stop || btnType == B_Cab)
	default:
		return false
	}
}

func Requests_clearAtCurrentFloor(e Elevator) Elevator {
	switch e.Config.ClearRequestVariant {
	case CV_All:
		for btn := 0; btn < N_BUTTONS; btn++ {
			e.Requests[e.Floor][btn] = false
		}
	case CV_InDirn:
		e.Requests[e.Floor][B_Cab] = false
		switch e.Dirn {
		case D_Up:
			if !requests_above(e) && !e.Requests[e.Floor][B_HallUp] {
				e.Requests[e.Floor][B_HallDown] = false
			}
			e.Requests[e.Floor][B_HallUp] = false
		case D_Down:
			if !requests_below(e) && !e.Requests[e.Floor][B_HallDown] {
				e.Requests[e.Floor][B_HallUp] = false
			}
			e.Requests[e.Floor][B_HallDown] = false
		case D_Stop:
			fallthrough
		default:
			e.Requests[e.Floor][B_HallUp] = false
			e.Requests[e.Floor][B_HallDown] = false
		}
	}
	return e
}
