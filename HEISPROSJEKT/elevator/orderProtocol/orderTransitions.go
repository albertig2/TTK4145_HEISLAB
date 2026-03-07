package orderProtocol

import (
	"HEISPROSJEKT/elevatorConfig"
	"HEISPROSJEKT/elevatorHardware"
)

type OrderTransition int

const (
	NoTransition OrderTransition = iota
	NoOrderToPending
	PendingToAssigned
	PendingToNoOrder
	AssignedToComplete
	AssignedToNoOrder
	CompleteToNoOrder
)

func CheckOrderTransitionStatusForHallRequests(
	system *elevatorHardware.ElevatorSystem,
	HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorHardware.OrderStatus,
	halldir int,
	floor int,
	alivePeers []string) OrderTransition {

	noordertopending := false
	pendingtonoorder := false
	pendingtoassigned := true
	assignedtocomplete := false
	completetonoorder := true

	ownHallStatus := HallRequestsForAllElevators[system.OwnId][floor][halldir]

	if ownHallStatus == elevatorHardware.NoOrder {
		for _, peerId := range alivePeers {
			peerHallStatus := HallRequestsForAllElevators[peerId][floor][halldir]
			if peerHallStatus == elevatorHardware.Completed {
				noordertopending = false
				break
			} else if peerHallStatus != elevatorHardware.NoOrder {
				noordertopending = true
			}
		}
	}

	for _, peerId := range alivePeers {
		peerHallStatus := HallRequestsForAllElevators[peerId][floor][halldir]
		if peerHallStatus != elevatorHardware.Pending {
			pendingtoassigned = false
		}
		if peerHallStatus == elevatorHardware.Completed && peerId != system.OwnId {
			pendingtonoorder = true
		}
		if peerHallStatus != elevatorHardware.NoOrder && peerId != system.OwnId {
			completetonoorder = false
		}
	}

	if ownHallStatus != elevatorHardware.Completed && completetonoorder {
		completetonoorder = false
	}

	if ownHallStatus == elevatorHardware.Assigned { //&& elevator_floorSensor() == floor
		assignedtocomplete = true
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

func GetAllHallRequestTransitions(system *elevatorHardware.ElevatorSystem, HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorHardware.OrderStatus, alivePeers []string) [elevatorConfig.N_FLOORS][2]OrderTransition {
	var transitions [elevatorConfig.N_FLOORS][2]OrderTransition
	for floor := range elevatorConfig.N_FLOORS {
		for _, halldir := range elevatorHardware.HallDirs {
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
	system *elevatorHardware.ElevatorSystem,
	CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorHardware.OrderStatus,
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
		if peerCabStatus != elevatorHardware.Pending {
			pendingtoassigned = false
		}
	}

	ownCabStatus := CabRequestsForAllElevators[system.OwnId][floor]
	if ownCabStatus == elevatorHardware.Assigned { // && elevator_floorSensor() == floor
		assignedtonoorder = true
	}

	if noordertopending {
		return NoOrderToPending
	} else if pendingtoassigned {
		return PendingToAssigned
	} else if assignedtonoorder {
		return AssignedToNoOrder
	}
	return NoTransition
}

func GetAllCabRequestTransitions(system *elevatorHardware.ElevatorSystem, CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorHardware.OrderStatus, alivePeers []string) [elevatorConfig.N_FLOORS]OrderTransition {
	var transitions [elevatorConfig.N_FLOORS]OrderTransition
	for floor := range elevatorConfig.N_FLOORS {
		transition := CheckOrderTransitionStatusForCabRequests(system, CabRequestsForAllElevators, floor, alivePeers)
		transitions[floor] = transition
	}
	return transitions
}

// Should update HallRequestsForAllElevators after this
func TransitionForHallRequestsByType(system *elevatorHardware.ElevatorSystem, hallRequestTransitions [elevatorConfig.N_FLOORS][2]OrderTransition, transitionType OrderTransition, alivePeers []string) {
	if transitionType == PendingToAssigned {
		transitionFromPendingToAssignedForHallRequests(system, hallRequestTransitions, alivePeers)
	} else {
		for floor := range elevatorConfig.N_FLOORS {
			for _, hallDir := range elevatorHardware.HallDirs {
				if hallRequestTransitions[floor][hallDir] == transitionType {
					var status elevatorHardware.OrderStatus
					switch transitionType {
					case NoOrderToPending:
						status = elevatorHardware.Pending
					case PendingToNoOrder:
						status = elevatorHardware.NoOrder
					case AssignedToComplete:
						status = elevatorHardware.Completed
						// Turn off lights and stuff here
					case CompleteToNoOrder:
						status = elevatorHardware.NoOrder
					default:
						status = system.HallRequests[floor][hallDir]
					}
					elevatorHardware.SetHallRequests(system, floor, hallDir, status)
				}
			}
		}
	}
}

// Called from within the TransitionForHallRequestsByType. Setting to private to avoid it being called directly
func transitionFromPendingToAssignedForHallRequests(system *elevatorHardware.ElevatorSystem, hallRequestTransitions [elevatorConfig.N_FLOORS][2]OrderTransition, alivePeers []string) {
	output := HallRequestAssigner(system, hallRequestTransitions, alivePeers)
	for floor := range elevatorConfig.N_FLOORS {
		for _, hallDir := range elevatorHardware.HallDirs {
			if (output)[system.OwnId][floor][hallDir] {
				elevatorHardware.SetHallRequests(system, floor, hallDir, elevatorHardware.Assigned)
				// Set lights and stuff here
			}
		}
	}
}

func TransitionForCabRequestsByType(system *elevatorHardware.ElevatorSystem, cabRequestTransitions [elevatorConfig.N_FLOORS]OrderTransition, transitionType OrderTransition, alivePeers []string) {
	for floor := range elevatorConfig.N_FLOORS {
		if cabRequestTransitions[floor] == transitionType {
			var status elevatorHardware.OrderStatus
			switch transitionType {
			case NoOrderToPending:
				status = elevatorHardware.Pending
			case PendingToAssigned:
				status = elevatorHardware.Assigned
				// Turn on lights and stuff here
			case AssignedToNoOrder:
				status = elevatorHardware.NoOrder
				// Turn off lights and stuff here
			default:
				status = system.States[system.OwnId].CabRequests[floor]
			}
			elevatorHardware.SetCabRequests(system, floor, status)
		}
	}
}

/*
func transitionHallRequests(system ElevatorSystem, HallRequestForAllElevators *map[string][elevatorConfig.N_FLOORS][2]orderStatus, floor int, halldir int, transition OrderTransition) {
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

// Not using so far but could

func set_OrderStatus(system *elevatorHardware.ElevatorSystem, floor int, halldir int, status elevatorHardware.OrderStatus) {
	system.HallRequests[floor][halldir] = status
}

/*
func OrderIsPendingForAllElevators(HallRequestsForAllIds map[string][elevatorConfig.N_FLOORS][2]orderStatus, halldir int, floor int, alivePeers []string) bool {
	for _, peerId := range alivePeers {
		if HallRequestsForAllIds[peerId][floor][halldir] != pending {
			return false
		}
	}
	return true
}

func OrderIsCompletedForAnotherElevators(HallRequestsForAllIds map[string][elevatorConfig.N_FLOORS][2]orderStatus, halldir int, floor int, alivePeers []string) bool {
	for _, peerId := range alivePeers {
		if peerId != system.OwnId {
			if HallRequestsForAllIds[peerId][floor][halldir] == completed {
				return true
			}
		}
	}
	return false
}

func OrderIsNoOrderForAllOtherElevators(HallRequestsForAllIds map[string][elevatorConfig.N_FLOORS][2]orderStatus, halldir int, floor int, alivePeers []string) bool {
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
