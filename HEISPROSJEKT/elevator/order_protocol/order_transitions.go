package orderProtocol

import (
	elevatorConfig "HEISPROSJEKT/elevator_config"
	"HEISPROSJEKT/synchronization"
)

type OrderTransition int

const (
	// Commmon for both hall and cab orders
	NoTransition OrderTransition = iota
	NoOrderToPending
	PendingToAssigned
	// Only for hall orders
	PendingToNoOrder
	PendingToPending
	AssignedToComplete
	CompleteToNoOrder
	// Only for cab orders
	AssignedToNoOrder
	UnknownToNoOrder
	UnknownToPending
	UnknownToAssigned
)

func transitioningAllHallOrders(peerView *elevatorConfig.PeerView, hallRequestTransitions [elevatorConfig.NumberOfFloors][2]OrderTransition, elevatorOrderChannels elevatorConfig.OrderChannels) {
	for floor := range elevatorConfig.NumberOfFloors {
		for _, hallDirection := range synchronization.HallDirections {
			var status elevatorConfig.OrderStatus
			switch hallRequestTransitions[floor][hallDirection] {
			case NoOrderToPending:
				status = elevatorConfig.Pending
			case PendingToNoOrder:
				status = elevatorConfig.NoOrder
				elevatorOrderChannels.ServicedPeerOrderChannel <- elevatorConfig.ButtonEvent{Floor: floor, Button: elevatorConfig.Button(hallDirection)}
			case PendingToPending:
				status = elevatorConfig.Pending
				elevatorOrderChannels.NewAssignedPeerOrderChannel <- elevatorConfig.ButtonEvent{Floor: floor, Button: elevatorConfig.Button(hallDirection)}
			case AssignedToComplete:
				status = elevatorConfig.Serviced
			case CompleteToNoOrder:
				status = elevatorConfig.NoOrder
			default:
				status = peerView.HallOrders[floor][hallDirection]
			}
			synchronization.SetHallOrders(peerView, floor, hallDirection, status)
		}
	}
	transitionFromPendingToAssignedForHallOrders(peerView, hallRequestTransitions, elevatorOrderChannels)
}

func transitionFromPendingToAssignedForHallOrders(peerView *elevatorConfig.PeerView, hallRequestTransitions [elevatorConfig.NumberOfFloors][2]OrderTransition, orderChannels elevatorConfig.OrderChannels) {
	output := hallOrderAssigner(peerView, hallRequestTransitions)
	for floor := range elevatorConfig.NumberOfFloors {
		for _, hallDirection := range synchronization.HallDirections {
			if (output)[peerView.OwnId][floor][hallDirection] {
				synchronization.SetHallOrders(peerView, floor, hallDirection, elevatorConfig.Assigned)
				orderChannels.NewAssignedOrderChannel <- elevatorConfig.ButtonEvent{Floor: floor, Button: elevatorConfig.Button(hallDirection)}
			} else {
				assignedToOther := false
				for _, peerId := range peerView.AlivePeers {
					if peerId != peerView.OwnId && output[peerId][floor][hallDirection] {
						assignedToOther = true
						break
					}
				}
				if assignedToOther {
					orderChannels.NewAssignedPeerOrderChannel <- elevatorConfig.ButtonEvent{Floor: floor, Button: elevatorConfig.Button(hallDirection)}
				}
			}
		}
	}
}

func transitioningAllCabOrders(peerView *elevatorConfig.PeerView, cabRequestTransitions [elevatorConfig.NumberOfFloors]OrderTransition, orderChannels elevatorConfig.OrderChannels) {
	for floor := range elevatorConfig.NumberOfFloors {
		var status elevatorConfig.OrderStatus
		switch cabRequestTransitions[floor] {
		case NoOrderToPending, UnknownToPending:
			status = elevatorConfig.Pending
		case PendingToAssigned, UnknownToAssigned:
			status = elevatorConfig.Assigned
			orderChannels.NewAssignedOrderChannel <- elevatorConfig.ButtonEvent{Floor: floor, Button: elevatorConfig.Cab}
		case AssignedToNoOrder, UnknownToNoOrder:
			status = elevatorConfig.NoOrder
		default:
			status = peerView.States[peerView.OwnId].CabOrders[floor]
		}
		synchronization.SetCabOrders(peerView, floor, status)
	}
}

