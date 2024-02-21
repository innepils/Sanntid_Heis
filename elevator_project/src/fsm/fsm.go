package fsm

import (
	"driver/con_load"
	"driver/elevator"
	"driver/elevator_io_types"
	Requests "driver/requests"
	"driver/timer"
	"fmt"
)

// AT: Synes dette (package.Datatype) er en bra lesbar syntax, men usikker på hvordan det egt funker
//
//	og det må være LIK syntax på alle modulene.
var (
	localElevator elevator.Elevator
	outputDevice  elevator_io_types.ElevOutputDevice
)

func init() {
	// AT: F.eks het denne egentlig elevator = elevator_unintialized().
	localElevator = elevator.UninitializedElevator()

	// AT: Tror denne henger sammen med con_load-modulen)
	con_load.LoadConfig("elevator.con")

	// AT(21.feb): GetOuputDevice eksisterir elevator_io.c men ikke go-filen vi fikk
	// outputDevice = elevator_io.GetOutputDevice()
}

func setAllLights(es *elevator.Elevator) {
	for floor := 0; floor < elevator_io_types.N_FLOORS; floor++ {
		for btn := 0; btn < elevator_io_types.N_BUTTONS; btn++ {
			btnType := elevator_io_types.Button(btn)
			outputDevice.RequestButtonLight(floor, btnType, es.Requests[floor][btn])
		}
	}
}

func FsmOnInitBetweenFloors() {
	outputDevice.MotorDirection(elevator_io_types.D_Down)
	localElevator.Dirn = elevator_io_types.D_Down
	localElevator.Behaviour = elevator.EB_Moving
}

// AT: Er/Skal ButtonType være definert i elevator?
func FsmOnRequestButtonPress(btnFloor int, btnType elevator_io_types.Button) {
	fmt.Printf("\n\n%s(%d, %s)\n", "FsmOnRequestButtonPress", btnFloor, elevator_io_types.Elevio_button_toString(btnType))
	localElevator.Print()

	switch localElevator.Behaviour {

	case elevator.EB_DoorOpen:
		if Requests.Requests_shouldClearImmediately(localElevator, btnFloor, btnType) {
			timer.TimerStart(localElevator.Config.DoorOpenDurationSec)
		} else {
			localElevator.Requests[btnFloor][btnType] = true
		} // AT: Go trenger visst ikke "breaks" i switch cases

	case elevator.EB_Moving:
		localElevator.Requests[btnFloor][btnType] = true

	case elevator.EB_Idle:
		localElevator.Requests[btnFloor][btnType] = true
		// AT: Her har chat smeltet sammen Idles "innvendige" fsm. Tror det bør endres..
		if localElevator.Behaviour == elevator.EB_Idle {
			pair := Requests.Requests_chooseDirection(localElevator)
			localElevator.Dirn = pair.Dirn
			localElevator.Behaviour = pair.Behaviour
			updateElevatorState(pair)
		}
	}

	setAllLights(&localElevator)
	fmt.Println("\nNew state:")
	localElevator.Print()
}

// AT: Dette var det som var "inne" i IDLE-staten i fsm-en her oppe^.
//
//	Så funksjonen der oppe kan vel erstattes med dette.
func updateElevatorState(pair Requests.DirnBehaviourPair) {
	switch pair.Behaviour {
	case elevator.EB_DoorOpen:
		outputDevice.DoorLight(true)
		timer.TimerStart(localElevator.Config.DoorOpenDurationSec)
		localElevator = Requests.Requests_clearAtCurrentFloor(localElevator)
	case elevator.EB_Moving:
		outputDevice.MotorDirection(localElevator.Dirn)
	case elevator.EB_Idle:
		// No action needed for Idle here
	}
}

func FsmOnFloorArrival(newFloor int) {
	fmt.Printf("\n\n%s(%d)\n", "FsmOnFloorArrival", newFloor)
	localElevator.Print()

	localElevator.Floor = newFloor

	outputDevice.FloorIndicator(newFloor)

	// AT: Dette var en enkel switch case i C. Virkemåten er identisk men
	// 	   switch-casen fører til bedre kodkvalitet/mer lesbart.
	if localElevator.Behaviour == elevator.EB_Moving && Requests.Requests_shouldStop(localElevator) {
		outputDevice.MotorDirection(elevator_io_types.D_Stop)
		outputDevice.DoorLight(true)
		localElevator = Requests.Requests_clearAtCurrentFloor(localElevator)
		timer.TimerStart(localElevator.Config.DoorOpenDurationSec)
		setAllLights(&localElevator)
		localElevator.Behaviour = elevator.EB_DoorOpen
	}

	fmt.Println("\nNew state:")
	localElevator.Print()
}

func FsmOnDoorTimeout() {
	fmt.Println("\n\nFsmOnDoorTimeout()")
	localElevator.Print()

	// AT: Dette er også VELDIG ulikt. Her var det i c en switch case inni en switch case (lesbart)
	//     Men her brukes altså BARE den ytre switch-casen (som en if)
	//     hvor den indre er "updateElevatorState" igjen.
	if localElevator.Behaviour == elevator.EB_DoorOpen {
		pair := Requests.Requests_chooseDirection(localElevator)
		localElevator.Dirn = pair.Dirn
		localElevator.Behaviour = pair.Behaviour
		updateElevatorState(pair)
	}

	fmt.Println("\nNew state:")
	localElevator.Print()
}
