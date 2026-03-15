package orderProtocol

import (
	"HEISPROSJEKT/elevatorConfig"
	"fmt"
	//"HEISPROSJEKT/elevatorHardware"
	"HEISPROSJEKT/synchronisation"
)

func InitializeOrderChannels() elevatorConfig.ElevatorOrderChannelStruckt {

	orderChannelse := elevatorConfig.ElevatorOrderChannelStruckt{

		NewRecievedOrderChannel: make(chan elevatorConfig.ButtonEvent),
		NewAssignedOrderChannel: make(chan elevatorConfig.ButtonEvent),
		ServicedOrderChannel:    make(chan elevatorConfig.ButtonEvent, 10),
	}

	return orderChannelse
}

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
	system *elevatorConfig.ElevatorSystem,
	HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus,
	halldir int,
	floor int,
	newOrders []elevatorConfig.ButtonEvent,
	servicedOrders []elevatorConfig.ButtonEvent) OrderTransition {

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
		for _, order := range servicedOrders {
			if order.Floor == floor && int(order.Button) == halldir {
				return AssignedToComplete
			}
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
func CheckOrderTransitionStatusForCabRequests(
	system *elevatorConfig.ElevatorSystem,
	CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus,
	floor int,
	newOrders []elevatorConfig.ButtonEvent,
	servicedOrders []elevatorConfig.ButtonEvent) OrderTransition {

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
		for _, order := range servicedOrders {
			if order.Floor == floor {
				return AssignedToNoOrder
			}
		}
	}

	return NoTransition
}

func GetAllHallRequestTransitions(system *elevatorConfig.ElevatorSystem, HallRequestsForAllElevators map[string][elevatorConfig.N_FLOORS][2]elevatorConfig.OrderStatus, newOrders []elevatorConfig.ButtonEvent, servicedOrders []elevatorConfig.ButtonEvent) [elevatorConfig.N_FLOORS][2]OrderTransition {
	var transitions [elevatorConfig.N_FLOORS][2]OrderTransition
	for floor := range elevatorConfig.N_FLOORS {
		for _, halldir := range synchronisation.HallDirections {
			transition := CheckOrderTransitionStatusForHallRequests(system, HallRequestsForAllElevators, halldir, floor, newOrders, servicedOrders)
			transitions[floor][halldir] = transition
		}
	}
	return transitions
}

func GetAllCabRequestTransitions(system *elevatorConfig.ElevatorSystem, CabRequestsForAllElevators map[string][elevatorConfig.N_FLOORS]elevatorConfig.OrderStatus, newOrders []elevatorConfig.ButtonEvent, servicedOrders []elevatorConfig.ButtonEvent) [elevatorConfig.N_FLOORS]OrderTransition {
	var transitions [elevatorConfig.N_FLOORS]OrderTransition
	for floor := range elevatorConfig.N_FLOORS {
		transition := CheckOrderTransitionStatusForCabRequests(system, CabRequestsForAllElevators, floor, newOrders, servicedOrders)
		transitions[floor] = transition
	}
	return transitions
}

func TransitioningAllHallRequests(system *elevatorConfig.ElevatorSystem, hallRequestTransitions [elevatorConfig.N_FLOORS][2]OrderTransition, elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt) {
	for floor := range elevatorConfig.N_FLOORS {
		for _, hallDir := range synchronisation.HallDirections {
			var status elevatorConfig.OrderStatus
			switch hallRequestTransitions[floor][hallDir] {
			case NoOrderToPending:
				status = elevatorConfig.Pending
			case PendingToNoOrder:
				status = elevatorConfig.NoOrder
			case AssignedToComplete:
				status = elevatorConfig.Completed
				//elevatorOrderChannels.ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: floor, Button: elevatorConfig.Button(hallDir)}
				// Turn off lights and stuff here
			case CompleteToNoOrder:
				status = elevatorConfig.NoOrder
			default:
				status = system.HallRequests[floor][hallDir]
			}
			synchronisation.SetHallRequests(system, floor, hallDir, status)
		}
	}
	transitionFromPendingToAssignedForHallRequests(system, hallRequestTransitions, elevatorOrderChannels)
}

// Called from within the TransitionForHallRequestsByType. Setting to private to avoid it being called directly
func transitionFromPendingToAssignedForHallRequests(system *elevatorConfig.ElevatorSystem, hallRequestTransitions [elevatorConfig.N_FLOORS][2]OrderTransition, elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt) {
	output := HallRequestAssigner(system, hallRequestTransitions)
	for floor := range elevatorConfig.N_FLOORS {
		for _, hallDir := range synchronisation.HallDirections {
			if (output)[system.OwnId][floor][hallDir] {
				synchronisation.SetHallRequests(system, floor, hallDir, elevatorConfig.Assigned)
				elevatorOrderChannels.NewAssignedOrderChannel <- elevatorConfig.ButtonEvent{Floor: floor, Button: elevatorConfig.Button(hallDir)}
				// Set lights and stuff here
			}
		}
	}
}

func TransitioningAllCabRequests(system *elevatorConfig.ElevatorSystem, cabRequestTransitions [elevatorConfig.N_FLOORS]OrderTransition, elevatorOrderChannels elevatorConfig.ElevatorOrderChannelStruckt) {
	for floor := range elevatorConfig.N_FLOORS {
		var status elevatorConfig.OrderStatus
		switch cabRequestTransitions[floor] {
		case NoOrderToPending, UnknownToPending:
			status = elevatorConfig.Pending
		case PendingToAssigned, UnknownToAssigned:
			// Turn on lights and stuff here
			status = elevatorConfig.Assigned

			elevatorOrderChannels.NewAssignedOrderChannel <- elevatorConfig.ButtonEvent{Floor: floor, Button: elevatorConfig.Cab}
			fmt.Println("New cab order assigned on floor ", floor)
		case AssignedToNoOrder, UnknownToNoOrder:
			status = elevatorConfig.NoOrder

			//elevatorOrderChannels.ServicedOrderChannel <- elevatorConfig.ButtonEvent{Floor: floor, Button: elevatorConfig.Cab}
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
