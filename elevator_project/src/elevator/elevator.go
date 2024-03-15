package elevator

import (
	"src/config"
	"src/elevator_io"
	"strings"
)

type (
	ElevatorBehaviour int
	RequestType       int
)

const (
	EB_Idle     ElevatorBehaviour = 0
	EB_DoorOpen ElevatorBehaviour = 1
	EB_Moving   ElevatorBehaviour = 2

	UndefinedRequest RequestType = -1
	NoRequest        RequestType = 0
	NewRequest       RequestType = 1
	ConfirmedRequest RequestType = 2
	CompletedRequest RequestType = 3
)

type Elevator struct {
	Floor     int
	Dirn      elevator_io.MotorDirection
	Behaviour ElevatorBehaviour
	Requests  [config.N_FLOORS][config.N_BUTTONS]bool
}

type HRAElevatorState struct {
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	Behavior    string `json:"behaviour"`
	CabRequests []bool `json:"cabRequests"`
}

func UninitializedElevator() Elevator {
	return Elevator{
		Floor:     -1,
		Dirn:      elevator_io.MD_Stop,
		Behaviour: EB_Idle,
		// Requests are received from assigner.
	}
}

func GetCabRequests(elevator Elevator) []bool {
	cabRequests := make([]bool, len(elevator.Requests))
	for floor, request := range elevator.Requests {
		cabRequests[floor] = request[elevator_io.BT_Cab]
	}
	return cabRequests
}

func SendLocalElevatorState(
	nodeID 						string,
	localElevator 				Elevator,
	ch_elevatorStateToAssigner 	chan<- map[string]HRAElevatorState,
	ch_elevatorStateToNetwork 	chan<- HRAElevatorState,
) {

	elevatorState := ElevToHRAElevatorState(nodeID, localElevator)

	ch_elevatorStateToAssigner <- elevatorState
	ch_elevatorStateToNetwork <- elevatorState[nodeID]
}

func SetAllButtonLights(requests [config.N_FLOORS][config.N_BUTTONS]RequestType) {
	for floor := range requests {
		for btn := range requests[floor] {
			if requests[floor][btn] == ConfirmedRequest {
				elevator_io.SetButtonLamp(elevator_io.ButtonType(btn), floor, true)
			} else {
				elevator_io.SetButtonLamp(elevator_io.ButtonType(btn), floor, false)
			}
		}
	}
}

// *** Functions for changing datatype ***

func ElevBehaviourToString(elevBehaviour ElevatorBehaviour) string {
	switch elevBehaviour {
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

func ElevDirnToString(elevDirection elevator_io.MotorDirection) string {
	switch elevDirection {
	case elevator_io.MD_Down:
		return "Down"
	case elevator_io.MD_Stop:
		return "Stop"
	case elevator_io.MD_Up:
		return "Up"
	default:
		return "Unknown"
	}
}

func ElevToHRAElevatorState(id string, localElevator Elevator) map[string]HRAElevatorState {
	return map[string]HRAElevatorState{
		id: {
			Floor:       localElevator.Floor,
			Direction:   strings.ToLower(ElevDirnToString(localElevator.Dirn)),
			Behavior:    strings.ReplaceAll(strings.ToLower(ElevBehaviourToString(localElevator.Behaviour)[3:]), "open", "Open"),
			CabRequests: GetCabRequests(localElevator),
		},
	}
}

// The Elevator_print() is currently not used,
// but remains included for potential debugging purposes,
// facilitating future maintenance or expansion efforts.
/*
func (e *Elevator) Elevator_print() {
	fmt.Println("  +--------------------+")
	fmt.Printf(
		"  |floor = %-2d          |\n"+
			"  |dirn  = %-12.12s|\n"+
			"  |behav = %-12.12s|\n",
		e.Floor,
		ElevDirnToString(e.Dirn),
		ElevBehaviourToString(e.Behaviour),
	)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := config.N_FLOORS - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < config.N_BUTTONS; btn++ {
			btnType := elevator_io.ButtonType(btn)
			if ((f == config.N_FLOORS-1) && (btnType == elevator_io.BT_HallUp)) ||
				(f == 0 && btnType == elevator_io.BT_HallDown) {
				fmt.Print("|     ")
			} else {
				if e.Requests[f][btn] {
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
*/