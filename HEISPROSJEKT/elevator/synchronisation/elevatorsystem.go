package synchronisation

import (
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
)

//"flag"

// flag.Int("id", 1, "Input id")

var HallDirections = [2]int{int(elevatorConfig.HallUp), int(elevatorConfig.HallDown)}

func SetBehavior(system *elevatorConfig.ElevatorSystem, b elevatorConfig.Behavior) {
	state := system.States[system.OwnId]
	state.Behavior = b
}

func SetFloor(system *elevatorConfig.ElevatorSystem, f int) {
	state := system.States[system.OwnId]
	state.Floor = f
}

func SetDirection(system *elevatorConfig.ElevatorSystem, dir elevatorConfig.Direction) {
	state := system.States[system.OwnId]
	state.Direction = dir
}

// Usikker på om jeg skal kalle det on eller off? eller en funksjon for på og en for av
func SetCabRequests(system *elevatorConfig.ElevatorSystem, f int, orderstatus elevatorConfig.OrderStatus) {
	state := system.States[system.OwnId]
	state.CabRequests[f] = orderstatus
}

// Usikker på om jeg skal kalle den up or down eller dont know
func SetHallRequests(system *elevatorConfig.ElevatorSystem, f int, halldir int, orderstatus elevatorConfig.OrderStatus) {
	system.HallRequests[f][halldir] = orderstatus
}

func InitializeElevatorSystem(system *elevatorConfig.ElevatorSystem, id string) {
	// To decide floor can just do the get_floor_sensor_signal() and initialize to that floor, but for now hardcoded
	system.OwnId = id
	system.HallRequests = [elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus{}
	system.States = make(map[string]*elevatorConfig.ElevatorState)
	currentFloor := 1 // Get floor sensor. (men helst ikke -1? så siste faktisk floor)
	system.States[id] = &elevatorConfig.ElevatorState{
		Behavior:    elevatorConfig.Idle,
		Floor:       currentFloor,
		Direction:   elevatorConfig.Stop,
		CabRequests: [elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus{},
	}

	initializeHallRequests(system)

	//timer

	//listen for broadcast
	//if recieved -> stop timer
	// else timeout AND RUN AS SINGLE ELEVATOR
	//<-timer.C

	initializeCabRequests(system)
}

func initializeHallRequests(system *elevatorConfig.ElevatorSystem) {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for _, halldir := range HallDirections {
			system.HallRequests[floor][halldir] = elevatorConfig.NoOrder
		}
	}
}

func initializeCabRequests(system *elevatorConfig.ElevatorSystem) {
	// Listen for other elevators to broadcast their view of your cab orders, and if you hear any, set your cab orders to be the combination of all of them (pending if any of them is pending or no order)
	// For each elevator you hear from, check all floors, if any of the floors have pending, set that floor to pending.
	/*
		for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
			if system.States[system.OwnId].CabRequests[floor] != Pending {
				system.States[system.OwnId].CabRequests[floor] = NoOrder
			}
		}*/
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		system.States[system.OwnId].CabRequests[floor] = elevatorConfig.Unknown
	}
}

// If only called from updatedElavatorSystemFromPeer, then I dont need the check for existence
// Only called if not existing
func addPeer(system *elevatorConfig.ElevatorSystem, peerSystem *elevatorConfig.ElevatorSystem) {
	CabRequests := peerSystem.States[peerSystem.OwnId].CabRequests
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		if CabRequests[floor] == elevatorConfig.Unknown {
			CabRequests[floor] = elevatorConfig.NoOrder
		}
	}

	peerState := peerSystem.States[peerSystem.OwnId]
	system.States[peerSystem.OwnId] = &elevatorConfig.ElevatorState{
		Behavior:    peerState.Behavior,
		Floor:       peerState.Floor,
		Direction:   peerState.Direction,
		CabRequests: CabRequests,
	}
}

