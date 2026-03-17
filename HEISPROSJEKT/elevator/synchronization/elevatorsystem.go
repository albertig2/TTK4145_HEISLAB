package synchronization

import (
	"HEISPROSJEKT/elevatorConfig"
	//"fmt"
)

//"flag"

// flag.Int("id", 1, "Input id")

var HallDirections = [2]int{int(elevatorConfig.HallUp), int(elevatorConfig.HallDown)}

func SetAlivePeers(peerView *elevatorConfig.PeerView, alivePeers []string) {
	peerView.AlivePeers = alivePeers
}

func SetBehavior(peerView *elevatorConfig.PeerView, behavior elevatorConfig.Behavior) {
	state := peerView.States[peerView.OwnId]
	state.Behavior = behavior
}

func SetFloor(peerView *elevatorConfig.PeerView, floor int) {
	state := peerView.States[peerView.OwnId]
	state.Floor = floor
}

func SetDirection(peerView *elevatorConfig.PeerView, direction elevatorConfig.Direction) {
	state := peerView.States[peerView.OwnId]
	state.Direction = direction
}

// Usikker på om jeg skal kalle det on eller off? eller en funksjon for på og en for av
func SetCabRequests(peerView *elevatorConfig.PeerView, floor int, orderStatus elevatorConfig.OrderStatus) {
	state := peerView.States[peerView.OwnId]
	state.CabRequests[floor] = orderStatus
}

// Usikker på om jeg skal kalle den up or down eller dont know
func SetHallRequests(peerView *elevatorConfig.PeerView, floor int, hallDirection int, orderStatus elevatorConfig.OrderStatus) {
	peerView.HallRequests[floor][hallDirection] = orderStatus
}

func InitializeElevatorSystem(peerView *elevatorConfig.PeerView, ownId string) {
	// To decide floor can just do the get_floor_sensor_signal() and initialize to that floor, but for now hardcoded
	peerView.AlivePeers = []string{ownId}
	peerView.OwnId = ownId
	peerView.HallRequests = [elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus{}
	peerView.States = make(map[string]*elevatorConfig.PeerState)
	unknownFloor := -1
	peerView.States[ownId] = &elevatorConfig.PeerState{
		Behavior:    elevatorConfig.Idle,
		Floor:       unknownFloor, //Just initializing to this, will be updated by the elevator fsm if not correct
		Direction:   elevatorConfig.Stop,
		CabRequests: [elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus{},
	}

	initializeHallRequests(peerView)

	//timer

	//listen for broadcast
	//if recieved -> stop timer
	// else timeout AND RUN AS SINGLE ELEVATOR
	//<-timer.C

	initializeCabRequests(peerView)
}

func initializeHallRequests(peerView *elevatorConfig.PeerView) {
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		for _, hallDirection := range HallDirections {
			peerView.HallRequests[floor][hallDirection] = elevatorConfig.NoOrder
		}
	}
}

func initializeCabRequests(peerView *elevatorConfig.PeerView) {
	// Listen for other elevators to broadcast their view of your cab orders, and if you hear any, set your cab orders to be the combination of all of them (pending if any of them is pending or no order)
	// For each elevator you hear from, check all floors, if any of the floors have pending, set that floor to pending.
	/*
		for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
			if system.States[system.OwnId].CabRequests[floor] != Pending {
				system.States[system.OwnId].CabRequests[floor] = NoOrder
			}
		}*/
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		peerView.States[peerView.OwnId].CabRequests[floor] = elevatorConfig.Unknown
	}
}

