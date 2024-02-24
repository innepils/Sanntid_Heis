package elevator

import (
	"driver/elevator_io_types"
	"fmt"
)

// SK: Elevator states
type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

// SK:
type ClearRequestVariant int

const (
	/* GIVEN: Assume everyone waiting for the elevator gets on the elevator, even if
	   they will be traveling in the "wrong" direction for a while */
	CV_all ClearRequestVariant = iota
	CV_InDirn
)

// Struct contain
type Elevator struct {
	Floor     int
	Dirn      elevator_io_types.Dirn
	Requests  [elevator_io_types.N_FLOORS][elevator_io_types.N_BUTTONS]bool
	Behaviour ElevatorBehaviour

	Config struct {
		ClearRequestVariant ClearRequestVariant
		DoorOpenDurationSec float64
	}
}

func EBToString(eb ElevatorBehaviour) string {
	switch eb {
	case EB_Idle:
		return "EB_Idle"
	case EB_DoorOpen:
		return "EB_DoorOpen"
	case EB_Moving:
		return "EB_Moving"
	default:
		return "EB_UNDEFINED"
	}
}

// Prints the state of the elevator
func (es *Elevator) Print() {
	fmt.Println("  +--------------------+")
	fmt.Printf(
		"  |floor = %-2d          |\n"+
			"  |dirn  = %-12.12s|\n"+
			"  |behav = %-12.12s|\n",
		es.Floor,
		elevator_io_types.Elevio_dirn_toString(es.Dirn), // Assuming DirnToString function exists
		EBToString(es.Behaviour),
	)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := elevator_io_types.N_FLOORS - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < elevator_io_types.N_BUTTONS; btn++ {
			btnType := elevator_io_types.Button(btn)
			if ((f == elevator_io_types.N_FLOORS-1) && (btnType == elevator_io_types.B_HallUp)) ||
				(f == 0 && btnType == elevator_io_types.B_HallDown) {
				fmt.Print("|     ")
			} else {
				if es.Requests[f][btn] != false {
					fmt.Print("|  #  ")
				} else {
					fmt.Print("|  -  ")
				}
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}

func UninitializedElevator() Elevator {
	return Elevator{
		Floor:     -1,
		Dirn:      elevator_io_types.D_Stop,
		Behaviour: EB_Idle,
		Config: struct {
			ClearRequestVariant ClearRequestVariant
			DoorOpenDurationSec float64
		}{
			ClearRequestVariant: CV_all,
			DoorOpenDurationSec: 3.0,
		},
	}
}

func GetCabRequests(elevator Elevator) []bool {
	// Create a new slice to store the last column elements
	cabRequests := make([]bool, len(elevator.Requests))

	// Loop through each row and access the last element
	for i, row := range elevator.Requests {
		cabRequests[i] = row[len(row)-1]
	}

	return cabRequests
}
