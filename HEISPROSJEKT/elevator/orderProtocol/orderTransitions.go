package orderProtocol

import (
	"HEISPROSJEKT/elevatorConfig"
	//"HEISPROSJEKT/elevatorHardware"
	"HEISPROSJEKT/synchronisation"
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
	UnknownToPending
	UnknownToNoOrder
	UnknownToAssigned
)

// Not using these at the moment
var cabTransitions = []OrderTransition{
	NoOrderToPending,
	PendingToAssigned,
	AssignedToNoOrder,
	UnknownToPending,
	UnknownToAssigned,
	UnknownToNoOrder,
}

var hallTransitions = []OrderTransition{
	NoOrderToPending,
	PendingToAssigned,
	PendingToNoOrder,
	AssignedToComplete,
	CompleteToNoOrder,
}

func CheckOrderTransitionStatusForHallRequests(
	system *synchronisation.ElevatorSystem,
	HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]synchronisation.OrderStatus,
	halldir int,
	floor int,
	alivePeers []string) OrderTransition {

	noordertopending := false
	pendingtonoorder := false
	pendingtoassigned := true
	assignedtocomplete := false
	completetonoorder := true

	ownHallStatus := HallRequestsForAllElevators[system.OwnId][floor][halldir]

	if ownHallStatus == synchronisation.NoOrder {
		for _, peerId := range alivePeers {
			peerHallStatus := HallRequestsForAllElevators[peerId][floor][halldir]
			if peerHallStatus == synchronisation.Completed {
				noordertopending = false
				break
			} else if peerHallStatus != synchronisation.NoOrder {
				noordertopending = true
			}
		}
	}

	for _, peerId := range alivePeers {
		peerHallStatus := HallRequestsForAllElevators[peerId][floor][halldir]
		if peerHallStatus != synchronisation.Pending {
			pendingtoassigned = false
		}
		if peerHallStatus == synchronisation.Completed && peerId != system.OwnId {
			pendingtonoorder = true
		}
		if peerHallStatus != synchronisation.NoOrder && peerId != system.OwnId {
			completetonoorder = false
		}
	}

	if ownHallStatus != synchronisation.Completed && completetonoorder {
		completetonoorder = false
	}

	if ownHallStatus == synchronisation.Assigned { //&& elevator_floorSensor() == floor
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

// Dont know if this should be different or what
// Cab order goes from no order to pending when you press the button
// Cab order goes from pending to assigned when all elevators have pending
// Cab order never goes from pending to no order for own orders
// Cab orders goes from assigned to no order when the elevator reaches the floor
func CheckOrderTransitionStatusForCabRequests(
	system *synchronisation.ElevatorSystem,
	CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]synchronisation.OrderStatus,
	floor int,
	alivePeers []string) OrderTransition {

	noordertopending := false
	pendingtoassigned := true
	assignedtonoorder := false
	unknowntonoorder := false
	unknowntopending := false
	unknowntoassigned := false

	if CabRequestsForAllElevators[system.OwnId][floor] == synchronisation.Unknown {
		allNoOrder := true
		allAssignedOrPending := true
		anyPending := false

		for _, peerId := range alivePeers {
			if peerId != system.OwnId {
				peerCabStatus := CabRequestsForAllElevators[peerId][floor]
				if peerCabStatus == synchronisation.Pending {
					anyPending = true
				}
				if peerCabStatus != synchronisation.NoOrder {
					allNoOrder = false
				}
				if peerCabStatus != synchronisation.Assigned && peerCabStatus != synchronisation.Pending {
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
		/*
			if elevator_requestButton(floor, B_Cab) {
				unknowntopending = true
			}
		*/
	}
	/*
		if elevator_requestButton(floor, B_Cab) {
			noordertopending = true
		}
	*/

	for _, peerId := range alivePeers {
		peerCabStatus := CabRequestsForAllElevators[peerId][floor]
		if peerCabStatus != synchronisation.Pending {
			pendingtoassigned = false
		}
	}

	ownCabStatus := CabRequestsForAllElevators[system.OwnId][floor]
	if ownCabStatus == synchronisation.Assigned { // && elevator_floorSensor() == floor
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

func GetAllHallRequestTransitions(system *synchronisation.ElevatorSystem, HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]synchronisation.OrderStatus, alivePeers []string) [elevatorConfig.N_FLOORS][2]OrderTransition {
	var transitions [elevatorConfig.N_FLOORS][2]OrderTransition
	for floor := range elevatorConfig.N_FLOORS {
		for _, halldir := range synchronisation.HallDirections {
			transition := CheckOrderTransitionStatusForHallRequests(system, HallRequestsForAllElevators, halldir, floor, alivePeers)
			transitions[floor][halldir] = transition
		}
	}
	return transitions
}

func GetAllCabRequestTransitions(system *synchronisation.ElevatorSystem, CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]synchronisation.OrderStatus, alivePeers []string) [elevatorConfig.N_FLOORS]OrderTransition {
	var transitions [elevatorConfig.N_FLOORS]OrderTransition
	for floor := range elevatorConfig.N_FLOORS {
		transition := CheckOrderTransitionStatusForCabRequests(system, CabRequestsForAllElevators, floor, alivePeers)
		transitions[floor] = transition
	}
	return transitions
}

func TransitioningAllHallRequests(system *synchronisation.ElevatorSystem, hallRequestTransitions [elevatorConfig.N_FLOORS][2]OrderTransition, alivePeers []string, elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt) {
	for floor := range elevatorConfig.N_FLOORS {
		for _, hallDir := range synchronisation.HallDirections {
			var status synchronisation.OrderStatus
			switch hallRequestTransitions[floor][hallDir] {
			case NoOrderToPending:
				status = synchronisation.Pending
			case PendingToNoOrder:
				status = synchronisation.NoOrder
			case AssignedToComplete:
				status = synchronisation.Completed
				elevatorOrderChannels.ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: floor, Button: elevatorConfig.Button(hallDir)}
				// Turn off lights and stuff here
			case CompleteToNoOrder:
				status = synchronisation.NoOrder
			default:
				status = system.HallRequests[floor][hallDir]
			}
			synchronisation.SetHallRequests(system, floor, hallDir, status)
		}
	}
	transitionFromPendingToAssignedForHallRequests(system, hallRequestTransitions, alivePeers, elevatorOrderChannels)
}

// Called from within the TransitionForHallRequestsByType. Setting to private to avoid it being called directly
func transitionFromPendingToAssignedForHallRequests(system *synchronisation.ElevatorSystem, hallRequestTransitions [elevatorConfig.N_FLOORS][2]OrderTransition, alivePeers []string, elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt) {
	output := HallRequestAssigner(system, hallRequestTransitions, alivePeers)
	for floor := range elevatorConfig.N_FLOORS {
		for _, hallDir := range synchronisation.HallDirections {
			if (output)[system.OwnId][floor][hallDir] {
				synchronisation.SetHallRequests(system, floor, hallDir, synchronisation.Assigned)
				elevatorOrderChannels.NewAssignedOrderChannel <- elevatorConfig.ButtonEvent{Floor: floor, Button: elevatorConfig.Button(hallDir)}
				// Set lights and stuff here
			}
		}
	}
}

func TransitioningAllCabRequests(system *synchronisation.ElevatorSystem, cabRequestTransitions [elevatorConfig.N_FLOORS]OrderTransition, elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt) {
	for floor := range elevatorConfig.N_FLOORS {
		var status synchronisation.OrderStatus
		switch cabRequestTransitions[floor] {
		case NoOrderToPending, UnknownToPending:
			status = synchronisation.Pending
		case PendingToAssigned, UnknownToAssigned:
			status = synchronisation.Assigned
			// Turn on lights and stuff here
			elevatorOrderChannels.NewAssignedOrderChannel <- elevatorConfig.ButtonEvent{Floor: floor, Button: elevatorConfig.Cab}
		case AssignedToNoOrder, UnknownToNoOrder:
			status = synchronisation.NoOrder
			elevatorOrderChannels.ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: floor, Button: elevatorConfig.Cab}
			// Turn off lights and stuff here
		default:
			status = system.States[system.OwnId].CabRequests[floor]
		}
		synchronisation.SetCabRequests(system, floor, status)
	}
}

// has to be made ... similar to the one above just checking cab requests instead of hall requests, and only for own elevator, since cab requests are not shared between elevators.
// might now have to have to functions, can probably just pass in cabrequests instead and handle the up and down stuff differently (just dont care if cab requests)
/*
Transitions:
pending -> assigned is dealt with through assigner
pending -> no order, when it sees one that has complete.
assigned -> complete is dealt with through the elevator itself, since it is the one that completes the order, so it can set the order to complete when it completes it.
complete -> no order, here we need all other elevators to go to no order when they see one that is complete. And the one with complete can go to noOrder when all other have no order instead of pending.

Probably need one that checks for
*/
