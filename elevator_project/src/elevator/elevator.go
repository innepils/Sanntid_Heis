package elevator

import (
	"driver/config"
	"driver/elevator_io"
	"fmt"
	"strings"
)

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

type Elevator struct {
	Floor               int
	Dirn                elevator_io.MotorDirection
	Requests            [config.N_FLOORS][config.N_BUTTONS]bool
	Behaviour           ElevatorBehaviour
	ClearRequestVariant config.ClearRequestVariant
}

type ElevatorState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests  [config.N_FLOORS][config.N_BUTTONS - 1]bool `json:"hallRequests"`
	ElevatorState map[string]ElevatorState                    `json:"states"`
}

func ElevBehaviourToString(eb ElevatorBehaviour) string {
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

func ElevDirnToString(d elevator_io.MotorDirection) string {
	switch d {
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

func ElevButtonToString(b elevator_io.ButtonType) string {
	switch b {
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
	fmt.Println("  +--------------------+\n")
}

func UninitializedElevator() Elevator {
	return Elevator{
		Floor:               -1,
		Dirn:                elevator_io.MD_Stop,
		Behaviour:           EB_Idle,
		ClearRequestVariant: config.SystemsClearRequestVariant,
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

// needs new name?
func ElevatorToHRAElevState(localElevator Elevator) map[string]ElevatorState {
	return map[string]ElevatorState{
		"self": ElevatorState{
			Behavior:    strings.ToLower(ElevBehaviourToString(localElevator.Behaviour)[3:]),
			Floor:       localElevator.Floor,
			Direction:   strings.ToLower(ElevDirnToString(localElevator.Dirn)),
			CabRequests: GetCabRequests(localElevator),
		},
	}
}

func SendLocalElevatorState(
	localElevator Elevator,
	ch_elevatorStateToAssigner chan map[string]ElevatorState,
	ch_elevatorStateToNetwork chan map[string]ElevatorState) {

	HRAElevState := ElevatorToHRAElevState(localElevator)

	ch_elevatorStateToAssigner <- HRAElevState
	ch_elevatorStateToNetwork <- HRAElevState
}

func SetAllButtonLights(requests [config.N_FLOORS][config.N_BUTTONS]int) {
	for i := range requests {
		for j := range requests[i] {
			if requests[i][j] == 2 {
				elevator_io.SetButtonLamp(elevator_io.ButtonType(j), i, true)
			} else {
				elevator_io.SetButtonLamp(elevator_io.ButtonType(j), i, false)
			}
		}
	}
}
