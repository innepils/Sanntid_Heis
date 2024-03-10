package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

const bufSize = 1024

func main() {
	// Replace with the number of your workspace
	workspaceNumber := 7
	sendPort := 20000 + workspaceNumber

	// Goroutine for receiving server IP
	go func() {
		addr := net.UDPAddr{
			Port: 30000,
			IP:   net.ParseIP("0.0.0.0"),
		}
		conn, err := net.ListenUDP("udp", &addr)
		if err != nil {
			fmt.Println(err)
			return
		}

		buffer := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Server IP: %s\n", string(buffer[:n]))

		defer conn.Close()
	}()

	// Goroutine for sending messages
	go func() {
		serverIP := "10.100.23.129"
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

		message := []byte("Hello from workspace " + strconv.Itoa(workspaceNumber))
		_, err = conn.Write(message)
		if err != nil {
			fmt.Println(err)
			return
		}

	}()

	// Go Read Routine
	go func() {
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
	}()

	// Wait for goroutines to finish
	time.Sleep(7 * time.Second)
}
