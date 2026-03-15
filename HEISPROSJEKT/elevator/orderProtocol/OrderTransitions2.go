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

	var otherAlivePeers []string
	for _, peerId := range system.AlivePeers {
		if peerId != system.OwnId {
			otherAlivePeers = append(otherAlivePeers, peerId)
		}
	}

	ownHallStatus := HallRequestsForAllElevators[system.OwnId][floor][halldir]

	switch ownHallStatus {
	case elevatorConfig.NoOrder:
		noordertopending := false
		// If any of the elevators are in completed, the new order will not be set, but then the person just have to press the button again I guess.
		for _, order := range newOrders {
			if order.Floor == floor && int(order.Button) == halldir {
				noordertopending = true
				break
			}
		}

		for _, peerID := range otherAlivePeers {
			peerHallStatus := HallRequestsForAllElevators[peerID][floor][halldir]
			if peerHallStatus == elevatorConfig.Completed {
				noordertopending = false
				break
			} else if peerHallStatus != elevatorConfig.NoOrder {
				noordertopending = true
			}
		}

		if noordertopending {
			return NoOrderToPending
		}
	case elevatorConfig.Pending:
		for _, peerID := range otherAlivePeers {
			peerHallStatus := HallRequestsForAllElevators[peerID][floor][halldir]
			if peerHallStatus == elevatorConfig.Completed {
				return PendingToNoOrder
			}
		}

		pendingtoassigned := true
		for _, peerID := range otherAlivePeers {
			peerHallStatus := HallRequestsForAllElevators[peerID][floor][halldir]
			if peerHallStatus != elevatorConfig.Pending {
				pendingtoassigned = false
			}
		}

		if pendingtoassigned {
			return PendingToAssigned
		}

	case elevatorConfig.Assigned:
		if servicedOrder.Floor == floor {
			return AssignedToComplete
		}
	case elevatorConfig.Completed:
		completetonoorder := true
		for _, peerID := range otherAlivePeers {
			peerHallStatus := HallRequestsForAllElevators[peerID][floor][halldir]
			if peerHallStatus != elevatorConfig.NoOrder {
				completetonoorder = false
				break
			}
		}
		if completetonoorder {
			return CompleteToNoOrder
		}
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

	var otherAlivePeers []string
	for _, peerId := range system.AlivePeers {
		if peerId != system.OwnId {
			otherAlivePeers = append(otherAlivePeers, peerId)
		}
	}

	ownCabStatus := CabRequestsForAllElevators[system.OwnId][floor]

	switch ownCabStatus {
	case elevatorConfig.Unknown:
		allNoOrder := true
		anyPending := false
		allAssignedOrPending := true

		for _, peerID := range otherAlivePeers {
			peerCabStatus := CabRequestsForAllElevators[peerID][floor]
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
			if order.Floor == floor && order.Button == elevatorConfig.Cab {
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
		for _, peerID := range otherAlivePeers {
			peerCabStatus := CabRequestsForAllElevators[peerID][floor]
			if peerCabStatus != elevatorConfig.Pending {
				pendingtoassigned = false
			}
		}
		if pendingtoassigned {
			return PendingToAssigned
		}

	case elevatorConfig.Assigned:
		if servicedOrder.Floor == floor {
			return AssignedToNoOrder
		}
	}

	return NoTransition
}
