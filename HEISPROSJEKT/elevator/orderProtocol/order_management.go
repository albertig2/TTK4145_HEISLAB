package orderProtocol

import (
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/synchronization"
)

func InitializeOrderChannels() elevatorConfig.OrderChannels {
	orderChannelse := elevatorConfig.OrderChannels{
		NewRecievedOrderChannel:     make(chan elevatorConfig.ButtonEvent, 10),
		NewAssignedOrderChannel:     make(chan elevatorConfig.ButtonEvent, 10),
		NewAssignedPeerOrderChannel: make(chan elevatorConfig.ButtonEvent, 10),
		ServicedOrderChannel:        make(chan elevatorConfig.ButtonEvent, 10),
		ServicedPeerOrderChannel:    make(chan elevatorConfig.ButtonEvent, 10),
	}
	return orderChannelse
}

func orderRutine(peerView *elevatorConfig.PeerView, hallRequestsForAllElevators *map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, cabRequestsForAllElevators *map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus, orderChannels elevatorConfig.OrderChannels, newHallOrders []elevatorConfig.ButtonEvent, newCabOrders []elevatorConfig.ButtonEvent, servicedHallOrders []elevatorConfig.ButtonEvent, servicedCabOrders []elevatorConfig.ButtonEvent) {
	HallRequestTransitions := getAllHallRequestTransitions(peerView, *hallRequestsForAllElevators, newHallOrders, servicedHallOrders)
	CabRequestTransitions := getAllCabRequestTransitions(peerView, *cabRequestsForAllElevators, newCabOrders, servicedCabOrders)
	transitioningAllHallRequests(peerView, HallRequestTransitions, orderChannels)
	transitioningAllCabRequests(peerView, CabRequestTransitions, orderChannels)
	(*hallRequestsForAllElevators)[peerView.OwnId] = peerView.HallRequests
	(*cabRequestsForAllElevators)[peerView.OwnId] = peerView.States[peerView.OwnId].CabRequests
}

func initializePeerView(peerView *elevatorConfig.PeerView, ownId string) (map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus) {
	synchronization.InitializePeerView(peerView, ownId)
	hallRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus)
	cabRequestsForAllElevators := make(map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus)
	hallRequestsForAllElevators[ownId] = peerView.HallRequests
	cabRequestsForAllElevators[ownId] = peerView.States[ownId].CabRequests
	return hallRequestsForAllElevators, cabRequestsForAllElevators
}

func appendOrderByType(hallOrders []elevatorConfig.ButtonEvent, cabOrders []elevatorConfig.ButtonEvent, order elevatorConfig.ButtonEvent) ([]elevatorConfig.ButtonEvent, []elevatorConfig.ButtonEvent) {
	if order.Button == elevatorConfig.HallUp || order.Button == elevatorConfig.HallDown {
		hallOrders = append(hallOrders, order)
	} else if order.Button == elevatorConfig.Cab {
		cabOrders = append(cabOrders, order)
	}
	return hallOrders, cabOrders
}

func filterValidPeersAndIncludeOwnId(peerView *elevatorConfig.PeerView, alivePeers []string) []string {
	validAlivePeers := []string{}
	for _, peerID := range alivePeers {
		if _, ok := peerView.States[peerID]; ok {
			validAlivePeers = append(validAlivePeers, peerID)
		}
	}
	if !synchronization.Contains(validAlivePeers, peerView.OwnId) {
		validAlivePeers = append(validAlivePeers, peerView.OwnId)
	}
	return validAlivePeers
}

func areAllFloorsValid(peerView *elevatorConfig.PeerView) bool {
	allFloorsValid := true
	for _, peerId := range peerView.AlivePeers {
		if peerView.States[peerId].Floor == -1 {
			allFloorsValid = false
			break
		}
	}
	return allFloorsValid
}

func initializeOrDrainOrders() ([]elevatorConfig.ButtonEvent, []elevatorConfig.ButtonEvent, []elevatorConfig.ButtonEvent, []elevatorConfig.ButtonEvent) {
	return make([]elevatorConfig.ButtonEvent, 0), make([]elevatorConfig.ButtonEvent, 0), make([]elevatorConfig.ButtonEvent, 0), make([]elevatorConfig.ButtonEvent, 0)
}
