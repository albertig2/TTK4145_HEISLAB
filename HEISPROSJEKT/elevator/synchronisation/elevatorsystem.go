package synchronisation

import (
	"HEISPROSJEKT/elevatorConfig"
)

//"flag"

// flag.Int("id", 1, "Input id")
type OrderStatus string

// Cab order only needs to go from no order to Pending to Completed, while hall orders also need Assigned, since they are Assigned to an elevator by the assigner
const (
	Unknown   OrderStatus = "unknown"
	NoOrder   OrderStatus = "no order"
	Pending   OrderStatus = "pending"
	Assigned  OrderStatus = "assigned"
	Completed OrderStatus = "completed"
)

var HallDirections = [2]int{int(elevatorConfig.HallUp), int(elevatorConfig.HallDown)}

type ElevatorState struct {
	Behavior    elevatorConfig.Behavior              `json:"behavior"`
	Floor       int                                  `json:"floor"`
	Direction   elevatorConfig.Direction             `json:"direction"`
	CabRequests [elevatorConfig.N_FLOORS]OrderStatus `json:"cabRequests"`
}

type ElevatorSystem struct {
	OwnId        string                                  `json:"ownId"`
	HallRequests [elevatorConfig.N_FLOORS][2]OrderStatus `json:"hallRequests"`
	States       map[string]*ElevatorState               `json:"states"`
}

func SetBehavior(system *ElevatorSystem, b elevatorConfig.Behavior) {
	state := system.States[system.OwnId]
	state.Behavior = b
}

func SetFloor(system *ElevatorSystem, f int) {
	state := system.States[system.OwnId]
	state.Floor = f
}

func SetDirection(system *ElevatorSystem, dir elevatorConfig.Direction) {
	state := system.States[system.OwnId]
	state.Direction = dir
}

// Usikker på om jeg skal kalle det on eller off? eller en funksjon for på og en for av
func SetCabRequests(system *ElevatorSystem, f int, orderstatus OrderStatus) {
	state := system.States[system.OwnId]
	state.CabRequests[f] = orderstatus
}

// Usikker på om jeg skal kalle den up or down eller dont know
func SetHallRequests(system *ElevatorSystem, f int, halldir int, orderstatus OrderStatus) {
	system.HallRequests[f][halldir] = orderstatus
}

func InitializeElevatorSystem(system *ElevatorSystem, id string) {
	// To decide floor can just do the get_floor_sensor_signal() and initialize to that floor, but for now hardcoded
	system.OwnId = id
	system.HallRequests = [elevatorConfig.N_FLOORS][2]OrderStatus{}
	system.States = make(map[string]*ElevatorState)
	currentFloor := 1 // Get floor sensor. (men helst ikke -1? så siste faktisk floor)
	system.States[id] = &ElevatorState{
		Behavior:    elevatorConfig.Idle,
		Floor:       currentFloor,
		Direction:   elevatorConfig.Stop,
		CabRequests: [elevatorConfig.N_FLOORS]OrderStatus{},
	}

	initializeHallRequests(system)

	//timer

	//listen for broadcast
	//if recieved -> stop timer
	// else timeout AND RUN AS SINGLE ELEVATOR
	//<-timer.C

	initializeCabRequests(system)
}

func initializeHallRequests(system *ElevatorSystem) {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for _, halldir := range HallDirections {
			system.HallRequests[floor][halldir] = NoOrder
		}
	}
}

func initializeCabRequests(system *ElevatorSystem) {
	// Listen for other elevators to broadcast their view of your cab orders, and if you hear any, set your cab orders to be the combination of all of them (pending if any of them is pending or no order)
	// For each elevator you hear from, check all floors, if any of the floors have pending, set that floor to pending.
	/*
		for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
			if system.States[system.OwnId].CabRequests[floor] != Pending {
				system.States[system.OwnId].CabRequests[floor] = NoOrder
			}
		}*/
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		system.States[system.OwnId].CabRequests[floor] = Unknown
	}
}

// If only called from updatedElavatorSystemFromPeer, then I dont need the check for existence
// Only called if not existing
func addPeer(system *ElevatorSystem, peerSystem *ElevatorSystem) {
	CabRequests := peerSystem.States[peerSystem.OwnId].CabRequests
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		if CabRequests[floor] == Unknown {
			CabRequests[floor] = NoOrder
		}
	}

	peerState := peerSystem.States[peerSystem.OwnId]
	system.States[peerSystem.OwnId] = &ElevatorState{
		Behavior:    peerState.Behavior,
		Floor:       peerState.Floor,
		Direction:   peerState.Direction,
		CabRequests: CabRequests,
	}
}

func updatePeer(system *ElevatorSystem, peerSystem *ElevatorSystem) {
	// Need to make sure that if I have info on the cabrequests of a peer and they are restarted, meaning that their caborders may be uninitialized
	peerSystemCabRequests := peerSystem.States[peerSystem.OwnId].CabRequests
	CabRequests := system.States[peerSystem.OwnId].CabRequests

	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		if peerSystemCabRequests[floor] != Unknown {
			CabRequests[floor] = peerSystemCabRequests[floor]
		}
	}

	peerState := peerSystem.States[peerSystem.OwnId]
	system.States[peerSystem.OwnId] = &ElevatorState{
		Behavior:    peerState.Behavior,
		Floor:       peerState.Floor,
		Direction:   peerState.Direction,
		CabRequests: CabRequests,
	}
}

func UpdateElevatorSystemWithPeer(system *ElevatorSystem, peerSystem *ElevatorSystem, HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]OrderStatus, CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]OrderStatus) {
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

func UpdateElevatorSystemWithSelf(system *ElevatorSystem, HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]OrderStatus, CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]OrderStatus) {
	HallRequestsForAllElevators[system.OwnId] = system.HallRequests
	CabRequestsForAllElevators[system.OwnId] = system.States[system.OwnId].CabRequests
}

func updateElevatorSystem(system *ElevatorSystem, hallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]OrderStatus, cabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]OrderStatus, receivedWorldView chan string) {
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
