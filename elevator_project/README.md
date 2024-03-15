# Elevator Project

The project runs **n** elevators in parallell across **m** floors using a peer to peer network and UDP broadcasting.

## Prerequisites
The program is meant to be run using Linux, which must be done to ensure correct functionality. 

The project is built using Go, which needs to be installed to compile and run the program. The latest version of Go can be installed from [the official Go website](https://go.dev/dl/).

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

The program will restart immediately if terminated using ctrl+c. In order to fully stop the running program, close the terminal or keep holding ctrl+x. 

*Note: If you want to manually input id and port when running the program, you can use 'go run main.go -id=any_id -port=elevator_port' in the command line. However, with the current implementation this will reduce the overall fault tolerance compared to using run.sh'.*

## Module descriptions

### Main

### Assigner
The assigner assigns requests to the local elevator. To do this it keeps track of incomming button presses, locally compleded orders and status of other elevators. The assigner also has a idle time out: if there are orders and our local elevator is idle for a "long" time it takes all orders.

### Backup
This saves the local cab calls to file and also supports extracting it from file again.

### Config
This package contains variables (system specifics) to be used by other packages. This makes them easily to access, as well as easy to modify if the user wishes to change some system specifics.

The package also contains two functions that determines the nodes personal ID and port.

### Cost
This utilizes the "[HallRequestAssigner](#hall-request-assigner)" to calculate witch requests our local elevator should serve.

### Deadlock
The modules only function ("Detector") detects if [FSM](#fsm), [assigner](#assigner), [peers](#network) or [heartbeat](#heartbeat) is stuck within a loop. If this is detected the program reboots.

### Elevator
This package defines the `Elevator` struct which holds the local elevators most important states, and is to be used by other packages, e.g. [assigner](#assigner) and [FSM](#FSM).

It also holds functions for initilaizing an elevator-object, setting the button-lights (according to the nodes agreed upon requests), sending the elevator-objects state, and functions for converting data-types 

### Elevator IO
Information can be found [here](https://github.com/TTK4145/driver-go).

### FSM
The FSM is event-driven, and after initializing the local elevator it checks for, and acts on following events:
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

### Network
The network-module consists of the handed out modules which includes [bcast](#bcast) for transmitting and recieving messages, [conn](#conn) for establishing the UDP connections, [localip](#localip) for finding the nodes local IP and [peers](#peers) to detect other peers on the network. We have expanded the functionality in [peers](#peers) to read the heartbeat and send a map of the external elevators to be used in [cost](#cost). 

Most of the documentation can be found [here](https://github.com/TTK4145/Network-go).

#### Heartbeat
This last networking module sets up the struct which is continuously broadcasted to the network, containing information about new hall requests and state from [assigner](#assigner) each local elevator. 

#### Peers
In the handed out peers.go we have added functionality to continuously update the alivePeers to be used in [cost](#cost). To avoid concurrency issues while reading and writing to the map both in peers and [assigner](#assigner), we serialize the maps into JSON using Marshal and Unmarshal. 

### Requests
This package takes care of logic regarding local requests, giving the options of checking where the requests are, and what resulting behaviour the elevator should have. All functions take in the local elevator by using pass-by-reference.

This package is based on [this](https://github.com/TTK4145/Project-resources/blob/master/elev_algo/requests.c) c-module, but translated and modified to fit the projects event-driven FSM. 
The most important changes are that:
- The elevator is passed-by-reference instead of making a duplicate copy.
- The function "AnnounceDirectionChange" is added to fulfill the project specifiactions.
- The function "ClearAtCurrentFloor" only clears one hall-button per time, to fulfill the project specifiactions.
- Also, "ClearImmidiately" was removed.