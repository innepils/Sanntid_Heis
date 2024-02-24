package main

import (
	"driver/cost"
	"driver/elevator"
	"driver/elevator_io"
	"driver/elevator_io_types"
	"fmt"
)

func main() {
	
	numFloors := 4

	elevator_io.Init("localhost:20007", numFloors)

	//var d elevator_io.MotorDirection = elevator_io.MD_Up
	//elevator_io.SetMotorDirection(d)

	drv_buttons := make(chan elevator_io.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevator_io.PollButtons(drv_buttons)
	go elevator_io.PollFloorSensor(drv_floors)
	go elevator_io.PollObstructionSwitch(drv_obstr)
	go elevator_io.PollStopButton(drv_stop)

}
