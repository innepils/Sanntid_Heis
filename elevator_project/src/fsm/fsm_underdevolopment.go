package fsm

import (
	"driver/con_load"
	"driver/elevator"
	"driver/elevator_io_types"
	"driver/requests"
	"driver/timer"
	"fmt"
)

// One single function for the Final State Machine, to be run as a goroutine from main
func Fsm(ch_arrivalFloor chan,
		 ch_buttonPressed,
		 ch_doorTimedOut,
		 ch_doorObstruction 
		) {
		// Do initializing
			localElevator := elevator.UninitializedElevator()
			
		// "For-Select" to supervise the different channels/events that changes the FSM
		for {
			select {
			case button <- ch_buttonPressed:

				fmt.Printf("\n\n%s(%d, %s)\n", functionName, btn_floor, elevio_button_toString(btn_type))
				elevator_print(localElevator)

				switch elevator.behaviour {
					case elevator.EB_DoorOpen:
						if Requests_shouldClearImmediately(elevator, btn_floor, btn_type) {
							TimerStart(localElevator.config.doorOpenDuration_s)
						} else {
							localElevator.Requests[btn_floor][btn_type] = 1
						}

					case elevator.EB_Moving:
						localElevator.Requests[btn_floor][btn_type] = 1

					case elevator.EB_Idle:
						localElevator.Requests[btn_floor][btn_type] = 1
						pair := Requests_chooseDirection(localElevator)
						localElevator.dirn = pair.dirn
						localElevator.behaviour = pair.behaviour

						switch pair.behaviour {
							case elevator.EB_DoorOpen:
								outputDevice.doorLight(1)
								TimerStart(localElevator.config.doorOpenDuration_s)
								elevator = Requests_clearAtCurrentFloor(elevator)

							case elevator.EB_Moving:
								outputDevice.motorDirection(localElevator.dirn)

							case elevator.EB_Idle:
								// No action needed
							}
					// setAllLights(localElevator) // "Denne skal PRØVE å sette på lys, mens dette må vel tas hånd om senere, etter UDP-greier"
   			 	}//switch e.behaviour


			case newFloor <- ch_arrivalFloor:
				
				fmt.Printf("\n\n%s(%d)\n", functionName, newFloor)
				elevator_print(localElevator)
				
				localElevator.floor = newFloor
				outputDevice.floorIndicator(localElevator.floor)
				
				switch localElevator.behaviour {
					case elevator.EB_Moving:
						if requests_shouldStop(localElevator) {
							outputDevice.motorDirection(D_Stop)
							outputDevice.doorLight(1)
							localElevator = Requests_clearAtCurrentFloor(localElevator)
							TimerStart(localElevator.config.doorOpenDuration_s)
							//setAllLights(localElevator) // DANGER?
							localElevator.behaviour = elevator.EB_DoorOpen
						}
					case elevator.EB_DoorOpen:
						// Should not be possible
					case elevator.EB_Idle:
						// Should not be possible

				}
			
			
			case: <- ch_doorTimedOut:

				fmt.Printf("\n\n%s()\n", functionName)
				elevator_print(localElevator)
				
				switch localElevator.behaviour {
					case elevator.EB_DoorOpen:
						pair := requests_chooseDirection(localElevator)
						localElevator.dirn = pair.dirn
						localElevator.behaviour = pair.behaviour
						
						switch localElevator.behaviour {
							case elevator.EB_DoorOpen:
								TimerStart(localElevator.config.doorOpenDuration_s)
								localElevator = requests_clearAtCurrentFloor(localElevator)
								//setAllLights(elevator)	// DANGER?

							case elevator.EB_Moving, elevator.EB_Idle:
								outputDevice.doorLight(0)
								outputDevice.motorDirection(localElevator.dirn)
							}
					}

			case: <- ch_doorObstruction:

				fmt.Printf("\n\n%s()\n", functionName)
				elevator_print(localElevator)

				switch localElevator.behaviour {
					case elevator.EB_DoorOpen:
						// RESET DOOR TIMER.
					case elevator.EB_Moving, elevator.EB_Idle:
						// Do nothing
				}

			case: // 


				// For every channel/event, check what state (behaviour) the elevator has
				// Events are:
					// New Floor / Arrival at floor
					// Button is pressed
					// Door-timeout
					// Obstruction
					// Stop-button (or is it included in the button-channel?)
					// New-order is confirmed?

				
			}//Select
		}//for
		


	fmt.Println("\nNew state:")
	elevator_print(localElevator)
}