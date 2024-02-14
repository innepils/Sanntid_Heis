package main

import (
	"elevio"
	"fmt"
	"fsm"
	"time"
)

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

func main() {
	
	fmt.Println("Started!")

	inputPollRate_ms := 25
	con_load.Load("elevator.con",
		con_load.Val("inputPollRate_ms", &inputPollRate_ms))

	input := elevator_io_device.GetInputDevice()

	if input.FloorSensor() == -1 {
		fsm.OnInitBetweenFloors()
	}

	for {
		// Request button
		prev := make([][]int, N_FLOORS)
		for i := range prev {
			prev[i] = make([]int, N_BUTTONS)
		}
		for f := 0; f < N_FLOORS; f++ {
			for b := 0; b < N_BUTTONS; b++ {
				v := input.RequestButton(f, b)
				if v != 0 && v != prev[f][b] {
					fsm.OnRequestButtonPress(f, b)
				}
				prev[f][b] = v
			}
		}

		// Floor sensor
		prevFloor := -1
		f := input.FloorSensor()
		if f != -1 && f != prevFloor {
			fsm.OnFloorArrival(f)
		}
		prevFloor = f

		// Timer
		if timer.TimedOut() {
			timer.Stop()
			fsm.OnDoorTimeout()
		}

		time.Sleep(time.Duration(inputPollRate_ms) * time.Millisecond)
	}
}
