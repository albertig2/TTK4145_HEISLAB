package orderProtocol

import elevatorConfig "HEISPROSJEKT/elevator_config"

/*
This file contains functions that determine state transitions for hall and cab orders
in the distributed elevator control system. Each function is called for an order in a
specific state and determines the transition from that state. The returned OrderTransition
values are used by the order state machine to update and synchronize order states across all peers.
*/

func getHallTransitionFromNoOrder(hallOrdersForAllElevators map[string][elevatorConfig.NumberOfFloors][elevatorConfig.NumberOfHallButtons]elevatorConfig.OrderStatus, hallDirection int, floor int, newOrders []elevatorConfig.ButtonEvent, otherAlivePeers []string) OrderTransition {
	// If any of the elevators are in completed, the new order will not be set, but then the person just have to press the button again I guess.
	noOrderToPending := false
	for _, order := range newOrders {
		if order.Floor == floor && int(order.Button) == hallDirection {
			noOrderToPending = true
			break
		}
	}

	for _, peerId := range otherAlivePeers {
		peerHallStatus := hallOrdersForAllElevators[peerId][floor][hallDirection]
		if peerHallStatus == elevatorConfig.Serviced {
			noOrderToPending = false
			break
		} else if peerHallStatus != elevatorConfig.NoOrder {
			noOrderToPending = true
		}
	}

	if noOrderToPending {
		return NoOrderToPending
	}

	return NoTransition
}

func getHallTransitionFromPending(hallOrdersForAllElevators map[string][elevatorConfig.NumberOfFloors][elevatorConfig.NumberOfHallButtons]elevatorConfig.OrderStatus, hallDirection int, floor int, otherAlivePeers []string) OrderTransition {
	for _, peerId := range otherAlivePeers {
		peerHallStatus := hallOrdersForAllElevators[peerId][floor][hallDirection]
		if peerHallStatus == elevatorConfig.Serviced {
			return PendingToNoOrder
		}
	}

	for _, peerId := range otherAlivePeers {
		peerHallStatus := hallOrdersForAllElevators[peerId][floor][hallDirection]
		if peerHallStatus == elevatorConfig.Assigned {
			return PendingToPending
		}
	}

	pendingToAssigned := true
	for _, peerId := range otherAlivePeers {
		peerHallStatus := hallOrdersForAllElevators[peerId][floor][hallDirection]
		if peerHallStatus != elevatorConfig.Pending {
			pendingToAssigned = false
		}
	}

	if pendingToAssigned {
		return PendingToAssigned
	}

	return NoTransition
}

func getHallTransitionFromAssigned(hallDirection int, floor int, servicedOrders []elevatorConfig.ButtonEvent) OrderTransition {
	for _, order := range servicedOrders {
		if order.Floor == floor && int(order.Button) == hallDirection {
			return AssignedToComplete
		}
	}

	return NoTransition
}

func getHallTransitionFromCompleted(hallOrdersForAllElevators map[string][elevatorConfig.NumberOfFloors][elevatorConfig.NumberOfHallButtons]elevatorConfig.OrderStatus, hallDirection int, floor int, otherAlivePeers []string) OrderTransition {
	completeToNoOrder := true
	for _, peerId := range otherAlivePeers {
		peerHallStatus := hallOrdersForAllElevators[peerId][floor][hallDirection]
		if peerHallStatus != elevatorConfig.NoOrder {
			completeToNoOrder = false
			break
		}
	}

	if completeToNoOrder {
		return CompleteToNoOrder
	}

	return NoTransition
}

func getCabTransitionFromUnknown(cabOrdersForAllElevators map[string][elevatorConfig.NumberOfFloors]elevatorConfig.OrderStatus, floor int, newOrders []elevatorConfig.ButtonEvent, otherAlivePeers []string) OrderTransition {
	allNoOrder := true
	anyPending := false
	allAssignedOrPending := true
	for _, peerId := range otherAlivePeers {
		peerCabStatus := cabOrdersForAllElevators[peerId][floor]
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

	return NoTransition
}

func getCabTransitionFromNoOrder(floor int, newOrders []elevatorConfig.ButtonEvent) OrderTransition {
	for _, order := range newOrders {
		if order.Floor == floor {
			return NoOrderToPending
		}
	}

	return NoTransition
}

func getCabTransitionFromPending(cabOrdersForAllElevators map[string][elevatorConfig.NumberOfFloors]elevatorConfig.OrderStatus, floor int, otherAlivePeers []string) OrderTransition {
	pendingtoassigned := true
	for _, peerId := range otherAlivePeers {
		peerCabStatus := cabOrdersForAllElevators[peerId][floor]
		if peerCabStatus != elevatorConfig.Pending {
			pendingtoassigned = false
		}
	}

	if pendingtoassigned {
		return PendingToAssigned
	}

	return NoTransition
}

func getCabTransitionFromAssigned(floor int, servicedOrders []elevatorConfig.ButtonEvent) OrderTransition {
	for _, order := range servicedOrders {
		if order.Floor == floor {
			return AssignedToNoOrder
		}
	}

	return NoTransition
}
