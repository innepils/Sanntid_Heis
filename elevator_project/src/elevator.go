package elevator

import (
	"fmt"
)

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

type ClearRequestVariant int

const (
	CV_all ClearRequestVariant = iota
	CV_InDirn
)

type Elevator struct {
	Floor int
	Dirn Dirn
	Request [N_FLOORS][N_BUTTONS]int
	Behaviour ElevatorBehaviour

	Config struct {
		ClearRequestVariant ClearRequestVariant
		DoorOpenDurationSec float64
	}
}

