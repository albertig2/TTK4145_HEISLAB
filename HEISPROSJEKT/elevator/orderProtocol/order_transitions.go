package orderProtocol

import (
	"HEISPROSJEKT/elevatorConfig"

	//"HEISPROSJEKT/elevatorHardware"
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

func checkOrderTransitionStatusForHallRequests(peerView *elevatorConfig.PeerView, hallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, hallDirection int, floor int, newOrders []elevatorConfig.ButtonEvent, servicedOrders []elevatorConfig.ButtonEvent) OrderTransition {
	var otherAlivePeers []string
	for _, peerId := range peerView.AlivePeers {
		if peerId != peerView.OwnId {
			otherAlivePeers = append(otherAlivePeers, peerId)
		}
	}

	transition := NoTransition
	ownHallStatus := hallRequestsForAllElevators[peerView.OwnId][floor][hallDirection]
	switch ownHallStatus {
	case elevatorConfig.NoOrder:
		transition = getHallTransitionFromNoOrder(hallRequestsForAllElevators, hallDirection, floor, newOrders, otherAlivePeers)
	case elevatorConfig.Pending:
		transition = getHallTransitionFromPending(hallRequestsForAllElevators, hallDirection, floor, otherAlivePeers)
	case elevatorConfig.Assigned:
		transition = getHallTransitionFromAssigned(hallDirection, floor, servicedOrders)
	case elevatorConfig.Completed:
		transition = getHallTransitionFromCompleted(hallRequestsForAllElevators, hallDirection, floor, otherAlivePeers)
	}
	return transition
}

func checkOrderTransitionStatusForCabRequests(peerView *elevatorConfig.PeerView, cabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus, floor int, newOrders []elevatorConfig.ButtonEvent, servicedOrders []elevatorConfig.ButtonEvent) OrderTransition {
	var otherAlivePeers []string
	for _, peerId := range peerView.AlivePeers {
		if peerId != peerView.OwnId {
			otherAlivePeers = append(otherAlivePeers, peerId)
		}
	}

	transition := NoTransition
	ownCabStatus := cabRequestsForAllElevators[peerView.OwnId][floor]
	switch ownCabStatus {
	case elevatorConfig.Unknown:
		transition = getCabTransitionFromUnknown(cabRequestsForAllElevators, floor, newOrders, otherAlivePeers)
	case elevatorConfig.NoOrder:
		transition = getCabTransitionFromNoOrder(floor, newOrders)
	case elevatorConfig.Pending:
		transition = getCabTransitionFromPending(cabRequestsForAllElevators, floor, otherAlivePeers)
	case elevatorConfig.Assigned:
		transition = getCabTransitionFromAssigned(floor, servicedOrders)
	}
	return transition
}

func getAllHallRequestTransitions(peerView *elevatorConfig.PeerView, hallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, newOrders []elevatorConfig.ButtonEvent, servicedOrders []elevatorConfig.ButtonEvent) [elevatorConfig.N_FLOORS][2]OrderTransition {
	var transitions [elevatorConfig.N_FLOORS][2]OrderTransition
	for floor := range elevatorConfig.N_FLOORS {
		for _, hallDirection := range synchronization.HallDirections {
			transition := checkOrderTransitionStatusForHallRequests(peerView, hallRequestsForAllElevators, hallDirection, floor, newOrders, servicedOrders)
			transitions[floor][hallDirection] = transition
		}
	}
	return transitions
}

func getAllCabRequestTransitions(peerView *elevatorConfig.PeerView, cabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus, newOrders []elevatorConfig.ButtonEvent, servicedOrders []elevatorConfig.ButtonEvent) [elevatorConfig.N_FLOORS]OrderTransition {
	var transitions [elevatorConfig.N_FLOORS]OrderTransition
	for floor := range elevatorConfig.N_FLOORS {
		transition := checkOrderTransitionStatusForCabRequests(peerView, cabRequestsForAllElevators, floor, newOrders, servicedOrders)
		transitions[floor] = transition
	}
	return transitions
}

func transitioningAllHallRequests(peerView *elevatorConfig.ElevatorSystem, hallRequestTransitions [elevatorConfig.N_FLOORS][2]OrderTransition, elevatorOrderChannels elevatorConfig.OrderChannels) {
	for floor := range elevatorConfig.N_FLOORS {
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
				status = elevatorConfig.Completed
			case CompleteToNoOrder:
				status = elevatorConfig.NoOrder
			default:
				status = peerView.HallRequests[floor][hallDirection]
			}
			synchronization.SetHallRequests(peerView, floor, hallDirection, status)
		}
	}
	transitionFromPendingToAssignedForHallRequests(peerView, hallRequestTransitions, elevatorOrderChannels)
}

// Called from within the TransitionForHallRequestsByType. Setting to private to avoid it being called directly
func transitionFromPendingToAssignedForHallRequests(peerView *elevatorConfig.PeerView, hallRequestTransitions [elevatorConfig.N_FLOORS][2]OrderTransition, orderChannels elevatorConfig.OrderChannels) {
	output := hallRequestAssigner(peerView, hallRequestTransitions)
	for floor := range elevatorConfig.N_FLOORS {
		for _, hallDirection := range synchronization.HallDirections {
			if (output)[peerView.OwnId][floor][hallDirection] {
				synchronization.SetHallRequests(peerView, floor, hallDirection, elevatorConfig.Assigned)
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

func transitioningAllCabRequests(peerView *elevatorConfig.PeerView, cabRequestTransitions [elevatorConfig.N_FLOORS]OrderTransition, orderChannels elevatorConfig.OrderChannels) {
	for floor := range elevatorConfig.N_FLOORS {
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
			status = peerView.States[peerView.OwnId].CabRequests[floor]
		}
		synchronization.SetCabRequests(peerView, floor, status)
	}
}
