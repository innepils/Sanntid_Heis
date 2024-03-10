package config

import (
	"driver/network/localip"
	"flag"
	"fmt"
	"os"
)

// ***** System specifications *****
const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

const DoorOpenDurationSec int = 3

// ***** Network-configuration *****
const DefaultPortPeer int = 22017
const DefaultPortBcast int = 22018
const DefaultPortBackup int = 22019

func InitializeConfig() (string, string) {
	var id, port string
	flag.StringVar(&id, "id", getDefaultID(), "ID of this peer")
	flag.StringVar(&port, "port", "15657", "Port of this peer")
	flag.Parse()
	return id, port
}

func getDefaultID() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println("Error obtaining local IP:", err)
		return "DISCONNECTED"
	}
	return fmt.Sprintf("peer_%s:%d", localIP, os.Getpid())
}
