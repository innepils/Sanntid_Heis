package main

import (
	"driver/elevator_io"
	"driver/fsm"
	"fmt"
)

func main() {

	numFloors := 4
	// comment
	elevator_io.Init("localhost:15657", numFloors)

	fsm.Init()

	var d elevator_io.MotorDirection = elevator_io.MD_Up
	//elevator_io.SetMotorDirection(d)

	drv_buttons := make(chan elevator_io.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevator_io.PollButtons(drv_buttons)
	go elevator_io.PollFloorSensor(drv_floors)
	go elevator_io.PollObstructionSwitch(drv_obstr)
	go elevator_io.PollStopButton(drv_stop)

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			elevator_io.SetButtonLamp(a.Button, a.Floor, true)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			if a == numFloors-1 {
				d = elevator_io.MD_Down
			} else if a == 0 {
				d = elevator_io.MD_Up
			}
			elevator_io.SetMotorDirection(d)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevator_io.SetMotorDirection(elevator_io.MD_Stop)
			} else {
				elevator_io.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevator_io.ButtonType(0); b < 3; b++ {
					elevator_io.SetButtonLamp(b, f, false)
				}
			}
		}
	}
}
