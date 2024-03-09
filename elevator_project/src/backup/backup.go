package backup

import (
	"driver/config"
	"driver/elevator_io"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	sendAddr       = "255.255.255.255:20007"
	receiveAddr    = config.DefaultPortBcastStr
	baseStatusMsg  = "heartbeat"
	heartbeatSleep = 1000
)

func SaveBackupToFile(filename string, status [4]bool) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	// Encode the array using a gob encoder
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(status)
	if err != nil {
		return
	}

	return
}

func LoadBackupFromFile(filename string, ch_buttonPressed chan elevator_io.ButtonEvent) {
	var data [4]bool

	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	// Decode the data using a gob decoder
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return
	}

	for i, element := range data {
		if element {
			ch_buttonPressed <- elevator_io.ButtonEvent{BtnFloor: i, BtnType: elevator_io.BT_Cab}
		}
	}

	return
}

func StartBackupProcess(id string, port string) {
	exec.Command("gnome-terminal", "--", "go", "run", "main.go", "-id="+id, "-port="+port).Run()
}

func BackupProcess(id string, port string) {
	localState := ""
	fmt.Println(localState)
	fmt.Printf("---------BACKUP PHASE---------\n")
	receiveUDPAddr, err := net.ResolveUDPAddr("udp", receiveAddr)
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
	conn.SetReadDeadline(time.Now().Add(heartbeatSleep * 5 * time.Millisecond))
	for {
		buffer := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(heartbeatSleep * 5 * time.Millisecond))
		n, _, err := conn.ReadFromUDP(buffer)

		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				fmt.Println("Backup did not receive heartbeat, becoming primary.")
				// This is where the backup takes over and becomes Primary
				conn.Close()
				StartBackupProcess(id, port)
				return
			} else {
				fmt.Println("Error reading from UDP:", err)
				return
			}
		}

		msg := string(buffer[:n])

		parts := strings.Split(msg, ";")

		// if len(parts) > 0 {
		// 	// Access the first element
		// 	firstElement := parts[0]
		// 	fmt.Println("First element:", firstElement)
		// } else {
		// 	fmt.Println("String is empty or doesn't contain any ';'")
		// }

		if parts[0] == id {
			localState = string(msg[2])
			conn.SetReadDeadline(time.Now().Add(heartbeatSleep * 5 * time.Millisecond))
		}
	}
}
