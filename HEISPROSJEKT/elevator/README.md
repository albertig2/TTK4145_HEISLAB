# The elevator project

This is group 51’s (Albertine Gjøs, Odin Sandli Mellingsæter and Mathilde Skaset-Haarr) solution to the distributed elevator project for spring 2026 (TTK4155). The project is solved using UDP with a peer to peer structure. Information is shared between nodes by regular broadcasts ( 30 hz update frequency). Each broadcast contains the state of the local node (direction, behavior and floor and alivelist) the order queue viewed from the nodes perspective, and the nodes cab orders. To be able to deal with order assignments, reassignment and prevent loss of cab orders, the orders are managed as a state machine. Together with the synchronization module, this limits double assignment, assures automatic reassignment of the orders after a disconnection and the restoration of cab orders after reconnection . The local elevator and its hardware are managed by another finite state machine. The modules communicate through channels. 
## Run the project
The project is run by calling:
```bash
go run main.go -id=  yourId -port= yourPort
```
In the terminal. Note that you must enter this from the correct folder. 
Both the id and port should be entered as an integer. If no id is provided the elevator takes the default value of id=1. If no port is provided, the port is set to its default value of port= 15657. The port specified in the terminal, is the TCP port used to connect to the elevator server or the simulator server.
#### `Example: Connect to the elevatorserver on port 15777, with id 5:`
```bash
go run main.go -id=5 -port=15777
```
## Global variables and configuration
All global variables can be found in the elevator_config folder. This includes the broadcast and peerUpdate ports, the number of floors and buttons, and the duration of the door open timer and motor failure detection timer. The file also contains all global structs and enums. Anny changes to these variables will affect the entire code, making changing the number of floors or the door open time easy, pain free and bug resistant. 
