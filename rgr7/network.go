package network

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

sendPort = "20007"


func UDPlistener(){
	//serverIP := "10.100.23.129"
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
		fmt.Printf("Server Response: %s\n", string(buffer[:n]))
	}
}


func UDPsend(){

}