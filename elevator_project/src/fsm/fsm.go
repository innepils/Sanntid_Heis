package fsm

import (
	"driver/config"
	"driver/elevator"
	"driver/elevator_io"
	"driver/requests"
	"fmt"
	"time"
)

// One single function for the Final State Machine, to be run as a goroutine from main
func Fsm(ch_arrivalFloor chan int,
	ch_localRequests chan [config.N_FLOORS][config.N_BUTTONS]bool,
	ch_doorObstruction chan bool,
	ch_stopButton chan bool,
	ch_completedRequests chan elevator_io.ButtonEvent,
	ch_elevatorStateToAssigner chan map[string]elevator.ElevatorState,
	ch_elevatorStateToNetwork chan elevator.ElevatorState,
) {

	// Initializing
	fmt.Printf("*****INITIALIZING ELEVATOR*****\n")
	localElevator := elevator.UninitializedElevator()

	// If elevator is between floors, run it downwards until a floor is reached.
	elevator_io.SetMotorDirection(elevator_io.MD_Down)
	newFloor := <-ch_arrivalFloor
	elevator_io.SetMotorDirection(elevator_io.MD_Stop)
	localElevator.Floor = newFloor

	// Initialize door
	elevator_io.SetDoorOpenLamp(false)
	doorTimer := time.NewTimer(time.Duration(config.DoorOpenDurationSec) * time.Second)
	prevObstruction := false

	elevator.SendLocalElevatorState(localElevator, ch_elevatorStateToAssigner, ch_elevatorStateToNetwork)

	// "For-Select" to supervise the different channels/events that changes the FSM
	for {
		select {
		case localRequests := <-ch_localRequests:
			fmt.Printf("Entered Local requests in FSM\n")

			localElevator.Requests = localRequests
			localElevator.Elevator_print()

			switch localElevator.Behaviour {

			case elevator.EB_DoorOpen:

				if requests.Requests_here(localElevator) {
					elevator_io.SetDoorOpenLamp(true)
					if prevObstruction {
						prevObstruction = <-ch_doorObstruction
					}
					doorTimer.Reset(time.Duration(config.DoorOpenDurationSec) * time.Second)
					localElevator = requests.Requests_clearAtCurrentFloor(localElevator, ch_completedRequests)
				}

			case elevator.EB_Idle:
				pair := requests.Requests_chooseDirection(localElevator)
				localElevator.Dirn = pair.Dirn
				localElevator.Behaviour = pair.Behaviour
				elevator.SendLocalElevatorState(localElevator, ch_elevatorStateToAssigner, ch_elevatorStateToNetwork)
				localElevator.Elevator_print()

				switch pair.Behaviour {
				case elevator.EB_DoorOpen:
					elevator_io.SetDoorOpenLamp(true)
					if prevObstruction {
						prevObstruction = <-ch_doorObstruction
					}
					doorTimer.Reset(time.Duration(config.DoorOpenDurationSec) * time.Second)
					localElevator = requests.Requests_clearAtCurrentFloor(localElevator, ch_completedRequests)
					elevator.SendLocalElevatorState(localElevator, ch_elevatorStateToAssigner, ch_elevatorStateToNetwork)

				case elevator.EB_Moving:
					elevator_io.SetMotorDirection(localElevator.Dirn)

				}
			} //switch e.behaviour*/

		case newFloor := <-ch_arrivalFloor:
			fmt.Printf("Entered new floor in FSM\n")
			localElevator.Elevator_print()

			localElevator.Floor = newFloor
			elevator_io.SetFloorIndicator(localElevator.Floor)

			switch localElevator.Behaviour {
			case elevator.EB_Moving:
				if requests.Requests_shouldStop(localElevator) {
					elevator_io.SetMotorDirection(elevator_io.MD_Stop)
					elevator_io.SetDoorOpenLamp(true)
					localElevator = requests.Requests_clearAtCurrentFloor(localElevator, ch_completedRequests)
					if prevObstruction {
						prevObstruction = <-ch_doorObstruction
					}
					doorTimer.Reset(time.Duration(config.DoorOpenDurationSec) * time.Second)
					localElevator.Behaviour = elevator.EB_DoorOpen
				}
			}

		// This channel automatically "transmits" when the timer times out.
		case <-doorTimer.C:
			fmt.Printf("Entered doorTimeout in FSM\n")

			localElevator.Elevator_print()

			switch localElevator.Behaviour {
			case elevator.EB_DoorOpen:

				if prevObstruction {
					prevObstruction = <-ch_doorObstruction
				}
				doorTimer.Reset(time.Duration(config.DoorOpenDurationSec) * time.Second)
				elevator_io.SetDoorOpenLamp(false)

				pair := requests.Requests_chooseDirection(localElevator)
				localElevator.Dirn = pair.Dirn
				localElevator.Behaviour = pair.Behaviour

				switch localElevator.Behaviour {
				case elevator.EB_Moving:
					elevator_io.SetMotorDirection(localElevator.Dirn)
				}
			}

		case obstruction := <-ch_doorObstruction:
			prevObstruction = obstruction

		case <-ch_stopButton:
			fmt.Printf("Entered Stop Button in FSM\n")

			localElevator.Elevator_print()

			switch localElevator.Behaviour {
			case elevator.EB_DoorOpen:
				if prevObstruction {
					prevObstruction = <-ch_doorObstruction
				}
				doorTimer.Reset(time.Duration(config.DoorOpenDurationSec) * time.Second)
				elevator_io.SetDoorOpenLamp(true)

			case elevator.EB_Moving:
				elevator_io.SetMotorDirection(elevator_io.MD_Stop)
				localElevator.Behaviour = elevator.EB_Idle
			}

			// Loops as long as something (true) is received on the stopbutton-channel.
			stopButtonPressed := true
			for stopButtonPressed {
				stopButtonPressed = false
				stopButtonPressed = <-ch_stopButton

			}
			switch localElevator.Behaviour {
			case elevator.EB_DoorOpen:
				if prevObstruction {
					prevObstruction = <-ch_doorObstruction
				}
				doorTimer.Reset(time.Duration(config.DoorOpenDurationSec) * time.Second)
				localElevator = requests.Requests_clearAtCurrentFloor(localElevator, ch_completedRequests)
			case elevator.EB_Idle:
				elevator_io.SetMotorDirection(localElevator.Dirn)
				localElevator.Behaviour = elevator.EB_Moving
			}

			localElevator.Elevator_print()

		default:
			// Do Nothing
		} //select
		elevator.SendLocalElevatorState(localElevator, ch_elevatorStateToAssigner, ch_elevatorStateToNetwork)
	} //For
} //Fsm
