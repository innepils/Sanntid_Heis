package main

import (
	"driver/elevator_io"
	"driver/elevator_io_types"
	"driver/network/localip"
	"flag"
	"fmt"
)

type ElevatorMessage struct {
	ID           string
	HallRequests bool
	state        int
}

func main() {

	/* Initialize elevator ID and port
	This section sets the elevators ID (anything) and port (of the running node/PC),
	which should be passed on in the command line using
	'go run main.go -id=any_id -port=port'
	*/
	var id string
	flag.StringVar(&id, "ID", "", "ID of this peer")
	var port string
	flag.StringVar(&port, "Port", "", "port of this peer")
	flag.Parse()

	// if no ID is given, use local IP address
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s", localIP)
	}

	// Initialize local elevator
	elevator_io.Init("localhost:"+port, elevator_io_types.N_FLOORS)

	// GOROUTINES network

	// Channel for recieving updates on the ID's of alive peers
	// peerUpdateCh := make(chan peers.PeerUpdate)

	/* Channel for enabling/disabling the transmitter after start.
	Can be used to signal that the node is "unavailable". */
	// peerTxEnable := make(chan bool)
	
	// Channels for sending and recieving

	ch_buttonPressed := make(chan elevator_io.ButtonEvent)
	ch_arrivalFloor := make(chan int)
	ch_doorTimedOut := make(chan bool)
	ch_doorObstruction := make(chan bool)
	ch_stopButton := make(chan bool)

	go elevator_io.PollButtons(ch_buttonPressed)
	go elevator_io.PollFloorSensor(ch_arrivalFloor)
	go elevator_io.PollFloorSensor(ch_doorTimedOut)
	go elevator_io.PollObstructionSwitch(ch_doorObstruction)
	go elevator_io.PollStopButton(ch_stopButton)
	

}
