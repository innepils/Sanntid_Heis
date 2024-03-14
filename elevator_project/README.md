# Elevator Project

The project runs **n** elevator in parallell across **m** floors using a peer to peer network and udp broadcasting.

## Setting up and running the project

The main program is run using run.sh while in it's directory.
Before running set the permissions using
```
chmod +x run.sh
```
Then run the main program with

```
./run.sh
```

To exit, close the terminal by holding down ctrl+c. 

## Module dscription

### Assigner

The assigner assigns requests to the local elevator. To do this it keeps track of incomming button presses, locally compleded orders and status of other elevators. It also has idle time out, if there are orders and our local elevator is idle for a "long" time it takes all orders.

### Backup

This saves the local cab calls to file and also supports extracting it from file again.

### Cost

This utilises the "[HallRequestAssigner](#hall-request-assigner)" to calculate witch requests our local elevator should serve.

### Deadlock detector

This detects if [FSM](#fsm), [assigner](#assigner), [peers](#network) or [heartbeat](#heartbeat) is stuck within a loop. If this is detected the program reboots.

### Elevator

### Elevator IO

### FSM

### Hall Request Assigner

The documentation can be found [here](https://github.com/TTK4145/Project-resources/tree/master/cost_fns/hall_request_assigner).

### Heartbeat
Sets up the struct which is broadcasted to the network, containing information about new hall requests and state from [assigner](#assigner) each local elevator. 

### Network

Most of the documentation can be found [here](https://github.com/TTK4145/Network-go).

In the handed out peers.go we have aded functionality to continuously update the alivePeers to be used in [cost](#cost). To avoid concurrency issues while reading and writing to the map both in peers and [assigner](#assigner)], we serialize the maps into JSON using Marshal and Unmarshal. 

### Requests
