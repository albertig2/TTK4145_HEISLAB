package main

import (
	"testing"
)

// go test -v -run TestInitialize
func TestInitialize(t *testing.T) {
	// Sholuld always update these maps with own state before checking for transitions or doing anything based on the status of the floors
	HallRequestsForAllElevators := make(map[string][N_FLOORS][2]orderStatus)
	CabRequestsForAllElevators := make(map[string][N_FLOORS]orderStatus)

	system1 := ElevatorSystem{}
	initialize(&system1, "1")
	setFloor(&system1, 3)
	setDirection(&system1, stop)
	setBehavior(&system1, doorOpen)
	setHallRequests(&system1, 2, hallDown, pending)

	system2 := ElevatorSystem{}
	initialize(&system2, "2")
	setFloor(&system2, 0)
	setDirection(&system2, stop)
	setBehavior(&system2, idle)
	setHallRequests(&system2, 3, hallUp, pending)

	system3 := ElevatorSystem{}
	initialize(&system3, "3")
	setFloor(&system3, 1)
	setDirection(&system3, up)
	setBehavior(&system3, moving)
	setCabRequests(&system3, 1, pending)
	setHallRequests(&system3, 3, hallUp, pending)

	updateElevatorSystemFromPeer(&system1, &system2, HallRequestsForAllElevators, CabRequestsForAllElevators)
	updateElevatorSystemFromPeer(&system1, &system3, HallRequestsForAllElevators, CabRequestsForAllElevators)

	HallRequestsForAllElevators["1"] = system1.HallRequests           // Always update these before checking for transitions or doing anything based on the status of the floors
	CabRequestsForAllElevators["1"] = system1.States["1"].CabRequests // Always update these before checking for transitions or doing anything based on the status of the floors

	t.Logf("HallRequest for elevator 1: %v", HallRequestsForAllElevators["1"])
	orderTransition := CheckOrderTransitionStatusForHallRequests(&system1, HallRequestsForAllElevators, hallUp, 3, []string{"1", "2", "3"})
	t.Logf("orderTransition: %v", orderTransition)
	setHallRequests(&system1, 3, hallUp, pending)
	HallRequestsForAllElevators["1"] = system1.HallRequests // Remember to update
	orderTransition2 := CheckOrderTransitionStatusForHallRequests(&system1, HallRequestsForAllElevators, hallUp, 3, []string{"1", "2", "3"})
	t.Logf("orderTransition2: %v", orderTransition2)

	hallRequestTransitions := GetAllHallRequestTransitions(&system1, HallRequestsForAllElevators, []string{"1", "2", "3"})
	hallRequestAssigner(&system1, hallRequestTransitions, []string{"1", "2", "3"})
	t.Logf("HallRequest for system1 before assigning: %v", system1.HallRequests[3][hallUp])
	TransitionForHallRequestsByType(&system1, hallRequestTransitions, pendingToAssigned, []string{"1", "2", "3"})
	t.Logf("HallRequest for system1 after assigning: %v", system1.HallRequests[3][hallUp])

	t.Logf("system json: %s", encodeElevatorSystem(&system1))
	transition := CheckOrderTransitionStatusForHallRequests(&system1, HallRequestsForAllElevators, hallUp, 2, []string{"1", "2", "3"})
	t.Logf("transition: %v", transition)
	t.Logf("HallRequests: %v", HallRequestsForAllElevators)
	t.Logf("CabRequests: %v", CabRequestsForAllElevators)
	//message := encodeElevatorSystem(&system)
	//fmt.Print(message)
	//system2 := decodeElevatorSystem(message)
	//initialize(&system2, 4)
	//setCabRequests(&system2, 4, 1, true)

	//t.Logf("system2 json: %s", toJsonString(system2))

}
