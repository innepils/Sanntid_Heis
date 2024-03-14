package fsm

import (
	"driver/config"
	"driver/elevator"
	"driver/elevator_io"
	"driver/requests"
	"fmt"
	"time"
)

// One single function for the Final State Machine of the elevator
func FSM(
	id 							string,
	ch_arrivalFloor 			<-chan int,
	ch_localRequests 			<-chan [config.N_FLOORS][config.N_BUTTONS]bool,
	ch_doorObstruction 			<-chan bool,
	ch_stopButton				<-chan bool,
	ch_completedRequests 		chan<- elevator_io.ButtonEvent,
	ch_elevatorStateToAssigner 	chan<- map[string]elevator.ElevatorState,
	ch_elevatorStateToNetwork 	chan<- elevator.ElevatorState,
	ch_FSMLifeLine 				chan<- int,
) {

	// Initializing
	fmt.Printf("***** INITIALIZING ELEVATOR *****\n")
	localElevator := elevator.UninitializedElevator()
	prevLocalElevator := localElevator
	prevObstruction := false

	// If elevator is between floors, run it downwards until a floor is reached.
	elevator_io.SetMotorDirection(elevator_io.MD_Down)
	newFloor := <-ch_arrivalFloor
	elevator_io.SetMotorDirection(elevator_io.MD_Stop)
	localElevator.Floor = newFloor
	elevator_io.SetFloorIndicator(localElevator.Floor)

	// Initialize door
	elevator_io.SetDoorOpenLamp(false)
	doorTimer := time.NewTimer(time.Duration(config.DoorOpenDurationSec) * time.Second)

	elevator.SendLocalElevatorState(id, localElevator, ch_elevatorStateToAssigner, ch_elevatorStateToNetwork)

	// "For-Select" to supervise the different channels/events that changes the FSM
	for {
		ch_FSMLifeLine <- 1
		select {
		case localRequests := <-ch_localRequests:
			localElevator.Requests = localRequests

			switch localElevator.Behaviour {
			case elevator.EB_Moving:
				//NOP
			case elevator.EB_DoorOpen:
				if requests.Here(&localElevator) && (localElevator.Dirn == elevator_io.MD_Stop) {
					requests.ClearAtCurrentFloor(&localElevator, ch_completedRequests)
					doorTimer.Reset(time.Duration(config.DoorOpenDurationSec) * time.Second)
				}

			case elevator.EB_Idle:
				if requests.Here(&localElevator) && (localElevator.Dirn == elevator_io.MD_Stop) {
					requests.ClearAtCurrentFloor(&localElevator, ch_completedRequests)
					elevator_io.SetDoorOpenLamp(true)
					doorTimer.Reset(time.Duration(config.DoorOpenDurationSec) * time.Second)
					localElevator.Behaviour = elevator.EB_DoorOpen
				} else {
					// See if requests are elsewhere
					requests.ChooseDirnAndBehaviour(&localElevator)
					if localElevator.Behaviour == elevator.EB_Moving {
						elevator_io.SetMotorDirection(localElevator.Dirn)
					}
				}
			} //switch e.behaviour

		case newFloor := <-ch_arrivalFloor:
			localElevator.Floor = newFloor
			elevator_io.SetFloorIndicator(localElevator.Floor)

			switch localElevator.Behaviour {
			case elevator.EB_Moving:
				if requests.ShouldStop(&localElevator) {
					elevator_io.SetMotorDirection(elevator_io.MD_Stop)
					requests.ClearAtCurrentFloor(&localElevator, ch_completedRequests)
					elevator_io.SetDoorOpenLamp(true)
					doorTimer.Reset(time.Duration(config.DoorOpenDurationSec) * time.Second)
					localElevator.Behaviour = elevator.EB_DoorOpen
				}
			case elevator.EB_DoorOpen, elevator.EB_Idle:
				//NOP
			}

		// This channel automatically "transmits" when the timer times out.
		case <-doorTimer.C:
			switch localElevator.Behaviour {
			case elevator.EB_Idle, elevator.EB_Moving:
				//NOP
			case elevator.EB_DoorOpen:
				// This "if" happens when hallButton in direction was cleared at new floor,
				// 	but should wait longer if hallButton in other direction should also be cleared:
				if requests.Here(&localElevator) {
					requests.AnnounceDirectionChange(&localElevator)
					requests.ClearAtCurrentFloor(&localElevator, ch_completedRequests)
					time.Sleep(time.Duration(config.DoorOpenDurationSec) * time.Second)
				}
				// Keeps the door open while obstruction is active
				for prevObstruction{
					ch_FSMLifeLine <- 1
					select {
					case prevObstruction = <-ch_doorObstruction:
							time.Sleep(time.Duration(config.DoorOpenDurationSec) * time.Second)
					default:
						//NOP
						}
				}
				elevator_io.SetDoorOpenLamp(false)

				//Decides further action:
				requests.ChooseDirnAndBehaviour(&localElevator)
				if localElevator.Behaviour == elevator.EB_Moving {
					elevator_io.SetMotorDirection(localElevator.Dirn)
				}
			}

		case obstruction := <-ch_doorObstruction:
			prevObstruction = obstruction

		case <-ch_stopButton:
			if localElevator.Behaviour == elevator.EB_Moving {
				elevator_io.SetMotorDirection(elevator_io.MD_Stop)
			}

			// Keeps the elevator and fsm stalled while stopButton is active
			stopButtonPressed := true
			for stopButtonPressed {
				ch_FSMLifeLine <- 1
				stopButtonPressed = false
				stopButtonPressed = <-ch_stopButton
			}

			// Makes sure the elevator keeps going when stopButton is no longer active.
			switch localElevator.Behaviour {
			case elevator.EB_Moving:
				elevator_io.SetMotorDirection(localElevator.Dirn)

			case elevator.EB_DoorOpen:
				doorTimer.Reset(time.Duration(config.DoorOpenDurationSec) * time.Second)
				localElevator.HoldDoorOpenIfObstruction(&prevObstruction, doorTimer, ch_doorObstruction)
			}
		default:
			// NOP
		} //select

		if prevLocalElevator != localElevator {
			prevLocalElevator = localElevator
			elevator.SendLocalElevatorState(id, localElevator, ch_elevatorStateToAssigner, ch_elevatorStateToNetwork)
			localElevator.Elevator_print()
		}
	} //For
} //FSM
