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

//	Mangler muligens channels for ordre-håndtering?
//  	Vi må se litt mer på hvordan det skal implementeres

//  Usikker på hva printen i starten av de forkjsellige casene gjør (utenom for button)

//*****************************************************************************


// One single function for the Final State Machine, to be run as a goroutine from main
func Fsm(ch_arrivalFloor chan int,
		 ch_buttonPressed chan elevator.button, //Usikker på denne elevator.button
		 ch_doorTimedOut chan bool,
		 ch_doorObstruction chan bool, 
		 ch_stopButton chan bool,
		) {
			
		// Do initializing
		localElevator := elevator.UninitializedElevator()
		doorTimer := time.NewTimer(time.Duration(localElevator.config.doorOpenDuration_s) * time.Second);
		
			
		// "For-Select" to supervise the different channels/events that changes the FSM
		for {
			select {
			case buttonPressed <- ch_buttonPressed:

				fmt.Printf("\n\n%s(%d, %s)\n", functionName, buttonPressed.btn_floor, elevio_button_toString(buttonPressed.btn_type))
				elevator.elevator_print(localElevator)

				switch localElevator.behaviour {
					case elevator.EB_DoorOpen:
						if requests.Requests_shouldClearImmediately(localElevator, buttonPressed.btn_floor, buttonPressed.btn_type) {
							doorTimer.Reset(time.Duration(localElevator.config.DoorOpenDuration) * time.Second)
						} else {
							localElevator.Requests[buttonPressed.btn_floor][buttonPressed.btn_type] = 1
						}

					case elevator.EB_Moving:
						localElevator.Requests[buttonPressed.btn_floor][buttonPressed.btn_type] = 1

					case elevator.EB_Idle:
						localElevator.Requests[buttonPressed.btn_floor][buttonPressed.btn_type] = 1
						pair := requests.Requests_chooseDirection(localElevator)
						localElevator.dirn = pair.dirn
						localElevator.behaviour = pair.behaviour

						switch pair.behaviour {
							case elevator.EB_DoorOpen:
								elevator_io.SetDoorOpenLamp(true)
								doorTimer.Reset(time.Duration(localElevator.config.DoorOpenDuration) * time.Second)
								localElevator = Requests_clearAtCurrentFloor(localElevator)

							case elevator.EB_Moving:
								elevator_io.SetMotorDirection(localElevator.dirn)

							case elevator.EB_Idle:
								// No action needed
							}
   			 	}//switch e.behaviour


			case newFloor <- ch_arrivalFloor:
				
				fmt.Printf("\n\n%s(%d)\n", functionName, newFloor)
				elevator.elevator_print(localElevator)
				
				localElevator.floor = newFloor
				elevator_io.SetFloorIndicator(localElevator.floor)	
				
				switch localElevator.behaviour {
					case elevator.EB_Moving:
						if request.Requests_shouldStop(localElevator) {
							elevator_io.SetMotorDirection(D_Stop)
							elevator_io.SetDoorOpenLamp(true)
							localElevator = requests.Requests_clearAtCurrentFloor(localElevator)
							doorTimer.Reset(time.Duration(localElevator.config.DoorOpenDuration) * time.Second)
							localElevator.behaviour = elevator.EB_DoorOpen
						}
					case elevator.EB_DoorOpen:
						// Should not be possible
					case elevator.EB_Idle:
						// Should not be possible

				}
			
				// This channel automatically "transmits" when the timer times out. 
			case: <- ch_doorTimer.C:

				fmt.Printf("\n\n%s()\n", functionName)
				elevator.elevator_print(localElevator)
				
				switch localElevator.behaviour {
					case elevator.EB_DoorOpen:
						pair := request.Requests_chooseDirection(localElevator)
						localElevator.dirn = pair.dirn
						localElevator.behaviour = pair.behaviour
						
						switch localElevator.behaviour {
							case elevator.EB_DoorOpen:
								doorTimer.Reset(time.Duration(localElevator.config.DoorOpenDuration) * time.Second)
								localElevator = request.Requests_clearAtCurrentFloor(localElevator)
								
							case elevator.EB_Moving, elevator.EB_Idle:
								elevator_io.SetDoorOpenLamp(false)
								elevator_io.SetMotorDirection(localElevator.dirn)
							}
					}

			case: <- ch_doorObstruction:

				fmt.Printf("\n\n%s()\n", functionName)
				elevator.elevator_print(localElevator)

				switch localElevator.behaviour {
					case elevator.EB_DoorOpen:
						doorTimer.Reset(time.Duration(localElevator.config.DoorOpenDuration) * time.Second)
					case elevator.EB_Moving, elevator.EB_Idle:
						// Do nothing, print message?
				}

			// Loops as long as something (true) is received on the stopbutton-channel.
			case: <- ch_stopButton:

				fmt.Printf("\n\n%s()\n", functionName)
				elevator.elevator_print(localElevator)

				switch localElevator.behaviour {
					case elevator.EB_DoorOpen:
						doorTimer.Reset(time.Duration(localElevator.config.DoorOpenDuration) * time.Second)
						elevator.io.SetDoorOpenLamp(true)
						
					case elevator.EB_Moving:
						elevator_io.SetMotorDirection(elevator.D_Stop)
						localElevator.behaviour = elevator.EB_Idle  // Might be wrong if Idle also means at a floor
						
					case elevator.EB_Idle:
						// Do nothing
				}

				stopButtonPressed := true
				for stopButtonPressed; {
					select {
						case stopButtonPressed:
							stopButtonPressed = false;	//Might not be needed if the opposite of a signal on the channel is "false"
							stopButtonPressed = <- ch_stopButton:	
					}
				}

			}// select
		}// for

	fmt.Println("\nNew state:")
	elevator.elevator_print(localElevator)
}