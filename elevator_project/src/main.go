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
	SenderID     string
	HallRequests [config.N_FLOORS][config.N_BUTTONS - 1]int
	State        elevator.ElevatorState
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
	//ch_BackupHeartbeat := make(chan string, 100)
	ch_peerUpdate := make(chan peers.PeerUpdate, 100)
	ch_peerTxEnable := make(chan bool, 100)
	ch_msgOut := make(chan HeartBeat, 100)
	ch_msgIn := make(chan HeartBeat, 100)
	ch_completedRequests := make(chan elevator_io.ButtonEvent, 100)
	ch_hallRequestsIn := make(chan [config.N_FLOORS][config.N_BUTTONS - 1]int, 100)
	ch_hallRequestsOut := make(chan [config.N_FLOORS][config.N_BUTTONS - 1]int, 100)
	ch_externalElevators := make(chan map[string]elevator.ElevatorState, 100)

	// Channels for local elevator
	ch_arrivalFloor := make(chan int, 100)
	ch_buttonPressed := make(chan elevator_io.ButtonEvent, 100)
	ch_localRequests := make(chan [config.N_FLOORS][config.N_BUTTONS]bool, 100)
	ch_doorObstruction := make(chan bool, 100)
	ch_stopButton := make(chan bool, 100)
	ch_elevatorStateToAssigner := make(chan map[string]elevator.ElevatorState, 100)
	ch_elevatorStateToNetwork := make(chan elevator.ElevatorState, 100)

	// Goroutines for sending and recieving messages
	//go bcast.Transmitter(config.DefaultPortBackup, ch_BackupHeartbeat)
	go bcast.Transmitter(config.DefaultPortBcast, ch_msgOut)
	go bcast.Receiver(config.DefaultPortBcast, ch_msgIn)

	go peers.Transmitter(config.DefaultPortPeer, id, ch_peerTxEnable)
	go peers.Receiver(config.DefaultPortPeer, ch_peerUpdate)

	// Backup goroutine
	go backup.LoadBackupFromFile("backup.txt", ch_buttonPressed)

	// elevator_io goroutines
	go elevator_io.PollButtons(ch_buttonPressed)
	go elevator_io.PollFloorSensor(ch_arrivalFloor)
	go elevator_io.PollObstructionSwitch(ch_doorObstruction)
	go elevator_io.PollStopButton(ch_stopButton)

	// Finite state machine goroutine
	go fsm.Fsm(
		ch_arrivalFloor,
		ch_localRequests,
		ch_doorObstruction,
		ch_stopButton,
		ch_completedRequests,
		ch_elevatorStateToAssigner,
		ch_elevatorStateToNetwork,
	)

	// Assigner goroutine
	go assigner.Assigner(
		ch_buttonPressed,
		ch_completedRequests,
		ch_localRequests,
		ch_hallRequestsIn,
		ch_hallRequestsOut,
		ch_elevatorStateToAssigner,
		ch_externalElevators,
	)

	// Send heartbeat to backup
	go backup.PrimaryProcess(id)

	// Send heartbeat to network incl. all info
	go func() {
		var HallRequests [config.N_FLOORS][config.N_BUTTONS - 1]int = <-ch_hallRequestsOut
		var State elevator.ElevatorState = <-ch_elevatorStateToNetwork

		for {
			select {
			case newHallRequests := <-ch_hallRequestsOut:
				HallRequests = newHallRequests
			case newState := <-ch_elevatorStateToNetwork:
				State = newState
			default:
				// NOP
			}
			HeartBeat := HeartBeat{
				SenderID:     id,
				HallRequests: HallRequests,
				State:        State,
			}
			ch_msgOut <- HeartBeat
			time.Sleep(100 * time.Millisecond)
			// fmt.Printf("\n Heartbeat sent\n")
		}
	}()

	fmt.Println("Started")

	AlivePeers := make(map[string]elevator.ElevatorState)
	for {
		select {
		case p := <-ch_peerUpdate:

			for _, peer := range p.Lost {
				if _, ok := AlivePeers[peer]; ok {
					delete(AlivePeers, peer)
				}
			}

			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-ch_msgIn:
			AlivePeers[a.SenderID] = a.State

			ch_hallRequestsIn <- a.HallRequests
			ch_externalElevators <- AlivePeers

			//fmt.Printf("Received: %#v\n", a)

		default:
			// NOP
		}

	}

}
