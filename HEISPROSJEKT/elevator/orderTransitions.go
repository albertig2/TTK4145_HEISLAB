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

func CheckOrderTransitionStatusForElevators(
	HallRequestsForAllIds map[int][N_FLOORS][2]orderStatus,
	halldir int,
	floor int,
	alivePeers []int) OrderTransition {

	pendingtoassigned := true
	pendingtonoorder := false
	completetonoorder := true
	assignedtocomplete := false

	for _, peerId := range alivePeers {
		if HallRequestsForAllIds[peerId][floor][halldir] != pending {
			pendingtoassigned = false
		}
		if HallRequestsForAllIds[peerId][floor][halldir] == completed && peerId != ownId {
			pendingtonoorder = true
		}
		if HallRequestsForAllIds[peerId][floor][halldir] != noOrder && peerId != ownId {
			completetonoorder = false
		}
	}

	if HallRequestsForAllIds[ownId][floor][halldir] != completed && completetonoorder {
		completetonoorder = false
	}

	if HallRequestsForAllIds[ownId][floor][halldir] == assigned && elevator_floorSensor() == floor {
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
