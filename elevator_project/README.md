# Elevator Project

The project runs **n** elevator in parallell across **m** floors using a peer to peer network and udp broadcasting.

## Prerequisites
The program is meant to be run using Linux, which must be done to ensure correct functionality. With small modifications the application can be expanded to work on all platforms. 

The project is built using Go, which needs to be installed to compile and run the project. The latest version of Go can be installed from [the official Go website](https://go.dev/dl/).



## Usage

The main program is run using run.sh while in it's directory.
Before running, set the permissions using
```
chmod +x run.sh
```
then run the main program with

```
./run.sh
```

To exit, close the terminal by holding down ctrl+c. 

## Module descriptions

### Assigner

The assigner assigns requests to the local elevator. To do this it keeps track of incomming button presses, locally compleded orders and status of other elevators. It also has idle time out, if there are orders and our local elevator is idle for a "long" time it takes all orders.

### Backup

This saves the local cab calls to file and also supports extracting it from file again.

### Cost

This utilizes the "[HallRequestAssigner](#hall-request-assigner)" to calculate witch requests our local elevator should serve.

### Deadlock detector

This detects if [FSM](#fsm), [assigner](#assigner), [peers](#network) or [heartbeat](#heartbeat) is stuck within a loop. If this is detected the program reboots.

### Elevator

### Elevator IO

Information can be found [here](https://github.com/TTK4145/driver-go).

### FSM

The FSM is event-driven, and after initializing the local elevator it checks for following events:
 - Recieved request from assigner
 - Arrival at new floor
 - Door timer timed out
 - Door obstructed
 - Stop-button pressed

- Arrival at new floor
- Recieved request from assigner
-

### Hall Request Assigner

Information can be found [here](https://github.com/TTK4145/Project-resources/tree/master/cost_fns/hall_request_assigner).

### Heartbeat
Sets up the struct which is broadcasted to the network, containing information about new hall requests and state from [assigner](#assigner) each local elevator. 

### Network

Most of the documentation can be found [here](https://github.com/TTK4145/Network-go).

In the handed out peers.go we have aded functionality to continuously update the alivePeers to be used in [cost](#cost). To avoid concurrency issues while reading and writing to the map both in peers and [assigner](#assigner), we serialize the maps into JSON using Marshal and Unmarshal. 

### Requests

This package is based on [this](https://github.com/TTK4145/Project-resources/blob/master/elev_algo/requests.c), but modified to fit the projects event-driven FSM. 
The most important changes are:
- The elevator is taken is passed-by-reference, to avoid uneccessary opying and making a duplicate ("pair") of the elvator.
- The function ClearAtCurrentFloor was modified so that if both hall-orders at a floor is active, they are not cleared at the same time, to fulfill the project specifiactions.
- The function "AnnounceDirectionChange" is added to fulfill the project specifiactions.
- The function ClearImmidiately was removed as it was no longer needed.
