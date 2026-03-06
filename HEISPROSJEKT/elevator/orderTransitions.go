package main

type OrderTransition int

const (
	noTransition OrderTransition = iota
	noOrderToPending
	pendingToAssigned
	pendingToNoOrder
	assignedToComplete
	assignedToNoOrder
	completeToNoOrder
)

func CheckOrderTransitionStatusForHallRequests(
	system *ElevatorSystem,
	HallRequestsForAllElevators map[string][N_FLOORS][2]orderStatus,
	halldir int,
	floor int,
	alivePeers []string) OrderTransition {

	noordertopending := false
	pendingtonoorder := false
	pendingtoassigned := true
	assignedtocomplete := false
	completetonoorder := true

	ownHallStatus := HallRequestsForAllElevators[system.OwnId][floor][halldir]

	if ownHallStatus == noOrder {
		for _, peerId := range alivePeers {
			peerHallStatus := HallRequestsForAllElevators[peerId][floor][halldir]
			if peerHallStatus == completed {
				noordertopending = false
				break
			} else if peerHallStatus != noOrder {
				noordertopending = true
			}
		}
	}

	for _, peerId := range alivePeers {
		peerHallStatus := HallRequestsForAllElevators[peerId][floor][halldir]
		if peerHallStatus != pending {
			pendingtoassigned = false
		}
		if peerHallStatus == completed && peerId != system.OwnId {
			pendingtonoorder = true
		}
		if peerHallStatus != noOrder && peerId != system.OwnId {
			completetonoorder = false
		}
	}

	if ownHallStatus != completed && completetonoorder {
		completetonoorder = false
	}

	if ownHallStatus == assigned && elevator_floorSensor() == floor {
		assignedtocomplete = true
	}

	if noordertopending {
		return noOrderToPending
	} else if pendingtoassigned {
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

func GetAllHallRequestTransitions(system *ElevatorSystem, HallRequestsForAllElevators map[string][N_FLOORS][2]orderStatus, alivePeers []string) [N_FLOORS][2]OrderTransition {
	var transitions [N_FLOORS][2]OrderTransition
	for floor := range N_FLOORS {
		for _, halldir := range HallDirs {
			transition := CheckOrderTransitionStatusForHallRequests(system, HallRequestsForAllElevators, halldir, floor, alivePeers)
			transitions[floor][halldir] = transition
		}
	}
	return transitions
}

// Dont know if this should be different or what
// Cab order goes from no order to pending when you press the button
// Cab order goes from pending to assigned when all elevators have pending
// Cab order never goes from pending to no order for own orders
// Cab orders goes from assigned to no order when the elevator reaches the floor
func CheckOrderTransitionStatusForCabRequests(
	system *ElevatorSystem,
	CabRequestsForAllElevators map[string][N_FLOORS]orderStatus,
	floor int,
	alivePeers []string) OrderTransition {

	noordertopending := false
	pendingtoassigned := true
	assignedtonoorder := false

	/*
		if elevator_requestButton(floor, B_Cab) {
			noordertopending = true
		}
	*/

	for _, peerId := range alivePeers {
		peerCabStatus := CabRequestsForAllElevators[peerId][floor]
		if peerCabStatus != pending {
			pendingtoassigned = false
		}
	}

	ownCabStatus := CabRequestsForAllElevators[system.OwnId][floor]
	if ownCabStatus == assigned && elevator_floorSensor() == floor {
		assignedtonoorder = true
	}

	if noordertopending {
		return noOrderToPending
	} else if pendingtoassigned {
		return pendingToAssigned
	} else if assignedtonoorder {
		return assignedToNoOrder
	}
	return noTransition
}

func GetAllCabRequestTransitions(system *ElevatorSystem, CabRequestsForAllElevators map[string][N_FLOORS]orderStatus, alivePeers []string) [N_FLOORS]OrderTransition {
	var transitions [N_FLOORS]OrderTransition
	for floor := range N_FLOORS {
		transition := CheckOrderTransitionStatusForCabRequests(system, CabRequestsForAllElevators, floor, alivePeers)
		transitions[floor] = transition
	}
	return transitions
}

// Should update HallRequestsForAllElevators after this
func TransitionForHallRequestsByType(system *ElevatorSystem, hallRequestTransitions [N_FLOORS][2]OrderTransition, transitionType OrderTransition, alivePeers []string) {
	if transitionType == pendingToAssigned {
		transitionFromPendingToAssignedForHallRequests(system, hallRequestTransitions, alivePeers)
	} else {
		for floor := range N_FLOORS {
			for _, hallDir := range HallDirs {
				if hallRequestTransitions[floor][hallDir] == transitionType {
					var status orderStatus
					switch transitionType {
					case noOrderToPending:
						status = pending
					case pendingToNoOrder:
						status = noOrder
					case assignedToComplete:
						status = completed
						// Turn off lights and stuff here
					case completeToNoOrder:
						status = noOrder
					}
					setHallRequests(system, floor, hallDir, status)
				}
			}
		}
	}
}

// Called from within the TransitionForHallRequestsByType. Setting to private to avoid it being called directly
func transitionFromPendingToAssignedForHallRequests(system *ElevatorSystem, hallRequestTransitions [N_FLOORS][2]OrderTransition, alivePeers []string) {
	output := hallRequestAssigner(system, hallRequestTransitions, alivePeers)
	for floor := range N_FLOORS {
		for _, hallDir := range HallDirs {
			if (output)[system.OwnId][floor][hallDir] {
				setHallRequests(system, floor, hallDir, assigned)
				// Set lights and stuff here
			}
		}
	}
}

func TransitionForCabRequestsByType(system *ElevatorSystem, cabRequestTransitions [N_FLOORS]OrderTransition, transitionType OrderTransition, alivePeers []string) {
	for floor := range N_FLOORS {
		if cabRequestTransitions[floor] == transitionType {
			var status orderStatus
			switch transitionType {
			case noOrderToPending:
				status = pending
			case pendingToAssigned:
				status = assigned
				// Turn on lights and stuff here
			case assignedToNoOrder:
				status = noOrder
				// Turn off lights and stuff here
			}
			setCabRequests(system, floor, status)
		}
	}
}

/*
func transitionHallRequests(system ElevatorSystem, HallRequestForAllElevators *map[string][N_FLOORS][2]orderStatus, floor int, halldir int, transition OrderTransition) {
	var status orderStatus
	switch transition {
	case noTransition:
		status = system.HallRequests[floor][halldir]
	case noOrderToPending:
		status = pending
		//Do work
	case pendingToAssigned:
		//Do work
		status = system.HallRequests[floor][halldir] // NOt necessarily taking status assigned for hall, but for cab yes
		//Do work
	case pendingToNoOrder:
		status = noOrder
		//Do work
	case assignedToComplete:
		status = completed
		//Do work
	case completeToNoOrder:
		status = noOrder
		//Do work
	}

	arr := (*HallRequestForAllElevators)[system.OwnId]
	arr[floor][halldir] = status
	(*HallRequestForAllElevators)[system.OwnId] = arr
	setHallRequests(&system, floor, halldir, status)
}
*/
//Need to make one for cab also....

// has to be made ... similar to the one above just checking cab requests instead of hall requests, and only for own elevator, since cab requests are not shared between elevators.
// might now have to have to functions, can probably just pass in cabrequests instead and handle the up and down stuff differently (just dont care if cab requests)
func set_OrderStatus(system *ElevatorSystem, floor int, halldir int, status orderStatus) {
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
