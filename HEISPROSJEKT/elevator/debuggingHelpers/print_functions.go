package debuggingHelpers

import (
	"HEISPROSJEKT/elevatorConfig"
	"Network-go/network/peers"
	"fmt"
)

func PrintLocalElvator(elevator elevatorConfig.Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf("  |%-6s = %-2d          |\n", "floor", elevator.Floor)
	fmt.Printf("  |%-6s = %-12.12s|\n", "dirn", elevatorConfig.DirectionToString(elevator.Direction))
	fmt.Printf("  |%-6s = %-12.12s|\n", "behav", elevatorConfig.BehaviorToString(elevator.Behavior))
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := elevatorConfig.N_FLOORS - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < elevatorConfig.N_BUTTONS; btn++ {
			if (f == elevatorConfig.N_FLOORS-1 && elevatorConfig.Button(btn) == elevatorConfig.HallUp) || (f == 0 && elevatorConfig.Button(btn) == elevatorConfig.HallDown) {
				fmt.Printf("|     ")
			} else {
				if elevator.LocalOrderQueue[f][btn] {
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

func orderStatusToSymbolToString(orderstatus elevatorConfig.OrderStatus) string {
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

func printHallOrderLine(orders [elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, orderTypeIndex int) {
	for index := 0; index < len(orders); index++ {
		orderstatus := orders[index][orderTypeIndex]
		fmt.Printf("%s", orderStatusToSymbolToString(orderstatus))
	}
	fmt.Printf(" |\n")
}

func printCabOrderLine(orders [elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus) {
	for index := 0; index < len(orders); index++ {
		orderstatus := orders[index]
		fmt.Printf("%s", orderStatusToSymbolToString(orderstatus))
	}
	fmt.Printf(" |\n")
}
func printCurrentWorkingElevatorFromPeerView(peerView elevatorConfig.PeerView) {
	workingNode := peerView.States[peerView.OwnId]
	hallOrders := peerView.HallRequests
	cabOrders := workingNode.CabRequests

	fmt.Println("+----------------------------+")
	fmt.Printf("|         ElevatorID: %v      |\n", peerView.OwnId)
	fmt.Println("+----------------------------+")
	fmt.Printf("| %-12s | %-12d|\n", "Floor", workingNode.Floor)
	fmt.Printf("| %-12s | %-12s|\n", "Direction", elevatorConfig.DirectionToString(workingNode.Direction))
	fmt.Printf("| %-12s | %-12s|\n", "Behavior", elevatorConfig.BehaviorToString(workingNode.Behavior))
	fmt.Println("+----------------------------+")
	fmt.Printf("| %-12s | %-12s|\n", "Floor", "1  2  3  4")
	fmt.Println("+----------------------------+")
	fmt.Printf("| %-12s |", "Up")
	printHallOrderLine(hallOrders, 0)
	fmt.Printf("| %-12s |", "Down")
	printHallOrderLine(hallOrders, 1)
	fmt.Printf("| %-12s |", "Cab")
	printCabOrderLine(cabOrders)
	fmt.Println("+----------------------------+")

}

func printPeerElevatorStates(peerView elevatorConfig.PeerView) {
	currentPeers := peerView.States

	for id, state := range currentPeers {
		if id != peerView.OwnId {
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
			printCabOrderLine(cabOrders)
			fmt.Println("+----------------------------+")
		}
		fmt.Printf("\n")

	}
}

func printDeadAndAliveElevators(peerView elevatorConfig.PeerView) {

	fmt.Printf("Alive elevators: ")
	for _, peerId := range peerView.AlivePeers {
		fmt.Printf("%s ", peerId)
	}
	fmt.Printf("\n")
	fmt.Printf("Lost elevators: ")

	lostOwnId := true
	for peerId := range peerView.States {
		idLost := true
		for _, aliveId := range peerView.AlivePeers {
			if aliveId == peerId {
				idLost = false
			}
			if aliveId == peerView.OwnId {
				lostOwnId = false
			}
		}
		if idLost {
			fmt.Printf("%s ", peerId)
		}
	}
	fmt.Printf("\n")
	if lostOwnId {
		fmt.Println("This elevator is currently offline")
	}
}

func PrintPeerViewUpdate(peerView elevatorConfig.PeerView) {
	fmt.Println("---------Start PeerView update------------")
	printDeadAndAliveElevators(peerView)
	printCurrentWorkingElevatorFromPeerView(peerView)
	fmt.Printf("\n")
	printPeerElevatorStates(peerView)
	fmt.Println("---------End PeerView update--------------")
}
