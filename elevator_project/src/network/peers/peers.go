package peers

import (
	"driver/config"
	"driver/elevator"
	"driver/heartbeat"
	"driver/network/conn"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"sort"
	"time"
)

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

const interval = 150 * time.Millisecond
const timeout = 500 * time.Millisecond

func Transmitter(port int, id string, transmitEnable <-chan bool) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			conn.WriteTo([]byte(id), addr)
		}
	}
}

func Receiver(port int, peerUpdateCh chan<- PeerUpdate) {

	var buf [1024]byte
	var p PeerUpdate
	lastSeen := make(map[string]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {
		updated := false

		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, _ := conn.ReadFrom(buf[0:])

		id := string(buf[:n])

		// Adding new connection
		p.New = ""
		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {
				p.New = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}

		// Removing dead connection
		p.Lost = make([]string, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}

		// Sending update
		if updated {
			p.Peers = make([]string, 0, len(lastSeen))

			for k := range lastSeen {
				p.Peers = append(p.Peers, k)
			}

			sort.Strings(p.Peers)
			sort.Strings(p.Lost)
			peerUpdateCh <- p
		}
	}
}

func Update(
	id string,
	ch_peerUpdate chan PeerUpdate,
	ch_msgIn chan heartbeat.HeartBeat,
	ch_hallRequestsIn chan [config.N_FLOORS][config.N_BUTTONS - 1]int,
	ch_externalElevators chan []byte) {

	alivePeers := make(map[string]elevator.ElevatorState)
	var prevHallRequests [config.N_FLOORS][config.N_BUTTONS - 1]int
	var prevAlivePeers map[string]elevator.ElevatorState
	for {
		//fmt.Println("Alive peers: ", alivePeers)
		select {
		case peers := <-ch_peerUpdate:
			for _, peer := range peers.Lost {
				if _, ok := alivePeers[peer]; ok {
					delete(alivePeers, peer)
					AlivePeersJson, _ := json.Marshal(alivePeers)
					ch_externalElevators <- AlivePeersJson
				}
			}

			fmt.Printf("\nPeer update:\n")
			fmt.Printf("  Peers:    %q\n", peers.Peers)
			fmt.Printf("  New:      %q\n", peers.New)
			fmt.Printf("  Lost:     %q\n", peers.Lost)

		case a := <-ch_msgIn:
			if a.SenderID != id {
				alivePeers[a.SenderID] = a.ElevatorState
				fmt.Println("Alive Peers: ", alivePeers)
				if prevHallRequests != a.HallRequests {
					prevHallRequests = a.HallRequests
					ch_hallRequestsIn <- prevHallRequests
				}
				if !reflect.DeepEqual(prevAlivePeers, alivePeers) {
					fmt.Println(alivePeers)
					prevAlivePeers = alivePeers
					AlivePeersJson, _ := json.Marshal(prevAlivePeers)
					ch_externalElevators <- AlivePeersJson
				}
				//fmt.Printf("Received: %#v\n", a)
			}
		default:
			// NOP
		}
	}

}

// FORSLAG TIL ENDRINEGER FRA CHAT (liste med punkter arkivert i chat)
// func Update(
// 	id string,
// 	ch_peerUpdate chan PeerUpdate,
// 	ch_msgIn chan heartbeat.HeartBeat,
// 	ch_hallRequestsIn chan [config.N_FLOORS][config.N_BUTTONS-1]int,
// 	ch_externalElevators chan []byte) {

// 	alivePeers := make(map[string]elevator.ElevatorState)
// 	var prevHallRequests [config.N_FLOORS][config.N_BUTTONS-1]int
// 	prevAlivePeers := make(map[string]elevator.ElevatorState)
// 	for {
// 		select {
// 		case peers := <-ch_peerUpdate:
// 			for _, peer := range peers.Lost {
// 				delete(alivePeers, peer)
// 			}
// 			for _, peer := range peers.New {
// 				// Initialize or update the peer's state as needed
// 				// alivePeers[peer] = someInitialState
// 			}
// 			alivePeersJSON, err := json.Marshal(alivePeers)
// 			if err != nil {
// 				// log.Println("Error marshalling alivePeers:", err) // Example logging
// 				continue
// 			}
// 			ch_externalElevators <- alivePeersJSON

// 		case a := <-ch_msgIn:
// 			if a.SenderID != id {
// 				alivePeers[a.SenderID] = a.ElevatorState
// 				if !reflect.DeepEqual(prevHallRequests, a.HallRequests) {
// 					prevHallRequests = a.HallRequests
// 					ch_hallRequestsIn <- prevHallRequests
// 				}
// 				if !reflect.DeepEqual(prevAlivePeers, alivePeers) {
// 					// Deep copy alivePeers to prevAlivePeers
// 					for k, v := range alivePeers {
// 						prevAlivePeers[k] = v
// 					}
// 					// Remove any keys in prevAlivePeers not in alivePeers
// 					for k := range prevAlivePeers {
// 						if _, ok := alivePeers[k]; !ok {
// 							delete(prevAlivePeers, k)
// 						}
// 					}
// 					alivePeersJSON, err := json.Marshal(alivePeers)
// 					if err != nil {
// 						// log.Println("Error marshalling alivePeers:", err) // Example logging
// 						continue
// 					}
// 					ch_externalElevators <- alivePeersJSON
// 				}
// 			}
// 		}
// 	}
// }