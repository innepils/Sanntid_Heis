package config

import (
	"src/network/localip"
	"flag"
	"fmt"
	"os"
)

const (
	// System specifications
	N_FLOORS           			int = 4
	N_BUTTONS            		int = 3
	DoorOpenDurationSec  		int = 3
	IdleTimeOutDurationSec 		int = 15
	DeadlockTimeOutDurationSec 	int = 10
	// Network-configuration
	DefaultPortPeer				int = 22017
	DefaultPortBcast  			int = 22018
	elevatorServerPort			string = "15657"
	HeartbeatSleepMillisec 		int = 100
	
)

func InitializeIDandPort() (string, string) {
	var nodeID, port string
	flag.StringVar(&nodeID, "id", getDefaultID(), "ID of this peer")
	flag.StringVar(&port, "port", elevatorServerPort, "Port of this peer")
	flag.Parse()
	return nodeID, port
}

func getDefaultID() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println("Error obtaining local IP:", err)
		return "DISCONNECTED"
	}
	return fmt.Sprintf("peer_%s:%d", localIP, os.Getpid())
}
