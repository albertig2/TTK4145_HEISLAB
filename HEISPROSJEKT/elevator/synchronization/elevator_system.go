package synchronization

import (
	"HEISPROSJEKT/elevatorConfig"
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

func SetCabRequests(peerView *elevatorConfig.PeerView, floor int, orderStatus elevatorConfig.OrderStatus) {
	state := peerView.States[peerView.OwnId]
	state.CabRequests[floor] = orderStatus
}

func SetHallRequests(peerView *elevatorConfig.PeerView, floor int, hallDirection int, orderStatus elevatorConfig.OrderStatus) {
	peerView.HallRequests[floor][hallDirection] = orderStatus
}

func InitializeElevatorSystem(peerView *elevatorConfig.PeerView, ownId string) {
	peerView.AlivePeers = []string{ownId}
	peerView.OwnId = ownId
	peerView.HallRequests = [elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus{}
	peerView.States = make(map[string]*elevatorConfig.PeerState)
	unknownFloor := -1
	peerView.States[ownId] = &elevatorConfig.PeerState{
		Behavior:    elevatorConfig.Idle,
		Floor:       unknownFloor,
		Direction:   elevatorConfig.Stop,
		CabRequests: [elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus{},
	}

	initializeHallRequests(peerView)
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
	for floor := 0; floor < elevatorConfig.N_FLOORS; floor++ {
		peerView.States[peerView.OwnId].CabRequests[floor] = elevatorConfig.Unknown
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

func UpdateLocalPeerViewWithPeer(localPeerView *elevatorConfig.PeerView, extarnalPeerView *elevatorConfig.PeerView, hallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, cabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus) {
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

func CopyPeerView(peerView *elevatorConfig.PeerView) *elevatorConfig.PeerView {
	copyPeerView := *peerView

	copyPeerView.AlivePeers = make([]string, len(peerView.AlivePeers))
	copy(copyPeerView.AlivePeers, peerView.AlivePeers)

	copyPeerView.HallRequests = peerView.HallRequests

	copyPeerView.States = make(map[string]*elevatorConfig.PeerState, len(peerView.States))
	for peerId, peerState := range peerView.States {
		stateCopy := *peerState
		copyPeerView.States[peerId] = &stateCopy
	}

	return &copyPeerView
}
