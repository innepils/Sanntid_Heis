package backup

import (
	"driver/config"
	"driver/elevator"
	"driver/elevator_io"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

func KillSelf(localID string, port string) { //unused
	StartBackupProcess(localID, port)
	panic("Program terminated")
}

func SaveBackupToFile(filename string, allRequests [config.N_FLOORS][config.N_BUTTONS]elevator.RequestType) {
	var cabRequests [config.N_FLOORS]bool
	for request := range allRequests {
		if allRequests[request][2] == elevator.ConfirmedRequest {
			cabRequests[request] = true
		} else {
			cabRequests[request] = false
		}
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(cabRequests)
	if err != nil {
		return
	}
}

func LoadBackupFromFile(filename string, ch_buttonPressed chan elevator_io.ButtonEvent) {
	var data [config.N_FLOORS]bool

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Eroor decoding data from backup")
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Eroor decoding data from backup")
	}

	for i, element := range data {
		if element {
			ch_buttonPressed <- elevator_io.ButtonEvent{BtnFloor: i, BtnType: elevator_io.BT_Cab}
		}
	}
}

func StartBackupProcess(localID string, port string) {
	exec.Command("gnome-terminal", "--", "go", "run", "main.go", "-id="+localID, "-port="+port).Run()
}

func ReportPrimaryAlive(localID string) {
	sendUDPAddr, err := net.ResolveUDPAddr("udp", config.BackupSendAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.DialUDP("udp", nil, sendUDPAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		msg := localID
		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Primary failed to send heartbeat:", err)
			return
		}
		time.Sleep(config.HeartbeatSleepSec * time.Second)
	}
}

func BackupProcess(localID string, port string) { //name change: ListenForPrimary ???
	localState := ""
	fmt.Println(localState)
	fmt.Printf("---------BACKUP PHASE---------\n")
	receiveUDPAddr, err := net.ResolveUDPAddr("udp", config.BackupReceiveAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.ListenUDP("udp", receiveUDPAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(config.HeartbeatSleepSec * 5 * time.Second))
	for {
		buffer := make([]byte, 1024)
		//conn.SetReadDeadline(time.Now().Add(heartbeatSleep * 2.5 * time.Millisecond))
		n, _, err := conn.ReadFromUDP(buffer)

		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				fmt.Println("Backup did not receive heartbeat, becoming primary.")
				// This is where the backup takes over and becomes Primary
				conn.Close()
				StartBackupProcess(localID, port)
				return
			} else {
				fmt.Println("Error reading from UDP:", err)
				return
			}
		}

		msg := string(buffer[:n])
		if msg == localID {
			fmt.Println("Primary is alive!")
			conn.SetReadDeadline(time.Now().Add(config.HeartbeatSleepSec * 2500 * time.Millisecond))
		}
	}
}
