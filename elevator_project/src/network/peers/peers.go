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

const (
	interval = 150 * time.Millisecond
	timeout = 500 * time.Millisecond
)

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
	nodeID 					string,
	ch_peerUpdate 		 	<-chan PeerUpdate,
	ch_msgIn 			 	<-chan heartbeat.HeartBeat,
	ch_hallRequestsIn 	 	chan<- [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType,
	ch_externalElevators	chan<- []byte,
	ch_peersDeadlock 	 	chan<- int,
){

	alivePeers := make(map[string]elevator.ElevatorState)
	var prevHallRequests [config.N_FLOORS][config.N_BUTTONS - 1]elevator.RequestType

	for {
		ch_peersDeadlock <- 1
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

		case heartbeat := <-ch_msgIn:
			if heartbeat.SenderID != nodeID {
				if !reflect.DeepEqual(alivePeers[heartbeat.SenderID], heartbeat.ElevatorState) {
					alivePeers[heartbeat.SenderID] = heartbeat.ElevatorState
					AlivePeersJson, _ := json.Marshal(alivePeers)
					ch_externalElevators <- AlivePeersJson
				}

				if prevHallRequests != heartbeat.HallRequests {
					prevHallRequests = heartbeat.HallRequests
					ch_hallRequestsIn <- prevHallRequests
				}
			}
		default:
			// NOP
		}
	}

}
