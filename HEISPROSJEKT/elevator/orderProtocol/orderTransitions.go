package orderProtocol

import (
	"HEISPROSJEKT/elevatorConfig"

	//"HEISPROSJEKT/elevatorHardware"
	"HEISPROSJEKT/synchronisation"
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

func checkOrderTransitionStatusForHallRequests(peerView *elevatorConfig.ElevatorSystem, hallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, hallDirection int, floor int, newOrders []elevatorConfig.ButtonEvent, servicedOrders []elevatorConfig.ButtonEvent) OrderTransition {
	var otherAlivePeers []string
	for _, peerId := range peerView.AlivePeers {
		if peerId != peerView.OwnId {
			otherAlivePeers = append(otherAlivePeers, peerId)
		}
	}

	ownHallStatus := hallRequestsForAllElevators[peerView.OwnId][floor][hallDirection]

	switch ownHallStatus {
	case elevatorConfig.NoOrder:
		noOrderToPending := false
		// If any of the elevators are in completed, the new order will not be set, but then the person just have to press the button again I guess.
		for _, order := range newOrders {
			if order.Floor == floor && int(order.Button) == hallDirection {
				noOrderToPending = true
				break
			}
		}
		for _, peerId := range otherAlivePeers {
			peerHallStatus := hallRequestsForAllElevators[peerId][floor][hallDirection]
			if peerHallStatus == elevatorConfig.Completed {
				noOrderToPending = false
				break
			} else if peerHallStatus != elevatorConfig.NoOrder {
				noOrderToPending = true
			}
		}
		if noOrderToPending {
			return NoOrderToPending
		}
	case elevatorConfig.Pending:
		for _, peerId := range otherAlivePeers {
			peerHallStatus := hallRequestsForAllElevators[peerId][floor][hallDirection]
			if peerHallStatus == elevatorConfig.Completed {
				return PendingToNoOrder
			}
		}
		pendingToAssigned := true
		for _, peerId := range otherAlivePeers {
			peerHallStatus := hallRequestsForAllElevators[peerId][floor][hallDirection]
			if peerHallStatus != elevatorConfig.Pending {
				pendingToAssigned = false
			}
		}
		if pendingToAssigned {
			return PendingToAssigned
		}
		for _, peerId := range otherAlivePeers {
			peerHallStatus := hallRequestsForAllElevators[peerId][floor][hallDirection]
			if peerHallStatus == elevatorConfig.Assigned {
				return PendingToPending
			}
		}
	case elevatorConfig.Assigned:
		for _, order := range servicedOrders {
			if order.Floor == floor && int(order.Button) == hallDirection {
				return AssignedToComplete
			}
		}
	case elevatorConfig.Completed:
		completeToNoOrder := true
		for _, peerId := range otherAlivePeers {
			peerHallStatus := hallRequestsForAllElevators[peerId][floor][hallDirection]
			if peerHallStatus != elevatorConfig.NoOrder {
				completeToNoOrder = false
				break
			}
		}
		if completeToNoOrder {
			return CompleteToNoOrder
		}
	}
	return NoTransition
}

func checkOrderTransitionStatusForCabRequests(peerView *elevatorConfig.ElevatorSystem, cabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus, floor int, newOrders []elevatorConfig.ButtonEvent, servicedOrders []elevatorConfig.ButtonEvent) OrderTransition {
	var otherAlivePeers []string
	for _, peerId := range peerView.AlivePeers {
		if peerId != peerView.OwnId {
			otherAlivePeers = append(otherAlivePeers, peerId)
		}
	}

	ownCabStatus := cabRequestsForAllElevators[peerView.OwnId][floor]

	switch ownCabStatus {
	case elevatorConfig.Unknown:
		allNoOrder := true
		anyPending := false
		allAssignedOrPending := true
		for _, peerId := range otherAlivePeers {
			peerCabStatus := cabRequestsForAllElevators[peerId][floor]
			if peerCabStatus != elevatorConfig.NoOrder {
				allNoOrder = false
			}
			if peerCabStatus == elevatorConfig.Pending {
				anyPending = true
			}
			if peerCabStatus != elevatorConfig.Assigned && peerCabStatus != elevatorConfig.Pending {
				allAssignedOrPending = false
			}
		}
		for _, order := range newOrders {
			if order.Floor == floor {
				return UnknownToPending
			}
		}
		if allNoOrder {
			return UnknownToNoOrder
		} else if allAssignedOrPending {
			return UnknownToAssigned
		} else if anyPending {
			return UnknownToPending
		}
	case elevatorConfig.NoOrder:
		for _, order := range newOrders {
			if order.Floor == floor {
				return NoOrderToPending
			}
		}
	case elevatorConfig.Pending:
		pendingtoassigned := true
		for _, peerId := range otherAlivePeers {
			peerCabStatus := cabRequestsForAllElevators[peerId][floor]
			if peerCabStatus != elevatorConfig.Pending {
				pendingtoassigned = false
			}
		}
		if pendingtoassigned {
			return PendingToAssigned
		}
	case elevatorConfig.Assigned:
		for _, order := range servicedOrders {
			if order.Floor == floor {
				return AssignedToNoOrder
			}
		}
	}

	return NoTransition
}

func getAllHallRequestTransitions(peerView *elevatorConfig.ElevatorSystem, hallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, newOrders []elevatorConfig.ButtonEvent, servicedOrders []elevatorConfig.ButtonEvent) [elevatorConfig.N_FLOORS][2]OrderTransition {
	var transitions [elevatorConfig.N_FLOORS][2]OrderTransition
	for floor := range elevatorConfig.N_FLOORS {
		for _, hallDirection := range synchronisation.HallDirections {
			transition := checkOrderTransitionStatusForHallRequests(peerView, hallRequestsForAllElevators, hallDirection, floor, newOrders, servicedOrders)
			transitions[floor][hallDirection] = transition
		}
	}
	return transitions
}

func getAllCabRequestTransitions(peerView *elevatorConfig.ElevatorSystem, cabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus, newOrders []elevatorConfig.ButtonEvent, servicedOrders []elevatorConfig.ButtonEvent) [elevatorConfig.N_FLOORS]OrderTransition {
	var transitions [elevatorConfig.N_FLOORS]OrderTransition
	for floor := range elevatorConfig.N_FLOORS {
		transition := checkOrderTransitionStatusForCabRequests(peerView, cabRequestsForAllElevators, floor, newOrders, servicedOrders)
		transitions[floor] = transition
	}
	return transitions
}

func transitioningAllHallRequests(peerView *elevatorConfig.ElevatorSystem, hallRequestTransitions [elevatorConfig.N_FLOORS][2]OrderTransition, elevatorOrderChannels elevatorConfig.OrderChannels) {
	for floor := range elevatorConfig.N_FLOORS {
		for _, hallDirection := range synchronisation.HallDirections {
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
			synchronisation.SetHallRequests(peerView, floor, hallDirection, status)
		}
	}
	transitionFromPendingToAssignedForHallRequests(peerView, hallRequestTransitions, elevatorOrderChannels)
}

// Called from within the TransitionForHallRequestsByType. Setting to private to avoid it being called directly
func transitionFromPendingToAssignedForHallRequests(peerView *elevatorConfig.ElevatorSystem, hallRequestTransitions [elevatorConfig.N_FLOORS][2]OrderTransition, orderChannels elevatorConfig.OrderChannels) {
	output := hallRequestAssigner(peerView, hallRequestTransitions)
	for floor := range elevatorConfig.N_FLOORS {
		for _, hallDirection := range synchronisation.HallDirections {
			if (output)[peerView.OwnId][floor][hallDirection] {
				synchronisation.SetHallRequests(peerView, floor, hallDirection, elevatorConfig.Assigned)
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

func transitioningAllCabRequests(peerView *elevatorConfig.ElevatorSystem, cabRequestTransitions [elevatorConfig.N_FLOORS]OrderTransition, orderChannels elevatorConfig.OrderChannels) {
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
		synchronisation.SetCabRequests(peerView, floor, status)
	}
}
