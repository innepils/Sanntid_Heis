package fsm

import (
	"driver/config"
	"driver/elevator"
	"driver/elevator_io"
	"driver/requests"
	"fmt"
	"time"
)

//*****************************************************************************
// 						*****	Status	*****
//	Mangler muligens channels for ordre-håndtering?
//  	Vi må se litt mer på hvordan det skal implementeres
//*****************************************************************************

// One single function for the Final State Machine, to be run as a goroutine from main
func Fsm(ch_arrivalFloor chan int,
	ch_buttonPressed chan elevator_io.ButtonEvent,
	ch_doorObstruction chan bool,
	ch_stopButton chan bool,
) {

	// Initializing
	localElevator := elevator.UninitializedElevator()
	doorTimer := time.NewTimer(time.Duration(config.DoorOpenDurationSec) * time.Second)
	elevator_io.SetMotorDirection(elevator_io.MD_Down)

	// Run the elevator to the bottom floor
	for {
		newFloor := <-ch_arrivalFloor

		if newFloor != 0 {
			elevator_io.SetMotorDirection(elevator_io.MD_Down)
		} else {
			elevator_io.SetMotorDirection(elevator_io.MD_Stop)
			localElevator.Floor = newFloor
			break
		}
	}

	// "For-Select" to supervise the different channels/events that changes the FSM
	for {
		select {
		case buttonPressed := <-ch_buttonPressed:

			switch localElevator.Behaviour {
			case elevator.EB_DoorOpen:
				if requests.Requests_shouldClearImmediately(localElevator, buttonPressed.BtnFloor, elevator_io.ButtonType(buttonPressed.BtnType)) {
					doorTimer.Reset(time.Duration(localElevator.Config.DoorOpenDurationSec) * time.Second)
				} else {
					localElevator.Requests[buttonPressed.BtnFloor][buttonPressed.BtnType] = true
				}

			case elevator.EB_Moving:
				localElevator.Requests[buttonPressed.BtnFloor][buttonPressed.BtnType] = true

			case elevator.EB_Idle:
				localElevator.Requests[buttonPressed.BtnFloor][buttonPressed.BtnType] = true
				pair := requests.Requests_chooseDirection(localElevator)
				localElevator.Dirn = pair.Dirn
				localElevator.Behaviour = pair.Behaviour

				switch pair.Behaviour {
				case elevator.EB_DoorOpen:
					elevator_io.SetDoorOpenLamp(true)
					doorTimer.Reset(time.Duration(localElevator.Config.DoorOpenDurationSec) * time.Second)
					localElevator = requests.Requests_clearAtCurrentFloor(localElevator)

				case elevator.EB_Moving:
					elevator_io.SetMotorDirection(elevator_io.MotorDirection(localElevator.Dirn))
					fmt.Printf("DirectionSet: %s\n", elevator.ElevDirnToString(elevator_io.MotorDirection(localElevator.Dirn)))

				case elevator.EB_Idle:
					// No action needed
				}
			} //switch e.behaviour

		case newFloor := <-ch_arrivalFloor:

			localElevator.Elevator_print()

			localElevator.Floor = newFloor
			elevator_io.SetFloorIndicator(localElevator.Floor)

			switch localElevator.Behaviour {
			case elevator.EB_Moving:
				if requests.Requests_shouldStop(localElevator) {
					elevator_io.SetMotorDirection(elevator_io.MD_Stop)
					elevator_io.SetDoorOpenLamp(true)
					localElevator = requests.Requests_clearAtCurrentFloor(localElevator)
					doorTimer.Reset(time.Duration(localElevator.Config.DoorOpenDurationSec) * time.Second)
					localElevator.Behaviour = elevator.EB_DoorOpen
				}
			case elevator.EB_DoorOpen:
				// Should not be possible
			case elevator.EB_Idle:
				// Should not be possible

			}

			// This channel automatically "transmits" when the timer times out.
		case <-doorTimer.C:

			localElevator.Elevator_print()

			switch localElevator.Behaviour {
			case elevator.EB_DoorOpen:
				pair := requests.Requests_chooseDirection(localElevator)
				localElevator.Dirn = pair.Dirn
				localElevator.Behaviour = pair.Behaviour

				switch localElevator.Behaviour {
				case elevator.EB_DoorOpen:
					doorTimer.Reset(time.Duration(localElevator.Config.DoorOpenDurationSec) * time.Second)
					localElevator = requests.Requests_clearAtCurrentFloor(localElevator)

				case elevator.EB_Moving, elevator.EB_Idle:
					elevator_io.SetDoorOpenLamp(false)
					elevator_io.SetMotorDirection(elevator_io.MotorDirection(localElevator.Dirn))
				}
			}

		case <-ch_doorObstruction:

			localElevator.Elevator_print()

			switch localElevator.Behaviour {
			case elevator.EB_DoorOpen:
				doorTimer.Reset(time.Duration(localElevator.Config.DoorOpenDurationSec) * time.Second)
			case elevator.EB_Moving, elevator.EB_Idle:
				// Do nothing, print message?
			}

		// Loops as long as something (true) is received on the stopbutton-channel.
		case <-ch_stopButton:

			localElevator.Elevator_print()

			switch localElevator.Behaviour {
			case elevator.EB_DoorOpen:
				doorTimer.Reset(time.Duration(localElevator.Config.DoorOpenDurationSec) * time.Second)
				elevator_io.SetDoorOpenLamp(true)

			case elevator.EB_Moving:
				elevator_io.SetMotorDirection(elevator_io.MD_Stop)
				localElevator.Behaviour = elevator.EB_Idle

			case elevator.EB_Idle:
				// Do nothing
			}

			stopButtonPressed := true
			for stopButtonPressed {
				stopButtonPressed = false //Might not be needed if the opposite of a signal on the channel is "false"
				stopButtonPressed = <-ch_stopButton

			}

			switch localElevator.Behaviour {
			case elevator.EB_DoorOpen:
				doorTimer.Reset(time.Duration(localElevator.Config.DoorOpenDurationSec) * time.Second)
			case elevator.EB_Idle:
				elevator_io.SetMotorDirection(localElevator.Dirn)
				localElevator.Behaviour = elevator.EB_Moving

			}
		} // select

		localElevator.Elevator_print()

	} // for
}
