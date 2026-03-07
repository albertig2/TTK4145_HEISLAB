package orderProtocol

import (
	"HEISPROSJEKT/elevatorHardware"
)

type OrderTransition int

const (
	noTransition OrderTransition = iota
	pendingToAssigned
	pendingToNoOrder
	assignedToComplete
	completeToNoOrder
)

func CheckOrderTransitionStatusForHallRequests(
	system *elevatorHardware.ElevatorSystem,
	HallRequestsForAllElevators map[string][elevatorHardware.N_FLOORS][2]elevatorHardware.OrderStatus,
	halldir int,
	floor int,
	alivePeers []string) OrderTransition {

	pendingtoassigned := true
	pendingtonoorder := false
	completetonoorder := true
	assignedtocomplete := false

	for _, peerId := range alivePeers {
		if HallRequestsForAllElevators[peerId][floor][halldir] != elevatorHardware.Pending {
			pendingtoassigned = false
		}
		if HallRequestsForAllElevators[peerId][floor][halldir] == elevatorHardware.Completed && peerId != system.OwnId {
			pendingtonoorder = true
		}
		if HallRequestsForAllElevators[peerId][floor][halldir] != elevatorHardware.NoOrder && peerId != system.OwnId {
			completetonoorder = false
		}
	}

	if HallRequestsForAllElevators[system.OwnId][floor][halldir] != elevatorHardware.Completed && completetonoorder {
		completetonoorder = false
	}

	if HallRequestsForAllElevators[system.OwnId][floor][halldir] == elevatorHardware.Assigned && elevatorHardware.Elevator_floorSensor() == floor {
		assignedtocomplete = true
	}

	if pendingtoassigned {
		return pendingToAssigned
	} else if pendingtonoorder {
		return pendingToNoOrder
	} else if assignedtocomplete {
		return assignedToComplete
	} else if completetonoorder {
		return completeToNoOrder
	}
	return noTransition
}

func CheckOrderTransitionStatusForCabRequests(
	system *elevatorHardware.ElevatorSystem,
	CabRequestsForAllElevators map[string][elevatorHardware.N_FLOORS]elevatorHardware.OrderStatus,
	floor int,
	alivePeers []string) OrderTransition {

	pendingtoassigned := true
	pendingtonoorder := false
	completetonoorder := true
	assignedtocomplete := false

	for _, peerId := range alivePeers {
		if CabRequestsForAllElevators[peerId][floor] != elevatorHardware.Pending {
			pendingtoassigned = false
		}
		if CabRequestsForAllElevators[peerId][floor] == elevatorHardware.Completed && peerId != system.OwnId {
			pendingtonoorder = true
		}
		if CabRequestsForAllElevators[peerId][floor] != elevatorHardware.NoOrder && peerId != system.OwnId {
			completetonoorder = false
		}
	}

	if CabRequestsForAllElevators[system.OwnId][floor] != elevatorHardware.Completed && completetonoorder {
		completetonoorder = false
	}

	if CabRequestsForAllElevators[system.OwnId][floor] == elevatorHardware.Assigned && elevatorHardware.Elevator_floorSensor() == floor {
		assignedtocomplete = true
	}

	if pendingtoassigned {
		return pendingToAssigned
	} else if pendingtonoorder {
		return pendingToNoOrder
	} else if assignedtocomplete {
		return assignedToComplete
	} else if completetonoorder {
		return completeToNoOrder
	}
	return noTransition
}

// has to be made ... similar to the one above just checking cab requests instead of hall requests, and only for own elevator, since cab requests are not shared between elevators.
// might now have to have to functions, can probably just pass in cabrequests instead and handle the up and down stuff differently (just dont care if cab requests)
func set_OrderStatus(system *elevatorHardware.ElevatorSystem, floor int, halldir int, status elevatorHardware.OrderStatus) {
	system.HallRequests[floor][halldir] = status
}

/*
func OrderIsPendingForAllElevators(HallRequestsForAllIds map[string][N_FLOORS][2]orderStatus, halldir int, floor int, alivePeers []string) bool {
	for _, peerId := range alivePeers {
		if HallRequestsForAllIds[peerId][floor][halldir] != pending {
			return false
		}
	}
	return true
}

func OrderIsCompletedForAnotherElevators(HallRequestsForAllIds map[string][N_FLOORS][2]orderStatus, halldir int, floor int, alivePeers []string) bool {
	for _, peerId := range alivePeers {
		if peerId != system.OwnId {
			if HallRequestsForAllIds[peerId][floor][halldir] == completed {
				return true
			}
		}
	}
	return false
}

func OrderIsNoOrderForAllOtherElevators(HallRequestsForAllIds map[string][N_FLOORS][2]orderStatus, halldir int, floor int, alivePeers []string) bool {
	for _, peerId := range alivePeers {
		if peerId != system.OwnId {
			if HallRequestsForAllIds[peerId][floor][halldir] != noOrder {
				return false
			}
		}
	}
	return true
}

*/
/*
Transitions:
pending -> assigned is dealt with through assigner
pending -> no order, when it sees one that has complete.
assigned -> complete is dealt with through the elevator itself, since it is the one that completes the order, so it can set the order to complete when it completes it.
complete -> no order, here we need all other elevators to go to no order when they see one that is complete. And the one with complete can go to noOrder when all other have no order instead of pending.

Probably need one that checks for
*/
