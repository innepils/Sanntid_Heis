package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		data, err := bufio.NewReader(conn).ReadString('\000')

		if err != nil {
			fmt.Println("Nå er det i alle fall feil:", err)
			return
		}
		fmt.Println(data)
		conn.Write([]byte(data + "\000"))
		time.Sleep(3 * time.Second)
	}
}

func main() {
	// Erstatt med den faktiske server-IPen og velg riktig port
	serverIP := "10.100.23.129" // Erstatt med serverens IP
	sendPort := "33546"         // Bruk "33546" for \0-terminerte meldinger, "34933" for fast størrelse

	sendAddr, err := net.ResolveTCPAddr("tcp", serverIP+":"+sendPort)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, sendAddr)
	if err != nil {
		fmt.Println("Feil ved tilkobling:", err)
		return
	}
	//defer conn.Close()

	fmt.Println("Koblet til serveren på", sendAddr)

	// Motta velkomstmelding
	message, err := bufio.NewReader(conn).ReadString('\000')

	if err != nil {
		fmt.Println("Feil ved mottak av velkomstmelding:", err)
		return
	}
	fmt.Print("Melding fra server: ", message)

	// Send meldinger
	_, err = conn.Write([]byte("Hello Dear TCP Server" + "\000"))
	fmt.Println("melding sendt...")

	if err != nil {
		fmt.Println("Feil ved sending av melding:", err)
		return
	}

	// Motta meldinger
	response, err := bufio.NewReader(conn).ReadString('\000')

	if err != nil {
		fmt.Println("Feil ved mottak av svar:", err)
		return
	}
	fmt.Println("Svar fra server: ", response)

	// Send melding for å be serveren koble tilbake
	myIP := "10.100.23.17"
	myPort := "20007"
	connectBackMessage := fmt.Sprintf("Connect to: %s:%s\000", myIP, myPort)
	_, err = conn.Write([]byte(connectBackMessage))

	if err != nil {
		fmt.Println("Feil ved sending av tilkoblingstilbake-melding:", err)
		return
	}

	fmt.Println("Venter på at serveren skal koble til...")

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+"20007")
	if err != nil {
		fmt.Println("feil i oppsett:", err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("feil i listen:", err)
		return
	}

	for {
		fmt.Println("in for loop")
		conn, err := listener.Accept()
		fmt.Println("listener accept")
		if err != nil {
			fmt.Println("Nå er det i alle fall feil:", err)
			return
		}

		go handleConnection(conn)

	}
	// Her kan du implementere logikk for å akseptere tilkoblingen fra serveren

	time.Sleep(3 * time.Second)

}
