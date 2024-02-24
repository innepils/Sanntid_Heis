package main

import (
	"driver/cost"
	"driver/elevator"
	"driver/elevator_io"
	"driver/elevator_io_types"
	"driver/elevator_io_types"
	"driver/fsm"
	"driver/network/bcast"
	"driver/network/conn"
	"driver/network/peers"
	"flag"
	"fmt"
	"os"
	"time"

)

func main() {

	/* Initialize elevator ID and port
	This section sets the elevators ID (anything) and port (of the running node/PC),
	which should be passed on in the command line using 
	'go run main.go -id=any_id -port=port'
	*/
	var id string
	var port string

	flag.StringVar(&id, "ID", "", "ID of this peer")
	flag.StringVar(%port, "Port", "", "port of this peer")
	flag.Parse()

	// if no ID is given, use local IP address
	// (legger også til process ID 'os.Getpid()', ikke helt sikker på hvorfor enda)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// if no port is given
	// if port == "" {
		
	// }

	// Initialize local elevator
	elevator_io.Init("localhost:" + port, elevator_io_types.N_FLOORS)

	// GOROUTINES network
	

	// Channel for recieving updates on the ID's of alive peers
	peerUpdateCh:= make(chan peers.PeerUpdate)
	
	/* Channel for enabling/disabling the transmitter after start.
	Can be used to signal that the node is "unavailable". */
	peerTxEnable := make(chan bool)

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
