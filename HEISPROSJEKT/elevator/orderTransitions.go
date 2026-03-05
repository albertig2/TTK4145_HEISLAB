package main

var ownId int = 1

type OrderTransition int

const (
	noTransition OrderTransition = iota
	pendingToAssigned
	pendingToNoOrder
	assignedToComplete
	completeToNoOrder
)

func CheckOrderTransitionStatusForHallRequests(
	HallRequestsForAllElevators map[int][N_FLOORS][2]orderStatus,
	halldir int,
	floor int,
	alivePeers []int) OrderTransition {

	pendingtoassigned := true
	pendingtonoorder := false
	completetonoorder := true
	assignedtocomplete := false

	for _, peerId := range alivePeers {
		if HallRequestsForAllElevators[peerId][floor][halldir] != pending {
			pendingtoassigned = false
		}
		if HallRequestsForAllElevators[peerId][floor][halldir] == completed && peerId != ownId {
			pendingtonoorder = true
		}
		if HallRequestsForAllElevators[peerId][floor][halldir] != noOrder && peerId != ownId {
			completetonoorder = false
		}
	}

	if HallRequestsForAllElevators[ownId][floor][halldir] != completed && completetonoorder {
		completetonoorder = false
	}

	if HallRequestsForAllElevators[ownId][floor][halldir] == assigned && elevator_floorSensor() == floor {
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
	CabRequestsForAllElevators map[int][N_FLOORS]orderStatus,
	floor int,
	alivePeers []int) OrderTransition {

	pendingtoassigned := true
	pendingtonoorder := false
	completetonoorder := true
	assignedtocomplete := false

	for _, peerId := range alivePeers {
		if CabRequestsForAllElevators[peerId][floor] != pending {
			pendingtoassigned = false
		}
		if CabRequestsForAllElevators[peerId][floor] == completed && peerId != ownId {
			pendingtonoorder = true
		}
		if CabRequestsForAllElevators[peerId][floor] != noOrder && peerId != ownId {
			completetonoorder = false
		}
	}

	if CabRequestsForAllElevators[ownId][floor] != completed && completetonoorder {
		completetonoorder = false
	}

	if CabRequestsForAllElevators[ownId][floor] == assigned && elevator_floorSensor() == floor {
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
func set_OrderStatus(system *ElevatorSystem, floor int, halldir int, status orderStatus) {
	system.HallRequests[floor][halldir] = status
}

/*
func OrderIsPendingForAllElevators(HallRequestsForAllIds map[int][N_FLOORS][2]orderStatus, halldir int, floor int, alivePeers []int) bool {
	for _, peerId := range alivePeers {
		if HallRequestsForAllIds[peerId][floor][halldir] != pending {
			return false
		}
	}
	return true
}

func OrderIsCompletedForAnotherElevators(HallRequestsForAllIds map[int][N_FLOORS][2]orderStatus, halldir int, floor int, alivePeers []int) bool {
	for _, peerId := range alivePeers {
		if peerId != ownId {
			if HallRequestsForAllIds[peerId][floor][halldir] == completed {
				return true
			}
		}
	}
	return false
}

func OrderIsNoOrderForAllOtherElevators(HallRequestsForAllIds map[int][N_FLOORS][2]orderStatus, halldir int, floor int, alivePeers []int) bool {
	for _, peerId := range alivePeers {
		if peerId != ownId {
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
