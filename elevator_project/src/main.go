package main

import (
	"driver/backup"
	"driver/config"
	"driver/elevator"
	"driver/elevator_io"
	"driver/fsm"
	"fmt"
)

type ElevMsg struct {
	ID           string
	HallRequests bool
	state        int
	Iter         int
}

func main() {
	/* Initialize elevator ID and port
	This section sets the elevators ID (anything) and port (of the running node/PC),
	which should be passed on in the command line using
	'go run main.go -id=any_id -port=port'
	*/
	var id string
	flag.StringVar(&id, "ID", "", "ID of this peer")

	if id == "" { // if no ID is given, use local IP address
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())

	}

	var port string
	flag.StringVar(&port, "port", "", "Port of this peer")
	flag.Parse()

	// Spawn backup
	backup.BackupProcess(id) //this halts the progression of the program while it is the backup
	fmt.Println("Primary started.")

	// Initialize local elevator
	elevator_io.Init("localhost:15657", config.N_FLOORS)

	// Assigner channels (Recieve updates on the ID's of of the peers that are alive on the network)
	ch_peerUpdate := make(chan peers.PeerUpdate)
	ch_peerTxEnable := make(chan bool)
	ch_msgOut := make(chan ElevMsg)
	ch_msgIn := make(chan ElevMsg)
	ch_localOrders := make(chan [config.N_FLOORS][config.N_BUTTONS]bool)
	ch_completedOrders := make(chan elevator_io.ButtonEvent)
	ch_hallRequests := make(chan [config.N_FLOORS][config.N_BUTTONS - 1]int)

	

	// Goroutines for sending and recieving messages
	go peers.Transmitter(config.GlobalPort, id, ch_peerTxEnable)
	go peers.Reciever(config.GlobalPort, ch_peerUpdate)
	go bcast.Transmitter(config.GlobalPort, ch_msgOut)
	go bcast.Reciever(config.GlobalPort, ch_msgIn)

	// Channels for local elevator
	ch_buttonPressed := make(chan elevator_io.ButtonEvent)
	ch_localOrders := make(chan [config.N_FLOORS][config.N_BUTTONS]bool)
	ch_doorObstruction := make(chan bool)
	ch_stopButton := make(chan bool)

	// Backup goroutine
	go backup.LoadBackupFromFile("status.txt", ch_buttonPressed)

	// Local elevator goroutines
	go elevator_io.PollButtons(ch_buttonPressed)
	go elevator_io.PollFloorSensor(ch_arrivalFloor)
	go elevator_io.PollObstructionSwitch(ch_doorObstruction)
	go elevator_io.PollStopButton(ch_stopButton)

	go fsm.Fsm(ch_arrivalFloor, ch_localOrders, ch_buttonPressed, ch_doorObstruction, ch_stopButton, ch_completedOrders)

	// Sending message
	go func() {
		ElevMsg := ElevMsg{"Hello from " + id, 0}
		for {
			ElevMsg.Iter++
			msgOut <- ElevMsg
			time.Sleep(100 * time.Millisecond())
		}
	}()

	go func() {
		for {
			event := <-ch_completedOrders
			fmt.Printf("Received event. Floor %d, Btn: %s\n", event.BtnFloor+1, elevator.ElevButtonToString(event.BtnType))
		}
	}()

	// Peer monitoring (for config/debug purposes)
	fmt.Println("Started")
	for {
		select {
		case p := <-ch_peerUpdate:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-helloRx:
			fmt.Printf("Received: %#v\n", a)
		}
	}
	for {
	}
}
