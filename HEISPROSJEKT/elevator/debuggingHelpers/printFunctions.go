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

	fmt.Println("--------New peer uppdate recieved----------")

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

	fmt.Println("---------End of peer update-------------- ")
}

// switch orderstatus {
// case elevatorConfig.NoOrder:
// 	fmt.Printf(" - ")
// case elevatorConfig.Pending:
// 	fmt.Printf(" ! ")
// case elevatorConfig.Assigned:
// 	fmt.Printf(" * ")
// case elevatorConfig.Completed:
// 	fmt.Printf(" ^ ")
// }

func OrderstatusToSymbol(orderstatus elevatorConfig.OrderStatus) string {
	switch orderstatus {
	case elevatorConfig.NoOrder:
		return " - "
	case elevatorConfig.Pending:
		return " ! "
	case elevatorConfig.Assigned:
		return " * "
	case elevatorConfig.Completed:
		return " ^ "
	case elevatorConfig.Unknown:
		return " ? "
	default:
		
		return "Invalid order staus, recieved: " + string(orderstatus)
	}

}

func printHallLine(orders [elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, orderTypeIndex int) {

	for index := 0; index < len(orders); index++ {
		orderstatus := orders[index][orderTypeIndex]
		fmt.Printf("%s", OrderstatusToSymbol(orderstatus))
	}
	fmt.Printf(" |\n")
}
func printCabLine(orders [elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus) {

	for index := 0; index < len(orders); index++ {
		orderstatus := orders[index]
		fmt.Printf("%s", OrderstatusToSymbol(orderstatus))

	}
	fmt.Printf(" |\n")
}
func PrintCurrentWorkingElevators(elevatorSystem elevatorConfig.ElevatorSystem) {
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
	fmt.Printf("| %-12s | %-12s|\n", "Floor", "1  2  3  4")
	fmt.Println("+----------------------------+")
	fmt.Printf("| %-12s |", "Up")
	printHallLine(hallOrders, 0)
	fmt.Printf("| %-12s |", "Down")
	printHallLine(hallOrders, 1)
	fmt.Printf("| %-12s |", "Cab")
	printCabLine(cabOrders)
	fmt.Println("+----------------------------+")
	//fmt.Println(divider)

	fmt.Printf("")

}
func PrintPeerElevators(elevatorSystem elevatorConfig.ElevatorSystem) {
	currentPeers := elevatorSystem.States

	for id, state := range currentPeers {
		if id != elevatorSystem.OwnId {
			cabOrders := state.CabRequests

			fmt.Println("+----------------------------+")
			fmt.Printf("|         ElevatorID: %v      |\n", id)
			fmt.Println("+----------------------------+")
			fmt.Printf("| %-12s | %-12d|\n", "Floor", state.Floor)
			fmt.Printf("| %-12s | %-12s|\n", "Direction", elevatorConfig.DirectionToString(state.Direction))
			fmt.Printf("| %-12s | %-12s|\n", "Behavior", elevatorConfig.BehaviorToString(state.Behavior))
			fmt.Println("+----------------------------+")
			fmt.Printf("| %-12s | %-12s|\n", "Floor", "1  2  3  4")
			fmt.Println("+----------------------------+")
			fmt.Printf("| %-12s |", "Cab")
			printCabLine(cabOrders)
			fmt.Println("+----------------------------+")
		}
		fmt.Printf("\n")

	}

}

func PrintElevatorSystem(elevatorSystem elevatorConfig.ElevatorSystem) {

	fmt.Println("---------Start System update-----------------")

	PrintCurrentWorkingElevators(elevatorSystem)
	fmt.Printf("\n")
	PrintPeerElevators(elevatorSystem)

	fmt.Println("---------End System update-------------------")

}

func TestPrintEElevatorSystem() {

	elevatorStatetest := elevatorConfig.ElevatorState{
		Behavior:    elevatorConfig.Idle,
		Direction:   elevatorConfig.Stop,
		CabRequests: [elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus{elevatorConfig.Pending, elevatorConfig.Assigned, elevatorConfig.NoOrder, elevatorConfig.Completed},
	}
	elevatorStatetest2 := elevatorConfig.ElevatorState{
		Behavior:    elevatorConfig.Idle,
		Direction:   elevatorConfig.Stop,
		CabRequests: [elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus{elevatorConfig.Pending, elevatorConfig.Assigned, elevatorConfig.NoOrder, elevatorConfig.Completed},
	}
	elevatorStatetest3 := elevatorConfig.ElevatorState{
		Behavior:    elevatorConfig.Idle,
		Direction:   elevatorConfig.Stop,
		CabRequests: [elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus{elevatorConfig.Pending, elevatorConfig.Assigned, elevatorConfig.NoOrder, elevatorConfig.Completed},
	}
	elevatortest := elevatorConfig.ElevatorSystem{
		OwnId:        "1",
		HallRequests: [elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus{{elevatorConfig.Pending, elevatorConfig.NoOrder}, {elevatorConfig.Completed, elevatorConfig.NoOrder}, {elevatorConfig.Assigned, elevatorConfig.NoOrder}, {elevatorConfig.Assigned, elevatorConfig.NoOrder}},
		States:       map[string]*elevatorConfig.ElevatorState{"1": &elevatorStatetest, "2": &elevatorStatetest2, "3": &elevatorStatetest3},
	}
	// elevatorConfig.PrintCurrentWorkingElevators(elevatortest)
	// elevatorConfig.PrintPeerElevators(elevatortest)
	PrintElevatorSystem(elevatortest)

}

// func printCurrentWorkingElevators(elevatorSystem ElevatorSystem ){
// 	workingNode := elevatorSystem.States[elevatorSystem.OwnId]
// 	// hallOrders := elevatorSystem.HallRequests
// 	// cabOrders := workingNode.CabRequests

// 	//divider := "+----------------------------+"
// 	fmt.Println(  "+----------------------------+" )
// 	fmt.Printf(   "|%-6s %v                        |\n", "ElevatorID: ", elevatorSystem.OwnId)
// 	fmt.Println(  "+----------------------------+" )
// 	fmt.Printf(   "|%-6s         | %-2d         |\n", "floor", workingNode.Floor)
// 	fmt.Printf(   "|%-6s         | %-12.12s     |\n", "Direction", elevatorConfig.DirectionToString(workingNode.Direction))
// 	fmt.Printf(   "|%-6s         | %-12.12s     |\n", "Behavior", elevatorConfig.BehaviorToString(workingNode.Behavior))
// 	fmt.Println(  "+----------------------------+" )
// 	fmt.Printf(  "| %-6s        |  1  2  3  4  |\n", "Floor")
// 	fmt.Printf (  "| %-6s        |", "Up" )
// 	fmt.Printf (  "| %-6s        |", "Down" )
// 	fmt.Printf (  "| %-6s        |", "Cab" )
// 	fmt.Println(  "+----------------------------+" )

// 	//fmt.Println(divider)

// 	fmt.Printf("")

// }
// func printPeerElevators(){

// }

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
/*
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
	*/
