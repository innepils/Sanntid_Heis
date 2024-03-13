package elevator

import (
	"driver/config"
	"driver/elevator_io"
	"fmt"
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

	None      RequestType = 0
	New       RequestType = 1
	Confirmed RequestType = 2
	Completed RequestType = 3
)

type Elevator struct {
	Floor     int
	Dirn      elevator_io.MotorDirection
	Requests  [config.N_FLOORS][config.N_BUTTONS]bool
	Behaviour ElevatorBehaviour
}

type ElevatorState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

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

func ElevButtonToString(buttonType elevator_io.ButtonType) string {
	switch buttonType {
	case elevator_io.BT_HallUp:
		return "HallUp"
	case elevator_io.BT_HallDown:
		return "HallDown"
	case elevator_io.BT_Cab:
		return "Cab"
	default:
		return "Unknown"
	}
}

func (es *Elevator) Elevator_print() {
	fmt.Println("  +--------------------+")
	fmt.Printf(
		"  |floor = %-2d          |\n"+
			"  |dirn  = %-12.12s|\n"+
			"  |behav = %-12.12s|\n",
		es.Floor,
		ElevDirnToString(es.Dirn),
		ElevBehaviourToString(es.Behaviour),
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
				if es.Requests[f][btn] {
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
		Dirn:      elevator_io.MD_Stop,
		Behaviour: EB_Idle,
	}
}

func GetCabRequests(elevator Elevator) []bool {
	cabRequests := make([]bool, len(elevator.Requests))
	for i, row := range elevator.Requests {
		cabRequests[i] = row[len(row)-1]
	}

	return cabRequests
}

func ElevToElevatorState(id string, localElevator Elevator) map[string]ElevatorState {
	return map[string]ElevatorState{
		id: {
			Behavior:    strings.ReplaceAll(strings.ToLower(ElevBehaviourToString(localElevator.Behaviour)[3:]), "open", "Open"),
			Floor:       localElevator.Floor,
			Direction:   strings.ToLower(ElevDirnToString(localElevator.Dirn)),
			CabRequests: GetCabRequests(localElevator),
		},
	}
}

func SendLocalElevatorState(
	id string,
	localElevator Elevator,
	ch_elevatorStateToAssigner chan<- map[string]ElevatorState,
	ch_elevatorStateToNetwork chan<- ElevatorState) {

	elevatorState := ElevToElevatorState(id, localElevator)
	ch_elevatorStateToAssigner <- elevatorState
	ch_elevatorStateToNetwork <- elevatorState[id]
}

func SetAllButtonLights(request [config.N_FLOORS][config.N_BUTTONS]RequestType) {
	for i := range request {
		for j := range request[i] {
			if request[i][j] == Confirmed {
				elevator_io.SetButtonLamp(elevator_io.ButtonType(j), i, true)
			} else {
				elevator_io.SetButtonLamp(elevator_io.ButtonType(j), i, false)
			}
		}
	}
}
