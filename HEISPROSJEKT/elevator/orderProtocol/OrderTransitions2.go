package orderProtocol

import (
	"HEISPROSJEKT/elevatorConfig"
)

func CheckOrderTransitionStatusForHallRequests2(
	system *elevatorConfig.ElevatorSystem,
	HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus,
	halldir int,
	floor int,
	newOrders []elevatorConfig.ButtonEvent,
	servicedOrder elevatorConfig.ButtonEvent) OrderTransition {

	noordertopending := false
	pendingtonoorder := false
	pendingtoassigned := true
	assignedtocomplete := false
	completetonoorder := true

	ownHallStatus := HallRequestsForAllElevators[system.OwnId][floor][halldir]

	if ownHallStatus == elevatorConfig.NoOrder {
		for _, peerId := range system.AlivePeers {
			peerHallStatus := HallRequestsForAllElevators[peerId][floor][halldir]
			if peerHallStatus == elevatorConfig.Completed {
				noordertopending = false
				break
			} else if peerHallStatus != elevatorConfig.NoOrder {
				noordertopending = true
			}
		}
	}

	for _, peerId := range system.AlivePeers {
		peerHallStatus := HallRequestsForAllElevators[peerId][floor][halldir]
		if peerHallStatus != elevatorConfig.Pending {
			pendingtoassigned = false
		}
		if peerHallStatus == elevatorConfig.Completed && peerId != system.OwnId {
			pendingtonoorder = true
		}
		if peerHallStatus != elevatorConfig.NoOrder && peerId != system.OwnId {
			completetonoorder = false
		}
	}

	if ownHallStatus != elevatorConfig.Completed && completetonoorder {
		completetonoorder = false
	}

	if ownHallStatus == elevatorConfig.Assigned && servicedOrder.Floor == floor {
		assignedtocomplete = true
	}

	if ownHallStatus == elevatorConfig.NoOrder {
		for _, order := range newOrders {
			if order.Floor == floor && int(order.Button) == halldir {
				noordertopending = true
				break
			}
		}
	}

	if noordertopending {
		return NoOrderToPending
	} else if pendingtoassigned {
		return PendingToAssigned
	} else if pendingtonoorder {
		return PendingToNoOrder
	} else if assignedtocomplete {
		return AssignedToComplete
	} else if completetonoorder {
		return CompleteToNoOrder
	}
	return NoTransition
}

// Are the priority of the transitions in correct order? or could noOrderToPending "overwrite" assigned?
// Dont know if this should be different or what
// Cab order goes from no order to pending when you press the button
// Cab order goes from pending to assigned when all elevators have pending
// Cab order never goes from pending to no order for own orders
// Cab orders goes from assigned to no order when the elevator reaches the floor
func CheckOrderTransitionStatusForCabRequests2(
	system *elevatorConfig.ElevatorSystem,
	CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus,
	floor int,
	newOrders []elevatorConfig.ButtonEvent,
	servicedOrder elevatorConfig.ButtonEvent) OrderTransition {

	noordertopending := false
	pendingtoassigned := true
	assignedtonoorder := false
	unknowntonoorder := false
	unknowntopending := false
	unknowntoassigned := false

	if CabRequestsForAllElevators[system.OwnId][floor] == elevatorConfig.Unknown {
		allNoOrder := true
		allAssignedOrPending := true
		anyPending := false

		for _, peerId := range system.AlivePeers {
			if peerId != system.OwnId {
				peerCabStatus := CabRequestsForAllElevators[peerId][floor]
				if peerCabStatus == elevatorConfig.Pending {
					anyPending = true
				}
				if peerCabStatus != elevatorConfig.NoOrder {
					allNoOrder = false
				}
				if peerCabStatus != elevatorConfig.Assigned && peerCabStatus != elevatorConfig.Pending {
					allAssignedOrPending = false
				}
			}
		}

		if allNoOrder {
			unknowntonoorder = true
		} else if allAssignedOrPending {
			unknowntoassigned = true
		} else if anyPending {
			unknowntopending = true
		}
		for _, order := range newOrders {
			if order.Floor == floor && order.Button == elevatorConfig.Cab {
				unknowntopending = true
				break
			}

		}
	}

	for _, order := range newOrders {
		if order.Floor == floor && order.Button == elevatorConfig.Cab {
			noordertopending = true
			break
		}

	}

	for _, peerId := range system.AlivePeers {
		peerCabStatus := CabRequestsForAllElevators[peerId][floor]
		if peerCabStatus != elevatorConfig.Pending {
			pendingtoassigned = false
		}
	}

	ownCabStatus := CabRequestsForAllElevators[system.OwnId][floor]
	if ownCabStatus == elevatorConfig.Assigned && servicedOrder.Floor == floor {
		assignedtonoorder = true
	}

	if noordertopending {
		return NoOrderToPending
	} else if pendingtoassigned {
		return PendingToAssigned
	} else if assignedtonoorder {
		return AssignedToNoOrder
	} else if unknowntonoorder {
		return UnknownToNoOrder
	} else if unknowntopending {
		return UnknownToPending
	} else if unknowntoassigned {
		return UnknownToAssigned
	}
	return NoTransition
}
