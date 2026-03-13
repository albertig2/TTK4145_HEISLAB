package debuggingHelpers

import (
	"HEISPROSJEKT/elevatorConfig"
	//"HEISPROSJEKT/synchronisation"
	"Network-go/network/peers"
	"fmt"
)

func Elevator_print(es elevatorConfig.Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf("  |%-6s = %-2d          |\n", "floor", es.Floor)
	fmt.Printf("  |%-6s = %-12.12s|\n", "dirn", elevatorConfig.DirectionToString(es.Direction))
	fmt.Printf("  |%-6s = %-12.12s|\n", "behav", elevatorConfig.BehaviorToString(es.Behavior))
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := elevatorConfig.N_FLOORS - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			if (f == elevatorConfig.N_FLOORS-1 && elevatorConfig.Button(btn) == elevatorConfig.HallUp) || (f == 0 && elevatorConfig.Button(btn) == elevatorConfig.HallDown) {
				fmt.Printf("|     ")
			} else {
				if es.Requests[f][btn] {
					fmt.Printf("|  #  ")
				} else {
					fmt.Printf("|  -  ")
				}
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}

func PrintPeerUpdate(peerUpdate peers.PeerUpdate) {

	fmt.Println("---------------------New peer uppdate recieved----------------------- ")

	fmt.Printf("Current alive peers: ")
	for peerIndex := 0; peerIndex < len(peerUpdate.Peers); peerIndex++ {
		fmt.Printf("%+v ", peerUpdate.Peers[peerIndex])
	}
	fmt.Printf("\n")

	if peerUpdate.New != "" {
		fmt.Printf("Elevator ID  %+v just joind the network \n", peerUpdate.New)
	}

	fmt.Printf("Peers considerd lost: ")
	for peerIndex := 0; peerIndex < len(peerUpdate.Lost); peerIndex++ {
		fmt.Printf("%+v ", peerUpdate.Lost[peerIndex])
	}
	fmt.Printf("\n")

	fmt.Println("---------------------End of peer update------------------------------ ")
}

type OrderStatus1 string

const (
	NoOrder   OrderStatus1 = "no order"
	Pending   OrderStatus1 = "pending"
	Assigned  OrderStatus1 = "assigned"
	Completed OrderStatus1 = "completed"
)

type ElevatorState1 struct {
	Behavior    elevatorConfig.Behavior               `json:"behavior"`
	Floor       int                                   `json:"floor"`
	Direction   elevatorConfig.Direction              `json:"direction"`
	CabRequests [elevatorConfig.N_FLOORS]OrderStatus1 `json:"cabRequests"`
}

type ElevatorSystem1 struct {
	OwnId        string                                   `json:"ownId"`
	HallRequests [elevatorConfig.N_FLOORS][2]OrderStatus1 `json:"hallRequests"`
	States       map[string]*ElevatorState1               `json:"states"`
}

func printHallLine(orders [elevatorConfig.N_FLOORS][2]OrderStatus1, orderTypeIndex int) {

	for index := 0; index < len(orders); index++ {
		orderstatus := orders[index][orderTypeIndex]
		switch orderstatus {
		case NoOrder:
			fmt.Printf(" - ")
		case Pending:
			fmt.Printf(" ! ")
		case Assigned:
			fmt.Printf(" * ")
		case Completed:
			fmt.Printf(" ^ ")
		}

	}
	fmt.Printf(" |\n")
}
func printCabLine(orders []OrderStatus1) {

	for index := 0; index < len(orders); index++ {
		orderstatus := orders[index]
		switch orderstatus {
		case NoOrder:
			fmt.Printf(" - ")
		case Pending:
			fmt.Printf(" ! ")
		case Assigned:
			fmt.Printf(" * ")
		case Completed:
			fmt.Printf(" ^ ")
		}

	}
	fmt.Printf(" |\n")
}
func PrintCurrentWorkingElevators(elevatorSystem ElevatorSystem1) {
	workingNode := elevatorSystem.States[elevatorSystem.OwnId]
	hallOrders := elevatorSystem.HallRequests
	cabOrders := workingNode.CabRequests

	// const boxWidth = 28
	// leftCol := 12
	// rightCol := boxWidth - leftCol - 4

	//divider := "+----------------------------+"
	fmt.Println("+----------------------------+")
	//fmt.Printf("| %-26s |\n", fmt.Sprintf("ElevatorID: %v", elevatorSystem.OwnId))
	fmt.Printf("|         ElevatorID: %v      |\n", elevatorSystem.OwnId)
	fmt.Println("+----------------------------+")
	fmt.Printf("| %-12s | %-12d|\n", "Floor", workingNode.Floor)
	fmt.Printf("| %-12s | %-12s|\n", "Direction", elevatorConfig.DirectionToString(workingNode.Direction))
	fmt.Printf("| %-12s | %-12s|\n", "Behavior", elevatorConfig.BehaviorToString(workingNode.Behavior))
	fmt.Println("+----------------------------+")
	fmt.Printf("| %-12s | %-12s|\n", "Floor", "1  2  3  4" )
	fmt.Printf("| %-12s |", "Up")
	printHallLine(hallOrders, 0)
	fmt.Printf("| %-12s |", "Down")
	printHallLine(hallOrders, 1)
	fmt.Printf("| %-12s |", "Cab")
	printCabLine(cabOrders[:])
	fmt.Println("+----------------------------+")
	//fmt.Println(divider)

	fmt.Printf("")

}
func printPeerElevators() {

}

// func PrintElevatorSystem(elevatorSystem synchronisation.ElevatorSystem){

// 	fmt.Println("---------Start System update-----------------")

// 	fmt.Println("---------End System update-------------------")

// }
/*

System worldview

noOrder: -
pending: !
Assigned: *
Completed: ^


+------------------+
|ElevatorID = 1    |
| Floor = 2        |
| Direction = down |
| Behavior = idle  |
+------------------+
|1     | dn |  up  |
|2     |    |      |
|3     |    |      |
|4     |    |      |
+------------------+


+---------------------------+
|     ElevatorID  1         |
+---------------------------+
| Floor       | 2           |
| Direction   | down        |
| Behavior    | Moving      |
+---------------------------+
| Floor       | 1  2  3  4  |
+---------------------------+
| Hall Up     | -  !  -  *  |
| Hall Down   | -  !  -  *  |
| Cab         | -  !  -  *  |
+---------------------------+

+----------------------------+
|     ElevatorID  1          |
+----------------------------+
| Floor       | 2            |
| Direction   | down         |
| Behavior    | Moving       |
+----------------------------+
| Floor       |  1  2  3  4  |
+----------------------------+
| Cab         |  -  !  -  *  |
+----------------------------+

+---------------------------+
|     ElevatorID  1         |
+---------------------------+
| Floor       | 2           |
| Direction   | down        |
| Behavior    | Moving      |
+---------------------------+
| Floor       | 1  2  3  4  |
+---------------------------+
| Cab         | -  !  -  *  |
+---------------------------+

compressed version:

ID: 2
Floor: 2, Dir: down, Behav: Idle

Hall UP :
Hall Down:
Cab

---------System update----------------

ElevatorID: 1 (Current Node)
+----------------------------+
|  Idle       |    #        |
| Floor       | 1  2  3  4  |
+---------------------------+
| Hall Up     | -  !  -  *  |
| Hall Down   | -  !  -  *  |
| Cab         | -  !  -  *  |
+---------------------------+

ElevatorID: 2
+----------------------------+
| Moving      |    #>       |
| Floor       | 1  2  3  4  |
+---------------------------+
| Hall Up     | -  !  -  *  |
| Hall Down   | -  !  -  *  |
| Cab         | -  !  -  *  |
+---------------------------+

ElevatorID: 3
+----------------------------+
| DoorOpen    |    #        |
| Floor       | 1  2  3  4  |
+---------------------------+
| Hall Up     | -  !  -  *  |
| Hall Down   | -  !  -  *  |
| Cab         | -  !  -  *  |
+---------------------------+

---------System update-----------------
*/

//these functions are just for testing of the fsm sync, delete when done

func InitializeOrderChannels() elevatorConfig.ElevatorOrderChannelStruckt {

	orderChannelse := elevatorConfig.ElevatorOrderChannelStruckt{

		NewRecievedOrderChannel: make(chan elevatorConfig.ButtonEvent),
		NewAssignedOrderChannel: make(chan elevatorConfig.ButtonEvent),
		ServicedOrderChannel:    make(chan elevatorConfig.ButtonEvent),
	}

	return orderChannelse
}

func MimicOrderAssignerAndSynch(orderChannelse elevatorConfig.ElevatorOrderChannelStruckt) {

	for {
		select {
		case newOrder := <-orderChannelse.NewRecievedOrderChannel:
			fmt.Println("New order Recieved: ", newOrder)
			orderChannelse.NewAssignedOrderChannel <- newOrder

		case servicedOrder := <-orderChannelse.ServicedOrderChannel:
			fmt.Println("New order serviced: ", servicedOrder)

		}

	}

}
