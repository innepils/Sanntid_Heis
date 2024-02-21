package fsm

import (
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
	localElevator.LoadConfiguration("elevator.con")

	// AT(21.feb): GetOuputDevice eksisterir elevator_io.c men ikke go-filen vi fikk
	outputDevice = elevator_io.GetOutputDevice()
}

func setAllLights(es *elevator.Elevator) {
	for floor := 0; floor < elevator.NFloors; floor++ {
		for btn := 0; btn < elevator.NButtons; btn++ {
			outputDevice.RequestButtonLight(floor, btn, es.Requests[floor][btn])
		}
	}
}

func FsmOnInitBetweenFloors() {
	outputDevice.MotorDirection(elevator_io_types.D_Down)
	elevator.Dirn = elevator_io_types.D_Down
	elevator.Behaviour = elevator.Moving
}

// AT: Er/Skal ButtonType være definert i elevator?
func FsmOnRequestButtonPress(btnFloor int, btnType elevator.ButtonType) {
	fmt.Printf("\n\n%s(%d, %s)\n", "FsmOnRequestButtonPress", btnFloor, btnType.ToString())
	elevator.Print()

	switch elevator.Behaviour {

	case elevator.DoorOpen:
		if Requests.Requests_shouldClearImmediately(elevator, btnFloor, btnType) {
			timer.TimerStart(elevator.Config.DoorOpenDurationS)
		} else {
			elevator.Requests[btnFloor][btnType] = true
		} // AT: Go trenger visst ikke "breaks" i switch cases

	case elevator.Moving:
		elevator.Requests[btnFloor][btnType] = true

	case elevator.Idle:
		elevator.Requests[btnFloor][btnType] = true
		// AT: Her har chat smeltet sammen Idles "innvendige" fsm. Tror det bør endres..
		if elevator.Behaviour == elevator.Idle {
			pair := Requests.Requests_chooseDirection(localElevator)
			elevator.Dirn = pair.Dirn
			elevator.Behaviour = pair.Behaviour
			updateElevatorState(pair)
		}
	}

	setAllLights(&localElevator)
	fmt.Println("\nNew state:")
	elevator.Print()
}

// AT: Dette var det som var "inne" i IDLE-staten i fsm-en her oppe^.
//
//	Så funksjonen der oppe kan vel erstattes med dette.
func updateElevatorState(pair Requests.DirnBehaviourPair) {
	switch pair.Behaviour {
	case elevator.DoorOpen:
		outputDevice.DoorLight(true)
		timer.TimerStart(elevator.Config.DoorOpenDurationS)
		elevator = Requests.Requests_clearAtCurrentFloor(elevator)
	case elevator.Moving:
		outputDevice.MotorDirection(elevator.Dirn)
	case elevator.Idle:
		// No action needed for Idle here
	}
}

func FsmOnFloorArrival(newFloor int) {
	fmt.Printf("\n\n%s(%d)\n", "FsmOnFloorArrival", newFloor)
	elevator.Print()

	elevator.Floor = newFloor

	outputDevice.FloorIndicator(newFloor)

	// AT: Dette var en enkel switch case i C. Virkemåten er identisk men
	// 	   switch-casen fører til bedre kodkvalitet/mer lesbart.
	if elevator.Behaviour == elevator.Moving && Requests.Requests_shouldStop(elevator) {
		outputDevice.MotorDirection(elevator_io_types.D_Stop)
		outputDevice.DoorLight(true)
		elevator = Requests.Requests_clearAtCurrentFloor(elevator)
		timer.TimerStart(elevator.Config.DoorOpenDurationS)
		setAllLights(&elevator)
		elevator.Behaviour = elevator.DoorOpen
	}

	fmt.Println("\nNew state:")
	elevator.Print()
}

func FsmOnDoorTimeout() {
	fmt.Println("\n\nFsmOnDoorTimeout()")
	elevator.Print()

	// AT: Dette er også VELDIG ulikt. Her var det i c en switch case inni en switch case (lesbart)
	//     Men her brukes altså BARE den ytre switch-casen (som en if)
	//     hvor den indre er "updateElevatorState" igjen.
	if elevator.Behaviour == elevator.DoorOpen {
		pair := Requests.Requests_chooseDirection(elevator)
		elevator.Dirn = pair.Dirn
		elevator.Behaviour = pair.Behaviour
		updateElevatorState(pair)
	}

	fmt.Println("\nNew state:")
	elevator.Print()
}
