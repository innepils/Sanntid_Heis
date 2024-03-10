package main

import (
	"driver/assigner"
	"driver/backup"
	"driver/config"
	"driver/elevator"
	"driver/elevator_io"
	"driver/fsm"
	"driver/network/bcast"
	"driver/network/peers"
	"fmt"
	"time"
)

type HeartBeat struct {
	ID           string
	HallRequests [config.N_FLOORS][config.N_BUTTONS - 1]int
	state        map[string]elevator.ElevatorState
	Iter         int
}

func main() {
	fmt.Printf("\n\n\n\n")
	fmt.Printf("************* NEW RUN *************")
	fmt.Printf("\n\n\n\n")

	/* Initialize elevator ID and port from command line:
	   'go run main.go -id=any_id -port=server_port' */
	id, port := config.InitializeConfig()

	// Spawn backup
	backup.BackupProcess(id, port) //this halts the progression of the program while it is the backup
	fmt.Println("Primary started.")

	// Initialize local elevator
	elevator_io.Init("localhost:"+port, config.N_FLOORS)
	fmt.Println("\n--- Initialized elevator " + id + " with port " + port + " ---\n")

	// Assigner channels (Recieve updates on the ID's of of the peers that are alive on the network)
	ch_BackupHeartbeat := make(chan string, 100)
	ch_peerUpdate := make(chan peers.PeerUpdate, 100)
	ch_peerTxEnable := make(chan bool, 100)
	ch_msgOut := make(chan HeartBeat, 100)
	ch_msgIn := make(chan HeartBeat, 100)
	ch_completedOrders := make(chan elevator_io.ButtonEvent, 100)
	ch_hallRequestsIn := make(chan [config.N_FLOORS][config.N_BUTTONS - 1]int, 100)
	ch_hallRequestsOut := make(chan [config.N_FLOORS][config.N_BUTTONS - 1]int)
	ch_externalElevators := make(chan map[string]elevator.ElevatorState, 100)

	// Channels for local elevator
	ch_arrivalFloor := make(chan int, 100)
	ch_buttonPressed := make(chan elevator_io.ButtonEvent, 100)
	ch_localOrders := make(chan [config.N_FLOORS][config.N_BUTTONS]bool, 100)
	ch_doorObstruction := make(chan bool, 100)
	ch_stopButton := make(chan bool, 100)
	ch_elevatorStateToAssigner := make(chan map[string]elevator.ElevatorState, 100)
	ch_elevatorStateToNetwork := make(chan map[string]elevator.ElevatorState, 100)

	// Goroutines for sending and recieving messages
	go bcast.Transmitter(config.DefaultPortBackup, ch_BackupHeartbeat)
	go bcast.Transmitter(config.DefaultPortBcast, ch_msgOut)
	go bcast.Receiver(config.DefaultPortBcast, ch_msgIn)

	go peers.Transmitter(config.DefaultPortPeer, id, ch_peerTxEnable)
	go peers.Receiver(config.DefaultPortPeer, ch_peerUpdate)

	// Backup goroutine
	go backup.LoadBackupFromFile("status.txt", ch_buttonPressed)

	// elevator_io goroutines
	go elevator_io.PollButtons(ch_buttonPressed)
	go elevator_io.PollFloorSensor(ch_arrivalFloor)
	go elevator_io.PollObstructionSwitch(ch_doorObstruction)
	go elevator_io.PollStopButton(ch_stopButton)

	// Finite state machine goroutine
	go fsm.Fsm(
		ch_arrivalFloor,
		ch_localOrders,
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
		ch_hallRequestsOut,
		ch_elevatorStateToAssigner,
		ch_externalElevators,
	)

	// Send heartbeat to network incl. all info
	go func() {
		HeartBeat := id
		for {
			ch_BackupHeartbeat <- HeartBeat
			time.Sleep(100 * time.Millisecond)
		}
	}()
	// Send heartbeat to network incl. all info
	go func() {
		for {
			HeartBeat := HeartBeat{
				ID:           id,
				HallRequests: <-ch_hallRequestsOut,
				state:        <-ch_elevatorStateToNetwork,
				Iter:         0,
			}
			HeartBeat.Iter++
			ch_msgOut <- HeartBeat
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// go func() {
	// 	for {
	// 		select {
	// 		//case event := <-ch_completedOrders:
	// 		//	fmt.Printf("Received event. Floor %d, Btn: %s\n", event.BtnFloor+1, elevator.ElevButtonToString(event.BtnType))

	// 		case <-ch_elevatorStateToNetwork:
	// 			//fmt.Printf("Received event from elevatorStateToNetWork\n")

	// 			//case <-ch_elevatorStateToAssigner:
	// 			//fmt.Printf("Received event from elevatorStateToAssigner\n")
	// 		}
	// 	}
	// }()

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
