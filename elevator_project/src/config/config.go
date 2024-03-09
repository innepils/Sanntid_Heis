package config

import (
	"driver/network/localip"
	"flag"
	"fmt"
	"os"
)

// ***** Network-configuration *****
const DefaultPortPeer int = 20017
const DefaultPortBcast int = 20018

// ***** System specifications *****
const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

type ClearRequestVariant int

const (
	CV_all    ClearRequestVariant = iota // Assumes customers enter the elevator even though its moving in the wrong direction
	CV_InDirn                            // Assumes customers only enter the elevator when its moving in the correct direction
)

const SystemsClearRequestVariant ClearRequestVariant = CV_InDirn
const DoorOpenDurationSec int = 3

/*
Initialize elevator ID and port
This section sets the elevators ID and port
which should be passed on in the command line using
'go run main.go -id=any_id -port=server_port'
*/
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
