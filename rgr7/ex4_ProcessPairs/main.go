// Combined Program using UDP communication
package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"time"
)

const (
	sendAddr       = "localhost:20009"
	receiveAddr    = "localhost:"
	heartbeatMsg   = "heartbeat"
	heartbeatSleep = 500
)

// Function to start a backup process that will become primary if needed.
func startBackupProcess() {
	exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
}

// The primary process sends heartbeats to the backup.
func primaryProcess() {
	sendUDPAddr, _ := net.ResolveUDPAddr("udp", sendAddr)
	conn, _ := net.DialUDP("udp", nil, sendUDPAddr)
	defer conn.Close()

	count := 1
	fmt.Printf("before for loop\n")
	for {
		fmt.Printf("for loop started\n")
		msg := heartbeatMsg + ":" + strconv.Itoa(count)
		_, err := conn.Write([]byte(msg))
		fmt.Printf("before if\n")
		if err != nil {
			fmt.Println("Primary failed to send heartbeat:", err)
			return
		}
		fmt.Printf("Primary count: %d\n", count)
		count++
		time.Sleep(heartbeatSleep * time.Millisecond)
	}
	fmt.Printf("ool\n")
}

// The backup process listens for heartbeats from the primary.
func backupProcess() {
	fmt.Printf("---------BACKUP STARTED---------")
	receiveUDPAddr, _ := net.ResolveUDPAddr("udp", receiveAddr)
	conn, _ := net.ListenUDP("udp", receiveUDPAddr)
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(heartbeatSleep * 5 * time.Millisecond))
		n, _, err := conn.ReadFromUDP(buffer)

		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				fmt.Println("Backup did not receive heartbeat, becoming primary.")
				// This is where the backup takes over and becomes Primary
				conn.Close()
				startBackupProcess()
				primaryProcess()
				return
			} else {
				fmt.Println("Error reading from UDP:", err)
				return
			}
		}

		msg := string(buffer[:n])
		if msg[:len(heartbeatMsg)] == heartbeatMsg {
			countStr := msg[len(heartbeatMsg)+1:]
			count, _ := strconv.Atoi(countStr)
			fmt.Printf("Backup received count: %d\n", count)
		}
	}
}

func main() {
	// if len(os.Args) > 1 && os.Args[1] == "primary" {
	// 	startBackupProcess()
	// 	primaryProcess()
	// } else {
	// 	backupProcess()
	// }
	backupProcess()

}
