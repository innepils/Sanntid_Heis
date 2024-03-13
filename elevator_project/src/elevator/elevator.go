package elevator

import (
	"driver/config"
	"driver/elevator_io"
	"fmt"
	"strings"
	"time"
)

type (
	ElevatorBehaviour int
	RequestType       int
)

const (
	EB_Idle     ElevatorBehaviour = 0
	EB_DoorOpen ElevatorBehaviour = 1
	EB_Moving   ElevatorBehaviour = 2

	NoOrder        RequestType = 0
	NewOrder       RequestType = 1
	ConfirmedOrder RequestType = 2
	CompletedOrder RequestType = 3
)

type Elevator struct {
	Floor     int
	Dirn      elevator_io.MotorDirection
	Behaviour ElevatorBehaviour
	Requests  [config.N_FLOORS][config.N_BUTTONS]bool
}

type ElevatorState struct {
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
	for i, row := range elevator.Requests {
		cabRequests[i] = row[len(row)-1]
	}
	return cabRequests
}

func SendLocalElevatorState(
	id 							string,
	localElevator 				Elevator,
	ch_elevatorStateToAssigner 	chan<- map[string]ElevatorState,
	ch_elevatorStateToNetwork 	chan<- ElevatorState) {

	elevatorState := ElevToElevatorState(id, localElevator)
	ch_elevatorStateToAssigner <- elevatorState
	ch_elevatorStateToNetwork <- elevatorState[id]
}

func SetAllButtonLights(requests [config.N_FLOORS][config.N_BUTTONS]RequestType) {
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

func (e *Elevator) HoldDoorOpenIfObstruction(
	prevObstruction *bool,
	doorTimer *time.Timer,
	ch_doorObstruction <-chan bool,
) {
	if *prevObstruction {
		fmt.Println("Door is obstructed")
		*prevObstruction = <-ch_doorObstruction
		doorTimer.Reset(time.Duration(config.DoorOpenDurationSec) * time.Second)
	}
}

func (e *Elevator) StallWhileStopButtonActive(ch_stopButton <-chan bool) {
	stopButtonPressed := true
	for stopButtonPressed {
		stopButtonPressed = false
		stopButtonPressed = <-ch_stopButton
	}
}

// Functions for changing datatype
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

func ElevToElevatorState(id string, localElevator Elevator) map[string]ElevatorState {
	return map[string]ElevatorState{
		id: {
			Floor:       localElevator.Floor,
			Direction:   strings.ToLower(ElevDirnToString(localElevator.Dirn)),
			Behavior:    strings.ReplaceAll(strings.ToLower(ElevBehaviourToString(localElevator.Behaviour)[3:]), "open", "Open"),
			CabRequests: GetCabRequests(localElevator),
		},
	}
}

// Printing
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
