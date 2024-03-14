# Elevator Project

The project runs **n** elevator in parallell across **m** floors using a peer to peer network and udp broadcasting.

## Setting up and running the project

The project is run using the run.sh file by writing

```
./run.sh
```

while in the directory where the file is located.

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

Information can be found [here](https://github.com/TTK4145/driver-go).

### FSM

The FSM is event-driven, and after initializing the local elevator it checks for following events:
 - Arrival at new floor
 - Recieved request from assigner
 - 

### Hall Request Assigner

Information can be found [here](https://github.com/TTK4145/Project-resources/tree/master/cost_fns/hall_request_assigner).

### Heartbeat

### Network

Most of the documentation can be found [here](https://github.com/TTK4145/Network-go).

### Requests
