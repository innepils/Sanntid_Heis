package fsm

import (
	"driver/config"
	"driver/elevator"
	"driver/elevator_io"
	"driver/requests"
	"time"
)

//*****************************************************************************
// 						*****	Status	*****
//	Mangler muligens channels for ordre-håndtering?
//  	Vi må se litt mer på hvordan det skal implementeres

//  Usikker på hva printen i starten av de forkjsellige casene gjør (utenom for button)
//  Denne://fmt.Printf("\n\n%s(%d, %s)\n", functionName, buttonPressed.btn_floor, elevio_button_toString(buttonPressed.btn_type)) //functionName should be a string, not sure why this is implemented
//*****************************************************************************

// One single function for the Final State Machine, to be run as a goroutine from main
func Fsm(ch_arrivalFloor chan int,
	ch_buttonPressed chan elevator_io.ButtonEvent, //Usikker på denne elevator.button
	ch_doorObstruction chan bool,
	ch_stopButton chan bool,
) {

	// Do initializing
	localElevator := elevator.UninitializedElevator()
	doorTimer := time.NewTimer(time.Duration(config.DoorOpenDurationSec) * time.Second)

	// "For-Select" to supervise the different channels/events that changes the FSM
	for {
		select {
		case buttonPressed := <-ch_buttonPressed:

			localElevator.Elevator_print()

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

				case elevator.EB_Idle:
					// No action needed
				}
			} //switch e.behaviour

		case newFloor := <-ch_arrivalFloor:

			//fmt.Printf("\n\n%s(%d)\n", functionName, newFloor)
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
				localElevator.Behaviour = elevator.EB_Idle // Might be wrong if Idle also means at a floor

			case elevator.EB_Idle:
				// Do nothing
			}

			stopButtonPressed := true
			for stopButtonPressed {
				stopButtonPressed = false //Might not be needed if the opposite of a signal on the channel is "false"
				stopButtonPressed = <-ch_stopButton

			}

		} // select
	} // for

	//fmt.Println("\nNew state:")
	//localElevator.Elevator_print()

}
