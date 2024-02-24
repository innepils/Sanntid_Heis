package network

import (
	"fmt"
	"net"
	"strconv"
)

var sendPort int = 20007

func UDPlistener(incommingMsgCh chan string) {
	//serverIP := "0.0.0.0" //Recieve from all (?)
	sendAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(sendPort))

	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.ListenUDP("udp", sendAddr)

	if err != nil {
		fmt.Println(err)
		return
	}
	buffer := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buffer)

		if err != nil {
			fmt.Println(err)
			return
		}
		var msgFromBuffer string = string(buffer[:n])
		fmt.Printf("Node Response: %s\n", msgFromBuffer)
		incommingMsgCh <- msgFromBuffer
	}
	// return
}

func UDPsend(msg string) {
	serverIP := "255.255.255.255"
	sendAddr, err := net.ResolveUDPAddr("udp", serverIP+":"+strconv.Itoa(sendPort))
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.DialUDP("udp", nil, sendAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	message := []byte(msg)
	_, err = conn.Write(message)

	if err != nil {
		fmt.Println(err)
		return
	}
}
