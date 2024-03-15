// Package main implements the elevator project. Check README.md
// Author: Currently disclosed

package main

import (
	"src/assigner"
	"src/backup"
	"src/config"
	"src/deadlock"
	"src/elevator"
	"src/elevator_io"
	"src/fsm"
	"src/network/bcast"
	"src/network/heartbeat"
	"src/network/peers"
)

func main() {

	nodeID, port := config.InitializeIDandPort()

	elevator_io.Init("localhost:" + port, config.N_FLOORS)

	// Request assigner channels (Recieve updates on the ID's of of the peers that are alive on the network)
	ch_peerUpdate 				:= make(chan peers.PeerUpdate, 1)
	ch_peerTxEnable 			:= make(chan bool, 1)
	ch_msgOut 					:= make(chan heartbeat.HeartBeat, 1)
	ch_msgIn 					:= make(chan heartbeat.HeartBeat, 1)
	ch_completedRequests 		:= make(chan elevator_io.ButtonEvent, 1)
	ch_hallRequestsIn 			:= make(chan [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType, 1)
	ch_hallRequestsOut 			:= make(chan [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType, 1)
	ch_externalElevators 		:= make(chan []byte, 1)

	// Channels for local elevator
	ch_arrivalFloor 			:= make(chan int, 1)
	ch_buttonPressed 			:= make(chan elevator_io.ButtonEvent, 1)
	ch_localRequests 			:= make(chan [config.N_FLOORS][config.N_BUTTONS]bool, 1)
	ch_doorObstruction 			:= make(chan bool, 1)
	ch_stopButton 				:= make(chan bool, 1)
	ch_elevatorStateToAssigner 	:= make(chan map[string]elevator.HRAElevatorState, 1)
	ch_elevatorStateToNetwork 	:= make(chan elevator.HRAElevatorState, 1)

	// Channels for deadlock detector
	ch_FSMDeadlock 				:= make(chan string, 1)
	ch_assignerDeadlock 		:= make(chan string, 1)
	ch_heartbeatDeadlock 		:= make(chan string, 1)
	ch_peersDeadlock 			:= make(chan string, 1)

	go backup.LoadBackupFromFile("backup.txt", ch_buttonPressed)

	// Goroutines for sending and recieving messages
	go bcast.Transmitter(config.DefaultPortBcast, ch_msgOut)
	go bcast.Receiver(config.DefaultPortBcast, ch_msgIn)

	go peers.Transmitter(config.DefaultPortPeer, nodeID, ch_peerTxEnable)
	go peers.Receiver(config.DefaultPortPeer, ch_peerUpdate)

	// elevator_io goroutines
	go elevator_io.PollButtons(ch_buttonPressed)
	go elevator_io.PollFloorSensor(ch_arrivalFloor)
	go elevator_io.PollObstructionSwitch(ch_doorObstruction)
	go elevator_io.PollStopButton(ch_stopButton)

	go fsm.FSM(
		nodeID,
		ch_localRequests,
		ch_arrivalFloor,
		ch_doorObstruction,
		ch_stopButton,
		ch_completedRequests,
		ch_elevatorStateToAssigner,
		ch_elevatorStateToNetwork,
		ch_FSMDeadlock,
	)

	go assigner.RequestAssigner(
		nodeID,
		ch_buttonPressed,
		ch_completedRequests,
		ch_elevatorStateToAssigner,
		ch_hallRequestsIn,
		ch_externalElevators,
		ch_hallRequestsOut,
		ch_localRequests,
		ch_assignerDeadlock,
	)
	
	go heartbeat.Send(
		nodeID,
		ch_hallRequestsOut,
		ch_elevatorStateToNetwork,
		ch_msgOut,
		ch_heartbeatDeadlock,
	)

	go peers.Update(
		nodeID,
		ch_peerUpdate,
		ch_msgIn,
		ch_hallRequestsIn,
		ch_externalElevators,
		ch_peersDeadlock,
	)

	go deadlock.Detector(
		ch_FSMDeadlock,
		ch_assignerDeadlock,
		ch_heartbeatDeadlock,
		ch_peersDeadlock,
	)

	select{}
}
