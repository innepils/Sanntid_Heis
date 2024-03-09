package config

import (
	"driver/network/localip"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// ***** Network-configuration *****
const DefaultPortPeer int = 22017
const DefaultPortBcast int = 22018
const DefaultPortBcastBackup int = DefaultPortBcast + 2

var DefaultPortBcastStr string = strconv.Itoa(DefaultPortBcast)
var DefaultPortBcastBackupStr string = strconv.Itoa(DefaultPortBcastBackup)

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
