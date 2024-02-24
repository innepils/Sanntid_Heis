package main

import (
	"driver/elevator"
	"driver/elevator_io"
	"driver/elevator_io_types"
	"driver/fsm"
	"fmt"
)

func main() {
	// Initialize system
	elevator_io.Init("localhost:15657", elevator_io_types.N_FLOORS)

	if (elevator_io.GetFloor() == -1){
		fsm.FsmOnInitBetweenFloors()
	} else {
		fsm.Init()
	}


	ch_buttons := make(chan elevator_io.ButtonEvent)
	ch_floors := make(chan int)
	ch_obstr := make(chan bool)
	ch_stop := make(chan bool)

	// 
	go elevator_io.PollButtons(ch_buttons)
	go elevator_io.PollFloorSensor(ch_floors)
	go elevator_io.PollObstructionSwitch(ch_obstr)
	go elevator_io.PollStopButton(ch_stop)

	
}
