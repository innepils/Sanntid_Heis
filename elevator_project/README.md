# Elevator Project

The project runs a n elevator across m floors using a peer to peer network and udp broadcasting.

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

This utilises the "HallRequestAssigner" to calculate witch requests our local elevator should serve.

### Deadlock detector

This detects if FSM, assigner, peers or heartbeat is stuck within a loop. If this is detected the program reboots.
