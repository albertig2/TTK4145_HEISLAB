# Debugging helpers

This module contains functions to ease the debugging process, specifically by providing a way to expose the state of the system through print functions. It contains functions to print the different structs used in the system, as well as some utility functions. There are three public functions, used to print peerUpdates,  the elevator object and print an update of the full peerView object and peer network info.

## Output examples

### PrintLocalElvator(elevator)
#### Input
``` go
elevator := elevatorConfig.Elevator{
    OwnId:           "1",
    Floor:           2,
    Direction:       elevatorConfig.down,
    LocalOrderQueue: [elevatorConfig.N_FLOORS][elevatorConfig.N_BUTTONS]bool{},
    Behavior:        elevatorConfig.Moving,
}
PrintLocalElvator(elevator)

```
#### Output
``` bash
  +--------------------+
  |floor  = 2          |
  |dirn   = down       |
  |behav  = moving     |
  +--------------------+
  |  | up  | dn  | cab |
  | 3|     |  -  |  -  |
  | 2|  -  |  -  |  -  |
  | 1|  -  |  -  |  -  |
  | 0|  -  |     |  -  |
  +--------------------+
```
### PrintPeerUpdate (Update peers.PeerUpdate)
#### Input

``` go
peerUpdate := peer.PeerUpdate{
    Peers [] string {"1", "2"},
    New   "1",
    Lost  []string {},
}
PrintPeerUpdate(peerUpdate)
```
#### Output

``` bash
--------New peer uppdate recieved----------
Current alive peers: 1 2
Elevator ID  1 just joind the network
Peers considerd lost:
---------End of peer update--------------
```
### PrintElevatorSystem(elevatorSystem elevatorConfig.PeerView)
#### Input
``` go
peerState1 := elevatorConfig.PeerState{
    Behavior:    elevatorConfig.Moving,
    Floor:    1,
    Direction:   elevatorConfig.down,
    CabRequests: [elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus{elevatorConfig.Assigned, elevatorConfig.Assigned, elevatorConfig.NoOrder, elevatorConfig.Assigned},
}
peerState2 := elevatorConfig.PeerState{
    Behavior:   elevatorConfig.down,
    Floor:       2,
    Direction:   elevatorConfig.Stop,
    CabRequests: [elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus{elevatorConfig.Assigned, elevatorConfig.Assigned, elevatorConfig.Assigned, elevatorConfig.Assigned},
}
peerState3 := elevatorConfig.PeerState{
    Behavior:    elevatorConfig.DoorOpen,
    Floor:        3,
    Direction:   elevatorConfig.down,
    CabRequests: [elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus{elevatorConfig.Assigned, elevatorConfig.Assigned, elevatorConfig.NoOrder, elevatorConfig.NoOrder},
}
peerView := elevatorConfig.PeerView{
    AlivePeers   ["3", "2"]
    OwnId:        "3",
    HallRequests: [elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus{{elevatorConfig.Pending, elevatorConfig.NoOrder}, {elevatorConfig.Assigned, elevatorConfig.NoOrder}, {elevatorConfig.Pending, elevatorConfig.NoOrder}, {elevatorConfig.NoOrder, elevatorConfig.Assigned}},
    States:       map[string]*elevatorConfig.PeerState{"1": &peerState1, "2": &peerState2, "3": &peerState3},
}

PrintElevatorSystem(peerView)

```
#### Output
``` bash
---------Start PeerView update------------

Alive elevators: 2 3
Lost elevators: 1
+----------------------------+
|         ElevatorID: 3      |
+----------------------------+
| Floor        | 3           |
| Direction    | down        |
| Behavior     | moving      |
+----------------------------+
| Floor        | 1  2  3  4  |
+----------------------------+
| Up           | *  !  *  -  |
| Down         | -  -  *  !  |
| Cab          | *  *  -  -  |
+----------------------------+


+----------------------------+
|         ElevatorID: 2      |
+----------------------------+
| Floor        | 2           |
| Direction    | down        |
| Behavior     | moving      |
+----------------------------+
| Floor        | 1  2  3  4  |
+----------------------------+
| Cab          | *  *  *  *  |
+----------------------------+

+----------------------------+
|         ElevatorID: 1      |
+----------------------------+
| Floor        | 1           |
| Direction    | down        |
| Behavior     | doorOpen    |
+----------------------------+
| Floor        | 1  2  3  4  |
+----------------------------+
| Cab          | *  -  *  *  |
+----------------------------+

---------End PeerView update--------------
```