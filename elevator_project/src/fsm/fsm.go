package fsm

import (
	"fmt"
	"elevator"
	"elevio"
	"requests"
	"timer"
)

// AT: Synes dette (package.Datatype) er en bra lesbar syntax, men usikker på hvordan det egt funker
//     og det må være LIK syntax på alle modulene.
var (
	elevator     elevator.Elevator
	outputDevice elevatorio.OutputDevice
)

func init() {
	// AT: F.eks het denne egentlig elevator = elevator_unintialized().
	elevator = elevator.Uninitialized()

	// AT: Tror denne henger sammen med con_load-modulen)
	elevator.LoadConfiguration("elevator.con")

	outputDevice = elevatorio.GetOutputDevice()
}

func setAllLights(es *elevator.Elevator) {
	for floor := 0; floor < elevator.NFloors; floor++ {
		for btn := 0; btn < elevator.NButtons; btn++ {
			outputDevice.RequestButtonLight(floor, btn, es.Requests[floor][btn])
		}
	}
}

func FsmOnInitBetweenFloors() {
	outputDevice.MotorDirection(elevatorio.Down)
	elevator.Dirn = elevatorio.Down
	elevator.Behaviour = elevator.Moving
}

// AT: Er/Skal ButtonType være definert i elevator?
func FsmOnRequestButtonPress(btnFloor int, btnType elevator.ButtonType) {
	fmt.Printf("\n\n%s(%d, %s)\n", "FsmOnRequestButtonPress", btnFloor, btnType.ToString())
	elevator.Print()

	switch elevator.Behaviour {
	
	case elevator.DoorOpen:
		if requests.ShouldClearImmediately(elevator, btnFloor, btnType) {
			timer.Start(elevator.Config.DoorOpenDurationS)
		} else {
			elevator.Requests[btnFloor][btnType] = true
		} // AT: Go trenger visst ikke "breaks" i switch cases
	
	case elevator.Moving:
		elevator.Requests[btnFloor][btnType] = true

	case elevator.Idle:
		elevator.Requests[btnFloor][btnType] = true
		// AT: Her har chat smeltet sammen Idles "innvendige" fsm. Tror det bør endres..
		//     Vet heller ikker hvor "updateElevatorState" kommer fra
		if elevator.Behaviour == elevator.Idle {
			pair := requests.ChooseDirection(elevator)
			elevator.Dirn = pair.Dirn
			elevator.Behaviour = pair.Behaviour
			updateElevatorState(pair)
		}
	}

	setAllLights(&elevator)
	fmt.Println("\nNew state:")
	elevator.Print()
}

// AT: Dette var det som var "inne" i IDLE-staten i fsm-en her oppe^.
// 	   Så funksjonen der oppe kan vel erstattes med dette.
func updateElevatorState(pair requests.DirnBehaviourPair) {
	switch pair.Behaviour {
	case elevator.DoorOpen:
		outputDevice.DoorLight(true)
		timer.Start(elevator.Config.DoorOpenDurationS)
		elevator = requests.ClearAtCurrentFloor(elevator)
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
	if elevator.Behaviour == elevator.Moving && requests.ShouldStop(elevator) {
		outputDevice.MotorDirection(elevatorio.Stop)
		outputDevice.DoorLight(true)
		elevator = requests.ClearAtCurrentFloor(elevator)
		timer.Start(elevator.Config.DoorOpenDurationS)
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
		pair := requests.ChooseDirection(elevator)
		elevator.Dirn = pair.Dirn
		elevator.Behaviour = pair.Behaviour
		updateElevatorState(pair)
	}

	fmt.Println("\nNew state:")
	elevator.Print()
}
