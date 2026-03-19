# Debugging helpers

This module contains functions to ease the debugging process, specifically by providing a way to expose the state of the system through print functions. It contains functions to print the different structs used in the system, as well as some utility functions. There are three public functions, used to print peerUpdates,  the elevator object and print an update of the full peerView object and peer network info.

## Output examples

### PrintLocalElvator(elevator)
### PrintPeerElevatorStates (peerView elevatorConfig.PeerView)
``` bash
--------New peer uppdate recieved----------
Current alive peers: 1 2
Elevator ID  1 just joind the network
Peers considerd lost:
---------End of peer update--------------
```
### PrintElevatorSystem(elevatorSystem elevatorConfig.PeerView) 
``` bash
---------Start System update-----------------
+----------------------------+
|         ElevatorID: 3      |
+----------------------------+
| Floor        | 1           |
| Direction    | stop        |
| Behavior     | idle        |
+----------------------------+
| Floor        | 1  2  3  4  |
+----------------------------+
| Up           | !  !  !  -  |
| Down         | -  -  -  -  |
| Cab          | -  -  -  -  |
+----------------------------+


+----------------------------+
|         ElevatorID: 1      |
+----------------------------+
| Floor        | 1           |
| Direction    | down        |
| Behavior     | moving      |
+----------------------------+
| Floor        | 1  2  3  4  |
+----------------------------+
| Cab          | -  -  -  -  |
+----------------------------+

+----------------------------+
|         ElevatorID: 2      |
+----------------------------+
| Floor        | 0           |
| Direction    | down        |
| Behavior     | doorOpen    |
+----------------------------+
| Floor        | 1  2  3  4  |
+----------------------------+
| Cab          | -  -  -  -  |
+----------------------------+
```