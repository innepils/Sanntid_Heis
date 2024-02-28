package backup

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

const (
	sendAddr       = "255.255.255.255:20007"
	receiveAddr    = ":20007"
	baseStatusMsg  = "heartbeat"
	heartbeatSleep = 1000
)

func StartBackupProcess() {
	exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
}

func BackupProcess(localID string) {
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
				StartBackupProcess()
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

		if parts[0] == localID {
			localState = string(msg[2])
			conn.SetReadDeadline(time.Now().Add(heartbeatSleep * 5 * time.Millisecond))
		}
	}
}