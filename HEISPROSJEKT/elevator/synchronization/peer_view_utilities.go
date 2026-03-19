package synchronization

/*
Handles the representation and manipulation of the shared elevator system state (PeerView).
Provides initialization, update, and merge logic for local and external peer data.
*/

import (
	elevatorConfig "HEISPROSJEKT/elevator_config"
)

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

func SetCabOrders(peerView *elevatorConfig.PeerView, floor int, orderStatus elevatorConfig.OrderStatus) {
	state := peerView.States[peerView.OwnId]
	state.CabOrders[floor] = orderStatus
}

func SetHallOrders(peerView *elevatorConfig.PeerView, floor int, hallDirection int, orderStatus elevatorConfig.OrderStatus) {
	peerView.HallOrders[floor][hallDirection] = orderStatus
}

func InitializePeerView(peerView *elevatorConfig.PeerView, ownId string) {
	peerView.AlivePeers = []string{ownId}
	peerView.OwnId = ownId
	peerView.HallOrders = [elevatorConfig.NumberOfFloors][2]elevatorConfig.OrderStatus{}
	peerView.States = make(map[string]*elevatorConfig.PeerState)
	unknownFloor := -1
	peerView.States[ownId] = &elevatorConfig.PeerState{
		Behavior:  elevatorConfig.Idle,
		Floor:     unknownFloor,
		Direction: elevatorConfig.Stop,
		CabOrders: [elevatorConfig.NumberOfFloors]elevatorConfig.OrderStatus{},
	}

	initializeHallRequests(peerView)
	initializeCabRequests(peerView)
}

func initializeHallRequests(peerView *elevatorConfig.PeerView) {
	for floor := 0; floor < elevatorConfig.NumberOfFloors; floor++ {
		for _, hallDirection := range HallDirections {
			peerView.HallOrders[floor][hallDirection] = elevatorConfig.NoOrder
		}
	}
}

func initializeCabRequests(peerView *elevatorConfig.PeerView) {
	for floor := 0; floor < elevatorConfig.NumberOfFloors; floor++ {
		peerView.States[peerView.OwnId].CabOrders[floor] = elevatorConfig.Unknown
	}
}

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
	cabRequests := externalPeerView.States[externalPeerView.OwnId].CabOrders
	for floor := 0; floor < elevatorConfig.NumberOfFloors; floor++ {
		if cabRequests[floor] == elevatorConfig.Unknown {
			cabRequests[floor] = elevatorConfig.NoOrder
		}
	}

	peerState := externalPeerView.States[externalPeerView.OwnId]
	peerView.States[externalPeerView.OwnId] = &elevatorConfig.PeerState{
		Behavior:  peerState.Behavior,
		Floor:     peerState.Floor,
		Direction: peerState.Direction,
		CabOrders: cabRequests,
	}
}

func updatePeer(localPeerView *elevatorConfig.PeerView, extarnalPeerView *elevatorConfig.PeerView) {
	peerSystemCabRequests := extarnalPeerView.States[extarnalPeerView.OwnId].CabOrders
	cabRequests := localPeerView.States[extarnalPeerView.OwnId].CabOrders

	for floor := 0; floor < elevatorConfig.NumberOfFloors; floor++ {
		if peerSystemCabRequests[floor] != elevatorConfig.Unknown {
			cabRequests[floor] = peerSystemCabRequests[floor]
		}
	}

	peerState := extarnalPeerView.States[extarnalPeerView.OwnId]
	localPeerView.States[extarnalPeerView.OwnId] = &elevatorConfig.PeerState{
		Behavior:  peerState.Behavior,
		Floor:     peerState.Floor,
		Direction: peerState.Direction,
		CabOrders: cabRequests,
	}
}

func UpdateLocalPeerViewWithPeer(localPeerView *elevatorConfig.PeerView, extarnalPeerView *elevatorConfig.PeerView, hallRequestsForAllElevators map[string][elevatorConfig.NumberOfFloors][2]elevatorConfig.OrderStatus, cabRequestsForAllElevators map[string][elevatorConfig.NumberOfFloors]elevatorConfig.OrderStatus) {
	if _, exists := localPeerView.States[extarnalPeerView.OwnId]; exists {
		updatePeer(localPeerView, extarnalPeerView)
	} else {
		addPeer(localPeerView, extarnalPeerView)
	}

	hallRequestsForAllElevators[extarnalPeerView.OwnId] = extarnalPeerView.HallOrders
	if _, exists := extarnalPeerView.States[localPeerView.OwnId]; exists {
		cabRequestsForAllElevators[extarnalPeerView.OwnId] = extarnalPeerView.States[localPeerView.OwnId].CabOrders
	}
}

func CopyPeerView(peerView *elevatorConfig.PeerView) *elevatorConfig.PeerView {
	copyPeerView := *peerView

	copyPeerView.AlivePeers = make([]string, len(peerView.AlivePeers))
	copy(copyPeerView.AlivePeers, peerView.AlivePeers)

	copyPeerView.HallOrders = peerView.HallOrders

	copyPeerView.States = make(map[string]*elevatorConfig.PeerState, len(peerView.States))
	for peerId, peerState := range peerView.States {
		stateCopy := *peerState
		copyPeerView.States[peerId] = &stateCopy
	}

	return &copyPeerView
}

func UpdateElevatorSystemFromElevator(elevator elevatorConfig.Elevator, peerView *elevatorConfig.PeerView) {
	SetBehavior(peerView, elevatorConfig.Behavior(elevator.Behavior))
	SetDirection(peerView, elevator.Direction)
	if elevator.Floor >= 0 && elevator.Floor < elevatorConfig.NumberOfFloors {
		SetFloor(peerView, elevator.Floor)
	}
}
