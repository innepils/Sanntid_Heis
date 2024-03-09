package main

import (
	"driver/assigner"
	"driver/backup"
	"driver/config"
	"driver/elevator"
	"driver/elevator_io"
	"driver/fsm"
	"driver/network/bcast"
	"driver/network/localip"
	"driver/network/peers"
	"flag"
	"fmt"
	"os"
	"time"
)

type HeartBeat struct {
	ID           string
	HallRequests [config.N_FLOORS][config.N_BUTTONS - 1]int
	state        map[string]elevator.ElevatorState
	Iter         int
}

func main() {
	/* Initialize elevator ID and port
	This section sets the elevators ID and port
	which should be passed on in the command line using
	'go run main.go -id=any_id -port=server_port'
	*/
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")

	if id == "" { // if no ID is given, use local IP address and process ID
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer_%s:%d", localIP, os.Getpid())
	}

	var port string
	flag.StringVar(&port, "port", "", "port of this peer")

	if port == "" { // if no port is given, use default port (15657)
		port = fmt.Sprint(15657)
	}
	flag.Parse()

	// Spawn backup
	backup.BackupProcess(id) //this halts the progression of the program while it is the backup
	fmt.Println("Primary started.")

	// Initialize local elevator
	elevator_io.Init("localhost:"+port, config.N_FLOORS)
	fmt.Println("\n--- Initialized elevator " + id + " with port " + port + " ---\n")

	// Assigner channels (Recieve updates on the ID's of of the peers that are alive on the network)
	ch_peerUpdate := make(chan peers.PeerUpdate, 100)
	ch_peerTxEnable := make(chan bool, 100)
	ch_msgOut := make(chan HeartBeat, 100)
	ch_msgIn := make(chan HeartBeat, 100)
	ch_completedOrders := make(chan elevator_io.ButtonEvent, 100)
	ch_hallRequestsIn := make(chan [config.N_FLOORS][config.N_BUTTONS - 1]int, 100)
	//ch_hallRequestsOut := make(chan [config.N_FLOORS][config.N_BUTTONS - 1]int)
	ch_externalElevators := make(chan map[string]elevator.ElevatorState, 100)

	// Goroutines for sending and recieving messages
	go bcast.Transmitter(config.DefaultPortBcast, ch_msgOut)
	go bcast.Receiver(config.DefaultPortBcast, ch_msgIn)

	go peers.Transmitter(config.DefaultPortPeer, id, ch_peerTxEnable)
	go peers.Receiver(config.DefaultPortPeer, ch_peerUpdate)

	// Channels for local elevator
	ch_arrivalFloor := make(chan int, 100)
	ch_buttonPressed := make(chan elevator_io.ButtonEvent, 100)
	ch_localOrders := make(chan [config.N_FLOORS][config.N_BUTTONS]bool, 100)
	ch_doorObstruction := make(chan bool, 100)
	ch_stopButton := make(chan bool, 100)
	ch_elevatorStateToAssigner := make(chan map[string]elevator.ElevatorState, 100)
	ch_elevatorStateToNetwork := make(chan map[string]elevator.ElevatorState, 100)
	//fmt.Printf("completed order-channel received in assign")

	// Backup goroutine
	go backup.LoadBackupFromFile("status.txt", ch_buttonPressed)

	// Local elevator goroutines
	go elevator_io.PollButtons(ch_buttonPressed)
	go elevator_io.PollFloorSensor(ch_arrivalFloor)
	go elevator_io.PollObstructionSwitch(ch_doorObstruction)
	go elevator_io.PollStopButton(ch_stopButton)

	// Finite state machine goroutine
	go fsm.Fsm(
		ch_arrivalFloor,
		ch_localOrders,
		ch_buttonPressed,
		ch_doorObstruction,
		ch_stopButton,
		ch_completedOrders,
		ch_elevatorStateToAssigner,
		ch_elevatorStateToNetwork,
	)

	// Assigner goroutine
	go assigner.Assigner(
		ch_buttonPressed,
		ch_completedOrders,
		ch_localOrders,
		ch_hallRequestsIn,
		ch_elevatorStateToAssigner,
		ch_externalElevators,
	)

	// Send heartbeat incl. all info
	go func() {
		HeartBeat := HeartBeat{"Hello from " + id, <-ch_hallRequestsIn, <-ch_elevatorStateToNetwork, 0}
		for {
			HeartBeat.Iter++
			ch_msgOut <- HeartBeat
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			select {
			//case event := <-ch_completedOrders:
			//	fmt.Printf("Received event. Floor %d, Btn: %s\n", event.BtnFloor+1, elevator.ElevButtonToString(event.BtnType))

			case <-ch_elevatorStateToNetwork:
				//fmt.Printf("Received event from elevatorStateToNetWork\n")

				//case <-ch_elevatorStateToAssigner:
				//fmt.Printf("Received event from elevatorStateToAssigner\n")
			}
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
		case a := <-ch_msgIn:
			fmt.Printf("Received: %#v\n", a)
		}
	}

	// select {}
}