func updatePeer(system *elevatorConfig.ElevatorSystem, peerSystem *elevatorConfig.ElevatorSystem) {
	// Need to make sure that if I have info on the cabrequests of a peer and they are restarted, meaning that their caborders may be uninitialized
	peerSystemCabRequests := peerSystem.States[peerSystem.OwnId].CabRequests
	CabRequests := system.States[peerSystem.OwnId].CabRequests

	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		if peerSystemCabRequests[floor] != elevatorConfig.Unknown {
			CabRequests[floor] = peerSystemCabRequests[floor]
		}
	}

	peerState := peerSystem.States[peerSystem.OwnId]
	system.States[peerSystem.OwnId] = &elevatorConfig.ElevatorState{
		Behavior:    peerState.Behavior,
		Floor:       peerState.Floor,
		Direction:   peerState.Direction,
		CabRequests: CabRequests,
	}
}

func UpdateElevatorSystemWithPeer(system *elevatorConfig.ElevatorSystem, peerSystem *elevatorConfig.ElevatorSystem, HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus) {
	if _, exists := system.States[peerSystem.OwnId]; exists {
		updatePeer(system, peerSystem)
	} else {
		addPeer(system, peerSystem)
	}

	HallRequestsForAllElevators[peerSystem.OwnId] = peerSystem.HallRequests
	if _, exists := peerSystem.States[system.OwnId]; exists {
		CabRequestsForAllElevators[peerSystem.OwnId] = peerSystem.States[system.OwnId].CabRequests
	}
}

func UpdateElevatorSystemWithSelf(system *elevatorConfig.ElevatorSystem, HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus) {
	HallRequestsForAllElevators[system.OwnId] = system.HallRequests
	CabRequestsForAllElevators[system.OwnId] = system.States[system.OwnId].CabRequests
}

func updateElevatorSystem(system *elevatorConfig.ElevatorSystem, hallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, cabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus, receivedWorldView chan string) {
	peerSystem := DecodeElevatorSystem(<-receivedWorldView)
	addPeer(system, &peerSystem)
	UpdateElevatorSystemWithPeer(system, &peerSystem, hallRequestsForAllElevators, cabRequestsForAllElevators)
	UpdateElevatorSystemWithSelf(system, hallRequestsForAllElevators, cabRequestsForAllElevators)
}

// I utgangspunktet har jeg en annen funksjon som fikser andre transisjoner ...., kanskje nok å sette HallRequest listen to the appropriate, og så finne derfifra hva man skal sette
// Hvor ofte skal man sjekke transisjoner? med en gang etter man har updated

// Funksjon for transisjoner mellom states (når man skal gå fra en state til en annen)
// Denne burde si hvilke transisjoner som skal gjøres (en pure function)
// Og så burde man ha en funskjon som utfører transjosjonen med påfølgende handlinger

// Må deale med transisjoner, når man skal sette pending? når man skal gå til de andre? Skal man gjøre det når man får inn fra andre (hvertfall pending?)
// Når transisjon så man kansje gjøre ting også så jeg har jo en pure en
// Må på et tidspunkt oppdatere HallRequest med egen id sin hallrequests også og cabRequests.

// SHoul maybe chnage hallup and halldown inputs to be button type instead of ints and then change inside the functions instead

func printHallLine(orders [elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, orderTypeIndex int) {

	for index := 0; index < len(orders); index++ {
		orderstatus := orders[index][orderTypeIndex]
		switch orderstatus {
		case elevatorConfig.NoOrder:
			fmt.Printf(" - ")
		case elevatorConfig.Pending:
			fmt.Printf(" ! ")
		case elevatorConfig.Assigned:
			fmt.Printf(" * ")
		case elevatorConfig.Completed:
			fmt.Printf(" ^ ")
		}

	}
	fmt.Printf(" |\n")
}
func printCabLine(orders []elevatorConfig.OrderStatus) {

	for index := 0; index < len(orders); index++ {
		orderstatus := orders[index]
		switch orderstatus {
		case elevatorConfig.NoOrder:
			fmt.Printf(" - ")
		case elevatorConfig.Pending:
			fmt.Printf(" ! ")
		case elevatorConfig.Assigned:
			fmt.Printf(" * ")
		case elevatorConfig.Completed:
			fmt.Printf(" ^ ")
		}

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
	printCabLine(cabOrders[:])
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
			printCabLine(cabOrders[:])
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
