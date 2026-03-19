package orderProtocol

import (
	elevatorConfig "HEISPROSJEKT/elevator_config"
	"HEISPROSJEKT/synchronization"
)

/*
This file contains utility functions for the order management and distribution protocol.
It includes functions for initializing order channels, processing order transitions, and managing the state of orders across the system.
*/

func orderRutine(peerView *elevatorConfig.PeerView, hallOrdersForAllElevators *map[string][elevatorConfig.NumberOfFloors][elevatorConfig.NumberOfHallButtons]elevatorConfig.OrderStatus, cabOrdersForAllElevators *map[string][elevatorConfig.NumberOfFloors]elevatorConfig.OrderStatus, orderChannels elevatorConfig.OrderChannels, newHallOrders []elevatorConfig.ButtonEvent, newCabOrders []elevatorConfig.ButtonEvent, servicedHallOrders []elevatorConfig.ButtonEvent, servicedCabOrders []elevatorConfig.ButtonEvent) {
	HallRequestTransitions := getAllHallRequestTransitions(peerView, *hallOrdersForAllElevators, newHallOrders, servicedHallOrders)
	CabRequestTransitions := getAllCabRequestTransitions(peerView, *cabOrdersForAllElevators, newCabOrders, servicedCabOrders)
	transitioningAllHallOrders(peerView, HallRequestTransitions, orderChannels)
	transitioningAllCabOrders(peerView, CabRequestTransitions, orderChannels)
	(*hallOrdersForAllElevators)[peerView.OwnId] = peerView.HallOrders
	(*cabOrdersForAllElevators)[peerView.OwnId] = peerView.States[peerView.OwnId].CabOrders
}

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

func initializePeerView(peerView *elevatorConfig.PeerView, ownId string) (map[string][elevatorConfig.NumberOfFloors][elevatorConfig.NumberOfHallButtons]elevatorConfig.OrderStatus, map[string][elevatorConfig.NumberOfFloors]elevatorConfig.OrderStatus) {
	synchronization.InitializePeerView(peerView, ownId)
	hallOrdersForAllElevators := make(map[string][elevatorConfig.NumberOfFloors][elevatorConfig.NumberOfHallButtons]elevatorConfig.OrderStatus)
	cabOrdersForAllElevators := make(map[string][elevatorConfig.NumberOfFloors]elevatorConfig.OrderStatus)
	hallOrdersForAllElevators[ownId] = peerView.HallOrders
	cabOrdersForAllElevators[ownId] = peerView.States[ownId].CabOrders
	return hallOrdersForAllElevators, cabOrdersForAllElevators
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