// If only called from updatedElavatorSystemFromPeer, then I dont need the check for existence
// Only called if not existing
func Contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func addPeer(peerView *elevatorConfig.PeerView, externalPeerView *elevatorConfig.PeerView) {
	if !Contains(peerView.AlivePeers, externalPeerView.OwnId) {
		peerView.AlivePeers = append(peerView.AlivePeers, externalPeerView.OwnId)
	}
	CabRequests := externalPeerView.States[externalPeerView.OwnId].CabRequests
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		if CabRequests[floor] == elevatorConfig.Unknown {
			CabRequests[floor] = elevatorConfig.NoOrder
		}
	}

	peerState := externalPeerView.States[externalPeerView.OwnId]
	peerView.States[externalPeerView.OwnId] = &elevatorConfig.PeerState{
		Behavior:    peerState.Behavior,
		Floor:       peerState.Floor,
		Direction:   peerState.Direction,
		CabRequests: CabRequests,
	}
}

func updatePeer(localPeerView *elevatorConfig.PeerView, extarnalPeerView *elevatorConfig.PeerView) {
	// Need to make sure that if I have info on the cabrequests of a peer and they are restarted, meaning that their caborders may be uninitialized
	peerSystemCabRequests := extarnalPeerView.States[extarnalPeerView.OwnId].CabRequests
	CabRequests := localPeerView.States[extarnalPeerView.OwnId].CabRequests

	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		if peerSystemCabRequests[floor] != elevatorConfig.Unknown {
			CabRequests[floor] = peerSystemCabRequests[floor]
		}
	}

	peerState := extarnalPeerView.States[extarnalPeerView.OwnId]
	localPeerView.States[extarnalPeerView.OwnId] = &elevatorConfig.PeerState{
		Behavior:    peerState.Behavior,
		Floor:       peerState.Floor,
		Direction:   peerState.Direction,
		CabRequests: CabRequests,
	}
}

func UpdateElevatorSystemWithPeer(localPeerView *elevatorConfig.PeerView, extarnalPeerView *elevatorConfig.PeerView, hallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, cabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus) {
	if _, exists := localPeerView.States[extarnalPeerView.OwnId]; exists {
		updatePeer(localPeerView, extarnalPeerView)
	} else {
		addPeer(localPeerView, extarnalPeerView)
	}

	hallRequestsForAllElevators[extarnalPeerView.OwnId] = extarnalPeerView.HallRequests
	if _, exists := extarnalPeerView.States[localPeerView.OwnId]; exists {
		cabRequestsForAllElevators[extarnalPeerView.OwnId] = extarnalPeerView.States[localPeerView.OwnId].CabRequests
	}
}
/*
func UpdateElevatorSystemWithSelf(peerView *elevatorConfig.PeerView, hallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, cabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus) {
	hallRequestsForAllElevators[peerView.OwnId] = peerView.HallRequests
	cabRequestsForAllElevators[peerView.OwnId] = peerView.States[peerView.OwnId].CabRequests
}
	*/

// Not used:
/*
func updateElevatorSystem(peerView *elevatorConfig.PeerView, hallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, cabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus, receivedWorldView chan string) {
	peerSystem := DecodeElevatorSystem(<-receivedWorldView)
	addPeer(peerView, &peerSystem)
	UpdateElevatorSystemWithPeer(peerView, &peerSystem, hallRequestsForAllElevators, cabRequestsForAllElevators)
	UpdateElevatorSystemWithSelf(peerView, hallRequestsForAllElevators, cabRequestsForAllElevators)
}
	*/

func CopyElevatorSystem(peerView *elevatorConfig.PeerView) *elevatorConfig.PeerView {
	copyPeerView := *peerView // shallow copy

	// Deep copy AlivePeers slice
	copyPeerView.AlivePeers = make([]string, len(peerView.AlivePeers))
	copy(copyPeerView.AlivePeers, peerView.AlivePeers)

	// Deep copy HallRequests (array, so direct copy is fine)
	copyPeerView.HallRequests = peerView.HallRequests

	// Deep copy States map
	copyPeerView.States = make(map[string]*elevatorConfig.PeerState, len(peerView.States))
	for k, v := range peerView.States {
		// Deep copy ElevatorState struct
		stateCopy := *v
		copyPeerView.States[k] = &stateCopy
	}

	return &copyPeerView
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
