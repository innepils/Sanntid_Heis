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

func (eb ElevatorBehaviour) String() string {
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

func (es Elevator) String() string {
    var b strings.Builder
    b.WriteString("  +--------------------+\n")
    b.WriteString(fmt.Sprintf(
        "  |floor = %-2d          |\n"+
            "  |dirn  = %-12s|\n"+
            "  |behav = %-12s|\n",
        es.Floor, elevio.DirnToString(es.Dirn), es.Behaviour.String(),
    ))
    b.WriteString("  +--------------------+\n")
    b.WriteString("  |  | up  | dn  | cab |\n")
    for f := N_FLOORS - 1; f >= 0; f-- {
        b.WriteString(fmt.Sprintf("  | %d", f))
        for btn := 0; btn < N_BUTTONS; btn++ {
            if (f == N_FLOORS-1 && btn == B_HallUp) || (f == 0 && btn == B_HallDown) {
                b.WriteString("|     ")
            } else {
                b.WriteString(es.Requests[f][btn] ? "|  #  " : "|  -  ")
            }
        }
        b.WriteString("|\n")
    }
    b.WriteString("  +--------------------+\n")
    return b.String()
}

func NewUninitializedElevator() Elevator {
    return Elevator{
        Floor: -1,
        Dirn:  D_Stop,
        Behaviour: EB_Idle,
        Config: struct {
            ClearRequestVariant ClearRequestVariant
            DoorOpenDurationS   float64
        }{
            ClearRequestVariant: CV_All,
            DoorOpenDurationS:   3.0,
        },
    }
}

