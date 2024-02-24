package main

import (
	"driver/elevator_io"
	"driver/elevator_io_types"
	"driver/fsm"
)

func main() {
	// Initialize system
	elevator_io.Init("localhost:15657", elevator_io_types.N_FLOORS)

	// GOROUTINES network

	// Channel for recieving updates on the id's of alive peers
	peerUpdateCh:= make(chan peers.PeerUpdate)

	// Channels for sending and recieving
	

	/* Fra utdelt

	ch_buttons := make(chan elevator_io.ButtonEvent)
	ch_floors := make(chan int)
	ch_obstr := make(chan bool)
	ch_stop := make(chan bool)

	go elevator_io.PollButtons(ch_buttons)
	go elevator_io.PollFloorSensor(ch_floors)
	go elevator_io.PollObstructionSwitch(ch_obstr)
	go elevator_io.PollStopButton(ch_stop)
	*/

}
