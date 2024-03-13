package main

import (
	"driver/assigner"
	"driver/backup"
	"driver/config"
	"driver/elevator"
	"driver/elevator_io"
	"driver/fsm"
	"driver/heartbeat"
	"driver/network/bcast"
	"driver/network/peers"
	"fmt"
)

func main() {
	fmt.Printf("\n\n************* NEW RUN *************\n\n")

	// Initialize elevator ID and port from command line: 'go run main.go -id=any_id -port=server_port'
	id, port := config.InitializeConfig()

	// Start backup process, halts the progression of the program while it is the backup
	backup.BackupProcess(id, port)
	fmt.Println("Primary started.")

	// Initialize local elevator
	elevator_io.Init("localhost:"+port, config.N_FLOORS)
	fmt.Println("\n--- Initialized local elevator " + id + " with port " + port + " ---\n")

	// Request assigner channels (Recieve updates on the ID's of of the peers that are alive on the network)
	ch_peerUpdate := make(chan peers.PeerUpdate, 100)
	ch_peerTxEnable := make(chan bool, 100)
	ch_msgOut := make(chan heartbeat.HeartBeat, 100)
	ch_msgIn := make(chan heartbeat.HeartBeat, 100)
	ch_completedRequests := make(chan elevator_io.ButtonEvent, 100)
	ch_hallRequestsIn := make(chan [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType, 100)
	ch_hallRequestsOut := make(chan [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType, 100)
	ch_externalElevators := make(chan []byte, 100)

	// Channels for local elevator
	ch_arrivalFloor := make(chan int, 100)
	ch_buttonPressed := make(chan elevator_io.ButtonEvent, 100)
	ch_localRequests := make(chan [config.N_FLOORS][config.N_BUTTONS]bool, 100)
	ch_doorObstruction := make(chan bool, 100)
	ch_stopButton := make(chan bool, 100)
	ch_elevatorStateToAssigner := make(chan map[string]elevator.ElevatorState, 5)
	ch_elevatorStateToNetwork := make(chan elevator.ElevatorState, 5)

	go backup.LoadBackupFromFile("backup.txt", ch_buttonPressed)
	
	// Goroutines for sending and recieving messages
	go bcast.Transmitter(config.DefaultPortBcast, ch_msgOut)
	go bcast.Receiver(config.DefaultPortBcast, ch_msgIn)

	go peers.Transmitter(config.DefaultPortPeer, id, ch_peerTxEnable)
	go peers.Receiver(config.DefaultPortPeer, ch_peerUpdate)

	// elevator_io goroutines
	go elevator_io.PollButtons(ch_buttonPressed)
	go elevator_io.PollFloorSensor(ch_arrivalFloor)
	go elevator_io.PollObstructionSwitch(ch_doorObstruction)
	go elevator_io.PollStopButton(ch_stopButton)

	// Finite state machine goroutine
	go fsm.Fsm(
		id,
		ch_arrivalFloor,
		ch_localRequests,
		ch_doorObstruction,
		ch_stopButton,
		ch_completedRequests,
		ch_elevatorStateToAssigner,
		ch_elevatorStateToNetwork,
	)

	// Assigner goroutine
	go assigner.RequestAssigner(
		id,
		ch_buttonPressed,
		ch_completedRequests,
		ch_elevatorStateToAssigner,
		ch_hallRequestsIn,
		ch_externalElevators,
		ch_hallRequestsOut,
		ch_localRequests,
	)

	go backup.ReportPrimaryAlive(id)

	go heartbeat.Send(
		id,
		ch_hallRequestsOut,
		ch_elevatorStateToNetwork,
		ch_msgOut,
	)

	go peers.Update(
		id,
		ch_peerUpdate,
		ch_msgIn,
		ch_hallRequestsIn,
		ch_externalElevators,
	)

	select {}

}