func checkOrderTransitionStatusForHallOrders(peerView *elevatorConfig.PeerView, hallOrdersForAllElevators map[string][elevatorConfig.NumberOfFloors][2]elevatorConfig.OrderStatus, hallDirection int, floor int, newOrders []elevatorConfig.ButtonEvent, servicedOrders []elevatorConfig.ButtonEvent) OrderTransition {
	var otherAlivePeers []string
	for _, peerId := range peerView.AlivePeers {
		if peerId != peerView.OwnId {
			otherAlivePeers = append(otherAlivePeers, peerId)
		}
	}

	transition := NoTransition
	ownHallStatus := hallOrdersForAllElevators[peerView.OwnId][floor][hallDirection]
	switch ownHallStatus {
	case elevatorConfig.NoOrder:
		transition = getHallTransitionFromNoOrder(hallOrdersForAllElevators, hallDirection, floor, newOrders, otherAlivePeers)
	case elevatorConfig.Pending:
		transition = getHallTransitionFromPending(hallOrdersForAllElevators, hallDirection, floor, otherAlivePeers)
	case elevatorConfig.Assigned:
		transition = getHallTransitionFromAssigned(hallDirection, floor, servicedOrders)
	case elevatorConfig.Serviced:
		transition = getHallTransitionFromCompleted(hallOrdersForAllElevators, hallDirection, floor, otherAlivePeers)
	}
	return transition
}

func checkOrderTransitionStatusForCabOrders(peerView *elevatorConfig.PeerView, cabOrdersForAllElevators map[string][elevatorConfig.NumberOfFloors]elevatorConfig.OrderStatus, floor int, newOrders []elevatorConfig.ButtonEvent, servicedOrders []elevatorConfig.ButtonEvent) OrderTransition {
	var otherAlivePeers []string
	for _, peerId := range peerView.AlivePeers {
		if peerId != peerView.OwnId {
			otherAlivePeers = append(otherAlivePeers, peerId)
		}
	}

	transition := NoTransition
	ownCabStatus := cabOrdersForAllElevators[peerView.OwnId][floor]
	switch ownCabStatus {
	case elevatorConfig.Unknown:
		transition = getCabTransitionFromUnknown(cabOrdersForAllElevators, floor, newOrders, otherAlivePeers)
	case elevatorConfig.NoOrder:
		transition = getCabTransitionFromNoOrder(floor, newOrders)
	case elevatorConfig.Pending:
		transition = getCabTransitionFromPending(cabOrdersForAllElevators, floor, otherAlivePeers)
	case elevatorConfig.Assigned:
		transition = getCabTransitionFromAssigned(floor, servicedOrders)
	}
	return transition
}

func getAllHallRequestTransitions(peerView *elevatorConfig.PeerView, hallOrdersForAllElevators map[string][elevatorConfig.NumberOfFloors][2]elevatorConfig.OrderStatus, newOrders []elevatorConfig.ButtonEvent, servicedOrders []elevatorConfig.ButtonEvent) [elevatorConfig.NumberOfFloors][2]OrderTransition {
	var transitions [elevatorConfig.NumberOfFloors][2]OrderTransition
	for floor := range elevatorConfig.NumberOfFloors {
		for _, hallDirection := range synchronization.HallDirections {
			transition := checkOrderTransitionStatusForHallOrders(peerView, hallOrdersForAllElevators, hallDirection, floor, newOrders, servicedOrders)
			transitions[floor][hallDirection] = transition
		}
	}
	return transitions
}

func getAllCabRequestTransitions(peerView *elevatorConfig.PeerView, cabOrdersForAllElevators map[string][elevatorConfig.NumberOfFloors]elevatorConfig.OrderStatus, newOrders []elevatorConfig.ButtonEvent, servicedOrders []elevatorConfig.ButtonEvent) [elevatorConfig.NumberOfFloors]OrderTransition {
	var transitions [elevatorConfig.NumberOfFloors]OrderTransition
	for floor := range elevatorConfig.NumberOfFloors {
		transition := checkOrderTransitionStatusForCabOrders(peerView, cabOrdersForAllElevators, floor, newOrders, servicedOrders)
		transitions[floor] = transition
	}
	return transitions
}
