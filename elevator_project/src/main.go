package main

import (
	"driver/backup"
	"driver/config"
	"driver/elevator"
	"driver/elevator_io"
	"driver/fsm"
	"fmt"
)

type ElevatorMessage struct {
	ID           string
	HallRequests bool
	state        int
	Iter         int
}

func main() {

	//var boolArray [4]bool = {false, true, false, true}
	//backup.SaveBackupToFile("status.txt", boolArray)

	/* Initialize elevator ID and port
	This section sets the elevators ID (anything) and port (of the running node/PC),
	which should be passed on in the command line using
	'go run main.go -id=any_id -port=port'
	*/ /*
			var id string
			flag.StringVar(&id, "ID", "", "ID of this peer")
			flag.Parse()

			// if no ID is given, use local IP address
			if id == "" {
				localIP, err := localip.LocalIP()
				if err != nil {

					fmt.Println(err)
					localIP = "DISCONNECTED"
				}
				id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())

			}/*

		//	backup.BackupProcess(id) //this halts the progression of the program while it is the backup
		//	fmt.Println("Primary started.")
			// Initialize local elevator
		//	localIPstring, _ := localip.LocalIP()
		//	elevator_io.Init(localIPstring+string(config.GlobalPort), elevator_io_types.N_FLOORS)


			// GOROUTINES network
		/*
			// We make a channel for receiving updates on the id's of the peers that are
			//  alive on the network
			peerUpdateCh := make(chan peers.PeerUpdate)

			//hannel for enabling/disabling the transmitter after start.
			//Can be used to signal that the node is "unavailable".
			peerTxEnable := make(chan bool)
			go peers.Transmitter(config.GlobalPort, id, peerTxEnable)
			go peers.Reciever(config.GlobalPort, peerUpdateCh)

			// Channels for sending and recieving
			msgTx := make(chan ElevatorMessage)
			msgRx := make(chan ElevatorMessage)

			go bcast.Transmitter(config.GlobalPort, msgTx)
			go bcast.Reciever(config.GlobalPort, msgRx)

			// example message
			go func() {
				testMsg := ElevatorMessage{"nice ID", true, 0, 1}
				for {
					ElevatorMessa/*ge.Iter++
					msgTx <- testMsg
					time.Sleep(1 * time.Second)
				}
			}()
	*/

	elevator_io.Init("localhost:15657", config.N_FLOORS)
	// Channels for sending and recieving
	ch_arrivalFloor := make(chan int)
	ch_buttonPressed := make(chan elevator_io.ButtonEvent)
	ch_localOrders := make(chan [config.N_FLOORS][config.N_BUTTONS]bool)
	ch_doorObstruction := make(chan bool)
	ch_stopButton := make(chan bool)
	ch_completedOrders := make(chan elevator_io.ButtonEvent)

	go backup.LoadBackupFromFile("status.txt", ch_buttonPressed)

	go elevator_io.PollButtons(ch_buttonPressed)
	go elevator_io.PollFloorSensor(ch_arrivalFloor)
	go elevator_io.PollObstructionSwitch(ch_doorObstruction)
	go elevator_io.PollStopButton(ch_stopButton)

	go fsm.Fsm(ch_arrivalFloor, ch_localOrders, ch_buttonPressed, ch_doorObstruction, ch_stopButton, ch_completedOrders)

	go func() {
		for {
			event := <-ch_completedOrders
			fmt.Printf("Received event. Floor %d, Btn: %s\n", event.BtnFloor+1, elevator.ElevButtonToString(event.BtnType))
		}
	}()

	for {

	}
}
