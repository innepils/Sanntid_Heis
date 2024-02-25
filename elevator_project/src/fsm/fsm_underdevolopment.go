package fsm

import (
	"driver/con_load"
	"driver/elevator"
	"driver/elevator_io_types"
	"driver/requests"
	"driver/timer"
	"fmt"
)

//*****************************************************************************
// 						*****	Status	*****

//	SetAllLights er foreløpig ikke implementert. Denne skal (forsøke) skru 
//  på alle lys pr ordre. Men vi må passe på at dette er for den synkroniserte 
//	ordre-matrisen.

//	Mangler muligens channels for ordre-håndtering?
//  	Vi må se litt mer på hvordan det skal implementeres

//  Door-obstruction-casen er ikke ferdig. Må implementere timer osv først.

//  Usikker på hva printen i starten av de forkjsellige casene gjør (utenom for button)

//  Også litt usikker på hvor btn_floor og btn_type skal komme fra. 
//  	Mulig det er en parameter av Button?

//			Ellers good, ser ryddigere ut, og skjønner virkemmåten nå.

//*****************************************************************************


// One single function for the Final State Machine, to be run as a goroutine from main
func Fsm(ch_arrivalFloor chan int,
		 ch_buttonPressed chan elevator.button, //Usikker på denne elevator.button
		 ch_doorTimedOut chan bool,
		 ch_doorObstruction chan bool 
		) {
			
		// Do initializing
			localElevator := elevator.UninitializedElevator()
			
		// "For-Select" to supervise the different channels/events that changes the FSM
		for {
			select {
			case button <- ch_buttonPressed:

				fmt.Printf("\n\n%s(%d, %s)\n", functionName, btn_floor, elevio_button_toString(btn_type))
				elevator_print(localElevator)

				switch localElevator.behaviour {
					case elevator.EB_DoorOpen:
						if Requests_shouldClearImmediately(localElevator, btn_floor, btn_type) {
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
								SetDoorOpenLamp(true)
								TimerStart(localElevator.config.doorOpenDuration_s)
								localElevator = Requests_clearAtCurrentFloor(localElevator)

							case elevator.EB_Moving:
								SetMotorDirection(localElevator.dirn)

							case elevator.EB_Idle:
								// No action needed
							}
					// setAllLights(localElevator) // "Denne skal PRØVE å sette på lys, mens dette må vel tas hånd om senere, etter UDP-greier"
   			 	}//switch e.behaviour


			case newFloor <- ch_arrivalFloor:
				
				fmt.Printf("\n\n%s(%d)\n", functionName, newFloor)
				elevator_print(localElevator)
				
				localElevator.floor = newFloor
				SetFloorIndicator(localElevator.floor)	
				
				switch localElevator.behaviour {
					case elevator.EB_Moving:
						if requests_shouldStop(localElevator) {
							SetMotorDirection(D_Stop)
							SetDoorOpenLamp(true)
							localElevator = Requests_clearAtCurrentFloor(localElevator)
							TimerStart(localElevator.config.doorOpenDuration_s)
							//setAllLights(localElevator)
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
								
								//setAllLights(elevator

							case elevator.EB_Moving, elevator.EB_Idle:
								SetDoorOpenLamp(false)
								SetMotorDirection(localElevator.dirn)
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
					// New Floor / Arrival at floor. DONE
					// Button is pressed			 DONE
					// Door-timeout					 DONE
					// Obstruction
					// Stop-button (or is it included in the button-channel?)
					// New-order is confirmed?

				
			}// select
		}// for
		


	fmt.Println("\nNew state:")
	elevator_print(localElevator)
}