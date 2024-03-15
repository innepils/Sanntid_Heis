# Elevator Project

The project runs **n** elevator in parallell across **m** floors using a peer to peer network and udp broadcasting.

## Setting up and running the project

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

The FSM is event-driven, and after initializing the local elevator it checks for, and acts on following events:
 - Recieved request from assigner
 - Arrival at new floor
 - Door timer timed out
 - Door obstructed
 - Stop-button pressed

### Hall Request Assigner

Information can be found [here](https://github.com/TTK4145/Project-resources/tree/master/cost_fns/hall_request_assigner).

### Heartbeat
Sets up the struct which is broadcasted to the network, containing information about new hall requests and state from [assigner](#assigner) each local elevator. 

### Network

Most of the documentation can be found [here](https://github.com/TTK4145/Network-go).

In the handed out peers.go we have added functionality to continuously update the alivePeers to be used in [cost](#cost). To avoid concurrency issues while reading and writing to the map both in peers and [assigner](#assigner), we serialize the maps into JSON using Marshal and Unmarshal. 

### Requests

This package takes care of logic regarding local requests, giving the options of checking where the requests are, and what resulting behaviour the elevator should have. All functions take in the local elevator by using pass-by-reference.

This package is based on [this](https://github.com/TTK4145/Project-resources/blob/master/elev_algo/requests.c) c-module, but translated and modified to fit the projects event-driven FSM. 
The most important changes are that:
- The elevator is passed-by-reference instead of making a duplicate copy.
- The function "AnnounceDirectionChange" is added to fulfill the project specifiactions.
- The function "ClearAtCurrentFloor" only clears one hall-button per time, to fulfill the project specifiactions.
- Also, "ClearImmidiately" was removed.