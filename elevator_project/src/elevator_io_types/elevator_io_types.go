package elevator_io_types

// Kan nok flyttes til config (bare at vi da må endre overalt ellers fra "elevator_io_types" til "config")
const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

type Dirn int

const (
	D_Down Dirn = -1
	D_Stop Dirn = 0
	D_Up   Dirn = 1
)

type Button int

const (
	B_HallUp   Button = 0
	B_HallDown Button = 1
	B_Cab      Button = 2
)

// Kan kanskje fjernes
type ElevInputDevice struct {
	FloorSensor   func() bool
	RequestButton func(int, Button) int
	StopButton    func() bool
	Obstruction   func() bool
}

// Kan kanskje fjernes
type ElevOutputDevice struct {
	FloorIndicator     func(int)
	RequestButtonLight func(int, Button, bool)
	DoorLight          func(bool)
	StopButtonLight    func(bool)
	MotorDirection     func(Dirn)
}

func Elevio_dirn_toString(d Dirn) string {
	switch d {
	case D_Down:
		return "Down"
	case D_Stop:
		return "Stop"
	case D_Up:
		return "Up"
	default:
		return "Unknown"
	}
}

func Elevio_button_toString(b Button) string {
	switch b {
	case B_HallUp:
		return "HallUp"
	case B_HallDown:
		return "HallDown"
	case B_Cab:
		return "Cab"
	default:
		return "Unknown"
	}
}
